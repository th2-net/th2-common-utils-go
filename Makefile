SRC_MAIN_PROTO_DIR=src/main/proto
GITHUB_TH2=github.com/th2-net

TH2_GRPC_COMMON=th2-grpc-common
TH2_GRPC_COMMON_URL=$(GITHUB_TH2)/$(TH2_GRPC_COMMON)@makefile

MODULE_NAME=th2-grpc
MODULE_DIR=$(MODULE_NAME)

PROTOBUF_VERSION=v1.5.2

PROTOC_VERSION=21.12

default: build

build:
	go vet ./...
	go build -v -race ./...

run-test:
	go test -v ./...