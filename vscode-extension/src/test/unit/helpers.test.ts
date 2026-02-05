import * as assert from 'assert';
import * as path from 'path';
import { getBinaryPath, LlamitConfig } from '../../helpers';

/**
 * Tests for helper functions that don't require complex mocking
 */
suite('Helper Functions', () => {
    suite('getBinaryPath', () => {
        let originalPlatform: NodeJS.Platform;
        let originalArch: string;

        setup(() => {
            originalPlatform = process.platform;
            originalArch = process.arch;
        });

        teardown(() => {
            Object.defineProperty(process, 'platform', { value: originalPlatform });
            Object.defineProperty(process, 'arch', { value: originalArch });
        });

        test('should return llamit-linux-amd64 on Linux x64', () => {
            Object.defineProperty(process, 'platform', { value: 'linux' });
            Object.defineProperty(process, 'arch', { value: 'x64' }); // Node uses x64 for amd64
            const extensionPath = '/home/user/.vscode/extensions/llamit';
            const result = getBinaryPath(extensionPath);

            assert.strictEqual(result, path.join(extensionPath, 'go-cli', 'bin', 'llamit-linux-amd64'));
        });

        test('should return llamit-linux-arm64 on Linux arm64', () => {
            Object.defineProperty(process, 'platform', { value: 'linux' });
            Object.defineProperty(process, 'arch', { value: 'arm64' });
            const extensionPath = '/home/user/.vscode/extensions/llamit';
            const result = getBinaryPath(extensionPath);

            assert.strictEqual(result, path.join(extensionPath, 'go-cli', 'bin', 'llamit-linux-arm64'));
        });

        test('should return llamit-darwin-amd64 on macOS x64', () => {
            Object.defineProperty(process, 'platform', { value: 'darwin' });
            Object.defineProperty(process, 'arch', { value: 'x64' });
            const extensionPath = '/Users/user/.vscode/extensions/llamit';
            const result = getBinaryPath(extensionPath);

            assert.strictEqual(result, path.join(extensionPath, 'go-cli', 'bin', 'llamit-darwin-amd64'));
        });

        test('should return llamit-darwin-arm64 on macOS arm64 (Apple Silicon)', () => {
            Object.defineProperty(process, 'platform', { value: 'darwin' });
            Object.defineProperty(process, 'arch', { value: 'arm64' });
            const extensionPath = '/Users/user/.vscode/extensions/llamit';
            const result = getBinaryPath(extensionPath);

            assert.strictEqual(result, path.join(extensionPath, 'go-cli', 'bin', 'llamit-darwin-arm64'));
        });

        test('should return llamit-windows-amd64.exe on Windows x64', () => {
            Object.defineProperty(process, 'platform', { value: 'win32' });
            Object.defineProperty(process, 'arch', { value: 'x64' });
            const extensionPath = 'C:\\Users\\user\\.vscode\\extensions\\llamit';
            const result = getBinaryPath(extensionPath);

            assert.strictEqual(result, path.join(extensionPath, 'go-cli', 'bin', 'llamit-windows-amd64.exe'));
        });

        test('should throw error for unsupported platform', () => {
            Object.defineProperty(process, 'platform', { value: 'sunos' });
            const extensionPath = '/path/to/ext';

            assert.throws(() => getBinaryPath(extensionPath), /Unsupported platform/);
        });
    });

    suite('LlamitConfig interface', () => {
        test('should enforce correct structure', () => {
            const config: LlamitConfig = {
                ollamaUrl: 'http://localhost:11434/api/generate',
                model: 'qwen2.5-coder:7b',
                commitFormat: 'conventional',
                customFormat: ''
            };

            assert.ok(config.ollamaUrl);
            assert.ok(config.model);
            assert.strictEqual(typeof config.ollamaUrl, 'string');
            assert.strictEqual(typeof config.model, 'string');
            assert.strictEqual(typeof config.commitFormat, 'string');
            assert.strictEqual(typeof config.customFormat, 'string');
        });

        test('should accept custom values', () => {
            const config: LlamitConfig = {
                ollamaUrl: 'https://custom-ollama.example.com/api/generate',
                model: 'llama2:13b',
                commitFormat: 'gitmoji',
                customFormat: 'My custom template'
            };

            assert.strictEqual(config.ollamaUrl, 'https://custom-ollama.example.com/api/generate');
            assert.strictEqual(config.model, 'llama2:13b');
            assert.strictEqual(config.commitFormat, 'gitmoji');
            assert.strictEqual(config.customFormat, 'My custom template');
        });
    });
});
