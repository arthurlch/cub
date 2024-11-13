.PHONY: help build build-all deps clean install uninstall verify-path

BINARY_NAME = cub
OUTPUT_DIR = bin

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)  # macOS
    INSTALL_DIR = /usr/local/bin
    SUDO_CMD = sudo
    PATH_CHECK := $(shell echo $$PATH | grep -q $(INSTALL_DIR) && echo "yes" || echo "no")
    SHELL_CONFIG = ~/.zshrc
else ifeq ($(UNAME_S),Linux)
    INSTALL_DIR = /usr/local/bin
    SUDO_CMD = sudo
    PATH_CHECK := $(shell echo $$PATH | grep -q $(INSTALL_DIR) && echo "yes" || echo "no")
    SHELL_CONFIG = ~/.bashrc
else ifeq ($(OS),Windows_NT)  # Windows
    BINARY_NAME := $(BINARY_NAME).exe
    INSTALL_DIR = C:\Windows\System32
    SUDO_CMD = 
    PATH_CHECK = yes  
    SHELL_CONFIG = 
endif

UNAME_M := $(shell uname -m)
ifeq ($(UNAME_M),x86_64)
    ARCH = amd64
else ifeq ($(UNAME_M),arm64)
    ARCH = arm64
else
    ARCH = $(UNAME_M)
endif

PLATFORMS = darwin-amd64 darwin-arm64 windows-amd64 linux-amd64
INSTALL_DIR_DARWIN = /usr/local/bin
INSTALL_DIR_LINUX = /usr/local/bin
INSTALL_DIR_WINDOWS = %USERPROFILE%\AppData\Local\Microsoft\WindowsApps
HOME_DIR = $(shell echo $$HOME)
HOME_BIN_DIR = $(HOME_DIR)/bin

# Detect OS and Architecture
ifeq ($(OS),Windows_NT)
	DETECTED_OS := windows
else
	DETECTED_OS := $(shell uname -s | tr A-Z a-z)
	DETECTED_ARCH := $(shell uname -m)
endif

all: build

help:
	@echo "Available targets:"
	@echo "  build             Build the binary for the current system ($(UNAME_S)-$(ARCH))"
	@echo "  build-all         Build binaries for all supported platforms"
	@echo "  deps              Install dependencies"
	@echo "  clean             Remove built binaries"
	@echo "  install           Install the binary to $(INSTALL_DIR)"
	@echo "  uninstall         Remove the binary from $(INSTALL_DIR)"
	@echo "  verify-path       Verify $(INSTALL_DIR) is in PATH"
	@echo "  help              Show this help message"

build: deps
	go build -o $(BINARY_NAME) ./cmd/cub

build-all: $(PLATFORMS)

$(PLATFORMS):
	@echo "Building for $@..."
	GOOS=$(word 1, $(subst -, ,$@)) GOARCH=$(word 2, $(subst -, ,$@)) \
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME)_$@ ./cmd/cub

deps:
	go mod tidy

clean:
ifeq ($(OS),Windows_NT)
	if exist $(BINARY_NAME) del /F $(BINARY_NAME)
	if exist $(OUTPUT_DIR) rmdir /S /Q $(OUTPUT_DIR)
else
	rm -f $(BINARY_NAME)
	rm -rf $(OUTPUT_DIR)
endif

verify-path:
ifneq ($(OS),Windows_NT)  
ifeq ($(PATH_CHECK),no)
	@echo "Warning: $(INSTALL_DIR) is not in your PATH!"
	@echo "Add the following line to your $(SHELL_CONFIG):"
	@echo "export PATH=\$$PATH:$(INSTALL_DIR)"
endif
endif

install: build verify-path
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
ifeq ($(OS),Windows_NT)
	@if not exist $(INSTALL_DIR) mkdir $(INSTALL_DIR)
	@copy /Y $(BINARY_NAME) $(INSTALL_DIR)
else
	@if [ ! -d $(INSTALL_DIR) ]; then \
		$(SUDO_CMD) mkdir -p $(INSTALL_DIR); \
	fi
	@$(SUDO_CMD) cp $(BINARY_NAME) $(INSTALL_DIR)/
	@$(SUDO_CMD) chmod 755 $(INSTALL_DIR)/$(BINARY_NAME)
endif
	@echo "Installation complete. You can now use '$(BINARY_NAME)' command from any directory."
ifneq ($(OS),Windows_NT)
ifeq ($(PATH_CHECK),no)
	@echo ""
	@echo "NOTE: Please run:"
	@echo "export PATH=\$$PATH:$(INSTALL_DIR)"
	@echo "Or add it to your shell's configuration file ($(SHELL_CONFIG))"
endif
endif

uninstall:
	@echo "Removing $(BINARY_NAME) from $(INSTALL_DIR)..."
ifeq ($(OS),Windows_NT)
	@del /F $(INSTALL_DIR)\$(BINARY_NAME)
else
	@$(SUDO_CMD) rm -f $(INSTALL_DIR)/$(BINARY_NAME)
endif
	@echo "Uninstallation complete."
