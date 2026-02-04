#!/bin/bash

# Llamit VS Code Extension - Test Setup Script
# This script installs dependencies and runs tests

set -e

echo "======================================"
echo "Llamit Extension Test Setup"
echo "======================================"
echo ""

# Check if we're in the right directory
if [ ! -f "package.json" ]; then
    echo "Error: package.json not found. Please run this script from the vscode-extension directory."
    exit 1
fi

# Install dependencies
echo "📦 Installing dependencies..."
npm install
echo "✅ Dependencies installed"
echo ""

# Compile TypeScript
echo "🔨 Compiling TypeScript..."
npm run compile
echo "✅ TypeScript compiled"
echo ""

# Run unit tests
echo "🧪 Running unit tests..."
npm run test:unit
echo "✅ Unit tests completed"
echo ""

# Ask if user wants to run integration tests
read -p "🤔 Run integration tests? (requires VS Code runtime, slower) [y/N]: " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🧪 Running integration tests..."
    npm test
    echo "✅ Integration tests completed"
else
    echo "ℹ️  Skipping integration tests"
fi

echo ""
echo "======================================"
echo "✅ Test setup complete!"
echo "======================================"
echo ""
echo "Available commands:"
echo "  npm run compile     - Compile TypeScript"
echo "  npm run watch       - Watch mode for development"
echo "  npm run test:unit   - Run unit tests (fast)"
echo "  npm test            - Run integration tests (slower)"
echo ""
