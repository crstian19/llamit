import * as vscode from 'vscode';

const UMAMI_ENDPOINT = 'https://umami.crstian.me/api/send';
const UMAMI_WEBSITE_ID = '10cfe8a0-4329-40a5-88c5-83b8ca2cc552';

export type TelemetryEventName = 'install' | 'update';

/**
 * Sends a single anonymous ping to the self-hosted Umami analytics instance.
 *
 * What is sent:
 *   - The event name: "install" (first time) or "update" (version changed)
 *   - The IDE app name (e.g. "Visual Studio Code", "VSCodium", "Cursor")
 *   - The extension version string
 *
 * What is NOT sent:
 *   - Any user identifiers, file paths, repository names, or PII
 *   - Commit messages, diffs, or any code content
 *   - IP address is not stored by Umami (configured server-side)
 *
 * The ping is skipped if:
 *   - vscode.env.isTelemetryEnabled is false (user disabled VS Code telemetry)
 *   - llamit.telemetry.enabled setting is false (user opted out)
 *   - The stored version matches the current version (no install/update occurred)
 */
export async function sendTelemetry(
    context: vscode.ExtensionContext,
    eventName: TelemetryEventName,
    version: string
): Promise<void> {
    // Respect VS Code's global telemetry setting
    if (!vscode.env.isTelemetryEnabled) {
        return;
    }

    // Respect the llamit-specific opt-out setting
    const config = vscode.workspace.getConfiguration('llamit');
    const telemetryEnabled = config.get<boolean>('telemetry.enabled', true);
    if (!telemetryEnabled) {
        return;
    }

    const payload = {
        payload: {
            website: UMAMI_WEBSITE_ID,
            url: '/llamit/activation',
            name: eventName,
            data: {
                client: vscode.env.appName,
                version: version
            }
        },
        type: 'event'
    };

    try {
        // Fire-and-forget: we send the ping but do not await or surface errors to the user.
        // Using native fetch (Node 18+, available in VS Code extension runtime).
        await fetch(UMAMI_ENDPOINT, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
            signal: AbortSignal.timeout(5000)
        });
    } catch {
        // Silently ignore network errors — telemetry must never break the extension
    }
}

/**
 * Checks whether a telemetry event should fire based on stored version state.
 * Returns the event name ('install' | 'update') or null if the version hasn't changed.
 */
export function getTelemetryEvent(
    context: vscode.ExtensionContext,
    currentVersion: string
): TelemetryEventName | null {
    const storedVersion = context.globalState.get<string>('llamit.lastTrackedVersion');

    if (!storedVersion) {
        return 'install';
    }

    if (storedVersion !== currentVersion) {
        return 'update';
    }

    return null;
}

/**
 * Persists the current version so subsequent activations won't re-fire the event.
 */
export async function persistTrackedVersion(
    context: vscode.ExtensionContext,
    version: string
): Promise<void> {
    await context.globalState.update('llamit.lastTrackedVersion', version);
}
