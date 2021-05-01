OWNER := fabianlee
PROJECT := k8s-downward-golang-web
VERSION := 1.0.0
OPV := $(OWNER)/$(PROJECT):$(VERSION)
WEBPORT := 8080:8080

# you may need to change to "sudo docker" if not a member of 'docker' group
DOCKERCMD := "docker"

BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
# unique id from last git commit
MY_GITREF := $(shell git rev-parse --short HEAD)


## builds docker image
docker-build:
	echo MY_GITREF is $(MY_GITREF)
	$(DOCKERCMD) build --build-arg MY_VERSION=$(VERSION) --build-arg MY_BUILDTIME=$(BUILD_TIME) -f Dockerfile -t $(OPV) .

## cleans docker image
clean:
	$(DOCKERCMD) image rm $(OPV) | true

## runs container in foreground, testing a couple of override values
docker-test-fg:
	$(DOCKERCMD) run -it -p $(WEBPORT) -e APP_CONTEXT=/golang-web/ -e MY_NODE_NAME=node1 --rm $(OPV)

## runs container in foreground, override entrypoint to use use shell
docker-test-cli:
	$(DOCKERCMD) run -it --rm --entrypoint "/bin/sh" $(OPV)

## run container in background
docker-run-bg:
	$(DOCKERCMD) run -d -p $(WEBPORT) --rm --name $(PROJECT) $(OPV)

## get into console of container running in background
docker-cli-bg:
	$(DOCKERCMD) exec -it $(PROJECT) /bin/sh

## tails $(DOCKERCMD)logs
docker-logs:
	$(DOCKERCMD) logs -f $(PROJECT)

## stops container running in background
docker-stop:
	$(DOCKERCMD) stop $(PROJECT)


## pushes to $(DOCKERCMD)hub
docker-push:
	$(DOCKERCMD) push $(OPV)

## pushes to kubernetes cluster
k8s-apply:
	sed -e 's/1.0.0/$(VERSION)/' golang-hello-world-web.yaml | kubectl apply -f -

k8s-delete:
	kubectl delete -f golang-hello-world-web.yaml

