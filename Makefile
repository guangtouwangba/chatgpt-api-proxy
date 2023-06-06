

build:
	@echo "Building..."
	@docker build -t proxy --build-arg ARCH=arm64v8 -f Dockerfile .

test:
	@echo "Testing..."
	@go test -v ./...

lint:
	@echo "Linting..."
	@golangci-lint run ./...