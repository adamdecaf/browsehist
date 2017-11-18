.PHONY: build check

check:
	go vet ./...
	go fmt ./...

test:
	go test -v ./...
