# Implementation Plan: Local Runtime Foundation

**Branch**: `001-foundation-runtime` | **Date**: 2026-10-01 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-foundation-runtime/spec.md`

## Summary

Deliver a single native Craftlane executable that detects the Git project at the current directory, runs a loopback-only project runtime, and serves a minimal dashboard. The runtime owns canonical project state, persists a bounded redacted event history, and exposes a versioned HTTP and WebSocket contract. It uses a thin Go command dispatcher, Go standard-library HTTP server plus a small WebSocket library, JSON snapshot/JSONL storage, and a React/TypeScript web shell built by contributors and embedded in the release binary.

## Technical Context

**Language/Version**: Go 1.25+ (pinned in `go.mod`); TypeScript 5.x for the contributor-built web shell

**Primary Dependencies**: Go standard library (`context`, `net/http`, `os/exec`, `os/signal`, `embed`); `github.com/coder/websocket`; React, Vite, and native browser Fetch/WebSocket

**Storage**: Per-project atomic JSON snapshot/runtime records and bounded append-only JSONL event segments; default retention is 10,000 events per project and configuration cannot lower it below 1,000 events; no database in this milestone

**Testing**: Go `testing`, `httptest`, temporary Git repositories, Go race detector; Vitest + React Testing Library; Playwright against the built executable

**Target Platform**: Local developer machines on macOS, Linux, and Windows; modern desktop browsers

**Project Type**: Native CLI and local web service with an embedded browser dashboard

**Performance Goals**: Project dashboard ready within 10 seconds in normal local conditions; ordered reconnect recovery for retained events; slow clients cannot block canonical state updates

**Constraints**: IPv4 loopback-only default; no Craftlane account/cloud/telemetry or required network access; explicit versioned contract; local paths constrained to approved roots; no generic shell endpoint; credentials redacted before persistence or presentation; every Go test/build target first produces embedded web assets; release binary has no Node.js/pnpm runtime dependency

**Scale/Scope**: One Git project per runtime, minimal dashboard and lifecycle/diagnostic commands, retained event recovery for normal development sessions; provider observability, Spec Kit parsing, task/worktree actions, dispatch, and 3D are deferred

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Pre-design review — PASS**

| Constitution obligation | Design response | Status |
|---|---|---|
| Productivity-first, keyboard-equivalent control | All delivered lifecycle operations are CLI commands; the dashboard is an optional status surface. | Pass |
| Local-first, account-free ownership | Runtime listens only on `127.0.0.1`; all state, history, and diagnostics remain local; no account, telemetry, or hosted service is designed. | Pass |
| Safe, explicit, reversible actions | Public routes are typed; no raw shell route exists. A runtime-control action requires a per-runtime secret, and only initiates graceful shutdown. Host/origin/path validation and pre-persistence redaction are mandatory. | Pass |
| Provider-agnostic core and versioned protocols | Project and event domain types remain provider-neutral; protocol schema version and additive evolution rules are defined before provider adapters. Core reducer is canonical state owner. | Pass |
| Spec-driven bounded delivery | The plan implements only `001-foundation-runtime`; later roadmap milestones are excluded. | Pass |
| Accessible 2D equivalence | This milestone provides a semantic HTML dashboard with labelled health and project state. No spatial interface is introduced. | Pass |
| Maintainable, testable open-source delivery | Domain packages are explicit, dependencies are minimal, and unit/integration/contract/browser tests cover state, security, embedding, and lifecycle behavior. | Pass |

**Post-design review — PASS**: The data model, contracts, quickstart, and source layout retain the same controls. No constitution exception or complexity justification is required.

## Project Structure

### Documentation (this feature)

```text
specs/001-foundation-runtime/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── cli.md
│   ├── http-api.yaml
│   └── websocket.md
└── tasks.md                 # Created later by $speckit-tasks
```

### Source Code (repository root)

```text
cmd/
└── craftlane/
    └── main.go                         # Command dispatch and process entry point

internal/
├── app/                                # Runtime composition and ownership
├── diagnostics/                        # Typed readiness checks and doctor output
├── events/                             # Reducer, replay, client hub, redaction boundary
├── project/                            # Canonical Git project detection and identity
├── protocol/                           # Go wire types and schema validation boundary
├── security/                           # Host/origin/token/path validation
├── server/
│   ├── assets/                         # Built frontend output embedded by the server package
│   └── server.go                       # HTTP, WebSocket, health, and SPA routes
└── storage/                            # Snapshot, runtime-control, and JSONL implementations

protocol/
├── schema/                             # Versioned JSON schemas; serialized contract source
└── README.md

web/
├── src/
│   ├── app/                            # App bootstrap and connection lifecycle
│   ├── dashboard/                      # Semantic project and health dashboard
│   ├── protocol/                       # Generated/validated wire representations
│   └── shared/                         # Small reusable presentational primitives
├── public/
├── package.json
├── pnpm-lock.yaml
├── tsconfig.json
└── vite.config.ts

tests/
├── fixtures/                           # Git and protocol fixtures
├── integration/                        # Executable/runtime lifecycle tests
└── e2e/                                # Playwright built-binary dashboard tests
```

**Structure Decision**: Use one Go module plus a separately built `web/` contributor application. The web build writes generated static assets beneath `internal/server/assets/`, allowing the Go server package to embed them without path traversal. Setup and Make targets must perform the web build before any Go test or build that compiles the embed package, preventing an absent-asset compile failure. `internal/` keeps implementation private until a stable provider SDK is intentionally introduced in a later milestone.
