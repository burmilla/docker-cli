---
title: "Docker commands"
description: "Docker's CLI command description and usage"
keywords: "Docker, Docker documentation, CLI, command line"
identifier: "smn_cli_guide"
---

<!-- This file is maintained within the docker/docker Github
     repository at https://github.com/docker/docker/. Make all
     pull requests against that repo. If you see this file in
     another repository, consider it read-only there, as it will
     periodically be overwritten by the definitive file. Pull
     requests which include edits to this file in other repositories
     will be rejected.
-->

# The Docker commands

This section contains reference information on using Docker's command line
client. Each command has a reference page along with samples. If you are
unfamiliar with the command line, you should start by reading about how to [Use
the Docker command line](cli.md).

## Commands by object

### Docker management commands

| Command | Description                                                        |
|:--------|:-------------------------------------------------------------------|
| [info](info.md) | Display system-wide information                            |
| [inspect](inspect.md)| Return low-level information on a container or image  |
| [version](version.md) | Show the Docker version information                  |


### Image commands

| Command | Description                                                        |
|:--------|:-------------------------------------------------------------------|
| [history](history.md) | Show the history of an image                         |
| [images](images.md) | List images                                            |
| [load](load.md) | Load an image from a tar archive or STDIN                  |
| [image prune](image_prune.md) | Remove unused images                         |
| [rmi](rmi.md) | Remove one or more images                                    |
| [save](save.md) | Save images to a tar archive                               |
| [tag](tag.md) | Tag an image into a repository                               |

### Container commands

| Command | Description                                                        |
|:--------|:-------------------------------------------------------------------|
| [attach](attach.md) | Attach to a running container                          |
| [container prune](container_prune.md) | Remove all stopped containers        |
| [cp](cp.md) | Copy files/folders from a container to a HOSTDIR or to STDOUT  |
| [create](create.md) | Create a new container                                 |
| [diff](diff.md) | Inspect changes on a container's filesystem                |
| [events](events.md) | Get real time events from the server                   |
| [exec](exec.md) | Run a command in a running container                       |
| [kill](kill.md) | Kill a running container                                   |
| [logs](logs.md) | Fetch the logs of a container                              |
| [port](port.md) | List port mappings or a specific mapping for the container |
| [ps](ps.md) | List containers                                                |
| [rename](rename.md) | Rename a container                                     |
| [restart](restart.md) | Restart a running container                          |
| [rm](rm.md) | Remove one or more containers                                  |
| [run](run.md) | Run a command in a new container                             |
| [start](start.md) | Start one or more stopped containers                     |
| [stop](stop.md) | Stop a running container                                   |
| [wait](wait.md) | Block until a container stops, then print its exit code    |

### Registry commands

| Command | Description                                                        |
|:--------|:-------------------------------------------------------------------|
| [pull](pull.md) | Pull an image or a repository from a Docker registry       |

### Network and connectivity commands

| Command | Description                                                        |
|:--------|:-------------------------------------------------------------------|
| [network connect](network_connect.md) | Connect a container to a network     |
| [network create](network_create.md) | Create a new network                   |
| [network disconnect](network_disconnect.md) | Disconnect a container from a network |
| [network inspect](network_inspect.md) | Display information about a network  |
| [network ls](network_ls.md) | Lists all the networks the Engine `daemon` knows about |
| [network prune](network_prune.md) | Remove all unused networks               |
| [network rm](network_rm.md) | Removes one or more networks                   |

### Shared data volume commands

| Command | Description                                                        |
|:--------|:-------------------------------------------------------------------|
| [volume create](volume_create.md) | Creates a new volume where containers can consume and store data |
| [volume inspect](volume_inspect.md) | Display information about a volume     |
| [volume ls](volume_ls.md) | Lists all the volumes Docker knows about         |
| [volume prune](volume_prune.md) | Remove all unused volumes                  |
| [volume rm](volume_rm.md) | Remove one or more volumes                       |
