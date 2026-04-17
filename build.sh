#!/bin/bash

set -e

PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$PROJECT_DIR"

VERSION=$(cat VERSION)
[ -z "$VERSION" ] && echo "❌ 无法获取版本号" && exit 1

echo "============================================"
echo "  Notepad v${VERSION} 一键打包"
echo "============================================"
echo ""

echo "📦 [1/3] 打包所有平台..."
./scripts/build-all.sh

echo ""
echo "📦 [2/3] 打包飞牛应用..."
./scripts/build-fnpack.sh

echo ""
echo "📦 [3/3] 构建 Docker 镜像..."
./scripts/build-docker.sh

echo ""
echo "============================================"
echo "  ✅ 全部完成!"
echo "============================================"
echo ""
echo "📁 release/v${VERSION}/"
ls -lh "release/v${VERSION}/"*.tar.gz 2>/dev/null | grep -v "all" | awk '{print "  平台包: " $NF " (" $5 ")"}'
ls -lh "release/v${VERSION}/"*.fpk 2>/dev/null | awk '{print "  飞牛包: " $NF " (" $5 ")"}'
echo ""
