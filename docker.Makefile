#
# github.com/docker/cli
#
# Makefile for developing using Docker
#

DEV_DOCKER_IMAGE_NAME = docker-cli-dev
MOUNTS = -v "$(CURDIR)":/go/src/github.com/docker/cli
VERSION = $(shell cat VERSION)
ENVVARS = -e VERSION=$(VERSION) -e GITCOMMIT

# build docker image (dockerfiles/Dockerfile.build)
.PHONY: build_docker_image
build_docker_image:
	docker build -t $(DEV_DOCKER_IMAGE_NAME) -f ./dockerfiles/Dockerfile.dev .

# build executable using a container
binary: build_docker_image
	docker run --rm $(ENVVARS) $(MOUNTS) $(DEV_DOCKER_IMAGE_NAME) make binary

build: binary


binary_arm64: build_docker_image
	docker run --rm $(ENVVARS) $(MOUNTS) $(DEV_DOCKER_IMAGE_NAME) make binary

build_arm64: binary_arm64

# clean build artifacts using a container
.PHONY: clean
clean: build_docker_image
	docker run --rm $(MOUNTS) $(DEV_DOCKER_IMAGE_NAME) make clean

# run go test
.PHONY: test
test: build_docker_image
	docker run --rm $(MOUNTS) $(DEV_DOCKER_IMAGE_NAME) make test

# start container in interactive mode for in-container development
.PHONY: dev
dev: build_docker_image
	docker run -ti $(MOUNTS) -v /var/run/docker.sock:/var/run/docker.sock $(DEV_DOCKER_IMAGE_NAME) ash

shell: dev

## generate man pages from go source and markdown
.PHONY: manpages
manpages: build_docker_image
	docker run -ti --rm $(MOUNTS) $(DEV_DOCKER_IMAGE_NAME) make manpages

## Generate documentation YAML files consumed by docs repo
.PHONY: yamldocs
yamldocs: build_docker_image
	docker run -ti --rm $(MOUNTS) $(DEV_DOCKER_IMAGE_NAME) make yamldocs
