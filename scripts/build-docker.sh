#!/bin/bash

set -e

PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
cd "$PROJECT_DIR"

VERSION=$(cat VERSION)
[ -z "$VERSION" ] && echo "❌ 无法获取版本号" && exit 1

APP_NAME="techfunway-notepad"
IMAGE_NAME="techfunways/notepad"
RELEASE_DIR="release/v${VERSION}"

[ ! -d "${RELEASE_DIR}/${APP_NAME}-${VERSION}-linux-amd64" ] && echo "❌ 请先运行 ./scripts/build-all.sh" && exit 1

echo "============================================"
echo "  Docker 多平台镜像构建 v${VERSION}"
echo "============================================"

# 清理旧镜像
docker rmi $(docker images "${IMAGE_NAME}" -q) 2>/dev/null || true

# 删除 .DS_Store
find "${RELEASE_DIR}" -name ".DS_Store" -delete 2>/dev/null || true

# 切换到 default context
docker context use default 2>/dev/null || true

echo ""
echo "🔨 构建多平台镜像..."
echo ""

cd "${RELEASE_DIR}"

# 创建 Dockerfile
cat > Dockerfile << EOF
FROM scratch

ARG TARGETARCH
ARG VERSION

COPY ${APP_NAME}-${VERSION}-linux-\${TARGETARCH}/notepad /app/notepad
COPY ${APP_NAME}-${VERSION}-linux-\${TARGETARCH}/www /app/www

VOLUME ["/app/data"]
EXPOSE 8904

ENV PORT=8904
ENV DB_PATH=/app/data/notepad.db
ENV DATA_DIR=/app/data
ENV WWW_DIR=/app/www

ENTRYPOINT ["/app/notepad"]
EOF

# 构建多平台镜像
docker buildx build \
    --builder default \
    --platform linux/amd64,linux/arm64 \
    --output type=docker \
    --build-arg VERSION=${VERSION} \
    -t "${IMAGE_NAME}:v${VERSION}" \
    -t "${IMAGE_NAME}:latest" \
    .

rm -f Dockerfile
cd "${PROJECT_DIR}"

echo ""
echo "✅ 完成!"
echo ""
echo "📦 本地镜像:"
docker images "${IMAGE_NAME}" --format "  {{.Repository}}:{{.Tag}}" | sort -u
echo ""
echo "💡 运行: docker run -p 8904:8904 -v ./data:/app/data ${IMAGE_NAME}:v${VERSION}"
