.PHONY: fix
fix: format lint

.PHONY: lint
lint:
	golangci-lint run

.PHONY: format
format:
	go fmt ./...

.PHONY: test
test:
	go test -v ./...
