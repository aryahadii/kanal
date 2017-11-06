ROOT := gitlab.com/arha/kanal
GO_VARS ?= CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GO ?= go
GIT ?= git
COMMIT := $(shell $(GIT) rev-parse HEAD)
VERSION ?= $(shell $(GIT) describe --tags ${COMMIT} 2> /dev/null || echo "$(COMMIT)")
BUILD_TIME := $(shell LANG=en_US date +"%F_%T_%z")
LD_FLAGS := -X $(ROOT).Version=$(VERSION) -X $(ROOT).Commit=$(COMMIT) -X $(ROOT).BuildTime=$(BUILD_TIME) -X $(ROOT).Title=kanal
DOCKER_IMAGE := registry.gitlab.com/arha/kanal

.PHONY: help clean update-dependencies dependencies docker push


help:
	@echo "Please use \`make <ROOT>' where <ROOT> is one of"
	@echo "  update-dependencies    to update glide.lock (refs to dependencies)"
	@echo "  dependencies           to install the dependencies"
	@echo "  kanald                 to build the binary"
	@echo "  clean                  to remove generated files"

clean:
	rm -rf kanald

update-dependencies:
	glide up

dependencies:
	glide install

kanald: *.go */*.go */*/*.go glide.lock
	$(GO_VARS) $(GO) build -o="kanald" -ldflags="$(LD_FLAGS)" $(ROOT)/cmd/kanal

docker: kanald Dockerfile
	docker build -t $(DOCKER_IMAGE):$(VERSION) .
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest

push:
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):latest
