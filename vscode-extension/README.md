# Llamit - AI-Powered Commit Messages

> ✨ **Fully vibecoded** - This extension was entirely developed using AI assistance.

Generate semantic commit messages instantly using your local Ollama LLM instance. No cloud services, no API keys - everything runs locally!

## Features

- 🚀 **Generate commit messages instantly** from staged changes
- 🔒 **Fully local** - uses your own Ollama instance
- 📝 **Conventional Commits** - follows standard commit message format
- ⚡ **Fast** - powered by a lightweight Go CLI
- 🎨 **Seamless integration** - works directly in VS Code's Source Control view

![Llamit Demo](https://raw.githubusercontent.com/crstian/llamit/main/assets/demo.gif)

## Prerequisites

Before using Llamit, you need:

1. **[Ollama](https://ollama.ai/)** installed and running locally
2. A compatible model downloaded (default: `qwen3-coder:30b`)

```bash
# Install Ollama, then pull a model:
ollama pull qwen3-coder:30b
```

## Usage

1. **Stage your changes** in Git
2. **Click the ✨ Llamit button** in the Source Control toolbar
3. **Review the generated commit message** and commit!

Or use the Command Palette:
- Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on Mac)
- Type: `Llamit: Generate Commit Message`

## Configuration

Customize Llamit in VS Code settings:

| Setting | Description | Default |
|---------|-------------|---------|
| `llamit.ollamaUrl` | Ollama API endpoint | `http://localhost:11434/api/generate` |
| `llamit.model` | Model to use | `qwen3-coder:30b` |

### Example Configuration

```json
{
  "llamit.ollamaUrl": "http://localhost:11434/api/generate",
  "llamit.model": "qwen3-coder:7b"
}
```

## Recommended Models

Any Ollama model works, but these are optimized for code:

- `qwen3-coder:30b` - Best quality (default) 🌟
- `qwen3-coder:7b` - Faster, lighter ⚡
- `codellama:13b` - Good balance
- `deepseek-coder:6.7b` - Fast and efficient

## How It Works

Llamit uses a two-component architecture:

1. **Go CLI**: Processes git diffs and communicates with Ollama
2. **VS Code Extension**: Integrates with the Source Control view

The extension:
1. Executes `git diff --cached` to get staged changes
2. Passes the diff to the Go CLI via stdin
3. The CLI sends it to your local Ollama instance
4. Returns a formatted commit message following Conventional Commits

## Requirements

- VS Code 1.85.0 or higher
- Ollama running locally
- Git repository

## Known Issues

- The extension requires Ollama to be running. If you see errors, make sure Ollama is started.
- Large diffs may take longer to process depending on your model size.

## Contributing

Found a bug or want to contribute? Check out the [GitHub repository](https://github.com/crstian/llamit)!

## Release Notes

### 0.1.0

Initial release of Llamit:
- Generate commit messages from staged changes
- Local Ollama integration
- Configurable models and endpoints
- Retry logic with exponential backoff

## License

MIT License - see [LICENSE](https://github.com/crstian/llamit/blob/main/LICENSE)

---

**Made with 🤖 and ✨ through vibecoding**

Enjoy using Llamit? [⭐ Star us on GitHub](https://github.com/crstian/llamit)
