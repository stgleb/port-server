# Port service

Port-service is a platform for managing information about ports

## Prerequisites

Up-to date docker, docker-compose and golang

## Install

### Install by `go get`

`go get github.com/stgleb/port-service`

## Build

build docker images

`make docker-build`

## Run

### Run from source

Run server

`make run-server`

Run client

`make run-client`
            
### Build and run in docker compose

```
make docker-build
make compose
```

## API

Client exposes REST api for getting information about ports

get info about port

`curl http://localhost:9091/port/AEAJM`

get info about ports

`curl http://localhost:9091/port`

## Development

If you want to make an update in grpc API install [protoc](http://google.github.io/proto-lens/installing-protoc.html) compiler 
and run command

`make proto`
  
Before commit to repo run `make verify` for verifying code

## Build and run all

`make all`