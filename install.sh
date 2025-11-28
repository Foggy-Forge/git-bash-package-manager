#!/usr/bin/env bash
set -e

# gbpm installer for Git Bash on Windows (and other Unix-like systems)

REPO="Foggy-Forge/git-bash-package-manager"
INSTALL_DIR="${GBPM_HOME:-$HOME/.gbpm}"
BIN_DIR="$INSTALL_DIR/bin"

echo "Installing gbpm..."

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
  linux*)
    OS="linux"
    ;;
  darwin*)
    OS="darwin"
    ;;
  mingw*|msys*|cygwin*)
    OS="windows"
    ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

case "$ARCH" in
  x86_64|amd64)
    ARCH="amd64"
    ;;
  arm64|aarch64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Get latest release
echo "Fetching latest release..."
LATEST_URL="https://api.github.com/repos/$REPO/releases/latest"
RELEASE_DATA=$(curl -sL "$LATEST_URL")
TAG=$(echo "$RELEASE_DATA" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$TAG" ]; then
  echo "Error: Could not determine latest release"
  exit 1
fi

echo "Latest version: $TAG"

# Construct download URL
if [ "$OS" = "windows" ]; then
  BINARY_NAME="gbpm-${OS}-${ARCH}.exe"
else
  BINARY_NAME="gbpm-${OS}-${ARCH}"
fi

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$TAG/$BINARY_NAME"

# Create directories
mkdir -p "$BIN_DIR"

# Download binary
echo "Downloading $BINARY_NAME..."
TEMP_FILE=$(mktemp)
if command -v curl >/dev/null 2>&1; then
  curl -sL -o "$TEMP_FILE" "$DOWNLOAD_URL"
elif command -v wget >/dev/null 2>&1; then
  wget -q -O "$TEMP_FILE" "$DOWNLOAD_URL"
else
  echo "Error: curl or wget is required"
  exit 1
fi

# Check if download was successful
if [ ! -s "$TEMP_FILE" ]; then
  echo "Error: Failed to download binary"
  rm -f "$TEMP_FILE"
  exit 1
fi

# Install binary
if [ "$OS" = "windows" ]; then
  TARGET="$BIN_DIR/gbpm.exe"
else
  TARGET="$BIN_DIR/gbpm"
fi

mv "$TEMP_FILE" "$TARGET"
chmod +x "$TARGET"

echo "✓ Installed gbpm to $TARGET"

# Check if bin directory is in PATH
if [[ ":$PATH:" != *":$BIN_DIR:"* ]]; then
  echo ""
  echo "⚠️  $BIN_DIR is not in your PATH"
  echo ""
  echo "Add the following line to your ~/.bashrc or ~/.bash_profile:"
  echo ""
  echo "    export PATH=\"$BIN_DIR:\$PATH\""
  echo ""
  
  # Offer to add to .bashrc
  if [ -f "$HOME/.bashrc" ]; then
    read -p "Would you like to add it to ~/.bashrc now? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
      echo "" >> "$HOME/.bashrc"
      echo "# gbpm" >> "$HOME/.bashrc"
      echo "export PATH=\"$BIN_DIR:\$PATH\"" >> "$HOME/.bashrc"
      echo "✓ Added to ~/.bashrc"
      echo "Run: source ~/.bashrc"
    fi
  fi
else
  echo "✓ $BIN_DIR is already in PATH"
fi

echo ""
echo "Installation complete! Run 'gbpm version' to verify."
