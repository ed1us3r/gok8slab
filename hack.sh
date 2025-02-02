#!/bin/bash

DIRECTORY="gok8slab"
LATEST_RELEASE_URL="https://api.github.com/repos/ed1us3r/gok8slab/releases/latest"
BINARY_NAME="gok8slab"

# Detect system architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Check for --install flag
INSTALL=false
if [[ "$1" == "--install" ]]; then
    INSTALL=true
fi

# Fetch the latest release version
tarball_url=$(curl -s $LATEST_RELEASE_URL | grep "browser_download_url" | grep "linux-$ARCH.tar.gz" | grep -v "md5" | cut -d '"' -f 4)

if [[ -z "$tarball_url" ]]; then
    echo "Failed to fetch the latest release tarball. Exiting."
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
    wget -q --show-progress "$tarball_url" -O "$BINARY_NAME.tar.gz"
    if ! command -v tar &> /dev/null; then
        echo "Error: tar command not found. Please install tar and try again."
        exit 1
    fi
    if ! tar -xzf "$BINARY_NAME.tar.gz"; then
        echo "Error: Failed to extract tarball. It may be corrupt."
        rm "$BINARY_NAME.tar.gz"
        exit 1
    fi
    rm "$BINARY_NAME.tar.gz"
    if [[ ! -f "$BINARY_NAME" ]]; then
        echo "Error: Extracted binary not found. Installation failed."
        exit 1
    fi
    mv "$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    echo "Binary installed at $INSTALL_DIR/$BINARY_NAME"
    exit 0
fi

# Skip user prompt if running in a non-interactive shell
if [[ -t 0 ]]; then
    echo "This script will create a directory called '$DIRECTORY' and download the latest source code. Continue? (y/n)"
    read -r RESPONSE

    if [[ "$RESPONSE" != "y" ]]; then
        echo "Operation canceled."
        exit 1
    fi
fi

# Create the directory
mkdir -p "$DIRECTORY"
cd "$DIRECTORY" || exit

# Download the source code as a zip file
wget -q --show-progress "$tarball_url" -O "$BINARY_NAME.tar.gz"

if ! command -v tar &> /dev/null; then
    echo "Error: tar command not found. Please install tar and try again."
    exit 1
fi

# Extract the file
if ! tar -xzf "$BINARY_NAME.tar.gz"; then
    echo "Error: Failed to extract tarball. It may be corrupt."
    rm "$BINARY_NAME.tar.gz"
    exit 1
fi
rm "$BINARY_NAME.tar.gz"

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

