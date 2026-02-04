# Llamit VS Code Extension Tests

This directory contains unit and integration tests for the Llamit extension.

## Test Structure

```
test/
├── unit/                    # Unit tests (fast, no VS Code runtime)
│   ├── helpers.test.ts      # Tests for utility functions
│   ├── gitDiff.test.ts      # Tests for git diff functionality
│   └── commitMessage.test.ts # Tests for commit message generation
├── suite/                   # Integration tests (require VS Code runtime)
│   ├── index.ts             # Test runner
│   └── extension.integration.test.ts
└── runTests.ts              # Entry point for integration tests
```

## Running Tests

### Unit Tests (Fast)
```bash
npm run test:unit
```

These tests run with Mocha directly and don't require VS Code runtime. They test individual functions in isolation.

### Integration Tests (Slower)
```bash
npm test
```

These tests run inside VS Code's Electron environment and test the extension as a whole, including:
- Extension activation
- Command registration
- Configuration defaults

## Writing New Tests

### Unit Tests

Place unit tests in `test/unit/`. They should:
- Test pure functions and logic
- Use mocks for external dependencies
- Run quickly without requiring VS Code runtime
- Follow the naming convention: `*.test.ts`

Example:
```typescript
import * as assert from 'assert';
import { myFunction } from '../../extension';

suite('My Function Tests', () => {
    test('should do something', () => {
        const result = myFunction(input);
        assert.strictEqual(result, expected);
    });
});
```

### Integration Tests

Place integration tests in `test/suite/`. They should:
- Test real extension behavior
- Use VS Code APIs
- Verify command registration, activation, etc.
- Follow the naming convention: `*.integration.test.ts`

## Test Coverage

Currently tested:
- ✅ Binary path resolution (cross-platform)
- ✅ Configuration loading
- ✅ Extension activation
- ✅ Command registration
- ✅ Function signatures and interfaces

Future improvements:
- Mock execFile for testing git diff execution
- Mock CLI binary execution
- Add error handling tests
- Test progress indicator behavior
