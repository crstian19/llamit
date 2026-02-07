#!/bin/bash

# Quick test script for smart git diff feature
# Run this in a test git repository

echo "=== Testing Llamit Smart Git Diff ==="
echo ""

# Setup test repo
TEST_DIR="/tmp/llamit-test-$(date +%s)"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

git init
git config user.name "Test User"
git config user.email "test@example.com"

echo "Test repo created at: $TEST_DIR"
echo ""

# Scenario 1: Only unstaged changes
echo "📝 Scenario 1: Only unstaged changes"
echo "content" > unstaged.txt
git status --short
echo "➡️  Open this folder in VS Code and run Llamit"
echo "   Expected: Should generate commit for unstaged.txt"
echo ""
read -p "Press Enter after testing scenario 1..."

# Scenario 2: Only staged changes
echo ""
echo "📝 Scenario 2: Only staged changes"
git add unstaged.txt
echo "more content" > staged.txt
git add staged.txt
git status --short
echo "➡️  Run Llamit again"
echo "   Expected: Should generate commit for both files (staged)"
echo ""
read -p "Press Enter after testing scenario 2..."

# Scenario 3: Mixed (staged + unstaged)
echo ""
echo "📝 Scenario 3: Mixed - staged takes priority"
git commit -m "initial commit"
echo "staged content" > file1.txt
git add file1.txt
echo "unstaged content" > file2.txt
git status --short
echo "➡️  Run Llamit again"
echo "   Expected: Should only mention file1.txt (staged), NOT file2.txt"
echo ""
read -p "Press Enter after testing scenario 3..."

# Scenario 4: No changes
echo ""
echo "📝 Scenario 4: No changes"
git add .
git commit -m "test commit"
git status --short
echo "➡️  Run Llamit again"
echo "   Expected: Should show 'No changes to commit'"
echo ""

echo "✅ All test scenarios complete!"
echo "Test repo location: $TEST_DIR"
echo "To clean up: rm -rf $TEST_DIR"
