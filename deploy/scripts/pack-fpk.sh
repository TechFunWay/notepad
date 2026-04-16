#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
VERSION=$(cat "$PROJECT_DIR/VERSION")
FPK_DIR="$PROJECT_DIR/deploy/fpk"
OUTPUT_DIR="$PROJECT_DIR/release/$VERSION"

mkdir -p "$OUTPUT_DIR"

echo "==> Building FPK package v${VERSION}..."

# Update version in manifest
sed -i.bak "s/version=.*/version=${VERSION}/" "$FPK_DIR/manifest" && rm -f "$FPK_DIR/manifest.bak"

# Update docker image version
sed -i.bak "s|image: wycto/notepad:.*|image: wycto/notepad:${VERSION}|" "$FPK_DIR/app/docker/docker-compose.yaml" && rm -f "$FPK_DIR/app/docker/docker-compose.yaml.bak"

# Make sure all cmd scripts are executable
chmod +x "$FPK_DIR/cmd/"*

# Copy icons if they exist
if [ -f "$PROJECT_DIR/deploy/icons/icon_64.png" ]; then
  cp "$PROJECT_DIR/deploy/icons/icon_64.png" "$FPK_DIR/ICON.PNG"
  cp "$PROJECT_DIR/deploy/icons/icon_64.png" "$FPK_DIR/app/ui/images/icon-64.png"
fi
if [ -f "$PROJECT_DIR/deploy/icons/icon_256.png" ]; then
  cp "$PROJECT_DIR/deploy/icons/icon_256.png" "$FPK_DIR/ICON_256.PNG"
  cp "$PROJECT_DIR/deploy/icons/icon_256.png" "$FPK_DIR/app/ui/images/icon-256.png"
fi

# Build FPK using fnpack
cd "$FPK_DIR"
fnpack build

# Move the generated .fpk to output directory
FPK_FILE=$(ls -t "$FPK_DIR"/*.fpk 2>/dev/null | head -1)
if [ -n "$FPK_FILE" ]; then
  mv "$FPK_FILE" "$OUTPUT_DIR/notepad_${VERSION}.fpk"
  echo "==> FPK package created: $OUTPUT_DIR/notepad_${VERSION}.fpk"
  echo "==> Size: $(du -h "$OUTPUT_DIR/notepad_${VERSION}.fpk" | cut -f1)"
else
  echo "==> ERROR: fnpack build did not produce a .fpk file"
  exit 1
fi
