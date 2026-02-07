import * as vscode from 'vscode';
import { execFile } from 'child_process';
import * as path from 'path';
import { getBinaryPath, LlamitConfig, generateCommitMessage, GitDiffResult, getGitDiffCascade } from './helpers';

// Helper to get the Git API from the built-in VS Code extension
export async function getGitAPI() {
    const gitExtension = vscode.extensions.getExtension('vscode.git');
    if (!gitExtension) {
        vscode.window.showErrorMessage('Git extension is not available.');
        return undefined;
    }
    if (!gitExtension.isActive) {
        await gitExtension.activate();
    }
    return gitExtension.exports.getAPI(1);
}

// Get configuration from VS Code settings
export function getConfiguration(): LlamitConfig {
    const config = vscode.workspace.getConfiguration('llamit');
    return {
        ollamaUrl: config.get<string>('ollamaUrl') || 'http://localhost:11434/api/generate',
        model: config.get<string>('model') || 'qwen2.5-coder:7b',
        commitFormat: config.get<string>('commitFormat') || 'conventional',
        customFormat: config.get<string>('customFormat') || ''
    };
}

export function activate(context: vscode.ExtensionContext) {
    let disposable = vscode.commands.registerCommand('llamit.generateCommit', async () => {
        const git = await getGitAPI();
        if (!git) {
            return;
        }

        if (git.repositories.length === 0) {
            vscode.window.showWarningMessage('No Git repository found.');
            return;
        }

        // Use the first repository for simplicity
        const repo = git.repositories[0];

        vscode.window.withProgress({
            location: vscode.ProgressLocation.SourceControl,
            title: 'Llamit: Generating commit message...',
        }, async () => {
            try {
                // 1. Get the git diff (staged changes with fallback to working directory)
                const gitPath = git.git.path;
                const repositoryRoot = repo.rootUri.fsPath;
                const diffResult = await getGitDiffCascade(gitPath, repositoryRoot);

                if (diffResult.isEmpty) {
                    vscode.window.showInformationMessage('No changes to commit.');
                    return;
                }

                // 2. Get configuration and binary path
                const config = getConfiguration();
                const binaryPath = getBinaryPath(context.extensionPath);

                // 3. Generate commit message
                const commitMessage = await generateCommitMessage(binaryPath, config, diffResult.diff);

                // 4. Update the commit message box
                repo.inputBox.value = commitMessage;

            } catch (e: any) {
                console.error(e);
                vscode.window.showErrorMessage(`Error generating commit message: ${e.message}`);
            }
        });
    });

    context.subscriptions.push(disposable);

    // Register command to open settings
    let openSettingsDisposable = vscode.commands.registerCommand('llamit.openSettings', () => {
        vscode.commands.executeCommand('workbench.action.openSettings', 'llamit');
    });

    context.subscriptions.push(openSettingsDisposable);
}

export function deactivate() { }
