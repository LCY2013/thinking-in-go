GOPATH:=$(shell go env GOPATH)
.PHONY: proto
proto:
	docker run --rm -v "$(shell pwd):/workspace" --workdir /workspace luochunyun/protoc:1.0.0 -I ./ --go_out=./ --micro_out=./ ./proto/payment-api/payment-api.proto

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o payment-api *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t payment-api:latest
