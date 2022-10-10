.PHONY: test test-cover vet lint

test:
	go test -v ./... -cover -timeout=0

test-cover:
	go test -v ./... -cover -covermode=count -coverprofile=coverage.out -timeout=0

vet:
	go vet -v ./...

lint:
	golangci-lint run
