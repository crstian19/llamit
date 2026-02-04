import * as assert from 'assert';
import * as vscode from 'vscode';

suite('Extension Integration Tests', () => {
    vscode.window.showInformationMessage('Running integration tests.');

    test('Extension should be present', () => {
        assert.ok(vscode.extensions.getExtension('llamit-dev.llamit'));
    });

    test('Extension should activate', async () => {
        const ext = vscode.extensions.getExtension('llamit-dev.llamit');
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

        assert.strictEqual(ollamaUrl, 'http://localhost:11434/api/generate');
        assert.strictEqual(model, 'qwen3-coder:30b');
    });
});
