ifndef PKGS
PKGS := $(shell go list ./... 2>&1 | grep -v 'vendor' | grep -v 'sanity')
endif

all: get-tools vendor-sync test lint vet docker-build compose

docker-build:
	# TODO

compose:
	# TODO

compose:
	docker-compose up

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go

verify: vendor-sync test lint vet

get-tools:
	go get -u golang.org/x/lint/golint
	go get github.com/golang/protobuf/protoc-gen-go

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

