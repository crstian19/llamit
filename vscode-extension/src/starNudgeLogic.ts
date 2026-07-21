import type { TelemetryEventName } from './telemetry';

/**
 * Pure decision logic for the star nudge / release notes prompts.
 *
 * These functions intentionally import NO vscode APIs so they can be exercised
 * by the fast unit test suite (tsx, no Electron runtime). The vscode side
 * effects live in starNudge.ts.
 */

/**
 * The star nudge is shown only on a fresh install, and only once ever
 * (subsequent installs on the same profile are already marked dismissed).
 */
export function shouldShowStarNudge(
    event: TelemetryEventName | null,
    dismissed: boolean
): boolean {
    return event === 'install' && !dismissed;
}

/**
 * Release notes are shown once per version, only when the activation was an
 * update (not a first install, not a plain restart on the same version).
 */
export function shouldShowReleaseNotes(
    event: TelemetryEventName | null,
    currentVersion: string,
    lastShownVersion: string | undefined
): boolean {
    return event === 'update' && lastShownVersion !== currentVersion;
}
