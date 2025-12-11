# Default binary name and install prefix
BINARY_NAME ?= wk
PREFIX ?= $(HOME)/.local
BINDIR ?= $(PREFIX)/bin

# Go settings
GOFILES := $(shell find . -name '*.go' -not -path './vendor/*')

.PHONY: all build lint test integration install clean

all: build

build:
	go build -o $(BINARY_NAME) ./cmd/wk

build-web:
	go build -o web ./cmd/web

lint:
	@echo "Running gofmt check..."
	@gofmt -l $(GOFILES) | tee /tmp/gofmt.out
	@if [ -s /tmp/gofmt.out ]; then \
		echo "gofmt found issues in the files above"; \
		exit 1; \
	fi
	@echo "Running go vet..."
	go vet ./...

test:
	go test ./...

integration:
	./scripts/integration_test.sh

install: build
	@mkdir -p "$(BINDIR)"
	@cp "$(BINARY_NAME)" "$(BINDIR)/$(BINARY_NAME)"
	@echo "Installed $(BINARY_NAME) to $(BINDIR)"

clean:
	rm -f $(BINARY_NAME)
