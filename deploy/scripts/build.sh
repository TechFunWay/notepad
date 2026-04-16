#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
VERSION=$(cat "$PROJECT_DIR/VERSION")

echo "============================================"
echo "  Notepad v${VERSION} - Full Build"
echo "============================================"

cd "$PROJECT_DIR"

# Step 1: Build frontend
echo ""
echo "==> [1/4] Building frontend..."
cd web && npm ci --silent && npm run build
cd "$PROJECT_DIR"

# Step 2: Copy frontend to Go embed dir
echo "==> [2/4] Preparing Go embed..."
rm -rf server/static/dist
cp -r web/dist server/static/dist

# Step 3: Cross-compile Go binaries
echo "==> [3/4] Cross-compiling Go binaries..."
mkdir -p "release/$VERSION"
cd server
BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "../release/$VERSION/notepad_linux_amd64" ./main.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -o "../release/$VERSION/notepad_linux_arm64" ./main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "../release/$VERSION/notepad_darwin_amd64" ./main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o "../release/$VERSION/notepad_darwin_arm64" ./main.go
cd "$PROJECT_DIR"

# Step 4: Build FPK
echo "==> [4/4] Building FPK package..."
"$SCRIPT_DIR/pack-fpk.sh"

# Summary
echo ""
echo "============================================"
echo "  Build Complete!"
echo "============================================"
echo ""
echo "Output directory: release/$VERSION/"
ls -lh "release/$VERSION/"
echo ""
