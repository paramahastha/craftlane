# Implementation Plan: Multi-Provider Agent Observability

**Branch**: `002-multi-provider-observability` | **Date**: 2026-10-02 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/002-multi-provider-observability/spec.md`

## Summary

Deliver read-only, local-first observability adapters for Claude Code, Codex, and OpenCode. Each adapter detects supported local lifecycle signals, presents an explicit preview before creating only its documented project-scoped integration artifact, projects raw provider data into a shared allowlisted envelope, and sends it through a generic authenticated local ingestion path. The existing reducer remains the sole owner of normalized agent state, durable event order, snapshot recovery, and browser delivery. The dashboard renders provider-neutral agent state and an explicit per-provider capability matrix; it never guesses that one provider has another provider's coverage.

## Technical Context

**Language/Version**: Go 1.25; TypeScript 5.x for the contributor-built dashboard; generated OpenCode plugin source compatible with the detected supported plugin API.

**Primary Dependencies**: Existing Go standard library and `github.com/coder/websocket`; existing React, Vite, Vitest, React Testing Library, and Playwright. No new production dependency.

**Storage**: Existing owner-only per-project JSON snapshot and bounded JSONL event log, plus one distinct owner-only observation secret and bounded processed-event fingerprints. No database, provider transcripts, or raw provider payload storage.

**Testing**: Go `testing`, `httptest`, temporary Git repositories, synthetic per-provider fixtures, integration tests, and race detector; Vitest/React Testing Library; Playwright against the built executable. Installed-provider smoke tests are optional and skipped with an explicit reason when unavailable.

**Target Platform**: Local macOS, Linux, and Windows developer machines; a current supported Claude Code, Codex, or OpenCode installation; modern desktop browser.

**Project Type**: Native CLI and local web service with an embedded dashboard plus provider-owned project integration artifacts.

**Performance Goals**: For each provider's supported lifecycle record, show a normalized dashboard update within 2 seconds in at least 19 of 20 local test sessions. Provider hook/plugin delivery is bounded and must never block agent execution. Reconnect recovers retained updates without duplicate application.

**Constraints**: Loopback-only; no account, hosted backend, telemetry, or provider transcript parsing; Claude Code, Codex, and OpenCode only; provider-neutral core; explicit preview/confirmation for every configuration write; strict Host and Origin validation for state-changing callback requests; provider-specific capability declaration; no raw prompts, commands, tool arguments/results, messages, paths outside the approved project, environment values, or secrets enter persistence/browser state; additive protocol changes only.

**Scale/Scope**: One Craftlane project/runtime and three initial adapters. This milestone observes only capabilities each provider reliably exposes: sessions, subagents, tools, permissions, failures, and terminal/idle states where available. It does not control agents, alter their permissions, dispatch work, parse Spec Kit, manage worktrees, or add spatial UI.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Pre-design review — PASS**

| Constitution obligation | Design response | Status |
|---|---|---|
| Productivity-first, keyboard-equivalent control | All setup, removal, status, and doctor actions are CLI operations. The dashboard adds semantic, keyboard-accessible situational awareness only. | Pass |
| Local-first, account-free ownership | Every adapter runs locally, uses loopback-only ingestion, and retains state locally. No Craftlane account, cloud service, telemetry, or provider transcript is required. | Pass |
| Safe, explicit, reversible actions | Every provider change is previewed then confirmed; removal is ownership-scoped. Callback authorization is distinct from runtime control; raw provider input never reaches generic state. | Pass |
| Provider-agnostic core and versioned protocols | A private adapter registry terminates provider formats. Core models/callback/reducer/UI use provider ID, capability, and normalized state only. New snapshot/event fields are additive v1 projections. | Pass |
| Spec-driven bounded delivery | This is the expanded first observability phase only. Provider controls, public SDK, worktrees, dispatch, and 3D remain deferred. | Pass |
| Accessible 2D equivalence | Semantic agent rows identify provider, state, reported/inferred source, safe tool status, parent relationship, history availability, and unavailable capabilities without colour-only meaning. | Pass |
| Maintainable, testable open-source delivery | The registry is intentionally private/minimal; adapters have isolated fixtures and setup tests; shared reducer/security/browser paths remain single-source and proportionately tested. | Pass |

**Post-design review — PASS**: The data model maintains one canonical reducer and safe browser projection. Provider-specific setup files never contain callback secrets, and uneven lifecycle coverage is exposed rather than hidden. No exception is required.

## Project Structure

### Documentation (this feature)

```text
specs/002-multi-provider-observability/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── cli.md
│   └── observability.md
└── tasks.md                 # Created later by $speckit-tasks
```

### Source Code (repository root)

```text
cmd/craftlane/main.go                          # observe command family and private bridge dispatch

internal/
├── app/                                       # runtime composition and observation-secret lifecycle
├── diagnostics/                               # provider-specific availability/setup/trust/capability checks
├── events/                                    # canonical reducer, agent snapshot state, deduplication, replay
├── project/                                   # Git-root identity and canonical path checks
├── protocol/                                  # provider-neutral agent/capability/wire types
├── providers/
│   ├── registry.go                             # private adapter lookup and shared setup/normalize boundary
│   ├── claude/                                # hook settings, bridge projection, detection, fixtures
│   ├── codex/                                 # hook settings/trust, bridge projection, detection, fixtures
│   └── opencode/                              # versioned plugin generation, projection, detection, fixtures
├── security/                                  # host/origin/token/path validators
├── server/                                    # generic allowlisted provider ingestion route
└── storage/                                   # observation secret, fingerprints, snapshot/event persistence

web/src/
├── dashboard/                                 # semantic agent list and capability presentation
└── protocol/                                  # agent-aware snapshot/event representations

tests/fixtures/
├── claude/
├── codex/
└── opencode/                                  # synthetic fixture inputs only
```

**Structure Decision**: Add a small private provider registry rather than expose an SDK or duplicate runtime paths. Each adapter owns detection, preview/apply/remove setup, raw-input projection, and normalization; it cannot mutate reducer state directly. The generic runtime owns authorization, validation, safe persistence, ordering, and delivery. A public provider API remains deferred until multiple adapters prove the boundary.

## Implementation Approach

1. Define a private `Adapter` boundary with provider ID, detection/capability report, setup preview/apply/remove, bridge projection, and normalized-update functions. The registry allowlists `claude`, `codex`, and `opencode`; it has no generic raw-provider passthrough.
2. Add provider-neutral agent, tool activity, observation setup/status, and capability projections. Extend the existing snapshot with optional `agents` and `providerStatuses`, and emit source-tagged additive `agent.updated` events. Existing clients safely ignore additions.
3. Extend the reducer—not a new store—to hydrate/persist agent state and a bounded fingerprint cache. One atomic operation validates a normalized update, deduplicates, applies terminal/inferred-state guards, appends a redacted event, saves the snapshot, then broadcasts. On restart, formerly active agents are `unknown`/stale until renewed evidence.
4. Add a generic `POST /api/v1/observe/{provider}` route. It verifies an allowlisted path, exact Host, and an Origin that is either absent for the private local bridge or exactly the runtime origin; it rejects every other Origin. It also validates body bound, envelope version, per-project observation token, project identity, and canonical path before adapter normalization. The separate observation token is never placed in settings, a URL, browser state, events, logs, or diagnostics.
5. Add private `craftlane observe bridge --provider <id>` handling. It consumes bounded standard input, projects/discards data at the provider boundary, reads the owner-only observation secret, makes a short loopback request, produces no model-facing output, and cannot execute arbitrary supplied commands.
6. Claude adapter: detect local executable/version and setup/capability state for CLI diagnostics, preview/merge only marked entries in `.claude/settings.local.json`, use fixed command hooks for supported lifecycle signals, and remove only owned entries. Its bridge discards prompts, tool input/output, assistant/permission text, transcript paths, commands, and environment values. Unsupported hook categories remain visible as unavailable.
7. Codex adapter: detect local executable/version, hook trust, and setup/capability state for CLI diagnostics; preview/merge only marked project hooks in `.codex/hooks.json` (or its documented project configuration location), warn that the artifact can be Git-tracked, and never alter `.gitignore`. Its lifecycle normalizer distinguishes turn completion from delayed/non-subagent session end and marks hosted/unhooked tools as partial coverage.
8. OpenCode adapter: detect local executable/version, supported plugin API, and setup/capability state for CLI diagnostics; after confirmation, write only a uniquely named generated project plugin under `.opencode/plugins/`, warn it may be Git-tracked, and remove only that exact owned file. Generate a pinned compatible plugin or report unavailable. The plugin sends only safe lifecycle/tool/permission/session-lineage metadata to the bridge; it never changes permission policy, reads shell environment, or forwards raw args/results/messages.
9. Expand doctor/status to show each provider's executable/version, setup state, trust/load state, callback self-test, project-artifact state, and available/unavailable capabilities with recovery actions. A missing/degraded provider warns only; it never disables the runtime or the other adapters.
10. Extend the dashboard to seed/upsert generic agents from snapshot/events and render accessible provider badges, explicit/inferred state source, safe tool summary/outcome, parent relation when supplied, prior-history marker, and per-provider coverage. Retain the existing reconnect/event-ID behaviour.
11. Add fixture, reducer, callback, setup/reversal, contract, integration, reconnect, redaction, dashboard, and optional installed-provider smoke coverage. Use fake IDs/path placeholders/secret canaries only; no real transcript or account data enters the repository.

## Complexity Tracking

No constitution violation requires justification.
