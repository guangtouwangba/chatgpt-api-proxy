

build:
	@echo "Building..."
	@go build -o proxy

test:
	@echo "Testing..."
	@go test -v ./...