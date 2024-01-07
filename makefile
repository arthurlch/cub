.PHONY: install deps clean

OS ?= darwin

deps:
	go mod tidy

install:
	GOOS=$(OS) go build -o cub

clean:
	rm -f cub cub_macos cub.exe
