.PHONY: install tests

install:
	go mod tidy && go mod vendor

tests:
	go test -covermode=set ./... -coverprofile=coverage.txt && go tool cover -func=coverage.txt