ifndef PKGS
PKGS := $(shell go list ./... 2>&1 | grep -v 'vendor' | grep -v 'sanity')
endif

all: get-tools vendor-sync test lint vet docker-build compose

compose:
	docker-compose up

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go

test:
	go test -v -race ./pkg/...

docker-build:
	docker build -t stgleb/port-client -f Dockerfile.client .
	docker build -t stgleb/port-server -f Dockerfile.server .

get-tools:
	go get -u golang.org/x/lint/golint
	go get github.com/golang/protobuf/protoc-gen-go
	go get github.com/golang/mock/mockgen

lint:
	for file in $(GO_FILES); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done

proto:
	protoc --go_out=plugins=grpc:api/ api/api.proto

vet:
	go vet $(PKGS)

vendor-sync:
	go mod tidy
	go mod download
	go mod vendor

