#!/bin/bash

set -e

PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
cd "$PROJECT_DIR"

VERSION=$(cat VERSION)
[ -z "$VERSION" ] && echo "❌ 无法获取版本号" && exit 1

APP_NAME="techfunway-notepad"
IMAGE_NAME="techfunways/notepad"
RELEASE_DIR="${PROJECT_DIR}/release/v${VERSION}"
BUILDER_NAME="${BUILDER_NAME:-notepad-multiarch}"
OCI_FILE="${RELEASE_DIR}/${APP_NAME}-v${VERSION}-multiarch.oci.tar"

# build-all.sh 压缩后即清理中间目录；Docker 构建需要 linux 平台目录，
# 缺失时从对应的 tar.gz 解包，用完在脚本末尾清理
DOCKER_TMP_DIRS=()
for arch in amd64 arm64; do
    DIR_NAME="${APP_NAME}-${VERSION}-linux-${arch}"
    if [ ! -d "${RELEASE_DIR}/${DIR_NAME}" ]; then
        TAR_FILE="${RELEASE_DIR}/${DIR_NAME}.tar.gz"
        [ ! -f "${TAR_FILE}" ] && echo "❌ 缺少 ${TAR_FILE}，请先运行 ./scripts/build-all.sh" && exit 1
        tar -xzf "${TAR_FILE}" -C "${RELEASE_DIR}"
        DOCKER_TMP_DIRS+=("${RELEASE_DIR}/${DIR_NAME}")
    fi
done

echo "============================================"
echo "  Docker 多平台 OCI 归档 v${VERSION}"
echo "============================================"

# 删除 .DS_Store
find "${RELEASE_DIR}" -name ".DS_Store" -delete 2>/dev/null || true

# 切换到 default context
docker context use default 2>/dev/null || true

# 合并 manifest 的 OCI 归档必须用 docker-container 驱动的 builder，
# 默认驱动不推 registry 导不出多平台合并结果
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

# 导出为本地 OCI 归档（不推送 registry），随发行目录分发
docker buildx build \
    --builder "${BUILDER_NAME}" \
    --platform linux/amd64,linux/arm64 \
    --output "type=oci,dest=${OCI_FILE}" \
    --build-arg VERSION=${VERSION} \
    -t "${IMAGE_NAME}:v${VERSION}" \
    -t "${IMAGE_NAME}:latest" \
    .

rm -f Dockerfile
cd "${PROJECT_DIR}"

# 清理解包出来的临时平台目录，发行目录只留最终产物
for tmp_dir in "${DOCKER_TMP_DIRS[@]}"; do
    rm -rf "${tmp_dir}"
done

# 重新切回 default builder（避免影响后续 docker build）
docker buildx use default 2>/dev/null || true

echo ""
echo "✅ 多平台 OCI 归档完成: ${OCI_FILE}"
echo ""
echo "  目标机器载入: docker load -i ${APP_NAME}-v${VERSION}-multiarch.oci.tar"
echo ""
