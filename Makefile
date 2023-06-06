
COVERAGE_THRESHOLD=0


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

test-coverage:
	go test -coverprofile=coverage.out ./...

coverage.out:
	@go test -coverprofile=coverage.out ./...

# 检查覆盖率是否达到阈值
check-coverage: coverage.out
	@go tool cover -func=coverage.out | grep total | awk '{ print substr($$3, 1, length($$3)-1) }' | { read cov; result=`echo "$$cov >= $(COVERAGE_THRESHOLD)" | bc -l`; if [ $$result -eq 0 ]; then echo Coverage is below $(COVERAGE_THRESHOLD)%; exit 1; fi; }


.PHONY: test coverage.outcheck-coverage

