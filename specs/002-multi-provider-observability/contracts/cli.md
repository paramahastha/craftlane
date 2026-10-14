# CLI Contract: Multi-Provider Observation

## Commands

| Command | Purpose | Mutation policy |
|---|---|---|
| `craftlane observe status` | Report safe setup/trust/load/callback/capability state for all three adapters. | Read-only. |
| `craftlane observe setup <claude|codex|opencode>` | Detect provider and print exact project-artifact diff, trust/load requirement, tracking warning, and secret-storage disclosure. | Read-only preview. |
| `craftlane observe setup <provider> --confirm` | Apply exactly the prior safe preview and initialize/retain the local observation secret. | Explicit provider-scoped mutation. |
| `craftlane observe remove <provider>` | Print owned entries/file and credential consequence. | Read-only preview. |
| `craftlane observe remove <provider> --confirm` | Remove only marked owned provider setup; remove shared secret only when no adapter needs it. | Explicit provider-scoped mutation. |
| `craftlane observe bridge --provider <provider>` | Private fixed bridge used by provider integration artifact. | Never a user workflow; no decision output. |

## Provider artifact rules

| Provider | Artifact | Required safety behavior |
|---|---|---|
| Claude Code | `.claude/settings.local.json` | Structurally merge/remove only marked Craftlane hooks; preserve unrelated settings. |
| Codex | `.codex/hooks.json` or documented project hook location | Structurally merge/remove only marked hooks; clearly warn that project configuration can be tracked; expose `needs_trust` until provider hook review succeeds. |
| OpenCode | `.opencode/plugins/craftlane-observability.<supported-extension>` | Write/remove only the uniquely owned version-compatible plugin; warn it can be tracked; do not edit permissions or general config. |

No command silently writes user-wide settings, adds an ignore rule, overwrites malformed/conflicting owned content, uses untrusted user input in a command, or exposes a callback secret. A preview is required before every write; repeated setup is idempotent; removal preserves unrelated configuration.
