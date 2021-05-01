package main

// https://ashishb.net/tech/docker-101-a-basic-web-server-displaying-hello-world/
// https://tutorialedge.net/golang/creating-simple-web-server-with-golang/
// https://stackoverflow.com/questions/47509272/how-to-set-package-variable-using-ldflags-x-in-golang-build
// https://blog.gopheracademy.com/advent-2017/kubernetes-ready-service/
// https://kubernetes.io/docs/tasks/inject-data-application/environment-variable-expose-pod-information/
// https://semaphoreci.com/community/tutorials/how-to-deploy-a-go-web-application-with-docker


import (
    "fmt"
    "log"
    "net/http"
    "os"
    "sync"
    "io/ioutil"
    "path/filepath"
    "strings"
)

// built into binary using ldflags
var Version string
var BuildTime string

// env keys set by k8s
var k8s_downward_env_list []string = []string { 
  "MY_NODE_NAME","MY_POD_NAME","MY_POD_IP","MY_POD_SERVICE_ACCOUNT",
  "MY_POD_LABEL_APP","MY_POD_ANNOTATION_AUTHOR","MY_POD_MEM_LIMIT_MB",
  "MY_POD_MEM_REQUEST_MB",
} 

// request count for this container
var counter int
var mutex = &sync.Mutex{}
func incrementCounter() {
    mutex.Lock()
    counter++
    mutex.Unlock()
}

func FilePathWalkDir(root string) ([]string, error) {
 var files []string
 err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
   // skip directories and files whose name starts with ".."
   if info.IsDir() || strings.HasPrefix(info.Name(),"..") {
     return nil
   }
   // Downward API files look like symbolic links
   if (info.Mode() & os.ModeSymlink)>0 {
     //log.Printf("found linked file %s",path)
     files = append(files, path )
   }
   return nil
 })
 return files, err
}


func StartWebServer() {
    log.Printf("build version/time: %s/%s", Version, BuildTime)

    // show Downward API env vars
    for _,keyName := range k8s_downward_env_list {
      log.Printf("ENV %s = %s\n", keyName, getenv(keyName,"empty") )
    }

    // show Downward API files
    _, derr := os.Stat("/etc/podinfo/")
    if os.IsNotExist(derr) {
      log.Printf("/etc/podinfo/ does not exist %v",derr)
    }else {
      myfiles, err := FilePathWalkDir("/etc/podinfo/")
      if err != nil {
        log.Printf("problem with FilePathWalkDir %v",err)
      }else {
        for _, file := range myfiles {
          log.Printf("FILE %s",file)
        } // each file
      }
    }

    // request handlers
    http.HandleFunc("/healthz", handleHealth)
    http.HandleFunc("/shutdown", handleShutdown)

    // APP_CONTEXT defaults to root
    appContext := getenv("APP_CONTEXT","/")
    log.Printf("app context: %s", appContext)
    http.HandleFunc(appContext, handleApp)

    port := getenv("PORT","8080")
    log.Printf("Starting web server on port %s", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        panic(err)
    }

}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type","application/json")
    fmt.Fprintf(w, "{\"health\":\"ok\", \"Version\":\"%s\", \"BuildTime\":\"%s\"}", Version, BuildTime )
}

func handleApp(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type","text/plain")

    // print main message with count to both stdout as well as response
    mainMsgFormat := "%d %s %s\n"
    log.Printf(mainMsgFormat, counter, r.Method, r.URL.Path)
    fmt.Fprintf(w, mainMsgFormat, counter, r.Method, r.URL.Path)

    // url path
    fmt.Fprintf(w, "path: %s\n", r.URL.Path)
    // 'Host' header is promoted to Request.Host field and removed from Header map
    fmt.Fprintf(w, "Host: %s\n", provideDefault(r.Host,"empty"))

    // env vars that are populated from kubernetes/docker environment
    for _,keyName := range k8s_downward_env_list {
      fmt.Fprintf(w, "ENV %s = %s\n", keyName, getenv(keyName,"empty") )
    }

    // show Downward API files
    _, derr := os.Stat("/etc/podinfo/")
    if os.IsNotExist(derr) {
      log.Printf("/etc/podinfo/ does not exist %v",derr)
    }else {
      myfiles, err := FilePathWalkDir("/etc/podinfo/")
      if err != nil {
        log.Printf("problem with FilePathWalkDir %v",err)
      }else {
        for _, file := range myfiles {

          data, ferr := ioutil.ReadFile(file)
          if ferr == nil {
            fmt.Fprintf(w,"FILE %s = %s\n",file,data)
          } else {
            log.Printf("error reading file %v",ferr)
          }

        } // each file

      }
    }

    incrementCounter()
}

// provide default for value
func provideDefault(value,defaultVal string) string {
  if len(value)==0 { 
    return defaultVal
  }
  return value
}
// pull from OS environment variable, provide default
func getenv(key, fallback string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    return value
}
// non-graceful and abrupt exit
func handleShutdown(w http.ResponseWriter, r *http.Request) {
    log.Printf("About to abruptly exit")
    os.Exit(0)
}

func main() {
    StartWebServer()
}
