import * as vscode from 'vscode';
import { execFile } from 'child_process';
import * as path from 'path';
import { getBinaryPath, LlamitConfig, generateCommitMessage, GitDiffResult, getGitDiffCascade, Repository } from './helpers';
import { sendTelemetry, getTelemetryEvent, persistTrackedVersion } from './telemetry';

async function getCurrentRepository(repositories: Repository[]): Promise<Repository | undefined> {
    if (repositories.length === 0) {
        return undefined;
    }

    if (repositories.length === 1) {
        return repositories[0];
    }

    // Always ask the user to select when there are multiple repositories
    const selected = await vscode.window.showQuickPick(
        repositories.map(repo => ({
            label: path.basename(repo.rootUri.fsPath),
            description: repo.rootUri.fsPath,
            repo: repo
        })),
        { placeHolder: 'Select a repository for commit message' }
    );

    return selected?.repo;
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
        model: config.get<string>('model') || 'qwen2.5-coder:7b',
        commitFormat: config.get<string>('commitFormat') || 'conventional',
        customFormat: config.get<string>('customFormat') || '',
        keepAlive: config.get<string>('keepAlive'),
        temperature: config.get<number>('temperature'),
        topK: config.get<number>('topK'),
        topP: config.get<number>('topP'),
        numCtx: config.get<number>('numCtx'),
        numPredict: config.get<number>('numPredict'),
        repeatPenalty: config.get<number>('repeatPenalty'),
        repeatLastN: config.get<number>('repeatLastN'),
        seed: config.get<number>('seed'),
        numGpu: config.get<number>('numGpu'),
        numThread: config.get<number>('numThread'),
        minP: config.get<number>('minP'),
        tfsZ: config.get<number>('tfsZ'),
        mirostat: config.get<number>('mirostat'),
        mirostatEta: config.get<number>('mirostatEta'),
        mirostatTau: config.get<number>('mirostatTau'),
        stop: config.get<string>('stop')
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

        // Get the current repository based on active editor or workspace context
        const repo = await getCurrentRepository(git.repositories as unknown as Repository[]);
        if (!repo) {
            vscode.window.showWarningMessage('Could not determine the active Git repository.');
            return;
        }

        vscode.window.withProgress({
            location: vscode.ProgressLocation.SourceControl,
            title: 'Llamit: Generating commit message...',
        }, async () => {
            try {
                // 1. Get the git diff (staged changes with fallback to working directory)
                // Use global git path with repository-specific working directory
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

    // Send anonymous telemetry on install or update only (not on every activation)
    const currentVersion = context.extension.packageJSON.version as string | undefined;
    if (currentVersion) {
        const event = getTelemetryEvent(context, currentVersion);
        if (event) {
            sendTelemetry(context, event, currentVersion).then(() => {
                persistTrackedVersion(context, currentVersion);
            }).catch(() => { /* ignore */ });
        }
    }
}

export function deactivate() { }
