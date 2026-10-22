#
# github.com/docker/cli
#
all: binary

# remove build artifacts
.PHONY: clean
clean:
	rm -rf ./build/* ./man/man[1-9] docs/yaml/gen

# run go test
# the "-tags daemon" part is temporary
.PHONY: test
test:
	./scripts/test/unit $(shell go list ./... | grep -v '/vendor/')

.PHONY: test-coverage
test-coverage:
	./scripts/test/unit-with-coverage $(shell go list ./... | grep -v '/vendor/')

.PHONY: binary
binary:
	@echo "WARNING: binary creates a Linux executable."
	./scripts/build/binary

## generate man pages from go source and markdown
.PHONY: manpages
manpages:
	scripts/docs/generate-man.sh

## generate documentation YAML files consumed by docs repo
.PHONY: yamldocs
yamldocs:
	scripts/docs/generate-yaml.sh
