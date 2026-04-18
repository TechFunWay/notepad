#!/bin/bash

set -e

PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
cd "$PROJECT_DIR"

VERSION=$(cat VERSION)
[ -z "$VERSION" ] && echo "❌ 无法获取版本号" && exit 1

APP_NAME="techfunway-notepad"
RELEASE_DIR="release/v${VERSION}"

[ ! -d "${RELEASE_DIR}/${APP_NAME}-${VERSION}-linux-amd64" ] && echo "❌ 请先运行 ./scripts/build-all.sh" && exit 1

echo "============================================"
echo "  飞牛 fnOS 打包 v${VERSION}"
echo "============================================"

cp fnpack/manifest fnpack/manifest.bak
sed -i '' "s/^version.*/version               = ${VERSION}/" fnpack/manifest

build_fnpack() {
    ARCH=$1
    PLATFORM=$2
    LABEL=$3
    echo -n "  📦 ${LABEL}.fpk... "
    
    DIR="${RELEASE_DIR}/.fnpack-${ARCH}"
    rm -rf "${DIR}"
    mkdir -p "${DIR}"
    
    # 复制飞牛模板，排除 .DS_Store
    find fnpack -name ".DS_Store" -delete 2>/dev/null || true
    cp -r fnpack/* "${DIR}/"
    sed -i '' "s/^platform.*/platform              = ${PLATFORM}/" "${DIR}/manifest"
    
    mkdir -p "${DIR}/app/server" "${DIR}/app/ui/images"
    cp "${RELEASE_DIR}/${APP_NAME}-${VERSION}-linux-${ARCH}/notepad" "${DIR}/app/server/"
    cp -r "${RELEASE_DIR}/${APP_NAME}-${VERSION}-linux-${ARCH}/www/"* "${DIR}/app/ui/"
    
    # Create ui config
    cat > "${DIR}/app/ui/config" << 'EOCONFIG'
{
    ".url": {
        "techfunway-notepad.Application": {
            "title": "记事本",
            "icon": "images/icon_{0}.png",
            "type": "url",
            "protocol": "http",
            "port": "${wizard_app_port}",
            "url": "/",
            "allUsers": false
        }
    }
}
EOCONFIG
    
    # 删除 .DS_Store
    find "${DIR}" -name ".DS_Store" -delete 2>/dev/null || true
    
    cd "${DIR}"
    fnpack build > /dev/null 2>&1
    cd "$PROJECT_DIR"
    
    mv "${DIR}/${APP_NAME}.fpk" "${RELEASE_DIR}/${APP_NAME}-v${VERSION}-${LABEL}.fpk"
    rm -rf "${DIR}"
    echo "✅"
}

echo ""
build_fnpack "arm64" "arm" "fnos-arm64"
build_fnpack "amd64" "x86" "fnos-amd64"

mv fnpack/manifest.bak fnpack/manifest

echo ""
echo "✅ 完成!"
ls -lh "${RELEASE_DIR}"/*.fpk | awk '{print "  " $NF " (" $5 ")"}'
