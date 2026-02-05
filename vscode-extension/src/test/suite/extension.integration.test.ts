import * as assert from 'assert';
import * as vscode from 'vscode';

suite('Extension Integration Tests', () => {
    vscode.window.showInformationMessage('Running integration tests.');

    test('Extension should be present', () => {
        assert.ok(vscode.extensions.getExtension('Crstian.llamit'));
    });

    test('Extension should activate', async () => {
        const ext = vscode.extensions.getExtension('Crstian.llamit');
        assert.ok(ext);

        await ext!.activate();
        assert.strictEqual(ext!.isActive, true);
    });

    test('Should register llamit.generateCommit command', async () => {
        const commands = await vscode.commands.getCommands();
        assert.ok(commands.includes('llamit.generateCommit'));
    });

    test('Configuration should have default values', () => {
        const config = vscode.workspace.getConfiguration('llamit');

        const ollamaUrl = config.get<string>('ollamaUrl');
        const model = config.get<string>('model');
        const commitFormat = config.get<string>('commitFormat');
        const customFormat = config.get<string>('customFormat');

        assert.strictEqual(ollamaUrl, 'http://localhost:11434/api/generate');
        assert.strictEqual(model, 'qwen2.5-coder:7b');
        assert.strictEqual(commitFormat, 'conventional');
        assert.strictEqual(customFormat, '');
    });
});
