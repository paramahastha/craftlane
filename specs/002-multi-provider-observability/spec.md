# Feature Specification: Multi-Provider Agent Observability

**Feature Branch**: `002-multi-provider-observability`

**Created**: 2026-10-02

**Status**: Draft

**Input**: User description: "Create the next bounded Craftlane milestone from `craftlane-plan-go.md`: add agent observability for Claude Code, Codex, and OpenCode. Craftlane should detect supported sessions and subagents using supported local lifecycle integrations, normalize activity into provider-agnostic agent states, display current tool activity and permission waits, and stream updates to the local UI without exposing secrets."

## Clarifications

### Session 2026-10-02

- Q: Should this milestone remain limited to a Claude Code adapter while requiring its normalized agent model to be reusable by future providers? → A: Superseded: the first phase must integrate Claude Code, Codex, and OpenCode while keeping provider-neutral core models.
- Q: What should Craftlane show when observation is enabled after an agent session is already running? → A: Create or update an agent from its first later lifecycle signal, and mark prior history as unavailable.
- Q: Must Claude Code, Codex, and OpenCode all expose the same observation features in this first phase? → A: Integrate all three; support each provider's available signals and clearly expose any unavailable capability.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - See Active Agent Work (Priority: P1)

A developer with Craftlane running can see active Claude Code, Codex, and OpenCode sessions and subagents for the current project without repeatedly checking terminals or refreshing the dashboard.

**Why this priority**: Immediate visibility into active agent work is the core value of this milestone and enables all later coordination features.

**Independent Test**: For each supported first-phase provider, start a session that creates a subagent when available, then verify that the local dashboard shows each agent, its relationship, current state, and latest activity as it changes.

**Acceptance Scenarios**:

1. **Given** Craftlane is running in a supported project and Claude Code, Codex, or OpenCode is available, **When** a supported provider session starts, **Then** Craftlane records and presents an agent for that session without a dashboard refresh.
2. **Given** an observed provider session starts a subagent, **When** the subagent lifecycle information is available, **Then** Craftlane presents the subagent separately and identifies its parent session when that relationship is available.
3. **Given** an observed agent performs supported tool activity, **When** the activity changes, **Then** Craftlane presents the current tool activity and a locally ordered update to connected dashboard clients.
4. **Given** a provider does not expose a supported observation category, **When** the developer views its integration status, **Then** Craftlane identifies that category as unavailable rather than implying equivalent coverage.

---

### User Story 2 - Respond to Permission and Completion States (Priority: P2)

A developer can quickly identify when an observed agent needs permission, has become idle, has completed, or has failed, and can distinguish reported state from a state inferred from activity.

**Why this priority**: Attention is costly; clear, trustworthy state helps a developer decide when to intervene without misinterpreting quiet activity as a failure.

**Independent Test**: Deliver representative lifecycle records for a permission request, inactivity, successful stop, and failure, then verify the normalized agent state, its source, and its latest relevant detail.

**Acceptance Scenarios**:

1. **Given** an observed agent explicitly requests permission, **When** Craftlane receives the request, **Then** the agent is marked as waiting for permission with the reported source and an actionable attention indicator.
2. **Given** an observed session remains alive but has no recent activity beyond the configured inactivity threshold, **When** no explicit idle state is available, **Then** Craftlane may mark it idle and labels that state as inferred.
3. **Given** an observed agent ends successfully or with failure, **When** Craftlane receives the terminal lifecycle information, **Then** it marks the agent completed or failed, preserves its last-known summary, and streams the change to connected clients.

---

### User Story 3 - Set Up and Diagnose Provider Observation Safely (Priority: P3)

A developer can explicitly enable observation for each available first-phase provider, understand the minimal changes Craftlane will make, and diagnose unavailable or degraded integration without preventing use of the local dashboard.

**Why this priority**: Each integration must respect local control and remain dependable across differing provider installations.

**Independent Test**: Run setup in projects with each first-phase provider available, unavailable, and already configured; inspect the proposed changes before approval; then run diagnostics and verify useful readiness and recovery results.

**Acceptance Scenarios**:

1. **Given** a project has no Craftlane integration for an available provider, **When** the developer requests setup, **Then** Craftlane detects relevant local prerequisites, explains the exact project-scoped changes, and waits for explicit confirmation before modifying provider configuration.
2. **Given** the developer confirms setup, **When** Craftlane adds its required project-scoped integration files or settings, **Then** it makes only the disclosed minimal changes and provides a reversible removal path.
3. **Given** a first-phase provider, its supported lifecycle integration, or the local callback path is unavailable, **When** the developer runs diagnostics, **Then** Craftlane identifies the affected condition, offers a practical next action, and keeps the rest of the local dashboard usable.

### Edge Cases

- Claude Code, Codex, or OpenCode is not installed, cannot be found, or does not support a required lifecycle integration: Craftlane shows that provider's observation as unavailable or degraded without treating the local runtime as failed.
- A lifecycle record is malformed, incomplete, duplicated, out of order, or from an unknown session: Craftlane rejects or safely records the issue without corrupting known agent state or exposing raw sensitive content.
- A subagent record has no parent relationship: Craftlane displays it as an independently observed agent and preserves the missing relationship as unavailable rather than inventing one.
- The local callback receives traffic without valid provider-specific authorization, from an unapproved host or Origin, or for another project: Craftlane rejects it and does not create agent state.
- A tool argument, provider payload, or diagnostic contains likely credentials: Craftlane redacts the sensitive value before storage, presentation, or delivery to a client.
- A first-phase provider evolves so that some lifecycle categories are unavailable: Craftlane detects available capabilities, observes the supported subset, and communicates the reduced coverage.
- Craftlane restarts while sessions are active: it preserves valid local observation history and reconnects or resumes observation where supported without presenting stale agents as confirmed active.
- Observation begins after a provider session is already running: Craftlane creates or updates the agent at its first subsequent supported lifecycle signal and clearly marks earlier activity as unavailable rather than reconstructing private or unstable provider history.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Craftlane MUST detect whether Claude Code, Codex, and OpenCode and their lifecycle integrations needed for observation are available for the current project, without requiring a Craftlane account.
- **FR-002**: Craftlane MUST provide an explicit, project-scoped setup action for each first-phase provider that describes the proposed changes and requires developer confirmation before creating or modifying provider-related settings or hook files.
- **FR-003**: Setup MUST make only the disclosed minimal changes, avoid overwriting unrelated user configuration, and provide a documented reversible removal path.
- **FR-004**: Craftlane MUST receive supported lifecycle information from each first-phase provider only through a local, authenticated provider callback path and MUST reject callback traffic that lacks valid provider-specific authorization, has an unapproved Host or Origin, or does not belong to the active local project.
- **FR-005**: Craftlane MUST feature-detect supported lifecycle categories for each first-phase provider and expose the resulting observation capabilities instead of assuming every provider installation provides the same events.
- **FR-006**: Craftlane MUST normalize supported Claude Code, Codex, and OpenCode lifecycle information into provider-agnostic agent records without requiring generic dashboard behavior to depend on provider-specific fields.
- **FR-007**: Each observed agent record MUST identify its local agent identity, provider, project, session when available, display name or usable fallback, normalized state, state source, most recent activity time, and current task or work context when supplied.
- **FR-008**: When lifecycle information supplies a parent/subagent relationship, Craftlane MUST preserve and present that relationship; when it does not, Craftlane MUST keep the relationship unspecified.
- **FR-009**: Craftlane MUST normalize and present supported current tool activity, including whether the activity started, completed, or failed, without presenting sensitive tool input or provider environment values.
- **FR-010**: Craftlane MUST represent working, waiting, waiting for permission, idle, completed, and failed states when supported evidence exists; it MUST distinguish explicitly reported states from derived states and expose the source or confidence of each state.
- **FR-011**: Craftlane MUST mark an agent as waiting for permission when supported lifecycle information explicitly indicates a permission request and surface it as requiring developer attention.
- **FR-012**: Craftlane MAY derive idle only for an observed live session with no recent activity beyond a configurable inactivity threshold; all derived idle states MUST be labeled as inferred and MUST NOT be treated as blocked or failed.
- **FR-013**: Craftlane MUST process duplicate, malformed, unknown, and out-of-order lifecycle records without corrupting the current state of previously observed agents.
- **FR-014**: Craftlane MUST publish normalized agent and tool updates through the existing versioned local event stream so connected dashboard clients receive changes without a manual refresh and reconnecting clients can recover retained updates in order.
- **FR-015**: Craftlane MUST redact likely credentials and secret-bearing values from provider payloads, tool activity, event records, logs, diagnostics, persisted state, and browser-visible data.
- **FR-016**: Craftlane MUST NOT send provider prompts, project source, provider environment values, or observation data to Craftlane-controlled infrastructure; observation data remains on the developer's machine by default.
- **FR-017**: Craftlane MUST retain only the local agent observation state and event history needed for current visibility and continuity, in accordance with the runtime's bounded local retention policy.
- **FR-018**: Craftlane MUST provide diagnostics for each first-phase provider that report availability, integration setup status, callback readiness, project-path validity, required local permissions, and detected observation capabilities, with a corrective next action for each unhealthy condition.
- **FR-019**: When a provider observation integration is unavailable or degraded, Craftlane MUST show the reason and continue to provide the local runtime, dashboard, and other provider integrations that do not depend on it.
- **FR-020**: Craftlane MUST use captured, sanitized lifecycle fixtures for Claude Code, Codex, and OpenCode to verify normalization for each supported lifecycle category.
- **FR-021**: When observation begins after a provider session is already running, Craftlane MUST create or update the agent from its first subsequent supported lifecycle signal and MUST identify prior activity as unavailable rather than attempting to reconstruct it from provider history.
- **FR-022**: Craftlane MUST integrate Claude Code, Codex, and OpenCode in this phase without requiring identical lifecycle coverage; it MUST normalize each provider's supported categories and explicitly identify unavailable categories in provider capability status and diagnostics.

### Key Entities

- **Provider Capability**: A locally detected observation function that a provider installation supports, such as sessions, subagents, tools, permissions, or terminal outcomes.
- **Agent**: A normalized representation of one observed provider session or subagent, with identity, provider, project, state, state source, activity summary, and optional parent relationship.
- **Agent Relationship**: A known parent-to-subagent association between observed agents; it remains unspecified when the provider does not supply it.
- **Tool Activity**: The current or most recent observed operation performed by an agent, including its lifecycle outcome and safely displayable summary.
- **Agent State**: The normalized current lifecycle condition of an agent, including whether it was explicitly reported or derived from available evidence.
- **Permission Request**: An explicit provider signal that an agent requires a developer decision before it can continue.
- **Observation Setup**: The project-scoped, developer-approved configuration that allows one first-phase provider's supported lifecycle information to reach the local Craftlane runtime.
- **Lifecycle Record**: One authenticated provider observation input that is validated, redacted, normalized, and converted into local agent state and event updates.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In 20 supported local test sessions per first-phase provider, Craftlane shows a newly started session in the dashboard within 2 seconds of receiving its lifecycle information in at least 19 sessions for each provider.
- **SC-002**: In fixture and end-to-end tests for Claude Code, Codex, and OpenCode covering every lifecycle category each provider supports, 100% of supported records produce the expected normalized agent or tool update.
- **SC-003**: In 20 observed sessions per first-phase provider that create supported subagents, Craftlane presents every provider-supplied parent/subagent relationship correctly and does not invent a relationship when the input omits one.
- **SC-004**: In 100 simulated reconnect cycles during agent and tool activity, a dashboard client recovers all retained agent updates in local sequence order with no duplicate event identities.
- **SC-005**: In security tests containing representative credentials in lifecycle records and tool inputs, 100% of detected secret-bearing values are absent from browser-visible state, local event history, diagnostics, and logs.
- **SC-006**: In tests of unauthorized, wrong-project, and malformed callback traffic, 100% of such inputs are rejected or safely isolated without changing an existing agent's valid state.
- **SC-007**: In usability testing with five developers who use at least one first-phase provider, at least four can identify an active agent, a subagent when supported, its current activity, and a permission wait when supported from the dashboard without opening a terminal.
- **SC-008**: In unavailable or degraded integration tests, every diagnostic names the affected condition and gives at least one corrective next action while the non-observation dashboard remains accessible.

## Assumptions

- This specification covers the expanded first observability phase: it ships adapters for Claude Code, Codex, and OpenCode. DeepSeek-based tools, Gemini, local agent CLIs not represented by these adapters, and custom providers are later integrations; the command dashboard, Spec Kit interpretation, worktrees, task dispatch, and spatial workspace are also out of scope.
- The existing local runtime foundation supplies the loopback-only workspace, authenticated local event delivery, bounded event history, and base dashboard shell required by this feature.
- Supported local lifecycle integrations are the preferred deterministic source of observation for each first-phase provider, and Craftlane observes only categories that the installed version explicitly supports.
- Feature parity is not assumed: each first-phase provider reports its available and unavailable observation capabilities so users can make informed coordination decisions.
- The feature is read-only with respect to agent execution: it observes lifecycle information but does not start, stop, message, approve, or otherwise control agents.
- A developer explicitly opts into project-scoped hook or integration setup for each provider; Craftlane never silently changes user-wide or project provider settings.
- Local observation data follows the runtime's existing bounded retention behavior; long-term archival, cloud synchronization, and team sharing are out of scope.
- If observation starts after a session begins, Craftlane establishes visibility from the first subsequent supported lifecycle signal; earlier session history is unavailable in this milestone.
- An inferred idle state uses a configurable local threshold with a reasonable default selected during planning; it is not proof that an agent is blocked, completed, or failed.
