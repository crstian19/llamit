import * as assert from 'assert';
import { shouldShowStarNudge, shouldShowReleaseNotes } from '../../starNudgeLogic';

/**
 * Unit tests for the star-nudge / release-notes decision logic.
 * These cover the pure predicates only; the vscode side effects in
 * starNudge.ts are exercised manually / in the integration suite.
 */
suite('Star Nudge Logic', () => {
    suite('shouldShowStarNudge', () => {
        test('shows on a fresh install that has not been dismissed', () => {
            assert.strictEqual(shouldShowStarNudge('install', false), true);
        });

        test('does not show once dismissed, even on install', () => {
            assert.strictEqual(shouldShowStarNudge('install', true), false);
        });

        test('does not show on update', () => {
            assert.strictEqual(shouldShowStarNudge('update', false), false);
        });

        test('does not show on a plain restart (no event)', () => {
            assert.strictEqual(shouldShowStarNudge(null, false), false);
        });
    });

    suite('shouldShowReleaseNotes', () => {
        test('shows on update when this version has not been shown yet', () => {
            assert.strictEqual(shouldShowReleaseNotes('update', '2.1.0', '2.0.0'), true);
        });

        test('shows on update when nothing was ever shown', () => {
            assert.strictEqual(shouldShowReleaseNotes('update', '2.1.0', undefined), true);
        });

        test('does not re-show the same version twice', () => {
            assert.strictEqual(shouldShowReleaseNotes('update', '2.1.0', '2.1.0'), false);
        });

        test('does not show on a fresh install', () => {
            assert.strictEqual(shouldShowReleaseNotes('install', '2.1.0', undefined), false);
        });

        test('does not show on a plain restart (no event)', () => {
            assert.strictEqual(shouldShowReleaseNotes(null, '2.1.0', '2.0.0'), false);
        });
    });
});
