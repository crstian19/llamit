# Platform-Specific Packaging Notes

## Current Status

The extension now uses platform-specific packaging, creating separate `.vsix` files for each platform.

### Known Limitation

Currently, each `.vsix` package includes **all** binary files (~16.5 MB) instead of just the platform-specific binary (~6 MB). This is because:

1. `vsce package --target` always runs the `vscode:prepublish` script
2. The prepublish script rebuilds **all** binaries regardless of the `VSCE_TARGET` environment variable
3. There's no way to conditionally exclude files from `.vsix` based on the target platform using `.vscodeignore`

### Why This Happens

The `vsce` tool doesn't pass the target platform to the prepublish hook, so even though we set `VSCE_TARGET` in our packaging script, the automatic prepublish hook doesn't have access to it and rebuilds all binaries.

### Impact

- Each `.vsix` is 16.5 MB instead of the ideal ~6 MB
- Users still get correct platform-specific functionality (runtime detection works)
- Download size is 2.75x larger than optimal, but still better than no platform-specific packaging
- The extension still works correctly on all platforms

### Potential Solutions (Future)

1. **Post-processing approach**: Unzip `.vsix`, remove unnecessary binaries, re-zip
2. **Multi-repo approach**: Separate extension packages per platform
3. **VSCE enhancement**: Request VSCE to pass `--target` to prepublish hooks
4. **Binary hosting**: Host binaries separately and download on first use

### Comparison

| Approach | Size per Platform | Total for 6 Platforms | User Download |
|----------|------------------|----------------------|---------------|
| Ideal platform-specific | ~6 MB | 36 MB | 6 MB |
| Current implementation | 16.5 MB | 99 MB | 16.5 MB |
| Previous universal | 28 MB | 28 MB | 28 MB |

While not optimal, the current approach is still an improvement over the universal package (16.5 MB vs 28 MB per download).

## How It Works

1. User installs extension from VS Code Marketplace
2. VS Code automatically selects the correct `.vsix` for their platform
3. All binaries are included, but only the correct one is used at runtime
4. The `helpers.ts` file detects the platform and selects the appropriate binary

## Testing

To test a specific platform package locally:

```bash
# Build a specific platform
VSCE_TARGET=linux-x64 npm run vscode:prepublish

# Install the .vsix
code --install-extension dist/llamit-0.3.0-linux-x64.vsix
```

## Build Process

The build script (`scripts/package-all-platforms.sh`) creates 6 separate packages:

1. `llamit-0.3.0-win32-x64.vsix` - Windows x64
2. `llamit-0.3.0-win32-arm64.vsix` - Windows ARM64
3. `llamit-0.3.0-linux-x64.vsix` - Linux x64
4. `llamit-0.3.0-linux-arm64.vsix` - Linux ARM64
5. `llamit-0.3.0-darwin-x64.vsix` - macOS Intel
6. `llamit-0.3.0-darwin-arm64.vsix` - macOS Apple Silicon

All packages are published to both VS Code Marketplace and Open VSX Registry.
