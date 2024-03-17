.PHONY: build deps clean build-macos-intel build-macos-arm build-windows build-linux install

# Default build uses the current system's OS and architecture.
build:
	GOOS=`go env GOOS` GOARCH=`go env GOARCH` go build -o ./cub ./cmd/cub

# Build for macOS (Intel).
build-macos-intel:
	GOOS=darwin GOARCH=amd64 go build -o ./cub_macos_intel ./cmd/cub

# Build for macOS (ARM).
build-macos-arm:
	GOOS=darwin GOARCH=arm64 go build -o ./cub_macos_arm ./cmd/cub

# Build for Windows.
build-windows:
	GOOS=windows GOARCH=amd64 go build -o ./cub.exe ./cmd/cub

# Build for Linux.
build-linux:
	GOOS=linux GOARCH=amd64 go build -o ./cub_linux ./cmd/cub

deps:
	go mod tidy

clean:
	rm -f ./cub ./cub_macos_intel ./cub_macos_arm ./cub.exe ./cub_linux

install: build
	mv ./cub /usr/local/bin/cub
