<div align="center">

# Llamit - AI-Powered Commit Messages

<img src="https://cdn.crstian.me/llamit.jpg" alt="Llamit Logo" width="200"/>

![License](https://img.shields.io/github/license/crstian19/llamit?style=for-the-badge&logo=unlicense&logoColor=white)
![Build Status](https://img.shields.io/github/actions/workflow/status/crstian19/llamit/publish.yml?style=for-the-badge&logo=githubactions&logoColor=white&label=Build)
![VS Code Marketplace Installs](https://img.shields.io/visual-studio-marketplace/i/Crstian.llamit?style=for-the-badge&label=VS%20Code%20Marketplace%20Downloads)
![Open VSX Downloads](https://img.shields.io/open-vsx/dt/Crstian/llamit?style=for-the-badge&logo=vscodium&logoColor=white&label=Open%20VSX%20downloads&color=blueviolet)
![VS Code](https://img.shields.io/badge/VS%20Code-1.85.0+-007ACC?style=for-the-badge&logo=visualstudiocode&logoColor=white)
![Go Version](https://img.shields.io/badge/Go-1.25.6-00ADD8?style=for-the-badge&logo=go&logoColor=white)

> ✨ **Fully vibecoded** - This project was entirely developed using AI assistance, showcasing the power of AI-driven development.

**Generate semantic commit messages using Ollama or any OpenAI-compatible endpoint.**

*Run it locally with Ollama or point it at the chat completion API you already use.*

</div>

## Highlights

- Generate commit messages directly from VS Code Source Control
- Support for `ollama` and `openai-compatible` APIs
- Full custom endpoint support through a configurable URL
- Conventional, Angular, Gitmoji, Karma, Semantic, Google, and custom formats
- Go CLI backend with retry logic and clean stdout/stderr separation

## What Changed in 2.0.0

`2.0.0` is the release where Llamit stops being Ollama-only.

- New settings:
  - `llamit.apiType`
  - `llamit.apiUrl`
  - `llamit.apiKey`
- Deprecated setting:
  - `llamit.ollamaUrl`
- Backward compatibility:
  - `llamit.ollamaUrl` still works as a fallback
  - the extension now shows a one-time migration prompt for users still on the old setting

## Prerequisites

- VS Code 1.85.0 or higher
- One of:
  - [Ollama](https://ollama.ai/) running locally
  - an OpenAI-compatible `/v1/chat/completions` endpoint
- A valid model name for the endpoint you choose

## Installation

### Option 1: VS Code Marketplace

1. Open VS Code
2. Go to Extensions
3. Search for `Llamit`
4. Install the extension

### Option 2: Build from Source

```bash
git clone https://github.com/crstian19/llamit.git
cd llamit

cd go-cli
go build -o cli main.go

cd ../vscode-extension
npm install
npm run vscode:prepublish
```

Then install the generated `.vsix` or launch the extension in development mode with `F5`.

## Configuration

Llamit is configured through VS Code extension settings.

### Ollama

```json
{
  "llamit.apiType": "ollama",
  "llamit.apiUrl": "http://localhost:11434/api/generate",
  "llamit.model": "qwen2.5-coder:7b",
  "llamit.commitFormat": "conventional"
}
```

### OpenAI-Compatible Endpoint

```json
{
  "llamit.apiType": "openai-compatible",
  "llamit.apiUrl": "https://api.openai.com/v1/chat/completions",
  "llamit.apiKey": "",
  "llamit.model": "gpt-4o-mini",
  "llamit.commitFormat": "conventional"
}
```

If `llamit.apiKey` is empty, the CLI falls back to:

1. `LLAMIT_API_KEY`
2. `OPENAI_API_KEY`

### Settings

- `llamit.apiType`: `ollama` or `openai-compatible`
- `llamit.apiUrl`: full API URL to call
- `llamit.apiKey`: optional API key for OpenAI-compatible endpoints
- `llamit.model`: model name to send to the provider
- `llamit.commitFormat`: commit message style
- `llamit.customFormat`: custom prompt template when `commitFormat=custom`

Advanced settings remain available. Some of them are Ollama-specific, while shared fields such as `temperature`, `topP`, `numPredict`, `seed`, and `stop` are forwarded when supported by the selected API type.

## Migration from `llamit.ollamaUrl`

If you already use Llamit with:

```json
{
  "llamit.ollamaUrl": "http://localhost:11434/api/generate"
}
```

migrate to:

```json
{
  "llamit.apiType": "ollama",
  "llamit.apiUrl": "http://localhost:11434/api/generate"
}
```

The extension now detects the deprecated setting and offers:

- `Migrate Settings`
- `Open Settings`
- `Don't Show Again`

## Usage

1. Stage changes, or leave them unstaged if you want the diff fallback behavior
2. Open Source Control
3. Run `Llamit: Generate Commit Message`
4. Review the generated message in the SCM input box
5. Commit as usual

## Architecture

Llamit consists of two components:

### Go CLI (`go-cli/`)

- reads Git diff from `stdin`
- builds the prompt from the selected commit format
- calls either:
  - Ollama `/api/generate`
  - an OpenAI-compatible `/v1/chat/completions` endpoint
- writes the commit message to `stdout`
- logs operational details to `stderr`

### VS Code Extension (`vscode-extension/`)

- integrates with the VS Code Git extension
- resolves the current repository
- uses staged changes first, then falls back to working tree changes
- spawns the bundled Go CLI for the current platform
- writes the generated message into the SCM input box

## Testing

### Go CLI

```bash
cd go-cli
go test -v -short
go test -v
```

Notes:

- `go test -v -short` runs the unit suite
- `go test -v` includes the live Ollama integration test

### VS Code Extension

```bash
cd vscode-extension
npm run test:unit
npm test
```

## Releases

Llamit uses automated CI/CD for releases.

- merge to `main`
- make sure `vscode-extension/package.json` has a new version
- update `vscode-extension/CHANGELOG.md`
- the workflow in [.github/workflows/publish.yml](./.github/workflows/publish.yml) creates the release and publishes marketplace packages

For maintainers, the full process is documented in [.github/RELEASE.md](./.github/RELEASE.md).

## Contributing

1. Fork the repository
2. Create a branch
3. Make changes
4. Run the relevant test suite
5. Open a pull request

See [CONTRIBUTING.md](./CONTRIBUTING.md) for the contributor workflow.
