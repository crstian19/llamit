# Release Process

This document explains how the automated release process works for the Llamit VS Code extension.

## Overview

Llamit uses GitHub Actions to automate the entire release process:
- ✅ Automatic releases on merge to `main`
- ✅ Platform-specific packaging (6 platforms)
- ✅ GitHub Release creation with all `.vsix` files
- ✅ Automatic publishing to VS Code Marketplace
- ✅ Automatic publishing to Open VSX Registry

## Workflows

### 1. Automatic Release (`.github/workflows/publish.yml`)

**Trigger**: Push to `main` branch (includes merges)

**What it does**:
1. Reads version from `vscode-extension/package.json`
2. Checks if a git tag for that version already exists
3. If tag exists → **skips release** (prevents duplicates)
4. If tag doesn't exist → proceeds with:
   - Compiles Go binaries for 6 platforms
   - Packages 6 platform-specific `.vsix` files
   - Creates git tag `vX.Y.Z`
   - Creates GitHub Release with all `.vsix` files
   - Publishes to Open VSX Registry (if token configured)
   - Publishes to VS Code Marketplace (if token configured)

**When to use**: Normal development workflow - just merge to `main` after bumping version

### 2. Manual Release (`.github/workflows/release.yml`)

**Trigger**: Manual GitHub Release creation

**What it does**:
- Packages extension for all platforms
- Uploads `.vsix` files to the existing release
- Publishes to marketplaces

**When to use**: Emergency releases or when automatic workflow fails

## How to Release a New Version

### Prerequisites

1. All changes are merged into `main`
2. Tests are passing
3. You're ready to publish

### Step 1: Update Version

Edit `vscode-extension/package.json`:

```json
{
  "version": "0.3.1"  // Bump to new version
}
```

Version format follows [Semantic Versioning](https://semver.org/):
- **Major** (1.0.0): Breaking changes
- **Minor** (0.1.0): New features, backward compatible
- **Patch** (0.0.1): Bug fixes, backward compatible

### Step 2: Update Changelog

Edit `vscode-extension/CHANGELOG.md`:

```markdown
## [0.3.1] - 2026-02-05

### Added
- New feature X

### Fixed
- Bug Y

### Changed
- Improved Z
```

### Step 3: Commit and Push

```bash
git add vscode-extension/package.json vscode-extension/CHANGELOG.md
git commit -m "chore: bump version to 0.3.1"
git push origin main
```

### Step 4: Wait for Automation

The `publish.yml` workflow will automatically:
1. Detect the new version (0.3.1)
2. Verify tag `v0.3.1` doesn't exist
3. Build all platform packages
4. Create GitHub Release
5. Publish to marketplaces

**Duration**: ~10-15 minutes

### Step 5: Verify Release

Check that everything succeeded:

1. **GitHub Release**: https://github.com/Crstian19/llamit/releases
   - ✅ Tag `v0.3.1` exists
   - ✅ Release has 6 `.vsix` files
   - ✅ Changelog is displayed

2. **VS Code Marketplace**: https://marketplace.visualstudio.com/items?itemName=Crstian.llamit
   - ✅ Version shows as 0.3.1
   - ✅ All platforms available

3. **Open VSX Registry**: https://open-vsx.org/extension/Crstian/llamit
   - ✅ Version shows as 0.3.1
   - ✅ All platforms available

## What Happens If...?

### You merge without bumping version?

**Nothing!** The workflow checks if the tag already exists. If you're still at version `0.3.0` and tag `v0.3.0` already exists, the workflow will skip the release entirely.

This is **intentional** to prevent duplicate releases.

### The workflow fails?

Check the Actions tab: https://github.com/Crstian19/llamit/actions

Common issues:
- **Build failure**: Fix the code and push again
- **Tag already exists**: Bump the version in `package.json`
- **Publishing failure**: Check that secrets are configured (see below)

### You need to fix something after release?

1. Bump to next patch version (e.g., 0.3.0 → 0.3.1)
2. Make your fixes
3. Follow the normal release process

**Never delete tags or releases** unless absolutely necessary.

## Secrets Configuration

The workflow can publish to marketplaces if these secrets are configured in GitHub:

**Settings → Secrets and variables → Actions → New repository secret**

### Required Secrets

| Secret | Purpose | Required? | How to get |
|--------|---------|-----------|------------|
| `GITHUB_TOKEN` | Create releases | ✅ Yes | Automatic (no setup needed) |
| `OPENVSX_TOKEN` | Publish to Open VSX | ⚠️  Optional | https://open-vsx.org/user-settings/tokens |
| `VSCE_PAT` | Publish to VS Code Marketplace | ⚠️  Optional | https://dev.azure.com/ → Personal Access Tokens |

### Getting OPENVSX_TOKEN

1. Go to: https://open-vsx.org/user-settings/tokens
2. Login with GitHub
3. Click "Create new token"
4. Copy the token (starts with `oy-`)
5. Add to GitHub Secrets as `OPENVSX_TOKEN`

### Getting VSCE_PAT

1. Go to: https://dev.azure.com/
2. Click user icon → Personal Access Tokens
3. Create new token with:
   - **Name**: Llamit Extension Publishing
   - **Organization**: All accessible organizations
   - **Scopes**: Custom defined → Marketplace → Publish
4. Copy the token
5. Add to GitHub Secrets as `VSCE_PAT`

### What if secrets are not configured?

The workflow will:
- ✅ Still create the GitHub Release with all `.vsix` files
- ⚠️  Skip publishing to Open VSX (if `OPENVSX_TOKEN` missing)
- ⚠️  Skip publishing to VS Code Marketplace (if `VSCE_PAT` missing)

You can manually publish later using:

```bash
cd vscode-extension
npm run publish:all
```

## Platform-Specific Packaging

The extension is packaged separately for each platform:

| Platform | Target | Binary | File Size |
|----------|--------|--------|-----------|
| Windows (x64) | `win32-x64` | `cli.exe` | ~5.5 MB |
| Windows (ARM64) | `win32-arm64` | `cli.exe` | ~5.5 MB |
| Linux (x64) | `linux-x64` | `cli` | ~5.5 MB |
| Linux (ARM64) | `linux-arm64` | `cli` | ~5.5 MB |
| macOS (Intel) | `darwin-x64` | `cli` | ~5.5 MB |
| macOS (Apple Silicon) | `darwin-arm64` | `cli` | ~5.5 MB |

**Benefits**:
- 📦 Users download only their platform (~5.5 MB vs ~28 MB universal)
- ⚡ 5x faster installation
- 💾 80% bandwidth savings
- 🎯 Follows VS Code best practices

**How it works**:
- VS Code automatically detects user's platform
- Marketplace serves the correct `.vsix` file
- No user configuration needed

## Manual Publishing (Emergency)

If automatic workflows fail, you can publish manually:

### 1. Build all platforms

```bash
cd vscode-extension
npm run package:all
```

This creates 6 `.vsix` files in `dist/`

### 2. Publish to Open VSX

```bash
cd vscode-extension
OPENVSX_TOKEN="your-token" ./scripts/publish-all-platforms.sh openvsx
```

### 3. Publish to VS Code Marketplace

```bash
cd vscode-extension
VSCE_PAT="your-token" ./scripts/publish-all-platforms.sh marketplace
```

## Rollback Procedure

If you need to undo a release:

### 1. Delete the release

```bash
gh release delete v0.3.1 --yes
```

Or via GitHub UI: Releases → Delete release

### 2. Delete the tag

```bash
git tag -d v0.3.1
git push --delete origin v0.3.1
```

### 3. Unpublish from marketplaces (if needed)

- **VS Code Marketplace**: Can't unpublish, but can deprecate
- **Open VSX**: Can unpublish via web UI

### 4. Fix and re-release

1. Fix the issue
2. Bump version to next patch (0.3.1 → 0.3.2)
3. Follow normal release process

## Monitoring

### GitHub Actions

View all workflow runs: https://github.com/Crstian19/llamit/actions

Each workflow shows:
- ✅ Success/failure status
- ⏱️  Duration
- 📋 Detailed logs for each step

### Email Notifications

GitHub sends email notifications for:
- ❌ Workflow failures
- ✅ Successful releases (if configured)

Configure in: GitHub Settings → Notifications

## Troubleshooting

### "Tag already exists" error

**Cause**: You're trying to release version X.Y.Z but tag `vX.Y.Z` already exists.

**Solution**: Bump the version in `package.json` to a new version.

### Publishing fails but release succeeds

**Cause**: GitHub Release created but marketplace publishing failed.

**Solution**:
1. Check secrets are configured
2. Manually publish using scripts (see Manual Publishing section)
3. The `.vsix` files are already in the GitHub Release

### Build fails on specific platform

**Cause**: Go compilation or packaging error for one platform.

**Solution**:
1. Check the logs in GitHub Actions
2. Test locally: `VSCE_TARGET=win32-x64 node scripts/build-go.js`
3. Fix the issue and push again

## Best Practices

### ✅ DO

- Bump version before merging to `main`
- Update `CHANGELOG.md` with each release
- Test locally before releasing: `npm run package:all`
- Use semantic versioning
- Keep releases small and focused

### ❌ DON'T

- Don't merge to `main` without bumping version (if you want a release)
- Don't manually create tags (let the workflow do it)
- Don't delete releases unless absolutely necessary
- Don't skip changelog updates
- Don't release breaking changes in patch versions

## FAQ

### Q: Can I release from a branch other than `main`?

**A**: No, the automatic workflow only triggers on `main`. You'd need to manually create a release.

### Q: How long does the release process take?

**A**: ~10-15 minutes total:
- Building: ~8 minutes (6 platforms in parallel)
- Publishing: ~2-3 minutes per marketplace

### Q: Can I test the workflow without actually releasing?

**A**: Not easily. The workflow is designed to be production-only. Test packaging locally instead:

```bash
cd vscode-extension
npm run package:all
```

### Q: What if I want to release a beta version?

**A**: Use a version like `0.4.0-beta.1` in `package.json`. The workflow will create a release, but you may want to mark it as "pre-release" in GitHub.

### Q: How do I know if secrets are configured?

**A**: Check the workflow logs. If publishing is skipped, you'll see:

```
⚠️  Open VSX Registry: Skipped (no OPENVSX_TOKEN configured)
```

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [VS Code Extension Publishing](https://code.visualstudio.com/api/working-with-extensions/publishing-extension)
- [Open VSX Publishing](https://github.com/eclipse/openvsx/wiki/Publishing-Extensions)
- [Semantic Versioning](https://semver.org/)

## Support

If you encounter issues with the release process:

1. Check GitHub Actions logs
2. Read this documentation
3. Open an issue: https://github.com/Crstian19/llamit/issues
