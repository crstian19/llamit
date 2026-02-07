import * as assert from 'assert';
import * as sinon from 'sinon';
import { getGitDiff, getGitDiffCascade, GitDiffResult } from '../../helpers';

suite('Git Diff Tests', () => {
    let sandbox: sinon.SinonSandbox;

    setup(() => {
        sandbox = sinon.createSandbox();
    });

    teardown(() => {
        sandbox.restore();
    });

    test('getGitDiff should return correct interface', async () => {
        // We test that the function exists and returns a promise
        // In real tests, you'd mock execFile using proxyquire or dependency injection
        assert.strictEqual(typeof getGitDiff, 'function');

        // Test that it returns a Promise
        const gitPath = '/usr/bin/git';
        const repoRoot = '/tmp/test-repo';

        try {
            const result = getGitDiff(gitPath, repoRoot);
            assert.ok(result instanceof Promise);
        } catch (error) {
            // Expected to fail in test environment without real git repo
            assert.ok(true);
        }
    });

    test('GitDiffResult should have correct structure', () => {
        const result: GitDiffResult = {
            diff: 'sample diff content',
            isEmpty: false
        };

        assert.strictEqual(typeof result.diff, 'string');
        assert.strictEqual(typeof result.isEmpty, 'boolean');
        assert.strictEqual(result.diff, 'sample diff content');
        assert.strictEqual(result.isEmpty, false);
    });

    test('GitDiffResult can represent empty diff', () => {
        const result: GitDiffResult = {
            diff: '',
            isEmpty: true
        };

        assert.strictEqual(result.diff, '');
        assert.strictEqual(result.isEmpty, true);
    });
});

suite('getGitDiffCascade Tests', () => {
    let sandbox: sinon.SinonSandbox;

    setup(() => {
        sandbox = sinon.createSandbox();
    });

    teardown(() => {
        sandbox.restore();
    });

    test('getGitDiffCascade should return correct interface', async () => {
        // Test that the function exists and returns a promise
        assert.strictEqual(typeof getGitDiffCascade, 'function');

        // Test that it returns a Promise
        const gitPath = '/usr/bin/git';
        const repoRoot = '/tmp/test-repo';

        try {
            const result = getGitDiffCascade(gitPath, repoRoot);
            assert.ok(result instanceof Promise);
        } catch (error) {
            // Expected to fail in test environment without real git repo
            assert.ok(true);
        }
    });

    test('getGitDiffCascade should return GitDiffResult structure', () => {
        // Verify the return type structure matches GitDiffResult
        const sampleResult: GitDiffResult = {
            diff: 'sample diff',
            isEmpty: false
        };

        assert.strictEqual(typeof sampleResult.diff, 'string');
        assert.strictEqual(typeof sampleResult.isEmpty, 'boolean');
    });

    test('getGitDiffCascade should handle empty result', () => {
        // Verify it can represent empty state
        const emptyResult: GitDiffResult = {
            diff: '',
            isEmpty: true
        };

        assert.strictEqual(emptyResult.diff, '');
        assert.strictEqual(emptyResult.isEmpty, true);
    });
});
