#!/bin/bash

# Release script for go-react-ink
# Usage: ./scripts/release.sh <version>

set -e

VERSION=${1:-"0.1.0"}
TAG="v${VERSION}"

echo "=== Preparing release ${TAG} ==="

# Check if tag exists
if git rev-parse "${TAG}" >/dev/null 2>&1; then
    echo "Tag ${TAG} already exists"
    exit 1
fi

# Run tests
echo "Running tests..."
go test ./...
if [ $? -ne 0 ]; then
    echo "Tests failed"
    exit 1
fi

# Update version in files
echo "Updating version in files..."
sed -i "s/version = \"[0-9.]*\"/version = \"${VERSION}\"/" cmd/gox/main.go
sed -i "s/\"version\": \"[0-9.]*\"/\"version\": \"${VERSION}\"/" tools/vscode-gox/package.json

# Build binaries
echo "Building binaries..."
go build -v -ldflags="-s -w" -o gox ./cmd/gox
cd tools/lsp-server && go build -v -ldflags="-s -w" -o ../../gox-lsp . && cd ..

# Create release archive
echo "Creating release archives..."
mkdir -p release

# Linux
GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w" -o release/gox-linux-amd64 ./cmd/gox
cd tools/lsp-server && GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w" -o ../../release/gox-lsp-linux-amd64 . && cd ..

# Windows
GOOS=windows GOARCH=amd64 go build -v -ldflags="-s -w" -o release/gox-windows-amd64.exe ./cmd/gox
cd tools/lsp-server && GOOS=windows GOARCH=amd64 go build -v -ldflags="-s -w" -o ../../release/gox-lsp-windows-amd64.exe . && cd ..

# macOS
GOOS=darwin GOARCH=amd64 go build -v -ldflags="-s -w" -o release/gox-darwin-amd64 ./cmd/gox
cd tools/lsp-server && GOOS=darwin GOARCH=amd64 go build -v -ldflags="-s -w" -o ../../release/gox-lsp-darwin-amd64 . && cd ..

# Create tarballs
cd release
tar -czvf gox-linux-amd64.tar.gz gox-linux-amd64 gox-lsp-linux-amd64
zip gox-windows-amd64.zip gox-windows-amd64.exe gox-lsp-windows-amd64.exe
tar -czvf gox-darwin-amd64.tar.gz gox-darwin-amd64 gox-lsp-darwin-amd64
cd ..

# Build VSCode extension
echo "Building VSCode extension..."
cd tools/vscode-gox
npm install
npm run compile
if command -v vsce &> /dev/null; then
    vsce package --out gox-${TAG}.vsix
fi
cd ..

echo "=== Release artifacts created ==="
ls -la release/

echo ""
echo "Next steps:"
echo "1. git add -A"
echo "2. git commit -m 'chore: bump version to ${VERSION}'"
echo "3. git tag ${TAG}"
echo "4. git push origin master"
echo "5. git push origin ${TAG}"
echo "6. Create GitHub release with release notes from CHANGELOG.md"
echo "7. Upload release assets from release/ directory"