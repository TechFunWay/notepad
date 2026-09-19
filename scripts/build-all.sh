#!/bin/bash

set -e

PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
cd "$PROJECT_DIR"

# 发行目录规范对齐 brick-game / reminders（见 reminders/scripts/build-all.sh）：
#   techfunway-notepad-<VERSION>-{linux,macos}-{amd64,arm64}.tar.gz
#   techfunway-notepad-<VERSION>-windows-amd64.zip（Windows 一律 zip）
#   docker-compose.yml（镜像 tag 钉死为当前版本）/ docker-compose.latest.yml
#   CHANGELOG.md / screenshots/ + screenshots-<VERSION>.zip
# 飞牛 .fpk 与 Docker 镜像由 build-fnpack.sh、build-docker.sh 补齐。

VERSION=$(cat VERSION)
[ -z "$VERSION" ] && echo "❌ 无法获取版本号" && exit 1
BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(git describe --always --dirty 2>/dev/null || echo "unknown")
LDFLAGS="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"
GOCACHE_DIR="${GOCACHE:-/tmp/notepad-go-cache}"

APP_NAME="techfunway-notepad"
RELEASE_DIR="release/v${VERSION}"

echo "============================================"
echo "  多平台打包 v${VERSION}"
echo "============================================"

mkdir -p "${RELEASE_DIR}"

echo ""
echo "📦 编译前端..."
cd web && npm install && npm run build && cd ..

compile() {
    GOOS=$1
    GOARCH=$2
    LABEL=$3
    echo -n "  📦 ${LABEL}... "
    
    DIR="${RELEASE_DIR}/${APP_NAME}-${VERSION}-${LABEL}"
    rm -rf "${DIR}"
    mkdir -p "${DIR}/www"

    # 复制前端构建产物到 server/static/dist 供 Go embed 使用
    cd "${PROJECT_DIR}"
    rm -rf server/static/dist
    cp -r web/dist server/static/dist

    cd server
    GOCACHE="${GOCACHE_DIR}" CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="${LDFLAGS}" -o "../${DIR}/notepad" .
    cd ..
    
    if [ "${GOOS}" = "windows" ]; then
        mv "${DIR}/notepad" "${DIR}/notepad.exe"
        echo '@echo off
cd /d "%~dp0"
notepad.exe %*' > "${DIR}/start.bat"
    else
        echo '#!/bin/bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"
chmod +x notepad 2>/dev/null || true
./notepad "$@"' > "${DIR}/start.sh"
        chmod +x "${DIR}/start.sh"
    fi
    
    cp -r web/dist/* "${DIR}/www/"

    # 复制使用说明文档（如果存在）
    for doc in release/v${VERSION}/README-*.md; do
        if [ -f "$doc" ]; then
            cp "$doc" "${DIR}/"
        fi
    done

    find "${DIR}" -name ".DS_Store" -delete 2>/dev/null || true

    echo "✅"
}

echo ""
echo "🔨 编译平台..."
echo ""
compile "linux" "amd64" "linux-amd64"
compile "linux" "arm64" "linux-arm64"
compile "darwin" "amd64" "macos-amd64"
compile "darwin" "arm64" "macos-arm64"
compile "windows" "amd64" "windows-amd64"

echo ""
echo "📦 创建压缩包..."
cd "${RELEASE_DIR}"
for dir in ${APP_NAME}-${VERSION}-*/; do
    name=$(basename "$dir")
    # 清掉旧格式的包，Windows 用 zip，其余用 tar.gz
    rm -f "${name}.tar.gz" "${name}.zip"
    if [[ "${name}" == *-windows-* ]]; then
        zip -qr "${name}.zip" "$dir" -x '*.DS_Store'
    else
        COPYFILE_DISABLE=1 tar -czf "${name}.tar.gz" --exclude='.DS_Store' "$dir"
    fi
    # 发行目录只留最终产物，未压缩的中间目录随手清掉
    rm -rf "$dir"
done
cd "$PROJECT_DIR"

# 更新日志、截图、compose 文件随发行目录分发，版本目录自包含
[ -f CHANGELOG.md ] && cp CHANGELOG.md "${RELEASE_DIR}/CHANGELOG.md"

SCREENSHOT_SRC="docs/images/v${VERSION}"
# 当前版本没有专属截图目录时，退回最新的版本截图目录（界面未变时沿用）
if [ ! -d "${SCREENSHOT_SRC}" ]; then
    SCREENSHOT_SRC=$(ls -d docs/images/v* 2>/dev/null | sort -V | tail -1 || true)
fi
if [ -n "${SCREENSHOT_SRC}" ] && [ -d "${SCREENSHOT_SRC}" ]; then
    rm -rf "${RELEASE_DIR}/screenshots"
    mkdir -p "${RELEASE_DIR}/screenshots"
    cp "${SCREENSHOT_SRC}"/*.png "${RELEASE_DIR}/screenshots/"
    find "${RELEASE_DIR}/screenshots" -name '.DS_Store' -delete 2>/dev/null || true
    (cd "${RELEASE_DIR}" && zip -qr "screenshots-v${VERSION}.zip" screenshots)
    echo "📸 已打包截图 screenshots-v${VERSION}.zip"
fi

# 版本文件中的镜像 tag 钉死为当前版本，latest 文件原样分发；
# 拷贝到目标主机后 docker compose up -d 即可运行
sed "s|techfunways/notepad:latest|techfunways/notepad:v${VERSION}|" docker-compose.yaml > "${RELEASE_DIR}/docker-compose.yml"
cp docker-compose.yaml "${RELEASE_DIR}/docker-compose.latest.yml"

echo ""
echo "✅ 完成!"
ls -lh "${RELEASE_DIR}"/ | awk 'NR>1 {print "  " $NF " (" $5 ")"}'
