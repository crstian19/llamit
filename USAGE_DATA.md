# Usage Data & Telemetry

Llamit collects minimal, anonymous telemetry to understand how many users install or update the extension and which editors they use. This helps prioritize development and compatibility work.

## What is collected

A single anonymous event is sent to a self-hosted [Umami](https://umami.is/) instance when you **install** or **update** the extension (not on every IDE open).

| Field | Value | Example |
|-------|-------|---------|
| Event name | `install` or `update` | `install` |
| Editor name | `vscode.env.appName` | `Visual Studio Code`, `VSCodium`, `Cursor` |
| Extension version | The installed version string | `1.6.0` |

## What is NOT collected

- No user identifiers, usernames, or account information
- No file paths, repository names, or project metadata
- No commit messages, diffs, or any code content
- No IP addresses (Umami is configured to not store them)
- No usage frequency (command invocations, feature usage, etc.)

## Where data is sent

Data is sent to a self-hosted Umami instance at `https://umami.crstian.me`. The data never leaves this server and is not shared with any third party.

## How to opt out

You have two ways to disable telemetry:

**Option 1 — Disable VS Code telemetry globally:**
Set `telemetry.telemetryLevel` to `off` in your VS Code settings. Llamit respects `vscode.env.isTelemetryEnabled`.

**Option 2 — Disable only Llamit telemetry:**
Set `llamit.telemetry.enabled` to `false` in your VS Code settings:

```json
{
  "llamit.telemetry.enabled": false
}
```

When either of these settings is `false`, no data is sent.
