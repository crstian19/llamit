import * as assert from 'assert';
import * as path from 'path';
import { getBinaryPath, LlamitConfig } from '../../helpers';

/**
 * Tests for helper functions that don't require complex mocking
 */
suite('Helper Functions', () => {
    suite('getBinaryPath', () => {
        let originalPlatform: string;
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
            Object.defineProperty(process, 'arch', { value: 'x64' });
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
                apiType: 'ollama',
                apiUrl: 'http://localhost:11434/api/generate',
                model: 'qwen2.5-coder:7b',
                commitFormat: 'conventional',
                customFormat: ''
            };

            assert.ok(config.apiType);
            assert.ok(config.apiUrl);
            assert.ok(config.model);
            assert.strictEqual(typeof config.apiType, 'string');
            assert.strictEqual(typeof config.apiUrl, 'string');
            assert.strictEqual(typeof config.model, 'string');
            assert.strictEqual(typeof config.commitFormat, 'string');
            assert.strictEqual(typeof config.customFormat, 'string');
        });

        test('should accept custom values', () => {
            const config: LlamitConfig = {
                apiType: 'openai-compatible',
                apiUrl: 'https://gateway.example.com/v1/chat/completions',
                apiKey: 'secret',
                model: 'gpt-4o-mini',
                commitFormat: 'gitmoji',
                customFormat: 'My custom template'
            };

            assert.strictEqual(config.apiType, 'openai-compatible');
            assert.strictEqual(config.apiUrl, 'https://gateway.example.com/v1/chat/completions');
            assert.strictEqual(config.apiKey, 'secret');
            assert.strictEqual(config.model, 'gpt-4o-mini');
            assert.strictEqual(config.commitFormat, 'gitmoji');
            assert.strictEqual(config.customFormat, 'My custom template');
        });

        test('should accept Ollama options', () => {
            const config: LlamitConfig = {
                apiType: 'ollama',
                apiUrl: 'http://localhost:11434/api/generate',
                model: 'qwen2.5-coder:7b',
                commitFormat: 'conventional',
                customFormat: '',
                keepAlive: '5m',
                temperature: 0.7,
                topK: 40,
                topP: 0.9,
                numCtx: 4096,
                numPredict: 512,
                repeatPenalty: 1.2,
                repeatLastN: 32,
                seed: 42,
                numGpu: 1,
                numThread: 4,
                minP: 0.05,
                tfsZ: 1.0,
                mirostat: 2,
                mirostatEta: 0.1,
                mirostatTau: 5.0,
                stop: 'END,STOP'
            };

            assert.strictEqual(config.keepAlive, '5m');
            assert.strictEqual(config.temperature, 0.7);
            assert.strictEqual(config.topK, 40);
            assert.strictEqual(config.topP, 0.9);
            assert.strictEqual(config.numCtx, 4096);
            assert.strictEqual(config.numPredict, 512);
            assert.strictEqual(config.repeatPenalty, 1.2);
            assert.strictEqual(config.repeatLastN, 32);
            assert.strictEqual(config.seed, 42);
            assert.strictEqual(config.numGpu, 1);
            assert.strictEqual(config.numThread, 4);
            assert.strictEqual(config.minP, 0.05);
            assert.strictEqual(config.tfsZ, 1.0);
            assert.strictEqual(config.mirostat, 2);
            assert.strictEqual(config.mirostatEta, 0.1);
            assert.strictEqual(config.mirostatTau, 5.0);
            assert.strictEqual(config.stop, 'END,STOP');
        });

        test('should accept optional Ollama options', () => {
            const config: LlamitConfig = {
                apiType: 'ollama',
                apiUrl: 'http://localhost:11434/api/generate',
                model: 'qwen2.5-coder:7b',
                commitFormat: 'conventional',
                customFormat: ''
            };

            assert.strictEqual(config.keepAlive, undefined);
            assert.strictEqual(config.temperature, undefined);
            assert.strictEqual(config.topK, undefined);
        });
    });
});
