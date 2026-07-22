import * as vscode from 'vscode';
import { TelemetryEventName } from './telemetry';
import { shouldShowStarNudge, shouldShowReleaseNotes } from './starNudgeLogic';

export const STAR_NUDGE_DISMISSED_KEY = 'llamit.starNudge.dismissed';
export const RELEASE_NOTES_SHOWN_KEY = 'llamit.releaseNotes.lastShownVersion';

const REPO_URL = 'https://github.com/crstian19/llamit';

/**
 * On a fresh install, shows a one-time, opt-in prompt inviting the user to star
 * the repo on GitHub. It never repeats: the dismissed flag is persisted before
 * the message is awaited, so even if activation happens again the toast is gone.
 *
 * This is a gentle nudge, not an automatic action — the user must click through
 * and star the repo themselves while logged in. Nothing is done on their behalf.
 */
export async function maybeShowStarNudge(
    context: vscode.ExtensionContext,
    event: TelemetryEventName | null
): Promise<void> {
    const dismissed = context.globalState.get<boolean>(STAR_NUDGE_DISMISSED_KEY, false);
    if (!shouldShowStarNudge(event, dismissed)) {
        return;
    }

    // Persist first so the nudge can only ever appear once, regardless of what
    // the user picks (or if they dismiss the notification without choosing).
    await context.globalState.update(STAR_NUDGE_DISMISSED_KEY, true);

    const selection = await vscode.window.showInformationMessage(
        'Thanks for installing Llamit! If it saves you time, a ⭐ on GitHub really helps the project.',
        '⭐ Star on GitHub',
        'Maybe later'
    );

    if (selection === '⭐ Star on GitHub') {
        await vscode.env.openExternal(vscode.Uri.parse(REPO_URL));
    }
}

/**
 * On an update to a new version, opens the bundled CHANGELOG in a rendered
 * Markdown preview so the user sees what's new. Shown at most once per version
 * (tracked by RELEASE_NOTES_SHOWN_KEY). If the Markdown preview command isn't
 * available (e.g. a stripped-down fork), it silently falls back to the GitHub
 * release page for the current version.
 */
export async function maybeShowReleaseNotes(
    context: vscode.ExtensionContext,
    event: TelemetryEventName | null,
    currentVersion: string,
    extensionPath: string
): Promise<void> {
    const lastShown = context.globalState.get<string>(RELEASE_NOTES_SHOWN_KEY);
    if (!shouldShowReleaseNotes(event, currentVersion, lastShown)) {
        return;
    }

    await context.globalState.update(RELEASE_NOTES_SHOWN_KEY, currentVersion);

    const changelogUri = vscode.Uri.joinPath(vscode.Uri.file(extensionPath), 'CHANGELOG.md');

    try {
        await vscode.commands.executeCommand('markdown.showPreview', changelogUri);
    } catch {
        // Markdown preview not available — fall back to the online release notes.
        await vscode.env.openExternal(
            vscode.Uri.parse(`${REPO_URL}/releases/tag/v${currentVersion}`)
        );
    }
}
