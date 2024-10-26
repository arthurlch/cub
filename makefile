.PHONY: help build build-all deps clean install

BINARY_NAME = cub
OUTPUT_DIR = bin
GOFLAGS = CGO_ENABLED=0

PLATFORMS = darwin-amd64 darwin-arm64 windows-amd64 linux-amd64

all: build

help:
	@echo "Available targets:"
	@echo "  build             Build the binary for the current system."
	@echo "  build-all         Build binaries for all supported platforms."
	@echo "  deps              Install dependencies."
	@echo "  clean             Remove built binaries."
	@echo "  install           Install the binary to /usr/local/bin."
	@echo "  help              Show this help message."

build: deps
	mkdir -p $(OUTPUT_DIR)
	$(GOFLAGS) go build -o $(OUTPUT_DIR)/$(BINARY_NAME) ./cmd/cub

build-all: $(PLATFORMS)

$(PLATFORMS):
	@echo "Building for $@..."
	GOOS=$(word 1, $(subst -, ,$@)) GOARCH=$(word 2, $(subst -, ,$@)) \
	$(GOFLAGS) go build -o $(OUTPUT_DIR)/$(BINARY_NAME)_$@ ./cmd/cub

deps:
	go mod tidy

clean:
	rm -rf $(OUTPUT_DIR)

install: build
	install -Dm755 $(OUTPUT_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
