BINARY_NAME := cmdr
MAIN_PACKAGE := .
INSTALL_DIR := /usr/local/bin

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "$(BINARY_NAME) built successfully."

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@echo "This operation requires root privileges. You may be prompted for your password."
	sudo cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	sudo chmod 0755 $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) installed successfully to $(INSTALL_DIR)."
	@echo "You can now run '$(BINARY_NAME)' from any directory."

uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_DIR)..."
	@echo "This operation requires root privileges. You may be prompted for your password."
	-sudo rm $(INSTALL_DIR)/$(BINARY_NAME) || true
	@echo "$(BINARY_NAME) uninstalled."

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
