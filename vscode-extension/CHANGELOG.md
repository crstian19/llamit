# Change Log

All notable changes to the Llamit extension will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
