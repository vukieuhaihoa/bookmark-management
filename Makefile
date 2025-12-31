.PHONY: run
start:
	go run ./cmd/api/main.go

.PHONY: mock-gen
mock-gen:
	go generate ./...

.PHONY: test
test:
	go test -v ./...