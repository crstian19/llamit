<div align="center">

# Llamit - AI-Powered Commit Messages

<img src="https://cdn.crstian.me/llamit.png" alt="Llamit Logo" width="200"/>

![License](https://img.shields.io/github/license/crstian19/llamit?style=for-the-badge&logo=unlicense&logoColor=white)
![Azure DevOps Installs](https://img.shields.io/visual-studio-marketplace/azure-devops/installs/total/Crstian.llamit?style=for-the-badge&logo=visualstudiocode&logoColor=white&label=installs)
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
- A compatible model (default: `qwen3-coder:30b`, but any model works)
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
  "llamit.model": "qwen3-coder:30b"
}
```

### Settings

- **`llamit.ollamaUrl`**: The Ollama API endpoint URL (default: `http://localhost:11434/api/generate`)
- **`llamit.model`**: The model to use for generation (default: `qwen3-coder:30b`)

### Recommended Models

Any Ollama model works, but these are optimized for code:
- `qwen3-coder:30b` - Best quality (default)
- `qwen3-coder:7b` - Faster, lighter
- `codellama:13b` - Good balance
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

## License

MIT License - see [LICENSE](LICENSE) file for details

## Acknowledgments

- Built with [Claude](https://claude.ai) - AI pair programming at its finest
- Powered by [Ollama](https://ollama.ai) - local LLM runtime
- Inspired by the need for better commit messages everywhere

---

**Made with 🤖 and ✨ through vibecoding**
