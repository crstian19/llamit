<div align="center">

# Llamit - AI-Powered Commit Messages

<img src="https://cdn.crstian.me/llamit.jpg" alt="Llamit Logo" width="200"/>

![License](https://img.shields.io/github/license/crstian19/llamit?style=for-the-badge&logo=unlicense&logoColor=white)
![Build Status](https://img.shields.io/github/actions/workflow/status/crstian19/llamit/publish.yml?style=for-the-badge&logo=githubactions&logoColor=white&label=Build)
![VS Code Marketplace Installs](https://img.shields.io/visual-studio-marketplace/i/Crstian.llamit?style=for-the-badge&logo=visualstudiocode&logoColor=white&label=VS%20Code%20Marketplace)
![Open VSX Downloads](https://img.shields.io/open-vsx/dt/Crstian/llamit?style=for-the-badge&logo=vscodium&logoColor=white&label=Open%20VSX&color=blueviolet)
![VS Code](https://img.shields.io/badge/VS%20Code-1.85.0+-007ACC?style=for-the-badge&logo=visualstudiocode&logoColor=white)
![Go Version](https://img.shields.io/badge/Go-1.25.6-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Ollama](https://img.shields.io/badge/Ollama-powered-black?style=for-the-badge&logo=ollama&logoColor=white)

> ✨ **Fully vibecoded** - This project was entirely developed using AI assistance, showcasing the power of AI-driven development.

**Generate semantic commit messages using your local Ollama LLM instance.**

*No cloud services, no API keys - everything runs locally.*

</div>

## Features

- 🚀 **Generate commit messages instantly** from staged changes
- 🔒 **Fully local** - uses your own Ollama instance
- 📝 **Conventional Commits** - follows standard commit message format
- ⚡ **Fast** - powered by a lightweight Go CLI
- 🎨 **VS Code integration** - seamless SCM toolbar button

## Prerequisites

- [Ollama](https://ollama.ai/) installed and running locally
- A compatible model (default: `qwen2.5-coder:7b`, but any model works)
- VS Code 1.85.0 or higher

## Installation

### Option 1: From VS Code Marketplace (Recommended)
1. Open VS Code
2. Go to Extensions view (`Ctrl+Shift+X` / `Cmd+Shift+X`)
3. Search for "Llamit"
4. Click **Install**

### Option 2: Build from Source
```bash
# Clone the repository
git clone https://github.com/crstian19/llamit.git
cd llamit

# Build the Go CLI
cd go-cli
go build -o cli main.go

# Build the VS Code extension
cd ../vscode-extension
npm install
npm run compile

# Package the extension
npx vsce package
# Install the generated .vsix file in VS Code
```
```

## Usage

1. **Start Ollama**: Make sure Ollama is running
   ```bash
   ollama serve
   ```

2. **Stage your changes**: Use Git to stage the files you want to commit
   ```bash
   git add .
   ```

3. **Generate commit message**:
   - Click the ✨ **Llamit** button in the Source Control toolbar, or
   - Open Command Palette (`Ctrl+Shift+P` / `Cmd+Shift+P`)
   - Run: `Llamit: Generate Commit Message`

4. **Review and commit**: The generated message appears in the commit input box. Review it and commit!

## Configuration

You can customize Llamit in VS Code settings:

```json
{
  "llamit.ollamaUrl": "http://localhost:11434/api/generate",
  "llamit.model": "qwen2.5-coder:7b",
  "llamit.commitFormat": "conventional",
  "llamit.customFormat": ""
}
```

### Settings

- **`llamit.ollamaUrl`**: The Ollama API endpoint URL (default: `http://localhost:11434/api/generate`)
- **`llamit.model`**: The model to use for generation (default: `qwen2.5-coder:7b`)
- **`llamit.commitFormat`**: The commit message format to use (default: `conventional`)
  - Available formats: `conventional`, `angular`, `gitmoji`, `karma`, `semantic`, `google`, `custom`
- **`llamit.customFormat`**: Custom format template (only used when `commitFormat` is set to `custom`)

### Commit Message Formats

Llamit supports multiple commit message formats to match your team's conventions:

#### Conventional Commits (Default)
```
feat(auth): add user login functionality

Implements OAuth2 authentication flow
```

#### Angular
```
feat(core): implement user authentication

- Add login service
- Add auth guard
- Update routing

Closes #123
```

#### Gitmoji
```
✨ feat(api): add new endpoint for user profiles

Implements GET /api/users/:id endpoint
```

#### Karma
```
feat(ui): add dark mode toggle

Implements theme switching functionality
```

#### Semantic
```
feat: implement user authentication system

Complete OAuth2 integration with JWT tokens
```

#### Google
```
Add user authentication system

Implements a complete authentication flow using OAuth2 and JWT tokens.
Includes login, logout, and token refresh functionality.
```

#### Custom Format
Set `llamit.commitFormat` to `custom` and provide your own template in `llamit.customFormat`:

```json
{
  "llamit.commitFormat": "custom",
  "llamit.customFormat": "Generate a simple commit message:\n<action>: <description>\n\nRules:\n1. Keep it under 50 characters\n2. Use imperative mood"
}
```

### Recommended Models

Any Ollama model works, but these are optimized for code:
- `qwen2.5-coder:7b` - Great balance of quality and speed (default)
- `qwen2.5-coder:14b` - Better quality, slower
- `codellama:13b` - Good alternative
- `deepseek-coder:6.7b` - Fast and efficient

## Architecture

Llamit consists of two components:

### 1. Go CLI (`go-cli/`)
A standalone command-line tool that:
- Reads git diffs from stdin
- Sends them to Ollama with a prompt template
- Returns a formatted commit message
- Implements retry logic with exponential backoff
- Handles errors gracefully

### 2. VS Code Extension (`vscode-extension/`)
A TypeScript extension that:
- Integrates with VS Code's Source Control view
- Executes `git diff --cached` to get staged changes
- Spawns the Go CLI as a subprocess
- Populates the commit message box with the result

## Development

### Running Tests

**Go CLI:**
```bash
cd go-cli
go test -v              # All tests
go test -v -short       # Unit tests only (skip integration)
```

**VS Code Extension:**
```bash
cd vscode-extension
npm run test:unit       # Fast unit tests
npm test               # Full integration tests
```

### File Structure

```
llamit/
├── go-cli/              # Go CLI binary
│   ├── main.go          # Core logic
│   ├── main_test.go     # Comprehensive tests
│   └── go.mod           # Go module file
├── vscode-extension/    # VS Code extension
│   ├── src/
│   │   ├── extension.ts # Extension entry point
│   │   └── test/        # Unit and integration tests
│   ├── package.json     # Extension manifest
│   └── tsconfig.json    # TypeScript config
└── CLAUDE.md           # AI assistant documentation
```

## Testing

Both components have comprehensive test coverage:

- **Go CLI**: 6 test cases covering success, errors, retries, and integration
- **VS Code Extension**: Unit tests + integration tests for all core functions

See [CLAUDE.md](./CLAUDE.md) for detailed testing information.

## Contributing

Contributions are welcome! This project was vibecoded, but that doesn't mean it can't be improved by humans too 😊

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'feat: add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

See [CONTRIBUTING.md](./CONTRIBUTING.md) for detailed guidelines.

## Releases

Llamit uses automated CI/CD for releases:
- ✅ Automatic releases on merge to `main`
- 📦 Platform-specific packaging (6 platforms)
- 🚀 Automatic publishing to VS Code Marketplace and Open VSX

For maintainers: See [.github/RELEASE.md](.github/RELEASE.md) for detailed release process documentation.

## License

MIT License - see [LICENSE](LICENSE) file for details

## Acknowledgments

- Built with [Claude](https://claude.ai) - AI pair programming at its finest
- Powered by [Ollama](https://ollama.ai) - local LLM runtime
- Inspired by the need for better commit messages everywhere

## Release Notes

### 0.2.2
- **Configurable Formats**: Added 6 predefined templates (Conventional, Angular, Gitmoji, Karma, Semantic, Google)
- **Custom Templates**: Support for user-defined commit message formats
- **Optimized Prompts**: Refined instructions for maximum conciseness and brevity
- **Markdown Cleanup**: Automatic removal of backticks and code blocks from LLM output
- **Automation**: Integrated Go CLI build into the extension lifecycle
- **UI improvements**: High-fidelity badges and CDN-based logo

### 0.1.0
- Initial release with local Ollama integration

---

**Made with 🤖 and ✨ through vibecoding**
