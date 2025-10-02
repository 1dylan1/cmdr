BINARY_NAME := cmdr
MAIN_PACKAGE := .
INSTALL_DIR := /usr/local/bin


USER_CONFIG_BASE_DIR := $(shell XDG_CONFIG_HOME=${HOME}/.config; echo "$$XDG_CONFIG_HOME")
USER_CONFIG_DIR := $(USER_CONFIG_BASE_DIR)/$(BINARY_NAME)
USER_CONFIG_FILE := /usr/share/cmdr/config.yaml

SOURCE_DEFAULT_USER_CONFIG := config.yaml.example

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "$(BINARY_NAME) built successfully."

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@echo "This operation requires root privileges. You may be prompted for your password."
	@sudo cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@sudo chmod 0755 $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) installed successfully to $(INSTALL_DIR)."
	@echo "You can now run '$(BINARY_NAME)' from any directory."
	@mkdir -p "$(USER_CONFIG_DIR)"
	@if [ ! -f "$(USER_CONFIG_FILE)" ]; then \
		install -m 0644 "$(SOURCE_DEFAULT_USER_CONFIG)" "$(USER_CONFIG_FILE)"; \
		echo "Default user config created at: $(USER_CONFIG_FILE). <-- Edit me!"; \
	else \
		echo "User config already exists at $(USER_CONFIG_FILE). Skipping"; \
	fi

uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_DIR)..."
	@echo "This operation requires root privileges. You may be prompted for your password."
	-sudo rm $(INSTALL_DIR)/$(BINARY_NAME) || true
	@echo "$(BINARY_NAME) uninstalled."
	@echo "$(USER_CONFIG_DIR) is NOT uninstalled"

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@echo "Clean complete."

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all       - Builds the executable (default)."
	@echo "  build     - Builds the executable."
	@echo "  install   - Installs the executable to $(INSTALL_DIR) (requires sudo)."
	@echo "  uninstall - Removes the executable from $(INSTALL_DIR) (requires sudo)."
	@echo "  clean     - Removes built executable and other artifacts."
	@echo "  help      - Display this help message."
