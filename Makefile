.PHONY: build run dev

# Get the current git commit hash
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')

build:
	go build -ldflags "-X 'github.com/hydrocode-de/datailama/internal/version.BuildTime=$(BUILD_TIME)' -X 'github.com/hydrocode-de/datailama/internal/version.GitCommit=$(GIT_COMMIT)'" -o datailama .

run: build
	./datailama serve

dev: build
	air serve

