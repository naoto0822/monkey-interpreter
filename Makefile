dep:
	go get ./...

test:
	go test ./...

build: test
	go build -o ./build/interpreter ./cmd/interpreter/main.go

lint:
	golint -set_exit_status $$(go list ./... | grep -v /vendor/)

.PHONY: dep test build lint
