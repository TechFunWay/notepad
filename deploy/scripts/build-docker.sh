#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
VERSION=$(cat "$PROJECT_DIR/VERSION")
BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(git -C "$PROJECT_DIR" rev-parse --short HEAD 2>/dev/null || echo "unknown")
DOCKER_REPO="wycto/notepad"

echo "==> Building Docker image v${VERSION}..."

cd "$PROJECT_DIR"

docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --build-arg VERSION="$VERSION" \
  --build-arg BUILD_TIME="$BUILD_TIME" \
  --build-arg GIT_COMMIT="$GIT_COMMIT" \
  -t "${DOCKER_REPO}:${VERSION}" \
  -t "${DOCKER_REPO}:latest" \
  .

echo "==> Docker image built: ${DOCKER_REPO}:${VERSION}"
