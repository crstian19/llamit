# CI/CD Setup Checklist

Use this checklist to enable automated releases for Llamit.

## Prerequisites

- [ ] Repository: https://github.com/Crstian19/llamit
- [ ] Admin access to repository settings
- [ ] Access to Open VSX account
- [ ] Access to Azure DevOps account (for VS Code Marketplace)

## Step 1: Configure GitHub Secrets

### 1.1 Open VSX Token

- [ ] Go to: https://open-vsx.org/user-settings/tokens
- [ ] Login with GitHub
- [ ] Click "Create new token"
- [ ] Copy the token (starts with `oy-`)
- [ ] Go to GitHub: https://github.com/Crstian19/llamit/settings/secrets/actions
- [ ] Click "New repository secret"
- [ ] Name: `OPENVSX_TOKEN`
- [ ] Value: Paste the token
- [ ] Click "Add secret"

### 1.2 VS Code Marketplace Token

- [ ] Go to: https://dev.azure.com/
- [ ] Click user icon → Personal Access Tokens
- [ ] Click "New Token"
- [ ] Fill in:
  - **Name**: `Llamit Extension Publishing`
  - **Organization**: All accessible organizations
  - **Expiration**: 1 year (or custom)
  - **Scopes**: Custom defined
    - [ ] Check: Marketplace → Publish
- [ ] Click "Create"
- [ ] Copy the token (save it somewhere safe - you can't see it again!)
- [ ] Go to GitHub: https://github.com/Crstian19/llamit/settings/secrets/actions
- [ ] Click "New repository secret"
- [ ] Name: `VSCE_PAT`
- [ ] Value: Paste the token
- [ ] Click "Add secret"

### 1.3 Verify Secrets

- [ ] Go to: https://github.com/Crstian19/llamit/settings/secrets/actions
- [ ] Verify you see:
  - `OPENVSX_TOKEN` (configured)
  - `VSCE_PAT` (configured)

## Step 2: Prepare Test Release

### 2.1 Update Version

- [ ] Open: `vscode-extension/package.json`
- [ ] Change version to: `0.3.1` (or next appropriate version)
- [ ] Save file

### 2.2 Update Changelog

- [ ] Open: `vscode-extension/CHANGELOG.md`
- [ ] Add new section at the top:
  ```markdown
  ## [0.3.1] - 2026-02-05

  ### Added
  - Automated CI/CD pipeline for releases
  - Platform-specific packaging for better performance

  ### Changed
  - Automated marketplace publishing on merge to main
  ```
- [ ] Save file

### 2.3 Commit Changes

```bash
cd /home/crstian/Dev/llamit
git add vscode-extension/package.json vscode-extension/CHANGELOG.md
git commit -m "chore: bump version to 0.3.1 and add CI/CD pipeline"
```

## Step 3: Test the Workflow

### 3.1 Push to Main

- [ ] Ensure you're on `main` branch: `git branch --show-current`
- [ ] Push changes: `git push origin main`

### 3.2 Monitor Workflow

- [ ] Go to: https://github.com/Crstian19/llamit/actions
- [ ] Find workflow: "Publish Extension on Main Merge"
- [ ] Click to view details
- [ ] Watch the progress:
  - [ ] `check-version` job completes
  - [ ] `build-and-package` jobs complete (6 parallel)
  - [ ] `create-release` job completes
  - [ ] `publish-marketplaces` job completes

### 3.3 Check for Errors

If any job fails:
- [ ] Click on the failed job
- [ ] Read the error logs
- [ ] Fix the issue
- [ ] Push fix to main (or rollback)

## Step 4: Verify Release

### 4.1 GitHub Release

- [ ] Go to: https://github.com/Crstian19/llamit/releases
- [ ] Verify new release exists: `v0.3.1`
- [ ] Check release has 6 `.vsix` files:
  - [ ] `llamit-v0.3.1-win32-x64.vsix`
  - [ ] `llamit-v0.3.1-win32-arm64.vsix`
  - [ ] `llamit-v0.3.1-linux-x64.vsix`
  - [ ] `llamit-v0.3.1-linux-arm64.vsix`
  - [ ] `llamit-v0.3.1-darwin-x64.vsix`
  - [ ] `llamit-v0.3.1-darwin-arm64.vsix`
- [ ] Check changelog is displayed in release notes

### 4.2 VS Code Marketplace

- [ ] Go to: https://marketplace.visualstudio.com/items?itemName=Crstian.llamit
- [ ] Verify version shows: `0.3.1`
- [ ] Check all platforms are available:
  - [ ] Windows x64
  - [ ] Windows ARM64
  - [ ] Linux x64
  - [ ] Linux ARM64
  - [ ] macOS x64
  - [ ] macOS ARM64
- [ ] Try installing in VS Code: `Extensions → Search "Llamit" → Install`

### 4.3 Open VSX Registry

- [ ] Go to: https://open-vsx.org/extension/Crstian/llamit
- [ ] Verify version shows: `0.3.1`
- [ ] Check all platforms are available
- [ ] Download one `.vsix` to test (optional)

## Step 5: Test Installation

### 5.1 Test in VS Code

- [ ] Open VS Code
- [ ] Go to Extensions (`Ctrl+Shift+X`)
- [ ] Search for "Llamit"
- [ ] Click Install
- [ ] Verify it installs the correct platform version
- [ ] Test the extension:
  - [ ] Stage some files: `git add .`
  - [ ] Click Llamit button in Source Control
  - [ ] Verify commit message generates

## Step 6: Cleanup (Optional)

### 6.1 Remove Test Files (if desired)

```bash
# The workflow files are permanent, but you can remove the summary if desired:
# rm .github/CI-CD-SUMMARY.md
# rm .github/SETUP-CHECKLIST.md

# Keep these files:
# - .github/workflows/publish.yml (REQUIRED)
# - .github/workflows/release.yml (REQUIRED)
# - .github/RELEASE.md (RECOMMENDED)
```

## Troubleshooting

### Issue: Workflow doesn't trigger

**Check**:
- [ ] Did you push to `main` branch?
- [ ] Does `.github/workflows/publish.yml` exist?
- [ ] Is the workflow file valid YAML?

**Fix**:
```bash
# Check workflow syntax
cd .github/workflows
cat publish.yml | grep -E "^[^ ]" # Should show valid YAML keys
```

### Issue: "Tag already exists" - workflow skips

**Check**:
- [ ] Did you bump the version in `package.json`?
- [ ] Does the tag already exist: `git tag | grep v0.3.1`?

**Fix**:
```bash
# Bump to a new version
cd vscode-extension
# Edit package.json and increase version
git add package.json
git commit -m "chore: bump version to 0.3.2"
git push origin main
```

### Issue: Build fails

**Check**:
- [ ] Does Go CLI compile locally?
- [ ] Are there syntax errors in the code?

**Fix**:
```bash
# Test local build
cd go-cli
go build -o cli main.go

# If it fails, fix the error and push again
```

### Issue: Publishing fails

**Check**:
- [ ] Are secrets configured correctly?
- [ ] Are tokens still valid (not expired)?
- [ ] Check workflow logs for specific error

**Fix**:
- Verify secrets in GitHub Settings
- Generate new tokens if expired
- Manually publish if needed:
  ```bash
  cd vscode-extension
  export OPENVSX_TOKEN="your-token"
  export VSCE_PAT="your-pat"
  npm run publish:all
  ```

## Success Indicators

You'll know the setup is complete when:
- ✅ Workflow runs successfully on push to main
- ✅ GitHub Release created automatically
- ✅ All 6 platform packages available
- ✅ Extension published to both marketplaces
- ✅ Users can install via Extensions UI
- ✅ No manual steps required

## Next Release

For the next release, you only need to:

1. Make your changes
2. Bump version in `package.json`
3. Update `CHANGELOG.md`
4. Commit and push to `main`
5. ☕ Relax while CI/CD does the work

## Resources

- **Workflow**: `.github/workflows/publish.yml`
- **Documentation**: `.github/RELEASE.md`
- **GitHub Actions**: https://github.com/Crstian19/llamit/actions
- **Releases**: https://github.com/Crstian19/llamit/releases

## Support

If you need help:
1. Read `.github/RELEASE.md` for detailed docs
2. Check workflow logs for errors
3. Open issue: https://github.com/Crstian19/llamit/issues

---

**Last Updated**: 2026-02-05
**Status**: Ready for first test release
