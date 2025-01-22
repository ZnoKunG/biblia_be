build:
	@go build -o bin/biblia_be cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/biblia_be