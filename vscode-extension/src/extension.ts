import * as vscode from 'vscode';
import { execFile } from 'child_process';
import * as path from 'path';

// Types for better testability
export interface LlamitConfig {
    ollamaUrl: string;
    model: string;
    commitFormat: string;
    customFormat: string;
}

export interface GitDiffResult {
    diff: string;
    isEmpty: boolean;
}

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
        model: config.get<string>('model') || 'qwen3-coder:30b',
        commitFormat: config.get<string>('commitFormat') || 'conventional',
        customFormat: config.get<string>('customFormat') || ''
    };
}

// Get the binary path for the current platform
export function getBinaryPath(extensionPath: string): string {
    const binaryName = process.platform === 'win32' ? 'cli.exe' : 'cli';
    return path.join(extensionPath, 'go-cli', binaryName);
}

// Execute git diff --cached to get staged changes
export function getGitDiff(gitPath: string, repositoryRoot: string): Promise<GitDiffResult> {
    return new Promise<GitDiffResult>((resolve, reject) => {
        execFile(gitPath, ['diff', '--cached'], { cwd: repositoryRoot }, (error, stdout, stderr) => {
            if (error) {
                // It's not a true error if stdout is empty, just means no staged changes.
                if (stdout.trim() === '') {
                    resolve({ diff: '', isEmpty: true });
                    return;
                }
                reject(new Error(stderr || error.message));
                return;
            }
            resolve({ diff: stdout, isEmpty: stdout.trim() === '' });
        });
    });
}

// Execute the Llamit CLI binary to generate commit message
export function generateCommitMessage(
    binaryPath: string,
    config: LlamitConfig,
    diff: string
): Promise<string> {
    return new Promise<string>((resolve, reject) => {
        const args = [
            '-ollama-url', config.ollamaUrl,
            '-model', config.model,
            '-format', config.commitFormat
        ];

        // Add custom template if format is 'custom' and template is provided
        if (config.commitFormat === 'custom' && config.customFormat) {
            args.push('-custom-template', config.customFormat);
        }

        const child = execFile(
            binaryPath,
            args,
            (error, stdout, stderr) => {
                if (error) {
                    console.error(`execFile error for llamit cli: ${error.message}`);
                    reject(new Error(stderr || error.message));
                    return;
                }
                resolve(stdout.trim());
            }
        );

        // Write the diff to the binary's stdin
        child.stdin?.write(diff);
        child.stdin?.end();
    });
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
                // 1. Get the staged diff
                const gitPath = git.git.path;
                const repositoryRoot = repo.rootUri.fsPath;
                const diffResult = await getGitDiff(gitPath, repositoryRoot);

                if (diffResult.isEmpty) {
                    vscode.window.showInformationMessage('No staged changes to commit.');
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
