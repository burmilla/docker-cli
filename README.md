# System Docker CLI

The `system-docker` client for BurmillaOS, forked from docker/cli 17.06 and
reduced to the commands the
[System Docker engine](https://github.com/burmilla/docker) still serves:
container, image, network, volume and system commands over a unix socket.
There is no build, push, login, swarm, plugin, stack or content trust.

## Building

```
make -f docker.Makefile binary
```

builds inside `golang:1.27.3-alpine` for the host architecture (amd64 or
arm64) and leaves a static `build/docker`.
Outside the image, with Go 1.27 and `GO111MODULE=off`, the tree must be at
`$GOPATH/src/github.com/docker/cli`; then `./scripts/build/binary` works.

`make -f docker.Makefile test` runs the unit tests.
