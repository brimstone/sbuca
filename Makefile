ifndef GOPATH
	GOPATH := ${PWD}/gopath
	export GOPATH
endif

all: sbuca docker-image docker-push

sbuca: *.go
	go get -v -d
	CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' -o sbuca

docker-image:
	docker build -t brimstone/sbuca .

docker-push:
	@[ -f ${HOME}/.dockercfg ] || docker login -e="${DOCKER_EMAIL}" -u="${DOCKER_USERNAME}" -p="${DOCKER_PASSWORD}"
	docker push brimstone/sbuca
