# AGENTS.md - Guidance for Agentic Coding Tools

This file provides guidance for AI coding agents working in this repository.

## Project Overview

Llamit generates commit messages using local Ollama LLM instances. It consists of two components:
- **Go CLI** (`go-cli/`): Command-line tool that reads git diffs from stdin and generates commit messages
- **VS Code Extension** (`vscode-extension/`): TypeScript extension integrating the Go CLI into VS Code

---

## Build, Lint, and Test Commands

### Go CLI (`go-cli/`)

```bash
# Build the CLI binary
cd go-cli && go build -o cli main.go

# Run all tests (including integration tests)
cd go-cli && go test -v

# Run only unit tests (skip integration test)
cd go-cli && go test -v -short

# Run a single test by name
cd go-cli && go test -v -run TestRun_EmptyInput

# Run with verbose output and cover profile
cd go-cli && go test -v -coverprofile=coverage.out

# Format code (go fmt)
cd go-cli && go fmt ./...

# Run go vet for static analysis
cd go-cli && go vet ./...

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o llamit-linux-amd64 main.go
GOOS=darwin GOARCH=amd64 go build -o llamit-darwin-amd64 main.go
GOOS=darwin GOARCH=arm64 go build -o llamit-darwin-arm64 main.go
GOOS=windows GOARCH=amd64 go build -o llamit-windows-amd64.exe main.go
```

### VS Code Extension (`vscode-extension/`)

```bash
# Install dependencies
cd vscode-extension && npm install

# Compile TypeScript
cd vscode-extension && npm run compile

# Watch mode for development
cd vscode-extension && npm run watch

# Run unit tests only (fast, no VS Code runtime)
cd vscode-extension && npm run test:unit

# Run all tests including integration (requires VS Code runtime)
cd vscode-extension && npm test

# Run a single unit test file
cd vscode-extension && npx mocha --ui tdd --require ts-node/register 'src/test/unit/helpers.test.ts'

# Package for publishing
cd vscode-extension && npm run vscode:prepublish
```

---

## Code Style Guidelines

### Go (go-cli/)

**Imports**
- Group imports: standard library first, then third-party packages
- Use explicit imports (no blank imports `_`)
- Format: `go fmt` handles automatic formatting

```go
import (
    "bytes"
    "encoding/json"
    "errors"
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
    "time"
)
```

**Naming Conventions**
- Use PascalCase for exported functions, types, and constants
- Use camelCase for unexported functions, variables, and constants
- Acronyms: keep uppercase (e.g., `HTTP`, `URL`, `API`)
- File names: lowercase with underscores (snake_case)

**Types and Interfaces**
- Define structs for API requests/responses with JSON tags
- Use meaningful type names (e.g., `OllamaGenerateRequest`)
- Export error variables as sentinel errors: `var ErrEmptyInput = errors.New(...)`

**Error Handling**
- Return errors with context using `fmt.Errorf("context: %w", err)`
- Check errors immediately after operations
- Log errors to stderr with `log.Printf("ERROR: ...")` before returning
- Use sentinel errors for expected error conditions

**Logging**
- Use `log` package for operational logging (stderr)
- Output commit message to stdout only
- Log format: `log.Ldate | log.Ltime | log.Lshortfile`

**Testing**
- Use `httptest.NewServer()` for HTTP mocking
- Skip integration tests with `testing.Short()`
- Use table-driven tests with `t.Run()` for multiple test cases
- Capture logs in tests with `bytes.Buffer`

### TypeScript (vscode-extension/)

**Imports**
- Use namespace imports for built-in modules: `import * as vscode from 'vscode'`
- Use named imports for external modules: `import { execFile } from 'child_process'`
- Order: external libraries first, then local modules

```typescript
import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import { execFile } from 'child_process';
import { getBinaryPath, getGitDiffCascade } from './helpers';
```

**Naming Conventions**
- Use PascalCase for types, interfaces, classes, and exported functions
- Use camelCase for variables, function parameters, and unexported functions
- Interfaces: prefix with `I` is optional (project uses descriptive names like `LlamitConfig`)

**Types**
- Always define explicit return types for exported functions
- Use interfaces for structured data (config objects, result types)
- Use `export interface` for types needed in tests

```typescript
export interface LlamitConfig {
    ollamaUrl: string;
    model: string;
    commitFormat: string;
    customFormat: string;
}
```

**Error Handling**
- Use try/catch for async operations
- Re-throw errors with descriptive messages: `throw new Error(context + ": " + error.message)`
- Show user-facing errors via `vscode.window.showErrorMessage()`

**VS Code Extension Patterns**
- Register commands with `vscode.commands.registerCommand()`
- Use `context.subscriptions.push(disposable)` for cleanup
- Use `vscode.window.withProgress()` for long-running operations
- Get Git API via `vscode.extensions.getExtension('vscode.git')`

**Testing**
- Unit tests: place in `src/test/unit/` with `.test.ts` suffix
- Use Mocha with TDD style (`describe`, `test`, `assert`)
- Use Sinon for mocking/spying
- Run with `npm run test:unit` (no VS Code runtime needed)

---

## Architecture Notes

### Communication Flow
1. VS Code extension calls `git diff --cached` for staged changes
2. Extension spawns Go CLI binary as subprocess
3. Git diff is written to CLI's stdin
4. CLI sends diff to Ollama API with prompt template
5. CLI returns commit message to stdout
6. Extension populates SCM input box with result

### Binary Distribution
The extension expects a pre-built binary at `vscode-extension/go-cli/bin/<platform-specific-binary>`. When developing:
1. Build the Go CLI first: `cd go-cli && go build -o cli main.go`
2. Copy/symlink to extension: `cp go-cli/cli vscode-extension/go-cli/bin/llamit-linux-amd64`

---

## Configuration

### CLI Flags
- `-ollama-url`: Ollama API endpoint (default: `http://localhost:11434/api/generate`)
- `-model`: Model name (default: `qwen2.5-coder:7b`)
- `-format`: Commit format (conventional, angular, gitmoji, karma, semantic, google, custom)
- `-custom-template`: Custom format template (used when format is 'custom')
- `-version`: Print version

### Extension Settings
- `llamit.ollamaUrl`: Full API endpoint URL
- `llamit.model`: Ollama model name
- `llamit.commitFormat`: Commit message format
- `llamit.customFormat`: Custom template (when commitFormat is 'custom')
