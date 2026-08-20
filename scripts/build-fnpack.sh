#!/bin/bash

set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

VERSION=$(tr -d '\n' < VERSION)
[ -z "$VERSION" ] && echo "❌ 无法获取版本号" && exit 1
BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"
GOCACHE_DIR="${GOCACHE:-/tmp/notepad-fnos-go-cache}"
APP_NAME="techfunway-notepad"
RELEASE_DIR="release/v${VERSION}"
WORK_DIR=$(mktemp -d "${TMPDIR:-/tmp}/notepad-fnpack.XXXXXX")
trap 'rm -rf "$WORK_DIR"' EXIT

command -v fnpack >/dev/null 2>&1 || { echo "❌ 未找到 fnpack"; exit 1; }
mkdir -p "$RELEASE_DIR"

echo "============================================"
echo "  仅构建飞牛 fnOS 安装包 v${VERSION}"
echo "============================================"

echo "📦 构建飞牛前端..."
VITE_FNOS_APP=true npm --prefix web run build -- --base="/app/${APP_NAME}/"
rm -rf server/static/dist
cp -R web/dist server/static/dist

build_fnpack() {
    local arch=$1 platform=$2 label=$3
    local binary="${WORK_DIR}/notepad-linux-${arch}"
    local package_dir="${WORK_DIR}/package-${arch}"

    echo "🔨 编译 linux/${arch}..."
    (
        cd server
        GOCACHE="$GOCACHE_DIR" CGO_ENABLED=0 GOOS=linux GOARCH="$arch" \
            go build -trimpath -ldflags="$LDFLAGS" -o "$binary" .
    )

    mkdir -p "$package_dir"
    cp -R fnpack/. "$package_dir/"
    rm -rf "$package_dir/app/server" "$package_dir/app/www"
    mkdir -p "$package_dir/app/server" "$package_dir/app/www"
    cp "$binary" "$package_dir/app/server/notepad"
    chmod +x "$package_dir/app/server/notepad"
    cp -R web/dist/. "$package_dir/app/www/"
    sed -e "s/^version.*/version               = ${VERSION#v}/" \
        -e "s/^platform.*/platform              = ${platform}/" \
        "$package_dir/manifest" > "$package_dir/manifest.tmp"
    mv "$package_dir/manifest.tmp" "$package_dir/manifest"
    find "$package_dir" -name '.DS_Store' -delete

    (cd "$package_dir" && fnpack build)
    mv "$package_dir/${APP_NAME}.fpk" "$RELEASE_DIR/${APP_NAME}-v${VERSION#v}-${label}.fpk"
    echo "✅ ${APP_NAME}-v${VERSION#v}-${label}.fpk"
}

build_fnpack arm64 arm fnos-arm64
build_fnpack amd64 x86 fnos-amd64

echo "✅ 飞牛安装包已输出到 ${RELEASE_DIR}/"
ls -lh "$RELEASE_DIR"/*.fpk
