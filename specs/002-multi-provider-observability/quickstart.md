# Quickstart: Validate Multi-Provider Agent Observability

## Prerequisites

- Built Craftlane executable and contributor web assets.
- Local Git repository.
- Any of Claude Code, Codex, or OpenCode installed locally. Tests must remain runnable with none installed.

## 1. Inspect provider readiness

```sh
craftlane start --no-open
craftlane doctor
craftlane observe status
```

Expected: loopback-only runtime is healthy. Each provider independently reports executable, setup/trust/load, callback, path, and capability status with a safe next action. Missing providers do not block the dashboard or other adapters.

## 2. Preview and confirm one provider at a time

```sh
craftlane observe setup claude
craftlane observe setup claude --confirm
craftlane observe setup codex
craftlane observe setup opencode
```

Expected: every first command is read-only and shows its exact project artifact. Claude uses local settings; Codex/OpenCode warn that their project artifact may be Git-tracked. Confirmation creates only the owned artifact/entries, never user-wide settings, credentials in provider configuration, or ignore rules. Codex additionally reports hook trust requirements. See [CLI contract](./contracts/cli.md).

## 3. Verify live observation and capability coverage

For each configured installed provider, start a session; perform a safe local tool operation; create a subagent or trigger permission only where the provider supports it.

Expected: the dashboard updates without refresh, showing provider, safe tool activity, state source, lineage when supplied, and attention status. Its capability view states unavailable categories rather than showing fabricated data. A session first seen after setup is marked as having unavailable prior history.

## 4. Verify recovery and safety

Disconnect/reconnect the dashboard repeatedly during fixture-driven agent activity, then restart Craftlane. Exercise wrong token/provider/project/path, oversized/malformed input, duplicate/reordered signals, synthetic credential canaries, and pre-existing unrelated provider configuration.

Expected: retained updates replay in local sequence without duplicate application; no unsafe input changes state; raw content/secrets are absent from storage, diagnostics, HTTP, WebSocket, and browser; revived pre-restart agents are not falsely marked active; setup/removal preserves unrelated configuration.

## 5. Remove safely

```sh
craftlane observe remove claude
craftlane observe remove claude --confirm
```

Expected: removal previews first, deletes only owned setup, retains unrelated provider settings, and leaves the normal Craftlane runtime/dashboard and other provider adapters working.

## Automated validation record

On 2026-10-02, the synthetic provider, reducer, storage, callback, bridge,
setup, dashboard, integration, race, and redaction suites passed. The full Go
suite passed with 100 reconnect-cycle coverage, and web typecheck, lint, unit,
and production-build checks passed. Installed-provider smoke tests and the
five-developer usability study remain manual follow-ups because they require
local provider installations and human participants.
