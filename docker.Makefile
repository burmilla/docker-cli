#
# github.com/docker/cli
#
# Makefile for developing using Docker
#

DEV_DOCKER_IMAGE_NAME = docker-cli-dev
LINTER_IMAGE_NAME = docker-cli-lint
CROSS_IMAGE_NAME = docker-cli-cross
MOUNTS = -v "$(CURDIR)":/go/src/github.com/docker/cli
VERSION = $(shell cat VERSION)
ENVVARS = -e VERSION=$(VERSION) -e GITCOMMIT

# build docker image (dockerfiles/Dockerfile.build)
.PHONY: build_docker_image
build_docker_image:
	docker build -t $(DEV_DOCKER_IMAGE_NAME) -f ./dockerfiles/Dockerfile.dev .

.PHONY: build_docker_image_arm64
build_docker_image_arm64:
	docker build -t $(DEV_DOCKER_IMAGE_NAME) -f ./dockerfiles/Dockerfile.dev.arm64 .

.PHONY: build_docker_image_mips64el
build_docker_image_mips64el:
	docker build -t $(DEV_DOCKER_IMAGE_NAME) -f ./dockerfiles/Dockerfile.dev.mips64el .

# build docker image having the linting tools (dockerfiles/Dockerfile.lint)
.PHONY: build_linter_image
build_linter_image:
	docker build -t $(LINTER_IMAGE_NAME) -f ./dockerfiles/Dockerfile.lint .

.PHONY: build_cross_image
build_cross_image:
	docker build -t $(CROSS_IMAGE_NAME) -f ./dockerfiles/Dockerfile.cross .


# build executable using a container
binary: build_docker_image
	docker run --rm $(ENVVARS) $(MOUNTS) $(DEV_DOCKER_IMAGE_NAME) make binary

build: binary


binary_arm64: build_docker_image_arm64
	docker run --rm $(ENVVARS) $(MOUNTS) $(DEV_DOCKER_IMAGE_NAME) make binary

binary_mips64el: build_docker_image_mips64el
	docker run --rm $(ENVVARS) $(MOUNTS) $(DEV_DOCKER_IMAGE_NAME) make binary

build_arm64: binary_arm64

build_mips64el: binary_mips64el

# clean build artifacts using a container
.PHONY: clean
clean: build_docker_image
	docker run --rm $(MOUNTS) $(DEV_DOCKER_IMAGE_NAME) make clean

# run go test
.PHONY: test
test: build_docker_image
	docker run --rm $(MOUNTS) $(DEV_DOCKER_IMAGE_NAME) make test

# build the CLI for multiple architectures using a container
.PHONY: cross
cross: build_cross_image
	docker run --rm $(ENVVARS) $(MOUNTS) $(CROSS_IMAGE_NAME) make cross

.PHONY: watch
watch: build_docker_image
	docker run --rm $(MOUNTS) $(DEV_DOCKER_IMAGE_NAME) make watch

# start container in interactive mode for in-container development
.PHONY: dev
dev: build_docker_image
	docker run -ti $(MOUNTS) -v /var/run/docker.sock:/var/run/docker.sock $(DEV_DOCKER_IMAGE_NAME) ash

shell: dev

# run linters in a container
.PHONY: lint
lint: build_linter_image
	docker run -ti $(MOUNTS) $(LINTER_IMAGE_NAME)

# download dependencies (vendor/) listed in vendor.conf, using a container
.PHONY: vendor
vendor: build_docker_image vendor.conf
	docker run -ti --rm $(MOUNTS) $(DEV_DOCKER_IMAGE_NAME) make vendor

dynbinary: build_cross_image
	docker run -ti --rm $(ENVVARS) $(MOUNTS) $(CROSS_IMAGE_NAME) make dynbinary

## generate man pages from go source and markdown
.PHONY: manpages
manpages: build_docker_image
	docker run -ti --rm $(MOUNTS) $(DEV_DOCKER_IMAGE_NAME) make manpages

## Generate documentation YAML files consumed by docs repo
.PHONY: yamldocs
yamldocs: build_docker_image
	docker run -ti --rm $(MOUNTS) $(DEV_DOCKER_IMAGE_NAME) make yamldocs
