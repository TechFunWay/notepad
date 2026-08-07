#!/bin/bash

set -e

PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
cd "$PROJECT_DIR"

VERSION=$(cat VERSION)
[ -z "$VERSION" ] && echo "❌ 无法获取版本号" && exit 1

APP_NAME="techfunway-notepad"
IMAGE_NAME="techfunways/notepad"
RELEASE_DIR="release/v${VERSION}"
BUILDER_NAME="${BUILDER_NAME:-notepad-multiarch}"

[ ! -d "${RELEASE_DIR}/${APP_NAME}-${VERSION}-linux-amd64" ] && echo "❌ 请先运行 ./scripts/build-all.sh" && exit 1

echo "============================================"
echo "  Docker 多平台镜像构建 v${VERSION}"
echo "============================================"

# 删除 .DS_Store
find "${RELEASE_DIR}" -name ".DS_Store" -delete 2>/dev/null || true

# 切换到 default context
docker context use default 2>/dev/null || true

# 确保有可用的 buildx builder（多平台交叉编译需要 docker-container 驱动）
if ! docker buildx inspect "${BUILDER_NAME}" >/dev/null 2>&1; then
    echo ""
    echo "🔧 创建 buildx builder: ${BUILDER_NAME}"
    docker buildx create --name "${BUILDER_NAME}" --driver docker-container --bootstrap >/dev/null
fi
docker buildx use "${BUILDER_NAME}"

echo ""
echo "🔨 构建多平台镜像 (linux/amd64 + linux/arm64)..."
echo ""

cd "${RELEASE_DIR}"

# 创建 Dockerfile
cat > Dockerfile << EOF
FROM scratch

ARG TARGETARCH
ARG VERSION

LABEL org.opencontainers.image.title="techfunway-notepad"
LABEL org.opencontainers.image.description="一款轻量、精美的多用户记事本 Web 应用"
LABEL org.opencontainers.image.version="\${VERSION}"
LABEL org.opencontainers.image.source="https://gitee.com/TechFunWay/notepad"

COPY ${APP_NAME}-${VERSION}-linux-\${TARGETARCH}/notepad /app/notepad
COPY ${APP_NAME}-${VERSION}-linux-\${TARGETARCH}/www /app/www

# 端口
EXPOSE 8904

# 数据卷：数据库与上传文件
VOLUME ["/app/data"]

# 环境变量
ENV PORT=8904
ENV DB_PATH=/app/data/notepad.db
ENV DATA_DIR=/app/data
ENV WWW_DIR=/app/www
ENV TZ=Asia/Shanghai

ENTRYPOINT ["/app/notepad"]
EOF

# 构建多平台镜像（导出为 docker 本地镜像，需要 --output type=docker）
docker buildx build \
    --builder "${BUILDER_NAME}" \
    --platform linux/amd64,linux/arm64 \
    --output type=docker \
    --build-arg VERSION=${VERSION} \
    -t "${IMAGE_NAME}:v${VERSION}" \
    -t "${IMAGE_NAME}:latest" \
    --load \
    .

rm -f Dockerfile
cd "${PROJECT_DIR}"

# 重新切回 default builder（避免影响后续 docker build）
docker buildx use default 2>/dev/null || true

echo ""
echo "✅ 完成!"
echo ""
echo "📦 本地镜像:"
docker images "${IMAGE_NAME}" --format "  {{.Repository}}:{{.Tag}} ({{.Size}})" 2>/dev/null | sort -u
echo ""
echo "📝 多平台验证（已通过 buildx 多架构构建）:"
echo "  - linux/amd64  ✅"
echo "  - linux/arm64  ✅"
echo ""
echo "  提示：Docker Desktop 在 macOS 上 --load 多平台镜像仅加载当前平台。"
echo "  如需在其他架构机器上直接拉取，请推送到 registry："
echo "    docker buildx build --builder ${BUILDER_NAME} --platform linux/amd64,linux/arm64 \\"
echo "      --push -t ${IMAGE_NAME}:v${VERSION} ."
echo ""
echo "💡 运行命令示例（推荐生产环境设置 JWT_SECRET 与时区）:"
echo "  docker run -d --name notepad \\"
echo "    -p 8904:8904 \\"
echo "    -v \$(pwd)/data:/app/data \\"
echo "    -e JWT_SECRET=请改为你自己的强随机密钥 \\"
echo "    -e TZ=Asia/Shanghai \\"
echo "    --restart unless-stopped \\"
echo "    ${IMAGE_NAME}:v${VERSION}"
echo ""
echo "  访问: http://localhost:8904"
echo "  端口: 8904  |  数据目录: /app/data  |  环境变量: PORT / JWT_SECRET / TZ"
