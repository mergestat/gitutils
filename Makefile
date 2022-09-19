.PHONY: test test-cover vet lint

test:
	go test -v ./... -cover

test-cover:
	go test -v ./... -cover -covermode=count -coverprofile=coverage.out

vet:
	go vet -v ./...

lint:
	golangci-lint run
