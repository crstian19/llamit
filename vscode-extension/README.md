<div align="center">

# Llamit - AI-Powered Commit Messages

<img src="https://cdn.crstian.me/llamit.png" alt="Llamit Logo" width="200"/>

![License](https://img.shields.io/github/license/crstian19/llamit?style=for-the-badge&logo=unlicense&logoColor=white)
![VS Code Installs](https://img.shields.io/visual-studio-marketplace/v/Crstian.llamit?style=for-the-badge&logo=visualstudiocode&logoColor=white&label=installs)
![VS Code](https://img.shields.io/badge/VS%20Code-1.85.0+-007ACC?style=for-the-badge&logo=visualstudiocode&logoColor=white)
![Ollama](https://img.shields.io/badge/Ollama-powered-black?style=for-the-badge&logo=ollama&logoColor=white)

> ✨ **Fully vibecoded** - This project was entirely developed using AI assistance, showcasing the power of AI-driven development.

**Generate semantic commit messages using your local Ollama LLM instance.**

*No cloud services, no API keys - everything runs locally.*

</div>

## Features

- 🚀 **Generate commit messages instantly** from staged changes
- 🔒 **Fully local** - uses your own Ollama instance
- 📝 **Conventional Commits** - follows standard commit message format (and many others!)
- ⚡ **Fast** - powered by a lightweight Go CLI
- 🎨 **VS Code integration** - seamless SCM toolbar button

## Prerequisites

- [Ollama](https://ollama.ai/) installed and running locally
- A compatible model (default: `qwen3-coder:30b`, but any model works)
- VS Code 1.85.0 or higher

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
  "llamit.model": "qwen3-coder:30b",
  "llamit.commitFormat": "conventional",
  "llamit.customFormat": ""
}
```

### Settings

- **`llamit.ollamaUrl`**: The Ollama API endpoint URL (default: `http://localhost:11434/api/generate`)
- **`llamit.model`**: The model to use for generation (default: `qwen3-coder:30b`)
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

#### Gitmoji
```
✨ feat(api): add new endpoint for user profiles

Implements GET /api/users/:id endpoint
```

#### Custom Format
Set `llamit.commitFormat` to `custom` and provide your own template in `llamit.customFormat`.

## Extension Settings

This extension contributes the following settings:

* `llamit.ollamaUrl`: URL of the Ollama generation API.
* `llamit.model`: Name of the Ollama model to use.
* `llamit.commitFormat`: Format of the generated commit message.
* `llamit.customFormat`: Custom template for the commit message.

## How It Works

Llamit consists of two components:
1. **Go CLI**: A fast, standalone binary that processes diffs and talks to Ollama.
2. **VS Code Extension**: A lightweight TypeScript bridge that integrates with the VS Code SCM view.

Everything runs on your machine. Your code never leaves your local environment.

## Release Notes

### 0.2.1
- Added configurable commit message formats (Conventional, Angular, Gitmoji, Karma, Semantic, Google)
- Added custom format template support
- Improved error handling and reliability
- Updated documentation and badges

### 0.1.0
- Initial release with Conventional Commits support

---

**Made with 🤖 and ✨ through vibecoding**

[GitHub Repository](https://github.com/crstian19/llamit)
