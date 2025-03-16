#!/bin/bash

echo "Installing Git hooks..."
cp hooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
echo "✅ Pre-commit hook installed successfully!"

# Install required Go tools
echo "📦 Installing gofumpt and golangci-lint..."
go install mvdan.cc/gofumpt@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Ensure GOPATH/bin is in PATH
export PATH=$PATH:$(go env GOPATH)/bin

echo "✅ Installation complete. Ensure you restart your shell if necessary."
