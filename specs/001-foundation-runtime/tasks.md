# Tasks: Local Runtime Foundation

**Input**: Design documents from specs/001-foundation-runtime/

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md, and contracts/

**Tests**: Required. The feature specification and Craftlane constitution require proportionate automated coverage for lifecycle, storage, protocol, security, browser behavior, and release packaging.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Establish the Go module, contributor web build, schemas, and reproducible developer commands.

- [X] T001 Create the Go module and initial runtime package layout in go.mod, cmd/craftlane/main.go, and internal/{app,diagnostics,events,project,protocol,security,server,storage}/.
- [X] T002 Initialize the React/TypeScript/Vite contributor application in web/package.json, web/tsconfig.json, web/vite.config.ts, and web/src/app/main.tsx; perform the initial web build into internal/server/assets/dist/ before any Go package embeds those assets.
- [X] T003 [P] Add repository ignore rules for generated web assets and local Craftlane state in .gitignore.
- [X] T004 [P] Add contributor build, test, lint, and embedded-asset targets in Makefile so every Go test/build target first runs the web build and fails if internal/server/assets/dist/index.html is absent.
- [X] T005 [P] Create the versioned protocol documentation and initial JSON Schema files in protocol/README.md, protocol/schema/common.schema.json, and protocol/schema/events.schema.json.
- [X] T006 Add pinned web toolchain metadata and baseline lint/test configuration in web/package.json, web/pnpm-lock.yaml, web/vitest.config.ts, and web/eslint.config.js after T002 creates the web project.
- [X] T007 [P] Create reusable temporary-Git-project and test-server helpers in tests/fixtures/gitrepo.go and tests/integration/runtime_test.go.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Implement provider-neutral canonical data, safe persistence, security boundaries, and runtime composition required by every story.

**CRITICAL**: Complete this phase before beginning user-story work.

- [X] T008 Define provider-neutral Project, RuntimeWorkspace, RuntimeRecord, DiagnosticCheck, Snapshot, and Event wire/domain types in internal/protocol/types.go and internal/project/project.go from data-model.md.
- [X] T009 Implement canonical Git-root detection and stable project identity using controlled cancellable Git commands in internal/project/detector.go and internal/project/detector_test.go.
- [X] T010 Implement atomic per-project runtime-record and snapshot storage with owner-only runtime-control storage in internal/storage/files.go, internal/storage/storage.go, and internal/storage/files_test.go.
- [X] T011 [P] Implement recursive copy-on-write secret redaction before every durable or browser-visible sink in internal/events/redact.go and internal/events/redact_test.go.
- [X] T012 [P] Implement exact loopback Host, WebSocket Origin, runtime-token, and approved-path validators in internal/security/validation.go and internal/security/validation_test.go.
- [X] T013 Implement the base HTTP router, embedded-asset serving, SPA fallback, and shared security middleware in internal/server/server.go and internal/server/server_test.go.
- [X] T014 Implement root-context ownership, bounded graceful shutdown, listener creation, runtime-record recovery, and dependency wiring in internal/app/runtime.go and internal/app/runtime_test.go.
- [X] T015 Add foundational integration coverage for non-Git rejection, atomic stale-record recovery, loopback-only binding, invalid Host rejection, invalid path rejection, redaction persistence, and hermetic no-egress startup in tests/integration/foundation_test.go.

**Checkpoint**: The core can safely identify one Git project, own a local runtime, serve embedded static assets, and reject invalid security inputs.

---

## Phase 3: User Story 1 - Open a Local Project Workspace (Priority: P1) MVP

**Goal**: A developer in a Git repository runs Craftlane and reaches a local dashboard that identifies the project and its health.

**Independent Test**: Build the executable, invoke it in a temporary Git repository with start --no-open, open the reported loopback URL, and verify the project name, root, and health are visible.

### Tests for User Story 1

- [X] T016 [P] [US1] Write CLI/start integration tests for default invocation, start --no-open, start --headless, explicit-port conflict handling, and loopback URL reporting; assert dashboard readiness within 10 seconds across 20 clean starts with at least 19 successes in tests/integration/start_test.go.
- [X] T017 [P] [US1] Write HTTP contract tests for GET /api/v1/health and GET /api/v1/project in tests/integration/project_api_test.go using contracts/http-api.yaml.
- [X] T018 [P] [US1] Write dashboard component tests for labelled project, root, and health rendering in web/src/dashboard/Dashboard.test.tsx.
- [X] T019 [P] [US1] Add Playwright built-binary test for a temporary Git project dashboard in tests/e2e/project-dashboard.spec.ts and playwright.config.ts.

### Implementation for User Story 1

- [X] T020 [US1] Implement command dispatch, default-to-start behavior, start options, and internal serve process launch in cmd/craftlane/main.go and internal/app/start.go.
- [X] T021 [US1] Implement safe GET /api/v1/health and GET /api/v1/project projections in internal/server/project_handlers.go against contracts/http-api.yaml.
- [X] T022 [US1] Implement the semantic dashboard shell, runtime API client, visible connection state, and accessible project/health presentation in web/src/app/App.tsx, web/src/dashboard/Dashboard.tsx, web/src/dashboard/dashboard.css, and web/src/protocol/client.ts.
- [X] T023 [US1] Implement best-effort browser opening that never changes successful startup status in internal/app/browser.go and internal/app/browser_test.go.
- [X] T024 [US1] Wire the production web build into Go embedding and verify static asset and SPA route behavior in internal/server/assets.go, internal/server/server.go, and internal/server/assets_test.go.
- [X] T025 [US1] Run and fix the User Story 1 Go, web, and Playwright test suites defined in tests/integration/start_test.go, web/src/dashboard/Dashboard.test.tsx, and tests/e2e/project-dashboard.spec.ts.

**Checkpoint**: craftlane delivers the MVP local project workspace from a released-style executable without runtime web-development tooling.

---

## Phase 4: User Story 2 - Manage the Local Workspace from the Command Line (Priority: P2)

**Goal**: A developer can inspect, open, diagnose, stop, and restart the local workspace through typed CLI actions.

**Independent Test**: Start without opening a browser, check status and doctor, open on request, restart, stop, then verify status reports no live matching runtime.

### Tests for User Story 2

- [X] T026 [P] [US2] Write CLI lifecycle integration tests for status, open, stop, and restart, including stale and stopped runtime results, in tests/integration/lifecycle_test.go.
- [X] T027 [P] [US2] Write diagnostics unit and integration tests for required Git/runtime/dashboard checks and optional Spec Kit/provider warnings in internal/diagnostics/diagnostics_test.go and tests/integration/doctor_test.go.
- [X] T028 [P] [US2] Write HTTP contract tests that reject missing, invalid, cross-origin, and secrets invalidated by shutdown, interruption, stale recovery, or restart in tests/integration/runtime_control_test.go.

### Implementation for User Story 2

- [X] T029 [P] [US2] Implement typed diagnostic checks, safe detail/action output, and CLI/health result rendering in internal/diagnostics/diagnostics.go and internal/diagnostics/render.go.
- [X] T030 [US2] Implement per-runtime control-secret lifecycle and authenticated graceful POST /api/v1/runtime/runtime/stop handler in internal/server/server.go and internal/app/commands.go.
- [X] T031 [US2] Implement status, open, stop, restart, and doctor command handlers with meaningful exit codes in internal/app/commands.go and cmd/craftlane/main.go.
- [X] T032 [P] [US2] Add optional Spec Kit and provider executable detection that degrades to warnings without blocking startup in internal/diagnostics/detectors.go.
- [X] T033 [US2] Run and fix the User Story 2 lifecycle, doctor, and runtime-control suites in tests/integration/lifecycle_test.go, tests/integration/doctor_test.go, and tests/integration/runtime_control_test.go.

**Checkpoint**: Every required lifecycle command works without the dashboard, reports a corrective action on failure, and cannot leak control credentials.

---

## Phase 5: User Story 3 - Receive Safe, Ordered Local Updates (Priority: P3)

**Goal**: A connected dashboard receives ordered normalized project updates, reconnects without missed retained updates, and resynchronizes explicitly when history has expired.

**Independent Test**: Connect a test client, record its sequence, disconnect, emit project updates, reconnect with the cursor, and verify exactly-once ordered recovery; repeat with an expired cursor and verify resync_required followed by a snapshot.

### Tests for User Story 3

- [X] T034 [P] [US3] Write protocol-schema and event-envelope compatibility tests for required fields and unknown optional fields in internal/protocol/events_test.go and tests/integration/protocol_contract_test.go.
- [X] T035 [P] [US3] Write reducer/storage tests for monotonic sequence assignment, redaction-before-append, the default 10,000-event and minimum 1,000-event retention limits, and snapshot consistency in internal/events/reducer_test.go and internal/storage/events_test.go.
- [X] T036 [P] [US3] Write WebSocket integration tests for subprotocol negotiation, exact Origin rejection, ordered replay, deduplication, resync, bootstrap races, slow-client isolation, and 100 reconnect cycles with no missed or duplicate retained events in tests/integration/websocket_test.go.
- [X] T037 [P] [US3] Write browser client tests for cursor persistence, duplicate suppression, reconnect, and resync handling in web/src/protocol/events.test.ts.

### Implementation for User Story 3

- [X] T038 [US3] Implement JSONL event append/read retention and snapshot cursor persistence in internal/storage/events.go and internal/storage/events_test.go.
- [X] T039 [US3] Implement the single-owner event reducer, normalized core project lifecycle events, and bounded client hub in internal/events/reducer.go, internal/events/hub.go, and internal/events/events.go.
- [X] T040 [US3] Implement GET /api/v1/events and the craftlane.events.v1 WebSocket handshake, replay, heartbeat, resync, and slow-client behavior in internal/server/events_handlers.go and internal/server/websocket.go.
- [X] T041 [US3] Implement dashboard live-event connection, snapshot application, cursor tracking, reconnection backoff, and degraded-state UI in web/src/protocol/events.ts, web/src/protocol/client.ts, and web/src/dashboard/Dashboard.tsx.
- [X] T042 [US3] Run and fix the User Story 3 protocol, reducer, storage, WebSocket, and browser-client suites in tests/integration/protocol_contract_test.go, tests/integration/websocket_test.go, and web/src/protocol/events.test.ts.

**Checkpoint**: The dashboard continues to present canonical runtime state through the versioned event contract, while disconnected or slow clients cannot affect runtime correctness.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Complete release-quality validation, documentation, and cross-platform safeguards without broadening milestone scope.

- [X] T043 [P] Add security regression coverage proving secrets never reach implemented JSONL, snapshots, HTTP/WS responses, or diagnostics in tests/integration/redaction_sinks_test.go.
- [X] T044 [P] Add contributor and user documentation for local-first behavior, CLI commands, security boundaries, and build/run prerequisites in README.md, docs/architecture/runtime.md, and docs/security/local-runtime.md.
- [X] T045 [P] Add CI jobs for Go formatting/vet/unit/race/integration tests, web lint/typecheck/unit/build, embedded-asset verification, and Playwright in .github/workflows/ci.yml.
- [X] T046 Add clean-install and built-binary smoke coverage for macOS, Linux, and Windows in .github/workflows/release-smoke.yml and tests/integration/binary_smoke_test.go.
- [ ] T047 Run the complete quickstart validation, contract checks, success-criteria audit, and a five-developer task-completion study; record timing, 20-start, 100-reconnect, offline-startup, and usability results plus any remediation in specs/001-foundation-runtime/quickstart.md and specs/001-foundation-runtime/checklists/requirements.md.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 — Setup**: Can start immediately. T001 and T002 unblock most implementation; T003–T005 and T007 may run in parallel once target directories exist. T006 depends on T002.
- **Phase 2 — Foundational**: Depends on T001–T002. T008–T010 establish canonical data, project detection, and persistence; T011–T012 can run in parallel; T013 depends on T008 and T011–T012; T014 depends on T009–T013; T015 validates the complete foundation.
- **Phase 3 — US1 (MVP)**: Depends on T015. Tests T016–T019 can run in parallel. T020–T024 follow the shared foundation; T025 is the MVP validation gate.
- **Phase 4 — US2**: Depends on T025 and the shared runtime. Tests T026–T028 can run in parallel. T029 and T032 may proceed in parallel; T030 depends on T010 and T012; T031 integrates them; T033 validates the story.
- **Phase 5 — US3**: Depends on T025 and T015. Tests T034–T037 can run in parallel. T038 precedes T039; T040 depends on T038–T039 and T012; T041 depends on T040; T042 validates the story.
- **Phase 6 — Polish**: T043–T045 can start after the covered functionality. T046 depends on the production build; T047 is the final feature gate.

### User Story Dependencies

- **US1 (P1)**: Requires only the foundational phase. It is the suggested MVP.
- **US2 (P2)**: Builds on the live runtime established by US1 but is independently testable entirely through the CLI.
- **US3 (P3)**: Builds on the runtime and dashboard from US1; its protocol, replay, and browser client validation are independently testable with fixture events.

### Parallel Opportunities

- Setup: T003, T004, T005, and T007 after directory scaffolding.
- Foundational: T011 and T012 after core types are available.
- US1 tests: T016, T017, T018, and T019.
- US2 tests: T026, T027, and T028; implementation T029 and T032.
- US3 tests: T034, T035, T036, and T037.
- Polish: T043, T044, and T045.

## Parallel Example: User Story 3

    Task: T034 Protocol compatibility tests in internal/protocol/events_test.go and tests/integration/protocol_contract_test.go
    Task: T035 Reducer/storage tests in internal/events/reducer_test.go and internal/storage/events_test.go
    Task: T036 WebSocket replay/security tests in tests/integration/websocket_test.go
    Task: T037 Browser event-client tests in web/src/protocol/events.test.ts

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete setup and foundational safety/runtime work through T015.
2. Complete T016–T025 for US1.
3. Run the US1 test suites and manually follow the first quickstart scenario.
4. Stop for review: the executable must open a loopback dashboard that identifies a Git project before adding lifecycle convenience or live events.

### Incremental Delivery

1. Complete US1 and validate the basic local workspace.
2. Complete US2 and validate CLI-only management and diagnostics.
3. Complete US3 and validate ordered live state/reconnect behavior.
4. Complete Phase 6 to certify the embedded, cross-platform release path.

### Task Summary

| Scope | Tasks |
|---|---:|
| Setup | 7 |
| Foundational | 8 |
| US1 — P1 MVP | 10 |
| US2 — P2 | 8 |
| US3 — P3 | 9 |
| Polish | 5 |
| **Total** | **47** |

## Notes

- Every checklist item has an ID, an exact file path, and a story label where required.
- [P] appears only where work can proceed on separate files without waiting on an unfinished task.
- No task introduces provider adapters, Spec Kit artifact mutation, task dispatch, worktree automation, or Three.js; those are explicitly deferred roadmap milestones.

---

## Phase 7: Convergence

- [X] T048 CRITICAL Implement the single-owner event lifecycle: redact input, assign a recovered monotonic sequence, durably append the event, atomically persist its snapshot cursor, and then broadcast normalized core lifecycle events per FR-007, FR-009–FR-012, Constitution III–IV (missing).
- [X] T049 Implement a detached internal `serve` process and health-confirmed runtime lifecycle so `start` returns after readiness, `restart` leaves a live replacement runtime, and startup/shutdown owns and releases all resources per FR-001, FR-002, plan: detached serve runtime (contradicts).
- [X] T050 Verify an existing runtime record against its health identity, classify stopped or unreachable records as stale, invalidate obsolete control secrets, and safely recover before start/status/stop/restart operations per FR-017 and edge case: stale runtime record (missing).
- [X] T051 Add implicit-default-port probing with clear conflict and recovery output while preserving explicit-port failure behavior per FR-018 and edge case: requested port in use (partial).
- [X] T052 Implement actual runtime, Git, embedded-dashboard, Spec Kit, and provider checks; expose the safe typed results consistently through health and `doctor`, with optional integrations as non-blocking warnings per FR-006, FR-018, FR-019 (partial).
- [X] T053 Complete bounded event storage and recovery: enforce default 10,000-event retention with a configurable 1,000-event minimum, restore sequence/snapshot state, and return explicit resynchronization for cursors older than retained history per FR-011 and data-model: Event persistence lifecycle (partial).
- [X] T054 Implement the documented `/api/v1/events` cursor semantics and `craftlane.events.v1` WebSocket endpoint, including exact Origin/subprotocol checks, ordered replay, snapshot/resync bootstrap, heartbeats, and bounded slow-client queues per FR-009–FR-011 and contracts/websocket.md (missing).
- [X] T055 Implement dashboard display of runtime origin and typed diagnostics plus live-event snapshot application, cursor persistence, duplicate suppression, reconnect backoff, resync, and degraded connection state per FR-005, US3/AC1–AC2 (partial).
- [X] T056 Centralize redaction before every supported persistence, diagnostic, HTTP, and WebSocket sink, including snapshots and event segments, and prove raw credential-shaped values cannot escape per FR-016 and Constitution III (partial).
- [X] T057 Add and pass the missing foundational, lifecycle, HTTP-contract, runtime-control, protocol, reducer/storage, WebSocket, dashboard, and Playwright built-binary suites, including 20-start and 100-reconnect acceptance exercises per SC-001–SC-006 and Constitution VII (missing).
- [X] T058 Require a pinned contributor web toolchain and a real web build for Go build/test verification; expand CI and add macOS/Linux/Windows clean-install built-binary smoke coverage per FR-020 and plan: embedded release binary (partial).
- [ ] T059 Execute the final quickstart, offline/startup timing, success-criteria, and five-developer usability validation; record results and remediation in the feature checklist and quickstart per SC-001–SC-007 (missing).

---

## Phase 8: Convergence

- [X] T060 Make the event journal and snapshot cursor crash-consistent so an event cannot become durable without a snapshot at its sequence; preserve the redact → sequence/ID → persist → snapshot → broadcast lifecycle and add rollback coverage per T048, FR-007, FR-009–FR-012.
- [X] T061 Replace whole-history event-file rewrites with bounded append-only JSONL event segments while retaining atomic recovery, configured retention, and ordered replay coverage per plan: Storage decision.

---

## Phase 9: Convergence

- [X] T062 Render the current typed diagnostic/health summary in `craftlane status`, including safe dependency availability and corrective actions where applicable, and add CLI coverage per FR-002 and contracts/cli.md (partial).
- [X] T063 Add a bounded project-change detector that emits normalized, redacted `project.changed` lifecycle events and prove detection, persistence, and delivery per FR-012 (missing).
- [X] T064 Apply valid normalized lifecycle events to dashboard state while retaining cursor persistence, duplicate suppression, reconnect, and resync behavior; add browser-client coverage per US3/AC1–AC2 and FR-009 (partial).
- [X] T065 Add explicit local event-retention configuration, plumb it into the runtime reducer, enforce the 1,000-event floor and 10,000-event default, and document/test it per FR-011 and data-model: Retention policy (missing).
- [X] T066 Add and consistently surface a typed project-detection diagnostic in health, doctor, and status results, including safe corrective guidance for unavailable detection, per FR-006 and contracts/cli.md (partial).
- [X] T067 Make event replay and snapshot capture a single atomic cursor boundary for HTTP and WebSocket bootstrap so clients cannot skip or duplicate events across concurrent appends; add race-focused contract coverage per FR-009–FR-011 and contracts/websocket.md (partial).

---

## Phase 10: Convergence

- [X] T068 Implement a bounded, coalesced project-change detector for an active runtime that emits normalized, redacted `project.changed` events and add persistence/delivery coverage per FR-012 and Constitution Operational Constraints (missing).
- [X] T069 Make detached runtime startup failure and interruption transactional: terminate an unready child, release listeners, invalidate control state, and persist a recoverable stopped or stale record on readiness timeout or initialization failure; add regression coverage per FR-017 and edge case: runtime interrupted during startup or shutdown (partial).

---

## Phase 11: Convergence

- [X] T070 Make detached startup cancellation and readiness-timeout cleanup transactional: retain and terminate the unready child, wait for its exit, clear control state, and persist a recoverable stopped or stale record; add regression coverage per FR-017 and edge case: runtime interrupted during startup or shutdown (partial).

---

## Phase 12: Convergence

- [ ] T071 Run the final five-developer unaided usability study with real Git-familiar participants; record individual and aggregate start/status/stop outcomes plus any remediation in specs/001-foundation-runtime/quickstart.md and specs/001-foundation-runtime/checklists/requirements.md per SC-007 and T047/T059 (missing).

---

## Phase 13: Convergence

- [ ] T072 Run and record the outstanding five-developer unaided usability study with real Git-familiar participants, including individual and aggregate start/status/stop outcomes and any remediation in specs/001-foundation-runtime/quickstart.md and specs/001-foundation-runtime/checklists/requirements.md per SC-007 (missing).
