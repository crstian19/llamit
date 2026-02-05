#!/bin/bash
set -e

# Usage: ./publish-all-platforms.sh [openvsx|marketplace|both]
TARGET="${1:-both}"

# Validate target
if [[ ! "$TARGET" =~ ^(openvsx|marketplace|both)$ ]]; then
  echo "Usage: $0 [openvsx|marketplace|both]"
  echo ""
  echo "Examples:"
  echo "  $0 openvsx      # Publish to Open VSX Registry only"
  echo "  $0 marketplace  # Publish to VS Code Marketplace only"
  echo "  $0 both         # Publish to both (default)"
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

# Publish to Open VSX Registry
if [[ "$TARGET" == "openvsx" ]] || [[ "$TARGET" == "both" ]]; then
  if [ -z "$OPENVSX_TOKEN" ]; then
    echo "⚠️  Warning: OPENVSX_TOKEN not set, skipping Open VSX Registry"
    echo "   Set it with: export OPENVSX_TOKEN=your_token_here"
    echo ""
  else
    echo "📤 Publishing to Open VSX Registry..."
    echo ""
    for vsix in dist/*.vsix; do
      echo "Publishing $(basename $vsix) to Open VSX..."
      npx ovsx publish "$vsix" -p "$OPENVSX_TOKEN"
      echo "✓ Published to Open VSX"
      echo ""
    done
    echo "✅ All platforms published to Open VSX Registry!"
    echo ""
  fi
fi

# Publish to VS Code Marketplace
if [[ "$TARGET" == "marketplace" ]] || [[ "$TARGET" == "both" ]]; then
  if [ -z "$VSCE_PAT" ]; then
    echo "⚠️  Warning: VSCE_PAT not set, skipping VS Code Marketplace"
    echo "   Set it with: export VSCE_PAT=your_token_here"
    echo ""
  else
    echo "📤 Publishing to VS Code Marketplace..."
    echo ""
    for vsix in dist/*.vsix; do
      echo "Publishing $(basename $vsix) to Marketplace..."
      npx @vscode/vsce publish --packagePath "$vsix" -p "$VSCE_PAT"
      echo "✓ Published to Marketplace"
      echo ""
    done
    echo "✅ All platforms published to VS Code Marketplace!"
    echo ""
  fi
fi

echo "🎉 Publishing complete!"
