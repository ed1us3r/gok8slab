#!/bin/env bash

DIRECTORY="gok8slab"
LATEST_RELEASE_URL="https://api.github.com/repos/ed1us3r/gok8slab/releases/latest"
ZIP_FILE="gok8slab-latest.tar.gz"
BINARY_NAME="gok8slab"

# Check for --install flag
INSTALL=false
if [[ "$1" == "--install" ]]; then
    INSTALL=true
fi

# Fetch the latest release download URL
REPO_URL=$(curl -s $LATEST_RELEASE_URL | grep "browser_download_url" | grep "zip" | cut -d '"' -f 4)

if [[ -z "$REPO_URL" ]]; then
    echo "Failed to fetch the latest release URL. Exiting."
    exit 1
fi

if $INSTALL; then
    # Determine the best installation directory
    for DIR in "$HOME/bin" "$HOME/.local/bin" "/usr/local/bin"; do
        if [[ -d "$DIR" && ":$PATH:" == *":$DIR:"* ]]; then
            INSTALL_DIR="$DIR"
            break
        fi
    done

    if [[ -z "$INSTALL_DIR" ]]; then
        echo "No suitable directory found in PATH. Please enter a directory to install the binary:"
        read -r INSTALL_DIR
        echo "Please add $INSTALL_DIR to your PATH to use the binary properly."
    fi

    mkdir -p "$INSTALL_DIR"
    wget -q --show-progress "$REPO_URL" -O "$BINARY_NAME.zip"
    if ! command -v unzip &> /dev/null; then
        echo "Error: unzip command not found. Please install unzip and try again."
        exit 1
    fi
    unzip -q "$BINARY_NAME.zip"
    rm "$BINARY_NAME.zip"
    if [[ ! -f "$BINARY_NAME" ]]; then
        echo "Error: Extracted binary not found. Installation failed."
        exit 1
    fi
    mv "$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    echo "Binary installed at $INSTALL_DIR/$BINARY_NAME"
    exit 0
fi

# Ask for user permission
echo "This script will create a directory called '$DIRECTORY' and download the latest source code. Continue? (y/n)"
read -r RESPONSE

if [[ "$RESPONSE" != "y" ]]; then
    echo "Operation canceled."
    exit 1
fi

# Create the directory
mkdir -p "$DIRECTORY"
cd "$DIRECTORY" || exit

# Download the source code as a zip file
wget -q --show-progress "$REPO_URL" -O "$ZIP_FILE"

if ! command -v unzip &> /dev/null; then
    echo "Error: unzip command not found. Please install unzip and try again."
    exit 1
fi

# Unzip the file
unzip -q "$ZIP_FILE"
rm "$ZIP_FILE"

# Display help text with README and compile instructions
if [[ -f "README.md" ]]; then
    echo -e "\n===== README.md CONTENT =====\n"
    cat README.md | head -n 20  # Display first 20 lines of README
    echo -e "\n============================\n"
else
    echo "README.md not found."
fi

# Print compile instructions
echo -e "\nTo compile the binary locally, run the following commands:"
echo "cd $DIRECTORY"
echo "make build  # If a Makefile exists, otherwise use 'go build .'"
echo "./gok8slab  # To execute the compiled binary"

# Markdown snippet for downloading and executing the script
echo -e "\n### Quick Install Command\n"
echo -e "\
```sh\
"
echo -e "curl -sL https://raw.githubusercontent.com/yourrepo/hack.sh | bash -s -- --install\
"
echo -e "```\
"

