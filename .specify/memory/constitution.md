<!--
Sync Impact Report
Version change: 1.0.0 → 1.1.0
Modified principles:
- I. Productivity and Keyboard-First Control → I. Productivity-First, Keyboard-Equivalent Control
- II. Local-First Ownership → II. Local-First, Account-Free Ownership
- III. Safe, Explicit, and Reversible Control → III. Safe, Explicit, and Reversible Actions
- IV. Provider-Agnostic Core and Stable Protocols → IV. Provider-Agnostic Core and Versioned Protocols
- V. Spec-Driven, Read-Only by Default → V. Spec-Driven, Bounded Delivery
- VI. Accessible Two-Dimensional Equivalence → VI. Accessible Two-Dimensional Equivalence
- VII. Maintainable Open-Source Delivery → VII. Maintainable, Testable Open-Source Delivery
Added sections: none.
Removed sections: none.
Follow-up TODOs: none.
-->

# Craftlane Constitution

## Core Principles

### I. Productivity-First, Keyboard-Equivalent Control

Every visual feature MUST improve situational awareness, task coordination, dependency
understanding, agent control, error visibility, workflow speed, or project comprehension. Frequently
used actions MUST have a keyboard-accessible standard-interface path. A spatial interaction MUST
never be the sole route to an important action, and a faster non-spatial path MUST remain available
when it requires fewer interactions. Craftlane reduces development coordination overhead; it is not a
game layer.

### II. Local-First, Account-Free Ownership

Craftlane MUST operate without a Craftlane-hosted backend or mandatory Craftlane account. By
default, its UI and API MUST bind only to loopback, and project state, logs, event history, and
operational metrics MUST remain local. LAN exposure, remote access, telemetry, cloud
synchronization, and sharing project content MUST be explicit opt-in capabilities with clear user
disclosure. This preserves user control, offline utility, and open-source independence.

### III. Safe, Explicit, and Reversible Actions

Craftlane MUST NOT silently expose local services, execute arbitrary UI-supplied shell commands,
modify Spec Kit artifacts, delete branches or worktrees, auto-merge code, collect telemetry, or send
project source to Craftlane-controlled infrastructure. State-changing actions MUST be typed,
target-validated, protected by the local runtime's action security, and initiated by clear user
action. Secrets MUST be redacted before persistence or presentation. Safety defaults take precedence
over automation convenience.

### IV. Provider-Agnostic Core and Versioned Protocols

The core domain MUST model providers, sessions, agents, subagents, tasks, tool activity,
permissions, worktrees, Spec Kit workflows, and projects without relying on a provider's event
format or terminology. Provider-specific formats MUST terminate at adapters, and declared
capabilities MUST govern available controls. The normalized wire protocol MUST be versioned;
established event semantics MUST not be silently redefined, and compatible evolution MUST prefer
additive fields. The core, not React or Three.js, owns canonical state.

### V. Spec-Driven, Bounded Delivery

Craftlane MUST treat the Spec Kit lifecycle as first-class: constitution, specification,
clarification, plan, checklist, tasks, analysis, implementation, and convergence. Integrations MUST
start read-only; artifact modification requires a deliberate user action. Work MUST be scoped to a
bounded roadmap milestone with testable requirements, explicit task dependencies, and measurable
outcomes. The master roadmap MUST NOT be implemented as a single feature.

### VI. Accessible Two-Dimensional Equivalence

All operational state and controls rendered in World Mode MUST be available in an accessible standard
interface. That interface MUST support keyboard navigation, visible focus, reduced motion,
non-color-only status communication, and screen-reader-compatible labels. The spatial world MUST be
a deterministic projection of canonical state and MUST NOT become its source of truth. Craftlane
MUST remain productive when Three.js is disabled or unsupported.

### VII. Maintainable, Testable Open-Source Delivery

Craftlane MUST favor understandable code, documented extension points, stable interfaces, testable
domain modules, minimal setup, and independently understandable boundaries. Dependencies and
abstraction layers MUST be justified by present product value; speculative flexibility is
insufficient. Changes to core state, adapters, storage, protocols, security boundaries, Git or
worktree behavior, and browser-facing behavior MUST include proportionate automated tests. Craftlane
MUST complement existing IDEs, terminals, Git tools, providers, and Spec Kit rather than replace
every development surface.

## Operational Constraints

- The UI MUST only project normalized core state; it MUST NOT own authoritative agent, task, or
  project state.
- Browser-visible data MUST NOT include credentials, provider environment variables, or arbitrary
  local file access. Local paths MUST be constrained to the project, known worktrees, or an
  explicitly selected directory.
- State-changing local endpoints MUST use runtime protection and strict host/origin validation; no
  generic shell endpoint or permissive cross-origin access is permitted.
- Write-heavy independent work MUST use isolated Git worktrees when supported. Early releases MUST
  NOT auto-merge, and worktree removal MUST require safety checks and user confirmation.
- Agent status and attention MUST distinguish ordinary progress, permission requests, waiting,
  blocking conditions, failures, and conflicts. Inferred state MUST remain distinguishable from
  explicit provider state.
- Project watching, Git refreshes, and event delivery MUST be bounded and coalesced. UI render loops
  MUST NOT drive repository polling or command execution.
- Released binaries MUST run without Node.js, npm, or pnpm at runtime; frontend assets MUST ship with
  the native executable.

## Development Workflow and Quality Gates

- Each roadmap milestone MUST follow the Spec Kit lifecycle before implementation: specification,
  clarification when needed, plan, requirements checklist, tasks, analysis, implementation, and
  convergence.
- Before implementation, specifications MUST have no unresolved critical clarification markers,
  plans MUST demonstrate constitution compliance, and task dependencies MUST be explicit.
- Reviews MUST verify local-only defaults, secret redaction, path and action validation, accessible
  two-dimensional equivalence, provider capability handling, protocol compatibility, and preservation
  of existing developer workflows.
- Releases MUST pass reproducible build, test, package, embedded-asset, and clean-install checks on
  supported platforms. User-visible configuration or protocol changes MUST include documentation and
  migration notes.

## Governance

This constitution supersedes conflicting implementation preferences, feature plans, and convenience
shortcuts. Every specification, plan, task set, review, and release decision MUST record and resolve
any material exception before work proceeds. Amendments require a written rationale, impact analysis
for active features and compatibility, an updated Sync Impact Report, and a semantic version bump:
major for removed or incompatible principles, minor for new or materially expanded governance, and
patch for clarifications that do not alter obligations. Compliance review MUST reject changes that
weaken local-first ownership, safety, accessibility, provider abstraction, canonical state ownership,
or protocol compatibility without an approved major amendment.

**Version**: 1.1.0 | **Ratified**: 2026-10-01 | **Last Amended**: 2026-10-01
