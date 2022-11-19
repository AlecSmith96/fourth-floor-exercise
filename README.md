# fourth-floor-exercise

A simple REST server that provides one endpoint for viewing details about a streamers videos.

## Table of contents
* [Getting Started](#getting-started)
* [Configuration Instructions](#configuration-instructions)

## Getting Started

This application provides a rest server using a variety of different technologies (see [Built With](#built-with) section for details).

* Clone the repo: `git clone git@github.com:AlecSmith96/fourth-floor-exercise.git`



### Local Development

The application can be run locally in two ways: 
* `make run-service` wraps the go run functionality.
* `make build` builds the binary which can be run from the `/dist` directory.

Once started the output should look like this:
```
fourth-floor-exercise$ make run-service
go run ./cmd --config dev-config.yaml
2022-11-19T15:21:37.495Z        info    adapters/loggingAdapter.go:55   Initialized logging adapter
2022-11-19T15:21:37.495Z        info    cmd/main.go:20  Starting the service
```

By default the service will run on port `8080`, although this is configurable (see [Configuration Instructions](#configuration-instructions)).

### Run Using Docker

The application provides convenient make commands for both building the docker image and running it. Configuration for the docker image will be pulled in from the same configuration file when running it locally.
* `make docker-build` - Builds the image with the tag `ff-exercise`
* `make docker-run` - Runs the image in detached mode, exposing the default port (8080)

This will build and run the image with the name `ff-exercise`:
```
fourth-floor-exercise$ docker container ls
CONTAINER ID   IMAGE                    COMMAND                  CREATED         STATUS                PORTS                                                                                                                                                                                                                                                    NAMES
224dcedc7e44   ff-exercise              "/build --config dev…"   8 seconds ago   Up 6 seconds          0.0.0.0:8080->8080/tcp, :::8080->8080/tcp
```

## Configuration Instructions
Configuration can be passed to the service through the `--config` argument. All configuration options come with default values to get you up and running faster, apart from the `auth` configuration. 

Below is a list of the configuration points for the service:
| Name                | Default Values | Valid Values |  Description |
| -----------------   | -- | --| ----------- |
| `rest.port`         | `8080` | `any integer` | The port the REST server will run on       |
| `rest.ginMode`      | `release` | `release \| debug` | The mode the gin logger will run on        |
| `logging.loglevel`  | `debug` | `debug \| info` | The loglevel of the logger        |
| `logging.encoding`  | `console` | `console \| json` | How the logger encodes its logs, `console` provides easy to read logs for when running locally in a terminal, `json` wraps all logs in json objects making it more suitable for production.        |
| `auth.clientId`     | `""` | `any string` | the clientId of the twitch app being used        |
| `auth.clientSecret` | `""` | `any string` | the clientSecret of the twitch app being used        |

## Built With
* [gin/gonic](https://github.com/gin-gonic/gin) - Provides HTTP router
* [knadh/koanf](https://github.com/knadh/koanf) - Config management library
* [uber-go/zap](https://github.com/uber-go/zap) - Structured, leveled logger for Go
* [google/wire](https://github.com/google/wire) - Compile-time dependency injection package for Go
