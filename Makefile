.PHONY: help build clean security security-scan deps-check code-scan install test

help:
	@echo "Doppel - Duplicate File Finder"
	@echo ""
	@echo "Usage:"
	@echo "  make build          Build the binary"
	@echo "  make install        Install to system"
	@echo "  make test           Run tests"
	@echo "  make security       Run all security scans"
	@echo "  make deps-check     Check dependencies for vulnerabilities"
	@echo "  make code-scan      Run code security analysis"
	@echo "  make clean          Remove build artifacts"

VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%d")
LDFLAGS := -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE) -X 'doppel/cmd.Version=$(VERSION)' -X 'doppel/cmd.Commit=$(COMMIT)' -X 'doppel/cmd.BuildDate=$(BUILD_DATE)'

build:
	@echo "Building doppel..."
	go build -ldflags "$(LDFLAGS)" -o doppel

install: build
	@echo "Installing doppel..."
	./install.sh

test:
	@echo "Running tests..."
	go test ./...

security: deps-check code-scan
	@echo ""
	@echo "✅ Security scans completed!"

deps-check:
	@echo "Checking dependencies for vulnerabilities..."
	@command -v govulncheck >/dev/null 2>&1 || { echo "Installing govulncheck..."; go install golang.org/x/vuln/cmd/govulncheck@latest; }
	@govulncheck ./...

code-scan:
	@echo "Running code security analysis..."
	@command -v gosec >/dev/null 2>&1 || { echo "Installing gosec..."; go install github.com/securego/gosec/v2/cmd/gosec@latest; }
	@gosec -conf .gosec.json ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -f doppel doppel-*
	rm -rf dist/
