GOPATH:=$(shell go env GOPATH)
.PHONY: proto
proto:
	docker run --rm -v "$(shell pwd):/workspace" --workdir /workspace luochunyun/protoc:1.0.0 -I ./ --go_out=./ --micro_out=./ ./proto/cart-api/cart-api.proto

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o cart-api *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t cart-api:latest
