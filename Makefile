default: build

build:
	go vet ./...
	go build -v -race ./...

run-test:
	go test -v ./...