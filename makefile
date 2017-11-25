.PHONY: build check

check:
	go list ./... | xargs -n1 go vet -v
	go fmt ./...

test: check
	CGO_ENABLED=1 go test -v ./...
