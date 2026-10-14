# Tasks: Multi-Provider Agent Observability

**Input**: Design documents from `specs/002-multi-provider-observability/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `quickstart.md`, and `contracts/`

**Tests**: Required. The specification and constitution require automated coverage for provider fixtures, state, storage, protocol, callback security, setup, and browser behavior.

**Organization**: Tasks are grouped by delivery increment. Provider setup intentionally precedes P1 because it is required to demonstrate the real live-observation user flow rather than only inject fixtures.

## Phase 1: Setup (Shared Infrastructure)

- [X] T001 Create private provider package and synthetic fixture layout in `internal/providers/claude/`, `internal/providers/codex/`, `internal/providers/opencode/`, `internal/providers/registry.go`, `tests/fixtures/claude/`, `tests/fixtures/codex/`, and `tests/fixtures/opencode/`.
- [X] T002 [P] Add fixture manifests with synthetic IDs, `{{PROJECT_ROOT}}` placeholders, and credential-shaped canaries in `tests/fixtures/claude/manifest.json`, `tests/fixtures/codex/manifest.json`, and `tests/fixtures/opencode/manifest.json`.
- [X] T003 [P] Document the adapter boundary and local-only data policy in `docs/architecture/runtime.md` and `docs/security/local-runtime.md`.

---

## Phase 2: Foundational (Blocking Prerequisites)

**⚠️ CRITICAL**: Complete this phase before provider or dashboard story work.

- [X] T004 Add provider-neutral agent, tool activity, provider status/capability, safe observation envelope, and additive snapshot types in `internal/protocol/types.go`, `protocol/schema/events.schema.json`, and `web/src/protocol/events.ts`.
- [X] T005 [P] Write protocol compatibility tests for additive agent/provider-status snapshots and source-tagged `agent.updated` events in `internal/protocol/events_test.go` and `tests/integration/protocol_contract_test.go`.
- [X] T006 Add owner-only observation-secret and bounded fingerprint persistence helpers in `internal/storage/observation.go` and `internal/storage/observation_test.go`.
- [X] T007 Define the private adapter registry, allowlisted provider IDs, setup/result contracts, and normalized-update boundary in `internal/providers/registry.go` and `internal/providers/registry_test.go`.
- [X] T008 Write reducer tests for agent hydration, duplicate suppression, terminal-state guards, explicit-over-inferred precedence, and stale-after-restart recovery in `internal/events/observation_test.go`.
- [X] T009 Extend the canonical reducer to atomically apply normalized observations, persist agent snapshots/fingerprints, emit `agent.updated`, and recover safe agent state in `internal/events/events.go` and `internal/events/reducer.go`.
- [X] T010 Write callback contract/security tests for allowlisted provider paths, exact Host, absent-or-exact Origin policy, method/body limits, token/project/path validation, malformed input, and no-mutation failures in `internal/server/observation_handlers_test.go`.
- [X] T011 Add the generic authenticated `POST /api/v1/observe/{provider}` callback route, exact-Origin guard, and runtime observation-secret lifecycle in `internal/server/observation_handlers.go`, `internal/server/server.go`, and `internal/app/runtime.go`.
- [X] T012 [P] Add bridge input-bound/projection/no-raw-output tests in `internal/providers/bridge_test.go`.
- [X] T013 Implement the private fixed `craftlane observe bridge --provider <provider>` command and bounded local callback client in `cmd/craftlane/main.go`, `internal/app/observe.go`, and `internal/providers/bridge.go`.
- [X] T014 Update v1 protocol documentation and browser snapshot/event type handling for optional agents/provider statuses and callback-Origin semantics in `protocol/README.md`, `specs/002-multi-provider-observability/contracts/observability.md`, and `web/src/protocol/client.ts`.
- [X] T015 Run foundational reducer, storage, protocol, callback, and bridge suites in `internal/events/observation_test.go`, `internal/storage/observation_test.go`, `internal/server/observation_handlers_test.go`, and `tests/integration/protocol_contract_test.go`.

**Checkpoint**: A safe, durable, provider-neutral ingestion path exists and enforces the constitution’s Host/Origin requirements.

---

## Phase 3: User Story 3 - Set Up and Diagnose Provider Observation Safely (Priority: P3, prerequisite for live MVP)

**Goal**: Enable a developer to configure a provider explicitly and safely before attempting the real P1 live-observation flow.

**Independent Test**: In temporary projects with absent, available, unsupported, untrusted, malformed, and preconfigured providers, confirm preview is non-mutating; confirmed setup changes only owned artifacts; removal preserves unrelated data; diagnostics offer recovery guidance.

- [X] T016 [P] [US3] Write Claude local-settings preview/merge/idempotency/remove tests with unrelated-hook preservation in `internal/providers/claude/setup_test.go`.
- [X] T017 [P] [US3] Write Codex project-hook preview/merge/trust/malformed-config/remove tests with tracking-warning coverage in `internal/providers/codex/setup_test.go`.
- [X] T018 [P] [US3] Write OpenCode generated-plugin version/preview/ownership/remove tests with tracking-warning coverage in `internal/providers/opencode/setup_test.go`.
- [X] T019 [P] [US3] Write CLI preview/confirm/status/remove and provider-isolation tests in `cmd/craftlane/main_test.go` and `tests/integration/observe_cli_test.go`.
- [X] T020 [P] [US3] Write provider diagnostic availability/setup/trust/callback/capability/action tests in `internal/diagnostics/observation_test.go`.
- [X] T021 [US3] Implement Claude executable/version, setup, and capability detection plus `.claude/settings.local.json` structural preview/apply/remove with exact owned-hook markers in `internal/providers/claude/adapter.go` and `internal/providers/claude/setup.go`.
- [X] T022 [US3] Implement Codex executable/version, trust, setup, and capability detection plus project-hook structural preview/apply/remove and tracked-artifact warning in `internal/providers/codex/adapter.go` and `internal/providers/codex/setup.go`.
- [X] T023 [US3] Implement OpenCode executable/version, plugin-API, setup, and capability detection plus generated-plugin preview/apply/remove with exact owned-file fingerprinting and tracked-artifact warning in `internal/providers/opencode/adapter.go` and `internal/providers/opencode/setup.go`.
- [X] T024 [US3] Add `observe setup`, `status`, and `remove` command routing with required `--confirm` and safe provider validation in `cmd/craftlane/main.go` and `internal/app/observe.go`.
- [X] T025 [US3] Implement independent provider diagnostics and doctor/status rendering from the prerequisite adapter detections for executable, setup, trust/load, callback, project path, and capability recovery actions in `internal/diagnostics/detectors.go`, `internal/diagnostics/diagnostics.go`, and `internal/app/commands.go`.
- [X] T026 [US3] Publish provider-status changes to the canonical reducer and dashboard through `internal/app/runtime.go`, `internal/events/events.go`, and `web/src/dashboard/Dashboard.tsx`.
- [X] T027 [US3] Add end-to-end setup/diagnostic/removal tests using temporary repositories and optional installed-provider smoke tests in `tests/integration/observe_setup_test.go` and `tests/e2e/provider-observation.spec.ts`.

**Checkpoint**: A developer can explicitly configure at least one provider and see its genuine setup/trust/capability state before live observation.

---

## Phase 4: User Story 1 - See Active Agent Work (Priority: P1) 🎯 MVP

**Goal**: Show live Claude Code, Codex, and OpenCode session/subagent/tool activity using one normalized dashboard model with honest capability disclosure.

**Independent Test**: Configure each available provider with Phase 3, start a supported session, then verify a live dashboard update for session/tool activity and supplied parent relation without refresh.

- [X] T028 [P] [US1] Write sanitized Claude hook projection/normalization fixture tests for sessions, subagents, and tool start/completion/failure in `internal/providers/claude/adapter_test.go` and `internal/providers/claude/testdata/`.
- [X] T029 [P] [US1] Write sanitized Codex hook projection/normalization fixture tests for sessions, subagents, tool activity, and partial local coverage in `internal/providers/codex/adapter_test.go` and `internal/providers/codex/testdata/`.
- [X] T030 [P] [US1] Write sanitized OpenCode plugin projection/normalization fixture tests for session status, parent lineage, and tool hooks in `internal/providers/opencode/adapter_test.go` and `internal/providers/opencode/testdata/`.
- [X] T031 [P] [US1] Write dashboard tests for snapshot seeding, `agent.updated` upsert, semantic provider/state/tool/parent rendering, and unavailable-capability display in `web/src/dashboard/Dashboard.test.tsx`.
- [X] T032 [US1] Implement Claude allowlisted lifecycle-hook projection/normalization using prerequisite provider detection in `internal/providers/claude/adapter.go` and `internal/providers/claude/hooks.go`.
- [X] T033 [US1] Implement Codex allowlisted lifecycle-hook projection/normalization using prerequisite provider detection in `internal/providers/codex/adapter.go` and `internal/providers/codex/hooks.go`.
- [X] T034 [US1] Implement OpenCode allowlisted session/tool/lineage projection/normalization using prerequisite provider detection in `internal/providers/opencode/adapter.go` and `internal/providers/opencode/plugin.go`.
- [X] T035 [US1] Register the three adapters, route bridge envelopes through their normalizers, and preserve provider source/version metadata in `internal/providers/registry.go` and `internal/app/observe.go`.
- [X] T036 [US1] Render accessible generic agent rows and provider capability status from snapshots/events in `web/src/dashboard/Dashboard.tsx` and `web/src/dashboard/dashboard.css`.
- [X] T037 [US1] Add per-provider fixture-driven callback/reconnect integration coverage for ordered live agent updates and parent relationships in `tests/integration/provider_observation_test.go` and `tests/integration/websocket_test.go`.

**Checkpoint**: A configured provider can show real safe agent work in the dashboard; the same view works for all three adapters without provider-specific UI logic.

---

## Phase 5: User Story 2 - Respond to Permission and Completion States (Priority: P2)

**Goal**: Present trustworthy permission, idle, failure, and terminal state while distinguishing reported evidence from inference and provider limits.

**Independent Test**: Feed each adapter permission, inactivity, turn/session/subagent ending, failure, duplicate, and reversed-order fixtures; confirm state/source correctness and no terminal agent revival.

- [X] T038 [P] [US2] Add table-driven cross-provider state-transition tests for reported permission, waiting, idle, terminal, and unavailable-terminal semantics in `internal/events/agent_state_test.go`.
- [X] T039 [P] [US2] Add provider fixture tests for Claude terminal evidence, Codex delayed/non-subagent session end, and OpenCode busy/idle/retry semantics in `internal/providers/claude/state_test.go`, `internal/providers/codex/state_test.go`, and `internal/providers/opencode/state_test.go`.
- [X] T040 [P] [US2] Add dashboard accessibility tests for permission attention, reported/inferred source, prior-history-unavailable, stale-after-restart, and non-colour-only states in `web/src/dashboard/Dashboard.test.tsx`.
- [X] T041 [US2] Implement the normalized agent-state transition table, configurable inferred-idle timer, and terminal/reordered-event guards in `internal/events/agent_state.go` and `internal/events/events.go`.
- [X] T042 [US2] Map provider-specific permission, failure, turn/session/subagent end, and idle signals to the shared state model and capability reasons in `internal/providers/claude/adapter.go`, `internal/providers/codex/adapter.go`, and `internal/providers/opencode/adapter.go`.
- [X] T043 [US2] Add attention-oriented state/source/history/staleness labels to the generic dashboard agent view in `web/src/dashboard/Dashboard.tsx` and `web/src/dashboard/dashboard.css`.
- [X] T044 [US2] Add state-focused integration coverage for permission updates, 100 reconnect cycles, restart recovery, duplicate/reversed events, and no false terminal parity in `tests/integration/agent_state_test.go` and `tests/integration/websocket_test.go`.

**Checkpoint**: Developers can distinguish requested permission, inferred idle, provider-reported failure, and confirmed terminal outcome without assuming unsupported parity.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T045 [P] Add redaction sink regressions scanning provider fixtures through JSONL, snapshots, diagnostics, HTTP, WebSocket, and browser output in `tests/integration/redaction_sinks_test.go`.
- [X] T046 [P] Update user/developer documentation for three-provider setup, capability differences, local artifacts, trust warnings, callback Origin policy, and removal in `README.md`, `docs/architecture/runtime.md`, and `docs/security/local-runtime.md`.
- [X] T047 [P] Add CI coverage for synthetic adapter fixtures, Go race/integration suites, web unit tests, and built-binary provider-observation smoke skip reporting in `.github/workflows/ci.yml`.
- [X] T048 Validate HTTP/WebSocket schema/contracts for additive agents/provider statuses and the callback Origin policy in `protocol/schema/events.schema.json`, `protocol/README.md`, and `specs/002-multi-provider-observability/contracts/observability.md`.
- [ ] T049 Run quickstart and success-criteria validation per provider, including 20-session latency exercises, 100 reconnect cycles, Origin/authentication negatives, and five-developer usability study; record outcomes/remediation in `specs/002-multi-provider-observability/quickstart.md` and `specs/002-multi-provider-observability/checklists/requirements.md`.

---

## Dependencies & Execution Order

- **Phase 1**: T001 starts immediately; T002–T003 can run in parallel after path creation.
- **Phase 2**: T004–T007 establish common contracts; T005, T006, and T008 follow T004. T009 depends on T004, T006, T008; T010 depends on T004, T007; T011 depends on T006, T007, T010; T012 follows T007; T013 depends on T011–T012; T014 depends on T004, T009; T015 validates the shared path.
- **Phase 3**: T016–T020 run in parallel; T021–T023 depend on matching tests and establish provider detection; T024 depends on T021–T023 and T013; T025 depends on T020–T023; T026 depends on T025/T009; T027 validates T024–T026.
- **Phase 4**: T028–T031 run in parallel; T032–T034 depend on matching tests and the corresponding prerequisite detection; T035 depends on T032–T034/T013; T036 depends on T014/T026; T037 validates T035–T036.
- **Phase 5**: T038–T040 run in parallel; T041 depends on T038; T042 depends on T039/T041; T043 depends on T040/T036; T044 validates T041–T043.
- **Phase 6**: T045–T047 can run after their relevant implementation; T048 depends on T014/T037; T049 is the final acceptance gate after T044–T048.

### Parallel Opportunities

- T002/T003; T005/T006/T008/T012; T016–T020; T021–T023; T028–T031; T032–T034; T038–T040; T045–T047.

## Implementation Strategy

1. Build and secure the shared ingestion pipeline (T001–T015).
2. Deliver explicit provider setup/diagnostics (T016–T027) so real sessions can be observed.
3. Deliver live multi-provider visibility (T028–T037), then attention/state semantics (T038–T044).
4. Run contract, security, performance, and usability acceptance validation (T045–T049).

## Task Summary

| Scope | Tasks |
|---|---:|
| Setup | 3 |
| Foundational | 12 |
| US3 setup prerequisite | 12 |
| US1 — P1 MVP | 10 |
| US2 — P2 | 7 |
| Polish | 5 |
| **Total** | **49** |

## Notes

- Every task uses the required checkbox, sequential ID, optional parallel marker, user-story label, and exact paths.
- Provider integrations are read-only: no task adds agent control, permission mutation, dispatch, worktree automation, cloud telemetry, or transcript parsing.

## Phase 7: Convergence

- [ ] T050 CRITICAL Sanitize or reject secret-bearing projected observation fields before reducer state, snapshot persistence, events, diagnostics, HTTP, WebSocket, and browser delivery per FR-015 and Constitution III (contradicts).
- [ ] T051 Replace generic provider setup markers and pass-through projection with version-gated Claude Code, Codex, and OpenCode lifecycle integrations that produce only allowlisted normalized envelopes; add sanitized fixtures for every supported category per FR-001, FR-005, FR-020, and FR-022 (missing).
- [ ] T052 Bind observation callbacks to the active runtime's canonical project identity/path and add reported-time ordering guards so malformed, wrong-project, and out-of-order records cannot alter valid agent state per FR-004 and FR-013 (partial).
- [ ] T053 Track confirmed-live session evidence and restrict configurable inferred-idle transitions to those sessions, preserving explicit state precedence and terminal guards per FR-010 and FR-012 (partial).
- [ ] T054 Publish provider-status updates through the versioned local event stream and render per-capability availability plus an accessible permission-attention indicator in the dashboard per FR-014, US1/AC4, and US2/AC1 (partial).
- [ ] T055 Expand provider diagnostics to independently check executable/version, lifecycle API, setup, trust/load, callback readiness, canonical project path, and each capability with a recovery action, while leaving unrelated providers usable per FR-018, FR-019, and FR-022 (partial).
- [ ] T056 Complete pending T049 acceptance validation, correcting unsupported quickstart claims and recording per-provider latency, replay, security, and usability evidence or remediation per SC-001, SC-002, SC-003, SC-004, SC-005, SC-006, SC-007, and SC-008 (missing).
