BINARY=quoctl

.PHONY: build run clean

build:
	go build -o $(BINARY) ./cmd/quoctl

run:
	go run ./cmd/quoctl --help

clean:
	@if command -v trash >/dev/null 2>&1; then trash $(BINARY) 2>/dev/null || true; else rm -f $(BINARY); fi
