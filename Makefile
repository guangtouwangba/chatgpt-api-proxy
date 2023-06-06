
dev:
	@echo "Build for dev..."
	@go build -o proxy

build:
	@echo "Building..."
	@docker build -t proxy --build-arg ARCH=arm64v8 -f Dockerfile .

test:
	@echo "Testing..."
	@go test -v ./...

lint:
	@echo "Linting..."
	@golangci-lint run ./...

coverage:
	@echo "Coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out


COVERAGE_THRESHOLD=80

