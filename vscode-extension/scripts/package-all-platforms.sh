#!/bin/bash
set -e

# Get version from package.json
VERSION=$(node -p "require('./package.json').version")

# Create directory for the .vsix files
mkdir -p dist

# Compile TypeScript first (shared across all platforms)
echo "Compiling TypeScript..."
npm run compile
echo ""

# Package for each platform
echo "Packaging for all platforms (version ${VERSION})..."
echo ""

# Define platforms
PLATFORMS=("win32-x64" "win32-arm64" "linux-x64" "linux-arm64" "darwin-x64" "darwin-arm64")
PLATFORM_NAMES=(
  "Windows x64"
  "Windows ARM64"
  "Linux x64"
  "Linux ARM64"
  "macOS Intel (x64)"
  "macOS Apple Silicon (ARM64)"
)

# Package each platform
for i in "${!PLATFORMS[@]}"; do
  TARGET="${PLATFORMS[$i]}"
  NAME="${PLATFORM_NAMES[$i]}"

  echo "Building and packaging ${NAME}..."

  # Build only the binary for this target
  VSCE_TARGET=$TARGET node scripts/build-go.js

  # Package with the target, skip the prepublish hook
  vsce package --target $TARGET -o dist/llamit-${VERSION}-${TARGET}.vsix --no-dependencies

  echo "✓ ${NAME} packaged"
  echo ""
done

echo "All packages created in dist/:"
ls -lh dist/*.vsix
