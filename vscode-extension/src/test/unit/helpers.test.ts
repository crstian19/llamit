import * as assert from 'assert';
import * as path from 'path';
import { getConfiguration, getBinaryPath, LlamitConfig } from '../../extension';

/**
 * Tests for helper functions that don't require complex mocking
 */
suite('Helper Functions', () => {
    suite('getBinaryPath', () => {
        let originalPlatform: NodeJS.Platform;

        setup(() => {
            originalPlatform = process.platform;
        });

        teardown(() => {
            Object.defineProperty(process, 'platform', { value: originalPlatform });
        });

        test('should return cli on Linux', () => {
            Object.defineProperty(process, 'platform', { value: 'linux' });
            const extensionPath = '/home/user/.vscode/extensions/llamit';
            const result = getBinaryPath(extensionPath);

            assert.strictEqual(result, path.join(extensionPath, 'go-cli', 'cli'));
        });

        test('should return cli on macOS', () => {
            Object.defineProperty(process, 'platform', { value: 'darwin' });
            const extensionPath = '/Users/user/.vscode/extensions/llamit';
            const result = getBinaryPath(extensionPath);

            assert.strictEqual(result, path.join(extensionPath, 'go-cli', 'cli'));
        });

        test('should return cli.exe on Windows', () => {
            Object.defineProperty(process, 'platform', { value: 'win32' });
            const extensionPath = 'C:\\Users\\user\\.vscode\\extensions\\llamit';
            const result = getBinaryPath(extensionPath);

            assert.strictEqual(result, path.join(extensionPath, 'go-cli', 'cli.exe'));
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
