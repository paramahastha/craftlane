# Research: Multi-Provider Agent Observability

## Decision: Use a private registry and one normalized observation pipeline

- **Decision**: Add a private adapter registry for Claude Code, Codex, and OpenCode. Adapters own detection, setup, raw projection, and normalization; a shared authenticated ingestion route and the existing reducer own validation, persistence, ordering, replay, and dashboard delivery.
- **Rationale**: This preserves provider-neutral canonical state and avoids three storage/event pipelines. It also prevents provider payload formats from reaching the UI or public protocol.
- **Alternatives considered**: A public SDK is premature before the three adapters validate the boundary. A generic raw webhook format was rejected because it would invite unsafe arbitrary payload persistence.

## Decision: Use explicit project-scoped integrations, with provider-specific safety disclosures

- **Decision**: Claude uses marked command hooks in `.claude/settings.local.json`; Codex uses marked project hooks in `.codex/hooks.json` or its documented project configuration location; OpenCode uses a uniquely named project plugin under `.opencode/plugins/`. Every write has an exact preview and needs confirmation. Removal is structural/ownership-scoped.
- **Rationale**: Each provider has a different official local integration surface. Codex/OpenCode artifacts may be Git-tracked, so setup must warn before the explicit write and must never silently add ignore rules. Settings/artifacts contain no callback secret.
- **Alternatives considered**: User-wide configuration was rejected because it violates project scope. Shared generated scripts and direct static-token callbacks were rejected because they broaden code/secret exposure.
- **Sources**: [Claude hooks](https://code.claude.com/docs/en/hooks), [Codex hooks](https://developers.openai.com/codex/hooks), [Codex advanced configuration](https://developers.openai.com/codex/config-advanced), [OpenCode plugins](https://opencode.ai/docs/plugins/).

## Decision: Capability-gate lifecycle coverage rather than require parity

- **Decision**: Each adapter declares available/unavailable session, subagent, tool, permission, failure, and terminal/idle capabilities. The UI and diagnostics expose the actual coverage. No provider is inferred to have Claude-equivalent semantics.
- **Rationale**: Codex `SessionEnd` may be delayed and does not cover subagents; OpenCode status exposes busy/idle/retry rather than a durable successful terminal signal; provider versions/load modes can limit hooks/plugins. Honest capability reporting is more useful than false completion state.
- **Alternatives considered**: Identical feature requirements were rejected as unreliable and likely to block delivery. Session-only support was rejected because it loses the core visibility value where richer local signals exist.
- **Sources**: [Codex hooks](https://developers.openai.com/codex/hooks), [OpenCode plugin API](https://github.com/anomalyco/opencode/blob/dev/packages/plugin/src/index.ts), [OpenCode session status schema](https://github.com/anomalyco/opencode/blob/dev/packages/schema/src/session-status-event.ts).

## Decision: Never send raw provider data beyond an adapter bridge

- **Decision**: Bridge envelopes contain only version, provider, opaque project/session/agent/correlation IDs, event kind, safe tool name/outcome, and safe timestamps. They discard prompts, messages, commands, transcript paths/content, environment values, tool arguments/results, permission metadata, and arbitrary fields. Existing redaction is a second sink defense.
- **Rationale**: All three providers may deliver private source/prompt/tool content in hook/plugin inputs. Redaction patterns alone cannot make arbitrary free text safe to persist or present.
- **Alternatives considered**: Parsing transcripts and forwarding full payloads were rejected for privacy, unstable formats, and unnecessary retention.
- **Sources**: [Claude hook common fields](https://code.claude.com/docs/en/hooks#common-input-fields), [Codex hooks](https://developers.openai.com/codex/hooks), [OpenCode plugins](https://opencode.ai/docs/plugins/).

## Decision: Use authenticated loopback ingestion with idempotent state guards

- **Decision**: A separate owner-only per-project observation secret authenticates a bounded `POST /api/v1/observe/{provider}` request. The reducer fingerprints provider event correlation data and rejects duplicates, malformed/wrong-project requests, and late nonterminal updates that would revive terminal state. Explicit provider state overrides inferred idle.
- **Rationale**: Command hooks/plugins may duplicate or reorder delivery. A local Craftlane sequence is delivery order only, not proof of provider causal order. The callback must not share runtime-control authority.
- **Alternatives considered**: Unauthenticated localhost ingestion, a static secret in provider settings, and last-write-wins updates were rejected for security and incorrect lifecycle state.

## Decision: Keep hooks/plugins non-controlling and bounded

- **Decision**: Adapters return no provider decision output, use fixed executable invocation/arguments, bounded input/body sizes, and short local delivery. Claude/OpenCode observation is fail-open; Codex uses its documented trusted hook path and bounded synchronous handling where terminal accounting needs reliability.
- **Rationale**: Craftlane is an observability layer, never an agent permission or execution controller. Codex requires exact hook review/trust; async hooks can lose events at terminal boundaries.
- **Alternatives considered**: Provider permission modification, retries that create uncontrolled duplicate work, and direct provider cloud telemetry were rejected.
- **Sources**: [Claude hook limitations](https://code.claude.com/docs/en/hooks#limitations), [Codex hook trust and limitations](https://developers.openai.com/codex/hooks), [OpenCode plugin migration](https://opencode.ai/v2/docs/build/plugins/migrate-v1).

## Decision: Fixtures are synthetic and provider-separated

- **Decision**: Maintain a versioned synthetic fixture set per adapter with fake IDs, `{{PROJECT_ROOT}}`, and credential-shaped canaries. Optional smoke tests run only when a provider is installed/configured; CI validates fixture behavior without accounts or network.
- **Rationale**: Provider payload captures can leak prompts, paths, credentials, or account data. Fixture-driven tests make unavailable provider installations a clear capability result rather than a build blocker.
- **Alternatives considered**: Captured transcripts and mandatory installed-provider end-to-end tests were rejected for privacy and reproducibility.
