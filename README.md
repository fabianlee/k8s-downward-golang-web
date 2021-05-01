# Summary
Golang http web server running by default on port 8080 that is intended for testing

Image is based on busybox:1.32.1-glibc, is about ~11Mb because it takes advantage of multi-stage building

# Environment variables

* PORT - listen port, defaults to 8080
* APP_CONTEXT - base context path of app, defaults to '/'

# Environment variables populated from Downward API
* MY_NODE_NAME - name of k8s node
* MY_POD_NAME - name of k8s pod
* MY_POD_IP - k8s pod IP
* MY_POD_SERVICE_ACCOUNT - service account of k8s pod

# Prerequisites
* make utility (sudo apt-get install make)

# Makefile targets
* docker-build (builds image)
* docker-run-fg (runs container in foreground, ctrl-C to exit)
* docker-run-bg (runs container in background)
