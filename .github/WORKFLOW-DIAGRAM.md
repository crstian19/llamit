# CI/CD Workflow Diagram

## Overview

This document provides a visual representation of the automated CI/CD pipeline.

## Workflow Trigger

```
┌─────────────────────────────────────────────────────────────┐
│                     Developer Actions                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ 1. Edit package.json (bump version)
                              │ 2. Update CHANGELOG.md
                              │ 3. git commit -m "chore: bump version"
                              │ 4. git push origin main
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│              GitHub: Push to main branch                     │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ Triggers: .github/workflows/publish.yml
                              │
                              ▼
```

## Job Flow

```
┌─────────────────────────────────────────────────────────────┐
│  JOB 1: check-version                                        │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│                                                               │
│  Steps:                                                       │
│  1. Checkout code                                             │
│  2. Read version from package.json → v0.3.1                  │
│  3. Check if git tag "v0.3.1" exists                          │
│                                                               │
│  Outputs:                                                     │
│  - version: "0.3.1"                                           │
│  - should_release: "true" | "false"                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ If should_release = false → STOP
                              │ If should_release = true → CONTINUE
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  JOB 2: build-and-package (Matrix Strategy)                 │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│                                                               │
│  Runs 6 jobs in PARALLEL:                                    │
│                                                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ win32-x64   │  │ win32-arm64 │  │ linux-x64   │         │
│  │             │  │             │  │             │         │
│  │ 1. Setup    │  │ 1. Setup    │  │ 1. Setup    │         │
│  │ 2. Build Go │  │ 2. Build Go │  │ 2. Build Go │         │
│  │ 3. Compile  │  │ 3. Compile  │  │ 3. Compile  │         │
│  │ 4. Package  │  │ 4. Package  │  │ 4. Package  │         │
│  │ 5. Upload   │  │ 5. Upload   │  │ 5. Upload   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│                                                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ linux-arm64 │  │ darwin-x64  │  │ darwin-arm64│         │
│  │             │  │             │  │             │         │
│  │ 1. Setup    │  │ 1. Setup    │  │ 1. Setup    │         │
│  │ 2. Build Go │  │ 2. Build Go │  │ 2. Build Go │         │
│  │ 3. Compile  │  │ 3. Compile  │  │ 3. Compile  │         │
│  │ 4. Package  │  │ 4. Package  │  │ 4. Package  │         │
│  │ 5. Upload   │  │ 5. Upload   │  │ 5. Upload   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│                                                               │
│  Artifacts Created:                                           │
│  - llamit-v0.3.1-win32-x64.vsix      (~5.5MB)               │
│  - llamit-v0.3.1-win32-arm64.vsix    (~5.5MB)               │
│  - llamit-v0.3.1-linux-x64.vsix      (~5.5MB)               │
│  - llamit-v0.3.1-linux-arm64.vsix    (~5.5MB)               │
│  - llamit-v0.3.1-darwin-x64.vsix     (~5.5MB)               │
│  - llamit-v0.3.1-darwin-arm64.vsix   (~5.5MB)               │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ Wait for all 6 jobs to complete
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  JOB 3: create-release                                       │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│                                                               │
│  Steps:                                                       │
│  1. Download all 6 .vsix artifacts                            │
│  2. Create git tag: v0.3.1                                    │
│  3. Push tag to GitHub                                        │
│  4. Extract changelog for v0.3.1                              │
│  5. Create GitHub Release:                                    │
│     - Title: v0.3.1                                           │
│     - Body: Changelog + download instructions                 │
│     - Assets: All 6 .vsix files                               │
│                                                               │
│  Outputs:                                                     │
│  - upload_url: "https://uploads.github.com/..."              │
└─────────────────────────────────────────────────────────────┘
                              │
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  JOB 4: publish-marketplaces                                 │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│                                                               │
│  Steps:                                                       │
│  1. Download all 6 .vsix artifacts                            │
│  2. Publish to Open VSX (if OPENVSX_TOKEN exists):           │
│     - llamit-v0.3.1-win32-x64.vsix                           │
│     - llamit-v0.3.1-win32-arm64.vsix                         │
│     - llamit-v0.3.1-linux-x64.vsix                           │
│     - llamit-v0.3.1-linux-arm64.vsix                         │
│     - llamit-v0.3.1-darwin-x64.vsix                          │
│     - llamit-v0.3.1-darwin-arm64.vsix                        │
│  3. Publish to VS Code Marketplace (if VSCE_PAT exists):     │
│     - (same 6 files)                                          │
│                                                               │
│  Results:                                                     │
│  ✅ Open VSX: Published (or ⚠️ Skipped if no token)         │
│  ✅ Marketplace: Published (or ⚠️ Skipped if no token)      │
└─────────────────────────────────────────────────────────────┘
                              │
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    ✅ RELEASE COMPLETE                       │
│                                                               │
│  Available at:                                                │
│  - GitHub: github.com/Crstian19/llamit/releases/v0.3.1      │
│  - Open VSX: open-vsx.org/extension/Crstian/llamit          │
│  - Marketplace: marketplace.visualstudio.com/.../llamit      │
└─────────────────────────────────────────────────────────────┘
```

## Job Dependencies

```
┌──────────────────┐
│  check-version   │
└────────┬─────────┘
         │
         │ depends on
         ▼
┌──────────────────────┐
│ build-and-package    │ (6 parallel jobs)
│  - win32-x64        │
│  - win32-arm64      │
│  - linux-x64        │
│  - linux-arm64      │
│  - darwin-x64       │
│  - darwin-arm64     │
└────────┬─────────────┘
         │
         │ depends on
         ▼
┌──────────────────┐
│  create-release  │
└────────┬─────────┘
         │
         │ depends on
         ▼
┌────────────────────┐
│ publish-marketplaces│
└────────────────────┘
```

## Decision Points

### Check Version Decision

```
┌─────────────────────────────────────┐
│ Read version from package.json      │
│ version = "0.3.1"                   │
└───────────────┬─────────────────────┘
                │
                ▼
┌─────────────────────────────────────┐
│ Check: Does git tag "v0.3.1" exist? │
└───────────────┬─────────────────────┘
                │
        ┌───────┴────────┐
        │                │
        ▼                ▼
    ┌──────┐         ┌──────┐
    │ YES  │         │  NO  │
    └──┬───┘         └───┬──┘
       │                 │
       │                 │
       ▼                 ▼
┌─────────────┐  ┌──────────────────┐
│ SKIP RELEASE│  │ CONTINUE RELEASE │
│ (prevents   │  │ (creates new     │
│  duplicates)│  │  release)        │
└─────────────┘  └──────────────────┘
```

### Publishing Decision

```
┌────────────────────────────────────────┐
│ Check: Is OPENVSX_TOKEN configured?   │
└───────────────┬────────────────────────┘
                │
        ┌───────┴────────┐
        │                │
        ▼                ▼
    ┌──────┐         ┌──────┐
    │ YES  │         │  NO  │
    └──┬───┘         └───┬──┘
       │                 │
       ▼                 ▼
┌─────────────┐  ┌─────────────────┐
│ PUBLISH TO  │  │ SKIP Open VSX   │
│ OPEN VSX    │  │ (warn in logs)  │
└─────────────┘  └─────────────────┘

┌────────────────────────────────────────┐
│ Check: Is VSCE_PAT configured?        │
└───────────────┬────────────────────────┘
                │
        ┌───────┴────────┐
        │                │
        ▼                ▼
    ┌──────┐         ┌──────┐
    │ YES  │         │  NO  │
    └──┬───┘         └───┬──┘
       │                 │
       ▼                 ▼
┌─────────────┐  ┌─────────────────┐
│ PUBLISH TO  │  │ SKIP Marketplace│
│ MARKETPLACE │  │ (warn in logs)  │
└─────────────┘  └─────────────────┘
```

## Timeline

```
Time (minutes)
│
0 ────┐
      │ check-version (30s)
      │
1 ────┼─────────────────────────────────────────────────────┐
      │                                                       │
      │ build-and-package (8 min, parallel)                 │
      │   ┌──────────────────────────────────────┐          │
      │   │  All 6 platforms build simultaneously │          │
      │   └──────────────────────────────────────┘          │
      │                                                       │
9 ────┼─────────────────────────────────────────────────────┘
      │
      │ create-release (1 min)
      │   - Create tag
      │   - Create release
      │   - Upload assets
      │
10 ───┼──────────┐
      │          │
      │ publish-marketplaces (3-5 min)
      │   - Publish to Open VSX
      │   - Publish to VS Code Marketplace
      │
15 ───┴──────────┘
      │
      ✅ DONE

Total: ~10-15 minutes
```

## Error Handling

### Build Failure

```
┌────────────────────────────┐
│ build-and-package job fails│
└──────────┬─────────────────┘
           │
           ▼
┌──────────────────────────────┐
│ Workflow STOPS               │
│ - No release created         │
│ - No tag created             │
│ - No publishing attempted    │
│ - Developer notified via     │
│   GitHub Actions UI/email    │
└──────────────────────────────┘
```

### Publishing Failure

```
┌────────────────────────────────┐
│ publish-marketplaces job fails │
└──────────┬─────────────────────┘
           │
           ▼
┌──────────────────────────────────┐
│ Workflow continues               │
│ - Release already created ✅     │
│ - Tag already created ✅         │
│ - .vsix files in GitHub ✅      │
│ - Can manually publish later     │
└──────────────────────────────────┘
```

## Security Flow

```
┌─────────────────────────────────────┐
│ Secrets (encrypted in GitHub)      │
│ - OPENVSX_TOKEN                     │
│ - VSCE_PAT                          │
│ - GITHUB_TOKEN (auto-generated)     │
└───────────────┬─────────────────────┘
                │
                │ Injected at runtime
                │ Never logged
                │
                ▼
┌─────────────────────────────────────┐
│ Workflow execution                  │
│ - Only visible in current job       │
│ - Masked in logs (****)             │
│ - Expires after job completion      │
└─────────────────────────────────────┘
```

## Platform Matrix Strategy

The build job uses GitHub Actions matrix strategy for parallel execution:

```yaml
strategy:
  matrix:
    target:
      - win32-x64
      - win32-arm64
      - linux-x64
      - linux-arm64
      - darwin-x64
      - darwin-arm64
```

This creates 6 parallel jobs that all run the same steps but with different platform targets.

## File Flow

```
Source Repository
    │
    ├─ vscode-extension/package.json ──────┐
    │  (version: "0.3.1")                   │
    │                                        │ Read version
    ├─ go-cli/main.go ─────────────────┐   │
    │                                    │   │
    │                           Build Go │   │
    │                            binaries│   │
    │                                    ▼   ▼
    │                            ┌──────────────┐
    │                            │ Build process│
    │                            └──────┬───────┘
    │                                   │
    │                         ┌─────────┴─────────┐
    │                         │                   │
    │                         ▼                   ▼
    │                    Go binaries        TypeScript
    │                    (6 platforms)      compiled
    │                         │                   │
    │                         └────────┬──────────┘
    │                                  │
    │                         Package into .vsix
    │                                  │
    │                         ┌────────┴────────┐
    │                         │                 │
    │                         ▼                 ▼
    ├─ dist/llamit-v0.3.1-win32-x64.vsix      ...
    ├─ dist/llamit-v0.3.1-win32-arm64.vsix    (x6)
    │  ...                                      │
    │                                           │
    │                              Upload to GitHub
    │                                           │
    │                         ┌─────────────────┤
    │                         │                 │
    │                         ▼                 ▼
    │                  GitHub Release    Marketplaces
    │                  (all 6 .vsix)     - Open VSX
    │                                     - VS Marketplace
    │
    └─ End users download correct platform version
```

## Summary

- **Total Jobs**: 4 (1 check + 6 parallel builds + 1 release + 1 publish)
- **Total Time**: ~10-15 minutes
- **Artifacts**: 6 platform-specific .vsix files (~5.5MB each)
- **Output**: GitHub Release + 2 marketplace publications
- **Automation**: 100% (zero manual steps required)

---

**Diagram Version**: 1.0
**Last Updated**: 2026-02-05
