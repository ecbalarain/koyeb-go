#!/bin/bash
# Build script for optimizing frontend assets for production

set -e

FRONTEND_DIR="cloudflare-pages-frontend"
BUILD_DIR="${FRONTEND_DIR}/dist"

echo "🚀 Starting frontend optimization build..."

# Clean previous build
if [ -d "$BUILD_DIR" ]; then
  echo "🧹 Cleaning previous build..."
  rm -rf "$BUILD_DIR"
fi

mkdir -p "$BUILD_DIR"

echo "📦 Copying static files..."
cp -r "${FRONTEND_DIR}"/*.html "${BUILD_DIR}/" 2>/dev/null || true
cp -r "${FRONTEND_DIR}"/*.json "${BUILD_DIR}/" 2>/dev/null || true
cp -r "${FRONTEND_DIR}"/admin "${BUILD_DIR}/" 2>/dev/null || true
cp "${FRONTEND_DIR}"/_headers "${BUILD_DIR}/" 2>/dev/null || true

echo "✨ Frontend files copied to ${BUILD_DIR}"
echo ""
echo "📝 Optimization recommendations:"
echo "  1. Use online tools to minify HTML:"
echo "     - https://www.toptal.com/developers/html-minifier"
echo "  2. Convert images to WebP format:"
echo "     - cwebp input.jpg -q 80 -o output.webp"
echo "  3. Add lazy loading to images:"
echo "     - <img loading=\"lazy\" src=\"image.jpg\" />"
echo "  4. Enable Cloudflare automatic minification:"
echo "     - JavaScript, CSS, HTML in Speed > Optimization"
echo "  5. Enable Cloudflare image optimization:"
echo "     - Polish (Lossless or Lossy)"
echo "     - WebP conversion"
echo "     - Check 'Mirage' for lazy loading"
echo ""
echo "✅ Build complete! Deploy ${BUILD_DIR} to Cloudflare Pages"
