.PHONY: build check

check:
	go list ./... | grep -v vendor | xargs -n1 go vet -v
	go list ./... | grep -v vendor | go fmt

test: check
	CGO_ENABLED=1 go test -v ./...
