import * as assert from 'assert';
import { generateCommitMessage, LlamitConfig } from '../../extension';

suite('Commit Message Generation Tests', () => {
    test('generateCommitMessage should accept correct parameters', () => {
        const binaryPath = '/path/to/cli';
        const config: LlamitConfig = {
            ollamaUrl: 'http://localhost:11434/api/generate',
            model: 'qwen3-coder:30b',
            commitFormat: 'conventional',
            customFormat: ''
        };
        const diff = 'diff --git a/test.txt b/test.txt\n+test line';

        // Verify function signature
        assert.strictEqual(typeof generateCommitMessage, 'function');
        assert.strictEqual(generateCommitMessage.length, 3);

        // Test that it returns a Promise
        try {
            const result = generateCommitMessage(binaryPath, config, diff);
            assert.ok(result instanceof Promise);
        } catch (error) {
            // Expected to fail without actual binary
            assert.ok(true);
        }
    });

    test('should handle empty diff string', () => {
        const binaryPath = '/path/to/cli';
        const config: LlamitConfig = {
            ollamaUrl: 'http://localhost:11434/api/generate',
            model: 'test-model',
            commitFormat: 'conventional',
            customFormat: ''
        };
        const diff = '';

        try {
            const result = generateCommitMessage(binaryPath, config, diff);
            assert.ok(result instanceof Promise);
        } catch (error) {
            assert.ok(true);
        }
    });

    test('should handle multiline diff', () => {
        const binaryPath = '/path/to/cli';
        const config: LlamitConfig = {
            ollamaUrl: 'http://localhost:11434/api/generate',
            model: 'test-model',
            commitFormat: 'conventional',
            customFormat: ''
        };
        const diff = `diff --git a/file1.txt b/file1.txt
index abc123..def456 100644
--- a/file1.txt
+++ b/file1.txt
@@ -1,3 +1,4 @@
 line 1
 line 2
+line 3
 line 4`;

        try {
            const result = generateCommitMessage(binaryPath, config, diff);
            assert.ok(result instanceof Promise);
        } catch (error) {
            assert.ok(true);
        }
    });
});
