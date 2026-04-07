import * as path from 'path';
import * as fs from 'fs';
import { execFile } from 'child_process';

// Types for better testability
export interface LlamitConfig {
    ollamaUrl: string;
    model: string;
    commitFormat: string;
    customFormat: string;
    keepAlive?: string;
    temperature?: number;
    topK?: number;
    topP?: number;
    numCtx?: number;
    numPredict?: number;
    repeatPenalty?: number;
    repeatLastN?: number;
    seed?: number;
    numGpu?: number;
    numThread?: number;
    minP?: number;
    tfsZ?: number;
    mirostat?: number;
    mirostatEta?: number;
    mirostatTau?: number;
    stop?: string;
}

export interface GitDiffResult {
    diff: string;
    isEmpty: boolean;
}

// Get the binary path for the current platform
export function getBinaryPath(extensionPath: string): string {
    const platform = process.platform;
    const arch = process.arch;

    let binaryName = '';

    if (platform === 'win32') {
        binaryName = 'llamit-windows-amd64.exe';
    } else if (platform === 'darwin') {
        binaryName = arch === 'arm64' ? 'llamit-darwin-arm64' : 'llamit-darwin-amd64';
    } else if (platform === 'linux') {
        binaryName = arch === 'arm64' ? 'llamit-linux-arm64' : 'llamit-linux-amd64';
    } else {
        throw new Error(`Unsupported platform: ${platform}-${arch}`);
    }

    const binaryPath = path.join(extensionPath, 'go-cli', 'bin', binaryName);

    // Make sure the binary is executable on Unix-like systems
    if (platform !== 'win32') {
        try {
            fs.chmodSync(binaryPath, 0o755);
        } catch (error) {
            console.error(`Failed to verify/set permissions for binary: ${error}`);
        }
    }

    return binaryPath;
}

// Helper function to execute git diff commands
function executeGitDiff(
    gitPath: string,
    repositoryRoot: string,
    args: string[]
): Promise<string> {
    return new Promise<string>((resolve, reject) => {
        execFile(gitPath, args, { cwd: repositoryRoot }, (error: Error | null, stdout: string, stderr: string) => {
            if (error && !stdout) {
                reject(new Error(stderr || error.message));
                return;
            }
            resolve(stdout);
        });
    });
}

/**
 * Gets git diff with smart fallback:
 * 1. Try staged changes (--cached)
 * 2. If empty, try working directory changes
 * 3. If both empty, return isEmpty: true
 */
export async function getGitDiffCascade(
    gitPath: string,
    repositoryRoot: string
): Promise<GitDiffResult> {
    try {
        // First: try staged changes
        const stagedDiff = await executeGitDiff(
            gitPath,
            repositoryRoot,
            ['diff', '--cached']
        );

        if (stagedDiff.trim() !== '') {
            return { diff: stagedDiff, isEmpty: false };
        }

        // Second: try working directory changes
        const workingDiff = await executeGitDiff(
            gitPath,
            repositoryRoot,
            ['diff']
        );

        if (workingDiff.trim() !== '') {
            return { diff: workingDiff, isEmpty: false };
        }

        // Both empty
        return { diff: '', isEmpty: true };

    } catch (error) {
        throw error;
    }
}

/**
 * Execute git diff --cached to get staged changes
 * @deprecated Use getGitDiffCascade() instead for smart fallback behavior
 */
export function getGitDiff(gitPath: string, repositoryRoot: string): Promise<GitDiffResult> {
    return new Promise<GitDiffResult>((resolve, reject) => {
        execFile(gitPath, ['diff', '--cached'], { cwd: repositoryRoot }, (error: Error | null, stdout: string, stderr: string) => {
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

        // Add keep-alive if provided
        if (config.keepAlive) {
            args.push('-keep-alive', config.keepAlive);
        }

        // Add Ollama options if provided
        if (config.temperature !== undefined && config.temperature !== 0) {
            args.push('-temperature', config.temperature.toString());
        }
        if (config.topK !== undefined && config.topK !== 0) {
            args.push('-top-k', config.topK.toString());
        }
        if (config.topP !== undefined && config.topP !== 0) {
            args.push('-top-p', config.topP.toString());
        }
        if (config.numCtx !== undefined && config.numCtx !== 0) {
            args.push('-num-ctx', config.numCtx.toString());
        }
        if (config.numPredict !== undefined && config.numPredict !== 0) {
            args.push('-num-predict', config.numPredict.toString());
        }
        if (config.repeatPenalty !== undefined && config.repeatPenalty !== 0) {
            args.push('-repeat-penalty', config.repeatPenalty.toString());
        }
        if (config.repeatLastN !== undefined && config.repeatLastN !== 0) {
            args.push('-repeat-last-n', config.repeatLastN.toString());
        }
        if (config.seed !== undefined && config.seed !== 0) {
            args.push('-seed', config.seed.toString());
        }
        if (config.numGpu !== undefined && config.numGpu !== 0) {
            args.push('-num-gpu', config.numGpu.toString());
        }
        if (config.numThread !== undefined && config.numThread !== 0) {
            args.push('-num-thread', config.numThread.toString());
        }
        if (config.minP !== undefined && config.minP !== 0) {
            args.push('-min-p', config.minP.toString());
        }
        if (config.tfsZ !== undefined && config.tfsZ !== 0) {
            args.push('-tfs-z', config.tfsZ.toString());
        }
        if (config.mirostat !== undefined && config.mirostat !== 0) {
            args.push('-mirostat', config.mirostat.toString());
        }
        if (config.mirostatEta !== undefined && config.mirostatEta !== 0) {
            args.push('-mirostat-eta', config.mirostatEta.toString());
        }
        if (config.mirostatTau !== undefined && config.mirostatTau !== 0) {
            args.push('-mirostat-tau', config.mirostatTau.toString());
        }
        if (config.stop) {
            args.push('-stop', config.stop);
        }

        const child = execFile(
            binaryPath,
            args,
            (error: Error | null, stdout: string, stderr: string) => {
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

export interface Repository {
    rootUri: {
        fsPath: string;
    };
    inputBox: {
        value: string;
    };
}
