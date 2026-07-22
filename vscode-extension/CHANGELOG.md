# Change Log

All notable changes to the Llamit extension will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

> ⭐ Enjoying Llamit? A star on [GitHub](https://github.com/crstian19/llamit) genuinely helps the project grow.

## [2.1.0] - 2026-07-21

### Added
- **First-install star prompt**: a one-time, opt-in notification inviting you to star the project on GitHub. It appears only once, never repeats, and performs no action on your account — clicking through opens the repo so you star it yourself.
- **Release notes on update**: updating to a new version opens the rendered CHANGELOG in a Markdown preview once per version, so you can see what changed.
- **Reset Notification Prompts command**: `Llamit: Reset Notification Prompts` clears the stored prompt state so the install and update prompts can be reviewed again.

## [2.0.1] - 2026-07-04

### Changed
- **Dependencies**: Updated dev dependencies and CI tooling (`sinon` to v22, `@types/sinon` to v22, `@vscode/test-electron` to v3, `@vscode/vsce` to v3.9.2, `mocha` to v11.7.6, `@types/node` to v25.9.4, `@types/vscode` to v1.125.0, and `actions/checkout` to v7).

## [2.0.0] - 2026-05-12

### Added
- **OpenAI-compatible endpoints**: Added support for `llamit.apiType`, `llamit.apiUrl`, and `llamit.apiKey` so Llamit can work with Ollama and OpenAI-compatible chat completion APIs.
- **Deprecated setting migration notice**: Added a one-time in-editor migration prompt for users still configured with `llamit.ollamaUrl`.

### Changed
- **Configuration model**: `llamit.ollamaUrl` remains supported as a deprecated fallback while users migrate to `llamit.apiUrl` and `llamit.apiType`.
- **Marketplace positioning**: Updated extension metadata and documentation to reflect multi-provider endpoint support.

### Breaking
- **Configuration migration**: Users should move from `llamit.ollamaUrl` to `llamit.apiUrl` and `llamit.apiType`. The old setting still works in 2.0.0 as a deprecated compatibility fallback, but future releases may remove it.

## [1.9.3] - 2026-04-25

### Changed
- **Dependencies**: Updated dev dependencies and build tooling (`sinon`, `@types/vscode`, `typescript`, `@vscode/vsce`, `softprops/action-gh-release`).

## [1.9.2] - 2026-04-10

### Fixed
- **Telemetry**: Add `activationEvents: onStartupFinished` so the extension activates at startup and the telemetry ping fires correctly on install/update.

## [1.9.1] - 2026-04-10

### Fixed
- **Telemetry**: Use `context.extension.packageJSON.version` instead of `vscode.extensions.getExtension()` to reliably read the extension version.

## [1.9.0] - 2026-04-10

### Fixed
- **Telemetry**: Use `Mozilla/5.0 (compatible; ...)` User-Agent format to correctly bypass Umami's bot filter.

## [1.8.0] - 2026-04-10

### Fixed
- **Telemetry**: Added `User-Agent` header to Umami requests to prevent bot filter rejections.

## [1.7.0] - 2026-04-10

### Added
- **Anonymous telemetry**: Sends a single ping to a self-hosted Umami instance on install or update only (never on every activation). Collects editor name and extension version — no PII, no code, no usage frequency. Respects `vscode.env.isTelemetryEnabled` and the new `llamit.telemetry.enabled` setting.
- **`llamit.telemetry.enabled` setting**: Opt-out toggle for Llamit-specific telemetry (default: `true`).
- **`USAGE_DATA.md`**: Documents exactly what is and isn't collected, and how to opt out.

### Updated
- **Dependency Updates**:
  - `@types/node` to v25.6.0
  - `sinon` to v21.1.0
  - `@types/vscode` to v1.115.0

## [1.6.0] - 2026-04-07

### Updated
- **Dependency Updates**:
  - `typescript` from v5 to v6
  - `mocha` from v10 to v11
  - `sinon` from v20 to v21
  - `@types/node` to v25.5.2
- **GitHub Actions**: Updated Node.js version in CI workflows from 20 to 24
- **TypeScript v6 compatibility**: Added explicit `types` in `tsconfig.json`, typed `execFile` callbacks

## [1.5.0] - 2026-03-20

### Updated
- **Dependency Updates**: Updated multiple dependencies:
  - `@vscode/vsce` from v2 to v3
  - `glob` to v13
  - `@types/node` to v25
  - `@types/glob` to v9
  - `@types/vscode` to v1.110.0
- **GitHub Actions Updates**: Updated workflow actions:
  - `actions/checkout` to v6
  - `actions/setup-go` to v6
  - `actions/setup-node` to v6
  - `actions/download-artifact` to v8
  - `actions/upload-artifact` to v7
- **Go Version**: Updated to Go 1.26

## [1.4.0] - 2026-03-20

## [1.3.0] - 2026-02-25

### Added
- **Updated Go CLI Binaries**: Latest binaries included with extension package

### Fixed
- Binary files updated to latest version

## [1.2.0] - 2026-02-22

### Added
- **Repository Selector**: New dropdown to select the correct Git repository in multi-workspace scenarios
- **Workspace-Specific Git Path**: Extension now uses repository-specific git path for accurate repository detection

### Fixed
- Incorrect Git repository detection in multi-folder workspaces

## [1.1.0] - 2026-02-22

### Added
- **Enhanced Settings Organization**: Settings now grouped with ESSENTIAL and ADVANCED labels
- **Multi-Workspace Git Detection**: Improved logic to detect the correct Git repository when working with multiple folders

### Fixed
- Git repository detection in complex workspace configurations

## [1.0.0] - 2026-02-17

### Added
- **Advanced Ollama Parameters**: Full support for customizing model behavior with 15+ new parameters:
  - `keepAlive`: Control model memory retention time
  - `temperature`, `topK`, `topP`: Fine-tune output creativity and determinism
  - `numCtx`, `numPredict`: Manage context window and token generation
  - `repeatPenalty`, `repeatLastN`: Reduce repetitive output
  - `seed`: Enable reproducible results
  - `numGpu`, `numThread`: Optimize hardware utilization
  - `minP`, `tfsZ`: Advanced sampling parameters
  - `mirostat`, `mirostatEta`, `mirostatTau`: Mirostat sampling algorithm
  - `stop`: Custom stop sequences

### Changed
- **Documentation**: Release notes moved to GitHub Releases page for better maintainability
- Extension now provides granular control over Ollama model parameters

### Notes
- 🎉 **First stable release** - Llamit is now production-ready!
- All features tested and stable
- Backward compatible with existing configurations

## [0.4.0] - 2026-02-07

### Added
- **Smart Git Diff with Fallback**: Extension now automatically handles both staged and unstaged changes
  - If staged changes exist, uses only those (maintains original behavior)
  - If no staged changes, automatically falls back to working directory changes
  - Prioritizes staged changes when both staged and unstaged changes are present
  - No longer requires manual staging before generating commit messages

### Changed
- Updated message from "No staged changes to commit" to "No changes to commit" for better clarity
- Enhanced `getGitDiffCascade()` function with intelligent diff selection logic

### Internal
- Added `executeGitDiff()` helper function for better code organization
- Deprecated `getGitDiff()` in favor of `getGitDiffCascade()` (backward compatible)
- Added comprehensive test coverage for cascade logic

## [0.3.1] - 2026-02-05

  ### Fixed
  - CI/CD pipeline with Node.js 20 compatibility
  - GitHub Actions permissions for automatic releases

  ### Added
  - Automated publishing to Open VSX and VS Code Marketplace
  
## [0.3.0] - 2026-02-05

### Added
- **Platform-Specific Packaging**: Extension now distributed as separate packages for each platform
- Support for Windows ARM64 architecture
- Automated packaging scripts (`npm run package:all`, `npm run publish:all`)
- GitHub Actions workflow for automated releases
- Platform-specific build optimization with `VSCE_TARGET` environment variable

### Changed
- **Breaking**: Distribution method changed from universal package to platform-specific packages
- Each platform now gets a dedicated `.vsix` file (win32-x64, win32-arm64, linux-x64, linux-arm64, darwin-x64, darwin-arm64)
- VS Code automatically selects the correct package for user's platform
- Improved build script to support targeted platform compilation

### Technical Details
- Package size: ~16.5 MB per platform (includes all binaries due to vsce limitations)
- Runtime platform detection remains unchanged and fully functional
- All 6 platforms supported: Windows (x64, ARM64), Linux (x64, ARM64), macOS (Intel, Apple Silicon)
- Follows VS Code's recommended distribution pattern for native binaries

## [0.2.2] - 2026-02-04

### Added
- **Configurable Formats**: Support for 6 predefined styles (Conventional, Angular, Gitmoji, Karma, Semantic, Google)
- **Custom Format Templates**: New setting `llamit.customFormat` for user-defined prompts
- **Advanced Post-processing**: Automated stripping of markdown backticks and artifacts from LLM output
- **Optimized Prompts**: Enhanced system instructions for precise, one-line commit messages
- **Marketplace Branding**: High-fidelity badges, refreshed logo, and professional description

### Changed
- Integrated Go CLI build process into extension packaging (`npm run package`)
- Improved retry logic for better resilience against transient Ollama errors
- Updated extension settings with descriptive enums and multi-line custom editor

### Fixed
- Outdated binary builds in extension packages
- Marker artifacts (```) appearing in commit messages
- Extension ID mapping in integration tests

## [0.1.0] - 2026-02-04

### Added
- Initial release of Llamit
- Generate commit messages from staged Git changes
- Integration with local Ollama LLM instances
- Configurable Ollama endpoint and model selection
- Source Control toolbar button for easy access
- Command Palette integration
- Retry logic with exponential backoff for network errors
- Comprehensive unit and integration tests
- Cross-platform support (Linux, macOS, Windows)
- Conventional Commits format support

### Features
- 🚀 One-click commit message generation
- 🔒 Fully local - no cloud services required
- ⚡ Fast Go-based CLI backend
- 📝 Semantic commit message formatting
- 🎨 Seamless VS Code integration

---

**Note**: This extension was fully vibecoded using AI assistance! 🤖✨
