# Contributing to Llamit

Thank you for your interest in contributing to Llamit! This project was vibecoded, but we welcome contributions from humans and AI alike. 🤖✨

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/yourusername/llamit.git
   cd llamit
   ```

3. **Set up the development environment**:
   ```bash
   # Install Go (1.25.6 or higher)
   # Install Node.js (18.x or higher)

   # Build the Go CLI
   cd go-cli
   go build -o cli main.go

   # Set up the VS Code extension
   cd ../vscode-extension
   npm install
   npm run compile
   ```

## Development Workflow

### Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the code style of the project

3. **Run tests**:
   ```bash
   # Go CLI tests
   cd go-cli
   go test -v

   # VS Code extension tests
   cd vscode-extension
   npm run test:unit
   npm test
   ```

4. **Commit with a semantic message**:
   ```bash
   git commit -m "feat: add amazing new feature"
   ```

   Use [Conventional Commits](https://www.conventionalcommits.org/):
   - `feat:` - New features
   - `fix:` - Bug fixes
   - `docs:` - Documentation changes
   - `test:` - Test additions or changes
   - `refactor:` - Code refactoring
   - `chore:` - Maintenance tasks

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Open a Pull Request** on GitHub

## Code Style

### Go
- Follow standard Go conventions (`gofmt`, `golint`)
- Write tests for new functionality
- Keep functions focused and well-documented
- Use meaningful variable names

### TypeScript
- Follow the existing code style
- Use strict TypeScript types
- Export functions that should be testable
- Add JSDoc comments for public APIs

## Testing

All contributions should include tests:

- **Go**: Add tests to `go-cli/main_test.go`
- **TypeScript**: Add unit tests to `vscode-extension/src/test/unit/`
- **Integration**: Add integration tests if modifying extension behavior

Run the test suite before submitting:
```bash
# Go tests
cd go-cli && go test -v

# TypeScript unit tests (fast)
cd vscode-extension && npm run test:unit

# Full integration tests
cd vscode-extension && npm test
```

## Areas for Contribution

Here are some ideas if you're looking for something to work on:

### Features
- [ ] Support for custom commit message templates
- [ ] Multi-repository support
- [ ] Commit message history/suggestions
- [ ] Support for other LLM providers (OpenAI, Anthropic, etc.)
- [ ] Commit message validation/linting

### Improvements
- [ ] Better error messages
- [ ] Progress indicators during generation
- [ ] Configurable retry logic
- [ ] Streaming support for faster responses
- [ ] Cache recent diffs to avoid regeneration

### Testing
- [ ] Increase test coverage
- [ ] Add end-to-end tests
- [ ] Mock-based testing for external dependencies
- [ ] Performance benchmarks

### Documentation
- [ ] Video tutorial
- [ ] More usage examples
- [ ] Troubleshooting guide
- [ ] Best practices for different model sizes

## Pull Request Guidelines

Before submitting a PR:

1. ✅ Tests pass locally
2. ✅ Code follows project style
3. ✅ Commit messages are semantic
4. ✅ PR description explains the change
5. ✅ Documentation is updated if needed

### PR Template

When opening a PR, please include:

```markdown
## Description
Brief description of the change

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
Describe how you tested your changes

## Checklist
- [ ] Tests pass
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
```

## Questions?

Feel free to open an issue for:
- Bug reports
- Feature requests
- Questions about the codebase
- Discussing potential changes

## Code of Conduct

Be respectful and constructive. We're all here to build something useful together, whether you're human or AI! 🤝

---

Thank you for contributing to Llamit! 🎉
