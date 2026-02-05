#!/bin/bash
set -e

if [ -z "$OPENVSX_TOKEN" ]; then
  echo "Error: OPENVSX_TOKEN environment variable not set"
  echo "Export your Open VSX token before running this script:"
  echo "  export OPENVSX_TOKEN=your_token_here"
  exit 1
fi

# Check if dist directory exists
if [ ! -d "dist" ]; then
  echo "Error: dist/ directory not found"
  echo "Run 'npm run package:all' first to create the .vsix files"
  exit 1
fi

# Check if there are any .vsix files
VSIX_COUNT=$(ls -1 dist/*.vsix 2>/dev/null | wc -l)
if [ "$VSIX_COUNT" -eq 0 ]; then
  echo "Error: No .vsix files found in dist/"
  echo "Run 'npm run package:all' first"
  exit 1
fi

echo "Found ${VSIX_COUNT} packages to publish"
echo ""

# Publish each .vsix to Open VSX
for vsix in dist/*.vsix; do
  echo "Publishing $(basename $vsix)..."
  npx ovsx publish "$vsix" -p "$OPENVSX_TOKEN"
  echo "✓ Published"
  echo ""
done

echo "All platforms published successfully to Open VSX!"
echo ""
echo "Note: To publish to Visual Studio Marketplace, use:"
echo "  vsce publish --packagePath dist/*.vsix"
