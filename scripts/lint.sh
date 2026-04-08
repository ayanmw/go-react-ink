#!/bin/bash

# Lint script for go-react-ink
# Usage: ./scripts/lint.sh [--fix]

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

FIX_MODE=false
if [ "$1" == "--fix" ]; then
    FIX_MODE=true
fi

echo "=== Go-React-Ink Lint Checker ==="
echo ""

# Check Go installation
echo -e "${YELLOW}Checking Go installation...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${RED}ERROR: Go is not installed${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi
GO_VERSION=$(go version)
echo -e "${GREEN}✓ $GO_VERSION${NC}"
echo ""

# Check gofmt
echo -e "${YELLOW}Checking gofmt...${NC}"
GOFMT_OUTPUT=$(gofmt -s -d .)
if [ -n "$GOFMT_OUTPUT" ]; then
    if [ "$FIX_MODE" = true ]; then
        echo -e "${YELLOW}Fixing formatting issues...${NC}"
        gofmt -s -w .
        echo -e "${GREEN}✓ Fixed formatting issues${NC}"
    else
        echo -e "${RED}ERROR: gofmt found issues:${NC}"
        echo "$GOFMT_OUTPUT"
        echo ""
        echo -e "${YELLOW}Run with --fix to auto-fix these issues${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✓ gofmt: No formatting issues${NC}"
fi
echo ""

# Check go vet
echo -e "${YELLOW}Running go vet...${NC}"
GO_VET_OUTPUT=$(go vet ./... 2>&1 || true)
if [ -n "$GO_VET_OUTPUT" ]; then
    echo -e "${RED}ERROR: go vet found issues:${NC}"
    echo "$GO_VET_OUTPUT"
    exit 1
else
    echo -e "${GREEN}✓ go vet: No issues${NC}"
fi
echo ""

# Check golint installation
echo -e "${YELLOW}Checking golint installation...${NC}"
if ! command -v golint &> /dev/null; then
    echo -e "${YELLOW}golint not found, installing...${NC}"
    go install golang.org/x/lint/golint@latest

    # Check if GOPATH/bin is in PATH
    GOPATH_BIN=$(go env GOPATH)/bin
    if [[ ":$PATH:" != *":$GOPATH_BIN:"* ]]; then
        echo -e "${YELLOW}Adding $GOPATH_BIN to PATH for this session${NC}"
        export PATH=$PATH:$GOPATH_BIN
    fi

    if ! command -v golint &> /dev/null; then
        echo -e "${RED}ERROR: Failed to install golint${NC}"
        echo "Please manually install: go install golang.org/x/lint/golint@latest"
        exit 1
    fi
    echo -e "${GREEN}✓ golint installed${NC}"
else
    echo -e "${GREEN}✓ golint: $(golint -version || echo 'installed')${NC}"
fi
echo ""

# Run golint
echo -e "${YELLOW}Running golint...${NC}"
GO_LINT_OUTPUT=$(golint ./... 2>&1 || true)

# Filter out acceptable warnings (F-keys and Ctrl keys in tcell)
FILTERED_LINT=$(echo "$GO_LINT_OUTPUT" | grep -v "KeyF[0-9]" | grep -v "KeyCtrl" || true)

if [ -n "$FILTERED_LINT" ]; then
    LINT_COUNT=$(echo "$FILTERED_LINT" | wc -l)
    echo -e "${YELLOW}golint found $LINT_COUNT warning(s):${NC}"
    echo "$FILTERED_LINT"
    echo ""
    echo -e "${YELLOW}Note: Some lint warnings may be acceptable. Review and fix as needed.${NC}"
    # Don't fail on lint warnings, just report
else
    echo -e "${GREEN}✓ golint: No issues${NC}"
fi
echo ""

# Final summary
echo "=== Summary ==="
if [ "$FIX_MODE" = true ]; then
    echo -e "${GREEN}All checks completed with fixes applied${NC}"
else
    echo -e "${GREEN}All checks completed${NC}"
fi
echo ""
echo "To auto-fix formatting issues, run: ./scripts/lint.sh --fix"