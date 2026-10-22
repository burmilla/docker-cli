#
# github.com/docker/cli
#
all: binary

# remove build artifacts
.PHONY: clean
clean:
	rm -rf ./build/*

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
