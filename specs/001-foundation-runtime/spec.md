# Feature Specification: Local Runtime Foundation

**Feature Branch**: `001-foundation-runtime`

**Created**: 2026-10-01

**Status**: Draft

**Input**: User description: "Generate the first bounded Craftlane feature from the master Craftlane plan: a local-first runtime, command-line entry point, local dashboard shell, event protocol, storage, project detection, localhost security, and diagnostics."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Open a Local Project Workspace (Priority: P1)

A developer working inside a Git repository runs Craftlane from the command line and is taken to a local dashboard that identifies the repository being managed.

**Why this priority**: This is the first-run promise of Craftlane and establishes the local workspace on which all later agent and workflow features depend.

**Independent Test**: From a Git repository, run the primary Craftlane command and verify that the browser dashboard opens, reports the repository identity and location, and is reachable only from the same machine.

**Acceptance Scenarios**:

1. **Given** a developer is in a valid Git repository and Craftlane is installed, **When** they run `craftlane`, **Then** Craftlane detects the project, starts a local workspace, and opens its dashboard unless opening the browser was explicitly disabled.
2. **Given** Craftlane has started for a project, **When** the developer opens the reported local address, **Then** the dashboard shows the detected project name, root location, and runtime health.
3. **Given** a developer uses the start command with a no-browser or headless option, **When** the service starts successfully, **Then** Craftlane reports the local address without attempting to open a browser.

---

### User Story 2 - Manage the Local Workspace from the Command Line (Priority: P2)

A developer can start, stop, restart, open, inspect, and diagnose their local Craftlane workspace without relying on the dashboard.

**Why this priority**: Command-line control keeps Craftlane compatible with existing terminal workflows and makes a usable runtime possible even when a browser is unavailable.

**Independent Test**: Start a workspace without opening a browser, use status to inspect it, open it on request, stop it, and verify that status reports it as stopped.

**Acceptance Scenarios**:

1. **Given** no workspace is running for the current project, **When** the developer runs the start command, **Then** Craftlane reports whether startup succeeded and how to access the workspace.
2. **Given** a workspace is running, **When** the developer runs status, **Then** Craftlane reports its running state, project identity, local address, and detected dependency health.
3. **Given** a workspace is running, **When** the developer runs stop or restart, **Then** Craftlane completes the requested lifecycle action cleanly and reports its result.
4. **Given** Craftlane or a project dependency is not ready, **When** the developer runs doctor, **Then** Craftlane identifies the affected check and provides a practical corrective action.

---

### User Story 3 - Receive Safe, Ordered Local Updates (Priority: P3)

A developer with the dashboard open sees new normalized project updates without manually refreshing, and can reconnect without losing updates produced while the dashboard was disconnected.

**Why this priority**: Reliable live state is the contract needed by later agent observability, Spec Kit, task, and worktree capabilities.

**Independent Test**: Connect a dashboard client, cause representative project lifecycle updates, disconnect it, reconnect it, and verify that it receives the missed updates once and in order.

**Acceptance Scenarios**:

1. **Given** a dashboard client is connected to a project workspace, **When** the runtime records a project update, **Then** the client receives a normalized update that identifies its schema version, project, source, type, time, and local order.
2. **Given** a dashboard client reconnects after an interruption, **When** it asks for updates after its last received order, **Then** it receives the missing updates in ascending order without duplicating already acknowledged updates.
3. **Given** a local website from a different origin attempts to perform a state-changing runtime operation, **When** it lacks the runtime's current authorization proof, **Then** Craftlane rejects the operation.

### Edge Cases

- The requested local port is in use: Craftlane either selects an available nearby local port or exits with a clear message that names the conflicting port and a recovery option.
- The command is run outside a Git repository: Craftlane does not accidentally manage an unrelated directory and explains that a Git project is required for this milestone.
- A prior runtime record exists but its process is no longer running: status identifies it as stale, and start can recover without leaving ambiguous state.
- A browser cannot be opened automatically: the successful start result still includes a copyable local address.
- The runtime is interrupted during startup or shutdown: it leaves no active listener or corrupt local state behind.
- A request contains an out-of-project or otherwise unapproved filesystem path: Craftlane denies it without exposing the path's contents.
- An event or diagnostic field appears to contain a credential: Craftlane redacts it before storing, displaying, or returning it.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Craftlane MUST provide a primary command-line entry point that, when run from a Git repository, detects the current project, starts its local workspace, and opens the dashboard by default.
- **FR-002**: Craftlane MUST provide start, stop, restart, open, status, and doctor commands, each with human-readable results and meaningful nonzero exit status on failure.
- **FR-003**: Craftlane MUST allow the start command to suppress automatic browser opening and to operate without an interactive browser.
- **FR-004**: Craftlane MUST run its workspace services on the loopback interface by default and MUST NOT expose them to a local network unless a future explicit opt-in feature authorizes it.
- **FR-005**: Craftlane MUST show a minimal dashboard that identifies the detected project, local workspace address, and current runtime health.
- **FR-006**: Craftlane MUST expose a local health result that reports runtime availability, project-detection status, Git availability, Spec Kit availability when detectable, and provider availability when detectable.
- **FR-007**: Craftlane MUST persist only the local runtime information needed to restore status, diagnostics, project identity, and event continuity; it MUST keep that information on the developer's machine.
- **FR-008**: Craftlane MUST not require a Craftlane account, a Craftlane-hosted backend, or network access for normal local startup and dashboard use.
- **FR-009**: Craftlane MUST make normalized project updates available to connected local dashboard clients in real time through a documented, versioned event contract.
- **FR-010**: Each emitted event MUST include a schema version, unique local event identity, monotonically increasing project-local sequence, timestamp, project identity, source identity, event type, and event payload.
- **FR-011**: Craftlane MUST retain a configurable, bounded event history of 10,000 events per project by default and no fewer than 1,000 events, so a reconnecting client can request updates after its last known sequence and receive them in order.
- **FR-012**: Craftlane MUST record and surface project lifecycle updates, including project detection and project changes, even before provider-specific observability is added.
- **FR-013**: Craftlane MUST validate the host of incoming local requests and the origin of live-update connections; it MUST NOT permit unrestricted cross-origin access.
- **FR-014**: Craftlane MUST require a short-lived, runtime-specific authorization proof for every state-changing local operation and MUST not expose a generic arbitrary-command execution operation.
- **FR-015**: Craftlane MUST validate all runtime-accessible filesystem locations against the current project, a known workspace, or an explicitly selected location before acting on them.
- **FR-016**: Craftlane MUST redact likely credentials and secret-bearing values before they are persisted, displayed, included in events, or included in diagnostics.
- **FR-017**: Craftlane MUST shut down cleanly when stopped or interrupted, releasing its local address and preserving only valid recoverable state.
- **FR-018**: Craftlane MUST provide clear startup, conflict, missing-dependency, and recovery messages that identify a next action a developer can take.
- **FR-019**: Craftlane MUST remain usable when optional provider or Spec Kit integrations are absent, showing their unavailable state without blocking the dashboard.
- **FR-020**: A released Craftlane executable MUST run the dashboard without requiring supplementary web-development tooling on the developer's machine.

### Key Entities

- **Project**: The Git repository currently managed by a Craftlane workspace, including its stable local identity, display name, root location, and detected capability health.
- **Runtime Workspace**: The local, project-scoped Craftlane service instance, including its lifecycle state, loopback address, authorization lifetime, and diagnostic state.
- **Runtime Record**: Persisted local information that lets Craftlane distinguish a running workspace from a stopped or stale one and provide useful status.
- **Event**: A versioned, normalized record of a project update with an ordered sequence, origin, type, time, and payload.
- **Event Source**: The origin that produced an event, including the core runtime or a future provider adapter, identified without exposing secret configuration.
- **Diagnostic Check**: A named readiness or health result with status, relevant local detail, and a recommended corrective action when unhealthy.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In a supported Git repository on a developer machine, a developer can reach a project-identifying Craftlane dashboard within 10 seconds of running the primary command under normal local conditions.
- **SC-002**: In acceptance testing across 20 consecutive clean starts, the primary command successfully detects the intended Git project and reports a usable loopback-only address in at least 19 runs.
- **SC-003**: In 100 simulated reconnect cycles with recorded project updates, clients recover all retained missed updates in sequence with no duplicate event identities.
- **SC-004**: In security acceptance tests, 100% of state-changing requests missing a valid runtime authorization proof and 100% of cross-origin live-update attempts are rejected.
- **SC-005**: In representative diagnostic failures for missing Git support, unavailable optional integrations, occupied ports, and stale runtime records, every failure report identifies the affected check and offers at least one corrective next action.
- **SC-006**: In a clean supported production installation without supplementary web-development tooling, a developer can start the local workspace and load the dashboard successfully.
- **SC-007**: In usability testing with five developers familiar with Git, at least four can start a workspace, find its status, and stop it without documentation or assistance.

## Assumptions

- This specification covers only roadmap milestone `001-foundation-runtime`; Claude observability, Spec Kit parsing, task graphs, worktrees, dispatch, and spatial visualization are later features.
- The first supported use case is a developer invoking Craftlane from within a local Git working tree; selecting arbitrary non-Git folders is out of scope.
- The local dashboard is intentionally minimal and serves as an operational shell rather than a full command dashboard.
- The local runtime chooses its standard local address by default; if that address is unavailable, it may use a nearby available local address and must clearly report the result.
- Event retention is configurable, defaults to 10,000 events per project, and may not be configured below 1,000 events; requests preceding retained history receive an explicit resynchronization result. Long-term archival, synchronization, and remote sharing are out of scope.
- Existing Craftlane configuration files, if present, may supply local preferences, but secrets and provider environment values are never presented to the browser.
- Browser launching is best-effort and does not determine whether startup succeeds.
