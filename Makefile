

build:
	@echo "Building..."
	@go build -o ./bin/proxy

test:
	@echo "Testing..."
	@go test -v ./...

lint:
	@echo "Linting..."
	@golangci-lint run ./...