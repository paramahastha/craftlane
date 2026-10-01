# Craftlane — End-to-End Build Plan

> **Project:** Craftlane  
> **Purpose:** Open-source, local-first spatial development environment for orchestrating AI coding agents  
> **Primary initial provider:** Claude Code  
> **Primary workflow integration:** GitHub Spec Kit  
> **Primary visual stack:** React + Three.js  
> **Execution model:** Localhost only by default  
> **Distribution target:** GitHub Releases + native single-binary CLI; package-manager distribution later  
> **Document role:** Master implementation roadmap to feed into Spec Kit feature-by-feature

---

## 0. Executive Summary

Craftlane is a local-first developer tool that gives engineers one place to understand, control, and coordinate multiple AI coding agents working on a software project.

The product should feel like a modern developer control plane with an optional spatial/3D layer, not like a game that happens to run code.

The core value proposition is:

- see every active AI agent and its current state;
- know which tasks are running, waiting, blocked, idle, completed, or failed;
- understand parent/subagent relationships;
- visualize Spec Kit artifacts and task dependencies;
- identify parallel work automatically;
- assign work to isolated Git worktrees;
- inspect agents without opening many terminals;
- control common actions quickly with keyboard-first UX;
- keep all orchestration and project state on the user's local machine;
- make the 3D environment useful for awareness, but never required for normal work.

A good one-sentence product definition is:

> **Craftlane is a local-first spatial development environment for orchestrating AI coding agents, visualizing Spec Kit workflows, and coordinating parallel software work.**

A shorter public tagline can be:

> **See your agents work. Control the flow. Ship faster.**

---

# 1. Product Principles

These principles are non-negotiable and should later be encoded into the Craftlane Spec Kit constitution.

## 1.1 Productivity before entertainment

Every visual feature must improve at least one of:

- situational awareness;
- task coordination;
- dependency understanding;
- agent control;
- error visibility;
- workflow speed;
- project comprehension.

A 3D interaction must never be the only way to complete an important action.

If a 3D interaction takes more actions than a command palette, keyboard shortcut, or normal UI control, the faster path must remain available.

## 1.2 Local-first

Craftlane must work locally without requiring a Craftlane-hosted backend.

Default runtime:

```text
Browser UI:
http://127.0.0.1:7331

Local API:
http://127.0.0.1:7332
```

The implementation may later serve the web UI and API from a single port, but the system design must preserve the logical separation between UI and runtime.

Default network binding must be:

```text
127.0.0.1
```

not:

```text
0.0.0.0
```

Remote/LAN exposure can be added later as an explicit opt-in feature.

## 1.3 No mandatory Craftlane account

Users should be able to:

```bash
craftlane
```

or:

```bash
craftlane start
```

inside a local Git repository and use the product.

No Craftlane login is required.

## 1.4 Provider-agnostic core

Claude Code is the first supported provider, but Craftlane's internal concepts must not be Claude-specific.

The core should understand:

```text
AgentProvider
AgentSession
Agent
Subagent
Task
ToolActivity
PermissionRequest
Worktree
SpecWorkflow
Project
```

rather than:

```text
ClaudeOnlySession
ClaudeOnlyTask
```

Future adapters should be possible for:

- Codex;
- OpenCode;
- Gemini;
- local agent CLIs;
- custom company agents.

## 1.5 Spec-driven development is first-class

Craftlane should not merely display Spec Kit files.

It should understand:

- constitution;
- specification;
- clarification state;
- technical plan;
- requirement checklist;
- tasks;
- parallel markers;
- dependencies;
- analyze findings;
- implementation state;
- convergence state.

The Spec Kit workflow should become visible and easier to operate.

## 1.6 Existing developer tools remain usable

Craftlane must never require users to stop using:

- their IDE;
- terminal;
- Git CLI;
- Claude Code directly;
- Spec Kit directly.

Craftlane is an orchestration and observability layer, not a replacement for every development surface.

## 1.7 Open-source friendly

Architecture must favor:

- understandable code;
- documented extension points;
- stable interfaces;
- testable modules;
- minimal setup;
- no proprietary cloud dependency.

## 1.8 Safety and control

Craftlane must not silently:

- expose localhost services to the network;
- run arbitrary shell commands from untrusted UI input;
- auto-merge code;
- delete branches/worktrees;
- modify Spec Kit artifacts without clear user action;
- send project source to Craftlane-owned infrastructure;
- collect telemetry by default.

---

# 2. Scope Strategy

Craftlane is a large product. Do **not** implement this master plan as one Spec Kit feature.

Use this file as the master roadmap, then create independent Spec Kit feature directories for each major milestone.

Recommended feature sequence:

```text
001-foundation-runtime
002-claude-observability
003-command-dashboard
004-speckit-integration
005-task-graph-and-worktrees
006-threejs-workspace
007-agent-control-and-dispatch
008-hardening-packaging-release
009-provider-sdk
010-insights-and-optimization
```

Each feature should go through:

```text
constitution
→ specify
→ clarify
→ plan
→ checklist
→ tasks
→ analyze
→ implement
→ converge
```

For major phases, implementation should be executed in bounded stages rather than one extremely large `/speckit.implement` run.

---

# 3. Target User

Primary user:

- software engineer using AI coding agents;
- uses Claude Code or similar tools;
- regularly runs multiple parallel agents;
- works with Git branches/worktrees;
- wants to understand agent progress without checking many terminals;
- wants more structured spec-driven development;
- values local control and developer tooling.

Secondary users:

- tech leads supervising agent-assisted feature development;
- open-source maintainers using multiple coding agents;
- small engineering teams experimenting with agentic development;
- developers using multiple AI providers.

---

# 4. Core User Problems

Craftlane should solve the following problems.

## 4.1 Too many terminals

Developers running multiple agents lose time switching between terminals to understand:

- which agent is still running;
- what each agent is doing;
- whether one needs permission;
- whether one is blocked;
- which task finished;
- which agent became idle.

## 4.2 Parallel work is hard to coordinate

Developers often need to mentally reason about:

- which tasks can run together;
- which agent should take which task;
- what dependencies block another task;
- which worktrees are safe to merge;
- where collisions may occur.

## 4.3 Spec Kit requires context switching

Spec Kit is powerful, but users still need to remember workflow stages and repeatedly open:

- `spec.md`;
- `plan.md`;
- `tasks.md`;
- checklists;
- analyze results;
- convergence results.

Craftlane should make that lifecycle visible.

## 4.4 Agent output is noisy

Raw logs are useful for debugging but poor for quick understanding.

Craftlane should provide concise operational summaries first and detailed logs second.

## 4.5 Agent attention is expensive

When several agents are active, the human should only be interrupted when attention is actually required.

Craftlane must separate:

```text
normal progress
informational completion
permission required
blocked work
failure
critical conflict
```

---

# 5. Product Modes

Craftlane should offer three UI modes.

## 5.1 Command Mode

Default mode.

Purpose:

- fast daily operation;
- compact agent overview;
- task overview;
- Spec Kit state;
- attention queue;
- command palette;
- keyboard-first actions.

This mode should remain useful even if Three.js is disabled.

## 5.2 World Mode

Three.js spatial view.

Purpose:

- visual understanding;
- project topology;
- parallel activity;
- dependencies;
- worktrees;
- agent movement;
- status recognition;
- exploration.

World Mode is optional.

## 5.3 Focus Mode

Focused view for one:

- agent;
- task;
- worktree;
- Spec Kit finding.

Purpose:

- deep inspection;
- current activity;
- logs;
- changed files;
- task details;
- decisions;
- review.

---

# 6. Product Theme and Visual Language

Craftlane uses a **modern software workshop / miniature engineering studio** theme.

Avoid:

- ocean/orca metaphors;
- realistic industrial factory simulation;
- mandatory avatar navigation;
- excessive game mechanics.

Use:

- clean miniature workspaces;
- low-poly robots;
- lanes;
- workbenches;
- blueprint stations;
- QA labs;
- review stations;
- integration gates;
- status lights;
- task routes.

Mapping:

| Craftlane World | Development Meaning |
|---|---|
| Robot | Agent |
| Smaller robot | Subagent |
| Workbench | Worktree |
| Blueprint station | Spec / Plan |
| Work order | Task |
| Development bay | Coding |
| QA lab | Tests |
| Inspection station | Review |
| Integration gate | Merge/integration |
| Lit route | Dependency satisfied |
| Closed gate | Dependency blocked |
| Yellow beacon | Human attention needed |
| Red beacon | Failure/block |
| Green completion light | Completed |
| Cargo/package | Task result |

Game-like elements should communicate real workflow state.

---

# 7. High-Level Architecture

```text
┌─────────────────────────────────────────────────────────────┐
│                    LOCAL DEVELOPER MACHINE                  │
│                                                             │
│  Claude Code / future agent providers                       │
│             │                                               │
│             ▼                                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │                 Provider Adapters                     │  │
│  │  Claude Adapter | future Codex | future OpenCode     │  │
│  └────────────────────────┬──────────────────────────────┘  │
│                           │                                 │
│                           ▼                                 │
│  ┌───────────────────────────────────────────────────────┐  │
│  │                  Craftlane Core                       │  │
│  │                                                       │  │
│  │ Agent Registry / Session Model                        │  │
│  │ Project Manager                                       │  │
│  │ Task Graph                                            │  │
│  │ Attention Engine                                      │  │
│  │ Worktree Manager                                      │  │
│  │ Spec Kit Adapter                                      │  │
│  │ Event Bus                                             │  │
│  │ Storage Adapter                                       │  │
│  └────────────────────────┬──────────────────────────────┘  │
│                           │                                 │
│                Local API + WebSocket                        │
│                           │                                 │
│         ┌─────────────────┼────────────────────┐            │
│         ▼                 ▼                    ▼            │
│       CLI            Command UI           Three.js          │
│                      React UI             World UI           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

Important rule:

> UI layers never become the source of truth for agent state.

The core owns normalized state.

The UI only projects that state.

---

# 8. Repository Structure

Craftlane should use a mixed Go + web repository rather than a TypeScript monorepo.

Recommended layout:

```text
craftlane/
├── cmd/
│   └── craftlane/
│       └── main.go
│
├── internal/
│   ├── app/
│   ├── project/
│   ├── agent/
│   ├── providers/
│   │   └── claude/
│   ├── events/
│   ├── protocol/
│   ├── tasks/
│   ├── speckit/
│   ├── git/
│   ├── worktree/
│   ├── storage/
│   ├── security/
│   ├── server/
│   └── diagnostics/
│
├── pkg/
│   └── provider/
│       └── provider.go
│
├── web/
│   ├── src/
│   │   ├── app/
│   │   ├── dashboard/
│   │   ├── world/
│   │   ├── agents/
│   │   ├── tasks/
│   │   ├── speckit/
│   │   ├── worktrees/
│   │   ├── command-palette/
│   │   ├── protocol/
│   │   └── shared/
│   ├── public/
│   ├── package.json
│   ├── pnpm-lock.yaml
│   ├── tsconfig.json
│   └── vite.config.ts
│
├── protocol/
│   ├── schema/
│   │   ├── events.schema.json
│   │   ├── api.schema.json
│   │   └── common.schema.json
│   └── README.md
│
├── tests/
│   ├── fixtures/
│   ├── integration/
│   └── e2e/
│
├── examples/
├── docs/
│   ├── architecture/
│   ├── providers/
│   ├── speckit/
│   ├── contributing/
│   └── security/
│
├── scripts/
├── .github/
├── go.mod
├── go.sum
├── Makefile
├── LICENSE
├── README.md
├── CONTRIBUTING.md
├── CODE_OF_CONDUCT.md
└── SECURITY.md
```

Use `internal/` for implementation details that should not become a public Go API.

Use `pkg/provider` only when a stable provider SDK is intentionally made public.

Avoid generic modules such as:

```text
utils.go
helpers.go
common.go
misc.go
```

Prefer domain-specific package and file names.

The web application is a separately built frontend, but its production assets are embedded into the Go executable.

---

# 9. Technology Decisions

Craftlane intentionally uses different technologies for the runtime and browser layers.

## Backend / Local Runtime

Use:

- Go;
- Go modules;
- standard library wherever practical;
- `context.Context` for lifecycle/cancellation;
- goroutines for concurrent watchers and provider streams;
- channels where they make event ownership clearer;
- `net/http` for the local API;
- a small mature WebSocket dependency if required;
- `os/exec` for controlled child processes;
- `os/signal` for graceful shutdown;
- `embed` for production web assets.

Why Go:

- single native executable;
- fast startup;
- low idle memory;
- good process management;
- natural concurrency model;
- straightforward cross-platform compilation;
- no Node.js runtime requirement for end users;
- strong fit for a long-running local developer daemon/CLI.

## Frontend

Use:

- TypeScript;
- React;
- Vite;
- Three.js;
- Zustand or another small client state layer if needed;
- native Fetch;
- WebSocket client;
- HTML/CSS for productivity UI and accessibility.

Three.js must only render a projection of Craftlane state.

React must not own canonical backend state.

## Protocol

Do not share implementation-language types directly between Go and TypeScript.

Use an explicit versioned wire contract:

```text
protocol/schema/*.schema.json
```

or OpenAPI plus JSON Schema where appropriate.

The contract is the source of truth for serialized messages.

Generate or validate:

```text
Go wire types
+
TypeScript wire types
```

from the protocol schema where practical.

Domain types inside Go do not need to exactly mirror browser view models.

## Development Tooling

Contributors need:

```text
Go toolchain
Node.js
pnpm
```

End users should need only:

```text
craftlane executable
```

Node.js and pnpm are build-time dependencies for the frontend, not runtime dependencies for released binaries.

## Testing

Backend:

- Go standard `testing`;
- table-driven unit tests;
- race detector;
- temporary Git repositories for integration tests;
- fixture-based provider adapter tests;
- `httptest` for HTTP APIs.

Frontend:

- Vitest for isolated TypeScript/React logic;
- React Testing Library where useful;
- Playwright for browser end-to-end tests.

Cross-layer:

- protocol compatibility tests;
- executable smoke tests;
- WebSocket reconnect tests;
- release-binary tests on supported operating systems.

## Build

Development:

```text
Go API/runtime
+
Vite dev server
```

Production:

```text
web build
    ↓
web/dist
    ↓
go:embed
    ↓
single craftlane executable
```

The Go binary serves the embedded web application and local API.

For development, the Vite server may proxy API/WebSocket traffic to the Go runtime.

## Packaging

Primary release artifacts:

```text
craftlane-darwin-arm64
craftlane-darwin-amd64
craftlane-linux-amd64
craftlane-linux-arm64
craftlane-windows-amd64.exe
```

Additional platforms can be added based on demand.

Distribution order:

1. GitHub Releases;
2. install script for macOS/Linux;
3. PowerShell installer for Windows;
4. Homebrew tap;
5. Winget/Scoop or equivalent later;
6. optional `go install` for users who already have Go.

Do not require npm for production installation.

## Storage

For v0.1 prioritize zero-friction single-binary installation.

Use:

- append-only JSONL event log;
- JSON snapshot/state files;
- bounded local history;
- atomic file replacement for snapshots.

Create a Go storage interface so SQLite can be added later without changing business logic.

Example:

```go
type Storage interface {
    LoadProjectState(ctx context.Context, projectID string) (*ProjectState, error)
    SaveProjectState(ctx context.Context, state ProjectState) error
    AppendEvent(ctx context.Context, event Event) error
    ReadEvents(ctx context.Context, query EventQuery) ([]Event, error)
}
```

Do not make v0.1 depend on CGO or native database compilation.

If SQLite is introduced later, prefer a distribution strategy that preserves simple cross-platform builds and a single executable.

---


# 9.1 Go Runtime Architecture

The Go runtime should be organized around explicit ownership and cancellation.

Recommended conceptual flow:

```text
main
  ↓
App
  ├── Project service
  ├── Provider supervisor
  ├── Spec Kit watcher
  ├── Git/worktree watcher
  ├── Event reducer/store
  ├── Attention engine
  ├── HTTP/WebSocket server
  └── Storage
```

Top-level startup owns one root context.

```go
ctx, cancel := signal.NotifyContext(
    context.Background(),
    os.Interrupt,
    syscall.SIGTERM,
)
defer cancel()
```

Long-running components receive that context.

Avoid detached goroutines with no owner.

## Event ownership

Provider adapters and watchers produce events.

Craftlane Core serializes canonical state transitions.

A useful initial model is:

```text
producer goroutines
      ↓
bounded event channel
      ↓
single reducer / state owner
      ↓
snapshot + event persistence
      ↓
WebSocket broadcast
```

This avoids many classes of shared-state races.

If later throughput requires sharding, change the model deliberately and add race/performance tests.

## Concurrency rules

- every goroutine must have a documented owner;
- blocking subprocesses must be cancellable;
- channels must have intentional buffering;
- slow WebSocket clients must not block canonical state updates;
- file watchers must debounce/coalesce bursts;
- shutdown waits for owned components with a bounded timeout;
- run `go test -race` in CI.

## Embedded frontend

Production builds should embed `web/dist`.

Conceptually:

```go
//go:embed web-dist/*
var webAssets embed.FS
```

The actual path may differ based on build layout.

The release pipeline must ensure web assets are built before the final Go binary.

The runtime serves:

```text
/
frontend assets

/api/v1/*
HTTP API

/ws
WebSocket

/internal/*
provider/local integration endpoints
```

Development may use separate ports:

```text
Vite:        127.0.0.1:7331
Go runtime:  127.0.0.1:7332
```

Production should prefer one localhost origin where practical:

```text
Craftlane:   127.0.0.1:7331
```

Using one origin simplifies security and installation.


---

# 10. Runtime Directories

Global:

```text
~/.craftlane/
├── config.json
├── projects.json
├── logs/
└── runtime/
```

Per repository:

```text
<repo>/.craftlane/
├── config.json
├── state/
├── sessions/
└── events/
```

Recommended `.gitignore` behavior:

```gitignore
.craftlane/state/
.craftlane/sessions/
.craftlane/events/
```

Allow:

```text
.craftlane/config.json
```

to be committed if users want shared project configuration.

Never place provider credentials into project-level committed config.

---

# 11. CLI Design

The installed native executable is the primary entry point:

```bash
craftlane
```

Running it inside a Git repository should detect the current project, start the localhost service, and open the UI.

Primary commands:

```bash
craftlane
craftlane start
craftlane start --no-open
craftlane start --headless
craftlane start --port 7331

craftlane stop
craftlane restart

craftlane open
craftlane status
craftlane doctor

craftlane agents
craftlane tasks
craftlane spec
craftlane worktrees

craftlane init
craftlane config
```

Later:

```bash
craftlane dispatch <task-id>
craftlane agent inspect <agent-id>
craftlane agent stop <agent-id>
craftlane agent message <agent-id> "..."
```

Developer/source build commands may include:

```bash
make dev
make test
make build
make release-snapshot
```

CLI expectations:

- useful output without web UI;
- script-friendly exit codes;
- JSON output option later;
- fast startup;
- graceful Ctrl+C handling;
- errors should explain corrective action;
- no Node.js process should be required for a production installation.

Example:

```text
Craftlane v0.1

Project: serverstreamgw
Root: /Users/me/projects/serverstreamgw

Claude Code: detected
Spec Kit: detected
Git: healthy

Agents: 4
Working: 2
Waiting: 1
Idle: 1

Workspace:
http://127.0.0.1:7331
```

---

# 12. Local Server Design

Craftlane runs a local service.

Responsibilities:

- project detection;
- provider adapter lifecycle;
- file watching;
- Spec Kit parsing;
- Git/worktree state;
- normalized agent state;
- event storage;
- WebSocket broadcasting;
- local command execution through safe internal actions.

Default:

```text
127.0.0.1 only
```

Requirements:

- fail if requested port is already occupied, or select a nearby port with clear output;
- expose health endpoint;
- graceful shutdown;
- persist only required state;
- never expose arbitrary filesystem access to the browser.

---

# 13. Security Model

Localhost is not automatically safe.

Craftlane must protect against malicious websites attempting to interact with local services.

Requirements:

## Network

- bind loopback only by default;
- validate `Host`;
- validate WebSocket `Origin`;
- no permissive `Access-Control-Allow-Origin: *`.

## Mutating API actions

Require an ephemeral per-runtime token for state-changing operations.

The token can be injected into the same-origin frontend at startup.

Mutating endpoints require:

```text
X-Craftlane-Token
```

A malicious cross-origin webpage should not be able to create that custom request successfully.

## Filesystem

Every path must be validated against:

- current project root;
- known worktree roots;
- explicit user-selected directories.

Never accept arbitrary raw filesystem paths from UI without validation.

## Shell/process execution

Do not provide a generic:

```text
POST /shell
```

endpoint.

Expose typed actions instead:

```text
dispatchTask()
stopAgent()
runSpecKitAnalyze()
createWorktree()
```

## Secrets

Redact likely credentials from:

- logs;
- event payloads;
- UI tool arguments.

At minimum detect:

```text
API_KEY
TOKEN
SECRET
PASSWORD
AUTHORIZATION
```

Provider environment variables should never be sent to browser state.

## Telemetry

Default:

```text
off
```

If telemetry is added later:

- opt-in;
- documented;
- no source code;
- no prompt contents by default;
- easy to disable.

---

# 14. Craftlane Protocol

All state updates should flow through a normalized event protocol.

Example envelope:

```go
type Event struct {
    SchemaVersion int             `json:"schemaVersion"`
    EventID       string          `json:"eventId"`
    Sequence      uint64          `json:"sequence"`
    Timestamp     time.Time       `json:"timestamp"`

    ProjectID string      `json:"projectId"`
    Source    EventSource `json:"source"`
    Type      string      `json:"type"`

    Payload json.RawMessage `json:"payload"`
}

type EventSource struct {
    Provider       string `json:"provider"`
    AdapterVersion string `json:"adapterVersion,omitempty"`
}
```

Core events:

```text
project.detected
project.changed

agent.discovered
agent.started
agent.activity
agent.waiting
agent.blocked
agent.idle
agent.completed
agent.failed
agent.stopped

subagent.started
subagent.completed
subagent.failed

tool.started
tool.completed
tool.failed

permission.requested
permission.resolved

task.discovered
task.updated
task.assigned
task.started
task.blocked
task.completed
task.failed

worktree.discovered
worktree.created
worktree.changed
worktree.removed

speckit.detected
speckit.phase_changed
speckit.tasks_changed
speckit.analysis_changed
speckit.convergence_changed

git.changed
tests.started
tests.completed

attention.created
attention.resolved
```

Every event should have a monotonically increasing local sequence number.

Clients should be able to reconnect and request missed events.

---

# 15. Agent State Model

Normalized states:

```text
DISCOVERED
STARTING
WORKING
WAITING
WAITING_PERMISSION
BLOCKED
IDLE
COMPLETED
FAILED
STOPPED
UNKNOWN
```

State definitions must be precise.

## WORKING

Recent agent or tool activity indicates active execution.

## WAITING_PERMISSION

Agent explicitly requires human permission.

## WAITING

Agent is waiting for human input or another external condition.

## BLOCKED

Progress cannot continue due to:

- dependency;
- failure;
- conflict;
- missing prerequisite.

Do not infer BLOCKED only because an agent has been quiet.

## IDLE

Session/process is alive but has no active work.

For providers without an explicit idle event, this may be a derived state using:

- process/session health;
- latest tool activity;
- configurable inactivity threshold.

Derived state must be labeled internally as inferred.

## COMPLETED

Assigned unit of work completed successfully.

## FAILED

Work ended due to failure.

---

# 16. Agent Model

Example:

```go
type Agent struct {
    ID        string
    Provider  string
    ProjectID string

    ParentAgentID string
    SessionID     string

    Role        string
    DisplayName string

    Status       AgentStatus
    StatusSource StatusSource

    CurrentTaskID string
    CurrentTool   *ToolActivity

    StartedAt      *time.Time
    LastActivityAt *time.Time

    WorktreeID string
    Model      string
    Summary    string
}
```

The UI should not depend on provider-specific fields.

Provider-specific raw data can exist under:

```text
metadata
```

but should not drive generic UI behavior.

---

# 17. Claude Code Adapter

Claude Code is the first provider integration.

Use Claude Code lifecycle hooks as the primary deterministic event source where supported.

Relevant categories include:

- session lifecycle;
- user prompt submission;
- tool pre/post events;
- permission requests;
- subagent lifecycle;
- stop events;
- notifications.

The adapter should feature-detect available hooks because Claude Code can evolve.

Do not hardcode assumptions that every installation supports identical hook events forever.

## Adapter responsibilities

- install or generate project hook configuration;
- receive hook payloads locally;
- normalize provider events;
- preserve parent/subagent relationships where available;
- update current tool state;
- detect permission waits;
- detect completion/failure;
- redact sensitive tool input;
- never modify user Claude settings without explicit action.

## Installation approach

Craftlane should support:

```bash
craftlane init
```

that:

1. detects Claude Code;
2. detects current project settings;
3. explains what Craftlane needs to add;
4. creates Craftlane hook scripts/config only after explicit confirmation;
5. keeps modifications minimal and reversible.

Provide:

```bash
craftlane doctor
```

to verify:

- Claude executable;
- hook setup;
- callback server;
- project paths;
- permissions.

## Hook endpoint

Example internal endpoint:

```text
POST /internal/providers/claude/hooks
```

This endpoint should only accept local provider adapter traffic and use a provider-specific secret/token.

---

# 18. Provider SDK

Do not wait until a second provider exists to define a clean adapter boundary.

The first implementation can live under `internal/providers`, while the interface remains small and intentionally designed for later promotion into `pkg/provider`.

Initial Go interface:

```go
type Adapter interface {
    ID() string
    Detect(ctx context.Context, project ProjectContext) (Detection, error)
    Start(ctx context.Context, sink EventSink) error
    Stop(ctx context.Context) error
    Capabilities() Capabilities
}

type Capabilities struct {
    ObserveSessions    bool
    ObserveSubagents   bool
    ObserveTools       bool
    ObservePermissions bool

    StartAgent  bool
    SendMessage bool
    StopAgent   bool
}
```

`Start` should respect `context.Context` cancellation.

Adapters emit normalized events into Craftlane Core through a narrow event sink instead of mutating global state directly.

Providers must advertise capabilities rather than Craftlane assuming feature parity.

Do not expose a stable public Go provider API before at least two adapters have validated the abstraction.

When the interface is promoted publicly, document compatibility guarantees.

---

# 19. Spec Kit Integration

Craftlane should detect:

```text
.specify/
specs/
```

and known Spec Kit artifacts.

Current supported lifecycle should conceptually map to:

```text
constitution
→ specify
→ clarify
→ plan
→ checklist
→ tasks
→ analyze
→ implement
→ converge
```

Not every step creates a unique standalone artifact, so Craftlane must distinguish between:

- explicit artifact presence;
- inferred workflow state;
- latest observed command/result.

## Files to understand

At minimum:

```text
constitution.md
spec.md
plan.md
tasks.md
checklists/*.md
.specify/extensions.yml
```

Do not assume hardcoded paths without checking current project layout.

## Read-only first

Initial Spec Kit integration should be read-only:

- detect files;
- parse tasks;
- derive stages;
- display checklists;
- display dependencies;
- show parallel task markers.

Mutation comes later via deliberate actions.

## Task parser

Create a resilient task parser.

It should extract:

```go
type SpecTask struct {
    ID    string
    Title string

    Completed      bool
    Parallelizable bool

    Phase     string
    UserStory string

    Dependencies []string

    SourceFile string
    SourceLine int
}
```

Parser behavior:

- preserve unknown text;
- do not rewrite task files;
- fail gracefully;
- report unsupported task syntax clearly.

## Spec Kit phase UI

Example:

```text
Constitution     ✓
Specification    ✓
Clarification    ✓
Plan             ✓
Checklist        ✓
Tasks            ✓
Analyze          ✓
Implementation   ●
Convergence      ○
```

Craftlane should provide the most likely next action:

```text
Recommended next:
Run implementation for Phase 4
```

or:

```text
Recommended next:
Resolve analyze findings before implementation
```

Craftlane should not automatically skip required human review.

---

# 20. Task Dependency Graph

Craftlane should convert Spec Kit tasks into a graph.

Node:

```text
Task
```

Edge:

```text
dependency
```

Each task has:

```text
NOT_READY
READY
ASSIGNED
RUNNING
BLOCKED
COMPLETED
FAILED
```

A task is READY when:

- all declared dependencies are completed;
- required quality gates are satisfied;
- no known blocking state exists.

Parallel opportunity:

```text
READY tasks > 1
AND
available agents > 1
```

Craftlane can recommend:

```text
3 independent tasks
3 idle agents
Dispatch in parallel?
```

Never automatically dispatch in the first public versions unless the user explicitly enables an automation policy.

---

# 21. Git and Worktree Model

Git is central to safe parallel implementation.

Craftlane must understand:

- repository root;
- current branch;
- worktrees;
- branch-to-worktree relationships;
- dirty state;
- untracked files;
- merge conflicts.

## Write-agent isolation

Default recommendation:

> one write-heavy independent task = one Git worktree.

Read-only research/review agents may share a repository view if safe.

## Worktree record

```go
type Worktree struct {
    ID     string
    Path   string
    Branch string

    TaskID   string
    AgentIDs []string

    Dirty      bool
    Conflicted bool
}
```

## Worktree operations

Initial:

```text
detect
inspect
open
```

Later:

```text
create
assign
archive
remove
```

Removal must never happen automatically.

Before removing a worktree, Craftlane must verify:

- no active agent;
- no uncommitted changes;
- user confirmation.

---

# 22. Attention Engine

The human should not monitor every agent continuously.

Craftlane should generate attention items.

Priority levels:

```text
INFO
ACTION_REQUIRED
WARNING
CRITICAL
```

Examples:

## INFO

- task completed;
- tests passed;
- agent became idle.

## ACTION_REQUIRED

- permission requested;
- agent asks a question;
- merge decision required.

## WARNING

- repeated test failure;
- agent retry loop;
- task taking unusually long;
- dependency unresolved.

## CRITICAL

- merge conflict;
- provider process crashed;
- corrupt Spec Kit state;
- security-sensitive action blocked.

Attention record:

```go
type AttentionItem struct {
    ID         string
    Priority   AttentionPriority
    SourceType string
    SourceID   string
    Title      string
    Summary    string
    CreatedAt  time.Time
    ResolvedAt *time.Time
}
```

---

# 23. Command Dashboard

Build before Three.js.

Default dashboard must answer:

```text
What is working?
What needs me?
What is blocked?
What can run next?
```

Top status bar:

```text
Agents 6
Working 3
Waiting 1
Blocked 1
Idle 1

Tasks 14/21
Tests Healthy
Spec Implementation
```

Main sections:

```text
Attention
Active agents
Ready tasks
Blocked tasks
Spec Kit progress
Worktrees
Recent completions
```

Clicking an agent should show:

```text
Name
Provider
Role
Status
Current task
Current tool
Elapsed time
Last activity
Worktree
Changed files if available
Parent/subagents
Recent events
```

Raw logs should be secondary.

---

# 24. Command Palette

Command palette is a first-class feature.

Shortcut:

```text
Cmd/Ctrl + K
```

Commands:

```text
Open agent
Open task
Open worktree
Show attention
Show blocked tasks
Show idle agents
Open Spec Kit
Run Spec Kit Analyze
Start task
Dispatch task
Stop agent
Message agent
Run tests
Open project
Open settings
Switch Command/World/Focus mode
```

Every frequently used workflow should be reachable from keyboard.

---

# 25. Keyboard UX

Suggested defaults:

```text
Cmd/Ctrl+K   command palette
1            Command Mode
2            World Mode
3            Focus Mode

A            Agents
T            Tasks
S            Spec Kit
W            Worktrees
G            Global/God view
Esc          Close panel/back
```

Avoid stealing shortcuts while user is typing in:

- terminal;
- text field;
- editor;
- search input.

---

# 26. Three.js World Architecture

The 3D world is a projection of normalized Craftlane state.

Do not let Three.js own project state.

Flow:

```text
Craftlane Core State
        ↓
World Projection
        ↓
World Entity Model
        ↓
Three.js Scene
```

Example:

```ts
interface WorldAgentEntity {
  id: string
  status: AgentStatus
  position: Vector3Like
  targetStation?: string
  attentionLevel?: string
}
```

World projection should be deterministic.

A state change:

```text
agent WORKING
```

becomes:

```text
agent entity moves to assigned station
active indicator turns on
```

A state change:

```text
WAITING_PERMISSION
```

becomes:

```text
agent stops
yellow attention beacon appears
```

---

# 27. Three.js Scene

Initial world should be intentionally small.

Use a miniature software workshop.

Core stations:

```text
Blueprint
Dispatch
Development
Testing
Review
Integration
Idle Dock
Attention Area
```

Do not build:

- giant open world;
- realistic city;
- vehicle physics;
- long movement paths;
- complex character controller.

Initial goal:

> understand 5–20 agents and their work at a glance.

---

# 28. Camera Modes

## Command camera

Default.

- isometric/top-down;
- pan;
- zoom;
- click entity;
- no avatar.

## Explore camera

Optional entertainment mode.

- WASD;
- free exploration;
- smooth camera.

Never require Explore Mode for normal actions.

## Focus camera

When inspecting a selected agent/task/worktree, smoothly frame the selected entity while keeping the normal detail panel.

---

# 29. Three.js Interaction

Use raycasting for:

- robot selection;
- task object selection;
- station selection;
- workbench/worktree selection.

On click:

- select entity;
- show React/HTML side panel.

Do not render large bodies of text as 3D text.

Use HTML for:

- logs;
- code;
- diffs;
- forms;
- buttons;
- task metadata;
- Spec Kit documents.

---

# 30. 3D Performance Requirements

Craftlane is a developer tool, so rendering must not consume resources needed by builds and agents.

Targets for a normal development laptop:

```text
Visible, active world:
target 60 FPS

Background/low activity:
reduce animation frequency

Hidden tab:
pause scene updates where possible

Idle CPU:
keep minimal

Agent capacity v1:
50 visible agent/subagent entities without major degradation
```

Use:

- `InstancedMesh` for repeated environment/robot parts where useful;
- LOD for distant geometry if needed;
- simple materials;
- compressed/small GLTF assets;
- render-on-demand for static states;
- paused animation loops in hidden tabs;
- limited dynamic lights;
- asset lazy loading.

Performance tests should measure:

- main-thread frame time;
- memory growth;
- CPU while idle;
- event burst handling.

---

# 31. Accessibility

World Mode is optional.

All information shown in 3D must also be available through accessible HTML.

Requirements:

- keyboard navigation;
- visible focus;
- reduced-motion support;
- status not communicated by color alone;
- ARIA labels;
- no flashing effects;
- screen-reader-compatible dashboard.

A user who cannot or does not want to use Three.js must still be able to use Craftlane productively.

---

# 32. Agent Control

Observation comes before control.

Phase 1:

```text
observe
inspect
notify
```

Phase 2:

```text
message
stop
resume where supported
```

Phase 3:

```text
dispatch
spawn
assign worktree
```

Every action must check provider capabilities.

Example UI:

```text
Agent: backend-agent
Status: WAITING

[Message]
[Open Worktree]
[Open Logs]
[Stop]
```

If provider does not support messaging:

```text
Messaging is not supported by this provider.
```

Do not simulate capabilities that do not exist.

---

# 33. Task Dispatch

Dispatch is a core productivity feature.

Flow:

```text
Select task
   ↓
Validate dependencies
   ↓
Choose or create worktree
   ↓
Choose agent/provider
   ↓
Build task context
   ↓
Start/assign agent
   ↓
Observe execution
   ↓
Collect completion state
```

Before dispatch, validate:

- task READY;
- worktree not conflicted;
- provider available;
- concurrency limit not exceeded.

Context should include only necessary artifacts.

Do not blindly inject the entire repository or every Spec Kit file.

---

# 34. Parallel Dispatch

Craftlane identifies runnable independent work.

Example:

```text
T021 READY
T022 READY
T023 depends T021
T024 READY
```

With three available agents:

```text
Parallel opportunity

Agent 1 → T021
Agent 2 → T022
Agent 3 → T024
```

User clicks:

```text
Dispatch All
```

Craftlane then tracks:

```text
dispatch group
task assignments
worktrees
agent sessions
completion state
```

No automatic merge.

---

# 35. Concurrency Controls

Project configuration:

```json
{
  "maxConcurrentAgents": 4,
  "maxWriteAgents": 3
}
```

Craftlane should distinguish:

- read-only agent;
- write agent.

This allows safe use of more research/review agents without increasing write conflicts.

---

# 36. Productivity Recommendations

Initial recommendations should be rule-based, not AI-based.

Examples:

```text
2 agents idle and 3 tasks READY.
```

```text
T017 and T018 have no dependency and can run in parallel.
```

```text
Agent A is waiting for permission.
```

```text
T022 cannot start until T019 finishes.
```

```text
Two write agents are currently touching the same worktree.
```

Only after enough execution history exists should Craftlane introduce intelligent recommendations.

---

# 37. Metrics

Collect operational metrics locally.

Useful:

- task duration;
- agent active time;
- waiting time;
- permission waiting time;
- failure count;
- retry count;
- tasks completed;
- parallel execution duration;
- dependency wait duration;
- test outcomes.

Avoid using these as productivity scores:

- lines of code;
- prompt count;
- token count;
- commit count;
- hours online.

Craftlane should optimize flow, not reward activity volume.

---

# 38. Gamification

Keep gamification subtle.

Positive outcomes:

```text
Spec converged
All tests green
Zero unresolved high findings
Parallel batch completed
No merge conflicts
```

Do not add XP systems in early versions.

No leaderboard.

No productivity score.

Game elements should help the developer understand the project.

---

# 39. API Design

Version routes:

```text
/api/v1
```

Read endpoints:

```text
GET /api/v1/health
GET /api/v1/project
GET /api/v1/agents
GET /api/v1/agents/:id
GET /api/v1/tasks
GET /api/v1/tasks/:id
GET /api/v1/worktrees
GET /api/v1/speckit
GET /api/v1/attention
GET /api/v1/events
```

Mutating endpoints later:

```text
POST /api/v1/agents/:id/message
POST /api/v1/agents/:id/stop

POST /api/v1/tasks/:id/dispatch

POST /api/v1/worktrees

POST /api/v1/speckit/analyze
POST /api/v1/speckit/implement
POST /api/v1/speckit/converge
```

Do not expose raw shell execution.

---

# 40. WebSocket Design

Endpoint:

```text
/ws
```

Client handshake should include protocol version.

Server sends:

```text
snapshot
event
heartbeat
error
```

On connection:

```text
client connects
→ server sends current snapshot
→ client subscribes to new events
```

For reconnect:

```text
lastSequence
```

allows client to request missing events when history remains available.

---

# 41. Client State

Frontend state should have clear domains:

```text
connection
project
agents
tasks
worktrees
speckit
attention
world
ui
```

Do not store provider raw data everywhere.

Normalize at the server/core boundary.

---

# 42. File Watching

Watch only required project paths.

Examples:

```text
.specify/
specs/
.git/
.craftlane/
```

Avoid watching the entire repository recursively unless needed.

Use debounce/coalescing for bursts.

Do not parse files on every filesystem event individually when several updates arrive together.

---

# 43. Git Status Watching

Git status can be expensive in large repositories.

Use:

- debouncing;
- bounded polling;
- event-triggered refresh;
- cached last result.

Do not execute `git status` every animation frame or every UI render.

---

# 44. Logging

Logs:

```text
~/.craftlane/logs/
```

Structured logging fields:

```text
timestamp
level
module
projectId
agentId
taskId
message
```

Levels:

```text
debug
info
warn
error
```

CLI:

```bash
craftlane doctor
```

should report relevant diagnostic log locations.

Secret redaction must happen before persistence.

---

# 45. Failure Handling

Craftlane should degrade gracefully.

Examples:

## Claude unavailable

```text
Claude Code not detected.
Craftlane dashboard remains available.
```

## Spec Kit unavailable

```text
Spec Kit not detected.
Agent/worktree monitoring still works.
```

## Three.js unsupported

Automatically use Command Mode.

## Hook failure

Show:

```text
Claude integration degraded.
Run: craftlane doctor
```

## Git unavailable

Disable Git/worktree features but keep provider observability where possible.

---

# 46. Configuration

Example `.craftlane/config.json`:

```json
{
  "version": 1,
  "port": 7331,
  "openBrowser": true,
  "defaultMode": "command",

  "agents": {
    "maxConcurrent": 4,
    "maxConcurrentWrite": 3
  },

  "world": {
    "enabled": true,
    "reducedMotion": false
  },

  "speckit": {
    "enabled": true
  }
}
```

Global config supplies defaults.

Project config overrides safe project-specific settings.

Environment variables can override runtime values.

---

# 47. Installation Experience

End users should install a native Craftlane executable.

Primary public installation should eventually be:

macOS/Linux:

```bash
curl -fsSL <official-install-script-url> | sh
```

Windows:

```powershell
# Official PowerShell installer from the Craftlane GitHub repository.
```

The actual release documentation must use the repository's real verified URLs before publication.

Users must also be able to manually download a binary from GitHub Releases.

Developer installation for users with Go installed:

```bash
go install <craftlane-module>/cmd/craftlane@latest
```

Do not make Node.js/npm/pnpm a production requirement.

First launch:

```text
Craftlane

✓ Git repository detected
✓ Claude Code detected
✓ Spec Kit detected

Claude integration is not installed for this project.

Craftlane needs project hooks to observe Claude sessions.

[Install integration]
[Continue without Claude integration]

Workspace:
http://127.0.0.1:7331
```

Do not silently modify project settings.

---

# 48. Public Open-Source Repository

Required files before first public release:

```text
README.md
LICENSE
CONTRIBUTING.md
CODE_OF_CONDUCT.md
SECURITY.md
CHANGELOG.md
```

Recommended license:

```text
Apache-2.0
```

Rationale:

- permissive;
- friendly to adoption;
- explicit patent grant.

This can be revisited before first public release.

Release source and compiled binaries should both be published from the same GitHub repository.

---

# 49. README Structure

README should explain the product in under 10 seconds.

Suggested opening:

```text
# Craftlane

Your AI development team, visible.

Craftlane is an open-source, local-first development environment for
orchestrating AI coding agents.

- Monitor parallel agents
- Visualize Spec Kit workflows
- Coordinate dependency-aware tasks
- Isolate work with Git worktrees
- Use an interactive Three.js workspace
- Run entirely on localhost
- Install as a native single executable
```

Quick Start should show:

```bash
craftlane
```

after the platform-specific installation step.

README sections:

```text
Demo
Why Craftlane
Quick Start
Installation
How It Works
Features
Spec Kit
Agent Providers
Security / Local-First
Architecture
Roadmap
Contributing
License
```

---

# 50. CI

GitHub Actions should separately validate Go and web code, then test the integrated executable.

Required jobs:

```text
Go format/vet
Go unit tests
Go race tests
Go integration tests

web install
web lint
web typecheck
web unit tests
web production build

protocol compatibility tests
embedded-assets build

Playwright browser tests
CLI smoke tests
release binary smoke tests
```

OS matrix for integration/release smoke tests:

```text
macOS
Ubuntu
Windows
```

CI should test the pinned Go version declared by the repository and the supported frontend Node version used for builds.

A production build should conceptually execute:

```text
pnpm --dir web install --frozen-lockfile
pnpm --dir web build
go test ./...
go build ./cmd/craftlane
```

Release CI must test the exact generated binaries, not only `go run`.

For each supported platform, verify:

```text
craftlane --version
craftlane doctor
craftlane start --no-open
health endpoint reachable
embedded web UI reachable
clean shutdown
```

---

# 51. Release Process

Use semantic versioning.

Early versions:

```text
0.1.x
0.2.x
...
```

Recommended release flow:

```text
tag
↓
GitHub Actions
↓
build web once
↓
embed web assets
↓
cross-compile native binaries
↓
checksum artifacts
↓
GitHub Release
```

Use a Go-oriented release tool such as GoReleaser if it simplifies reproducible cross-platform releases; keep the configuration understandable and avoid unnecessary release complexity.

Later publish package-manager manifests/taps from the same release artifacts.

Release gates:

```text
all CI green
Go race tests green
web production build green
embedded asset smoke test green
Windows/macOS/Linux binary smoke tests green
checksums generated
README current
migration notes present if config/protocol changed
```

---

# 52. Protocol Compatibility

Craftlane UI and server may eventually update independently.

Version the protocol from the beginning.

Rules:

- new optional fields are preferred;
- do not silently redefine existing event semantics;
- clients should ignore unknown optional event payload fields;
- incompatible protocol changes require schema version bump.

---

# 53. Testing Strategy

## Unit tests

Core state:

```text
agent transitions
task readiness
dependency calculation
attention rules
world projection
config merging
event sequence
redaction
```

## Adapter tests

Claude hook fixtures:

```text
session start
tool start
tool complete
permission
subagent start
subagent stop
failure
session stop
```

Do not test provider parsing from remembered terminal text.

Use captured fixture payloads.

## Spec Kit tests

Fixtures:

```text
valid tasks
parallel tasks
completed tasks
multiple phases
malformed tasks
missing plan
checklists
convergence tasks
```

## Git tests

Create temporary repositories and test:

```text
branch detection
worktree detection
dirty worktree
conflict state
create worktree
safe deletion checks
```

## API tests

Test:

```text
loopback binding
token protection
CORS/origin restrictions
path validation
event reconnect
```

## Browser tests

Playwright:

```text
dashboard loads
agent state updates
attention appears
command palette
task selection
Spec Kit view
mode switching
```

## Three.js tests

Do not rely only on pixel snapshots.

Test:

```text
world projection logic
entity creation/removal
selection mapping
camera mode state
animation state transitions
```

Add a small number of visual screenshots for major layouts.

---

# 54. Performance Testing

Add reproducible scenarios:

```text
1 project / 5 agents
1 project / 20 agents
1 project / 50 agents
10 worktrees
500 task events
high-frequency tool event burst
```

Measure:

```text
event processing latency
WebSocket backlog
React render count
Three.js frame time
memory growth
idle CPU
```

Performance regression tests should be added before large visual complexity.

---

# 55. Observability for Craftlane Itself

Expose local diagnostics:

```bash
craftlane status
craftlane doctor
```

Health API should show:

```json
{
  "status": "ok",
  "project": true,
  "git": true,
  "speckit": true,
  "providers": {
    "claude": "connected"
  }
}
```

Provide a support bundle command later:

```bash
craftlane doctor --report
```

Report must redact:

- source content;
- credentials;
- personal filesystem details where possible.

---

# 56. Implementation Milestones

## Milestone 0 — Product Constitution and Repo Bootstrap

Goal:

Create an implementation foundation without product behavior.

Deliver:

- repository;
- Apache-2.0 license;
- Go module;
- `cmd/craftlane`;
- `internal/` package layout;
- web TypeScript/React/Vite application;
- pnpm used only for web development/build;
- protocol schema directory;
- Makefile or equivalent contributor commands;
- CI;
- Go and web test tooling;
- formatting/linting;
- README skeleton;
- CONTRIBUTING;
- SECURITY;
- initial Spec Kit constitution.

Development exit criteria:

```text
go test ./...
go vet ./...
pnpm --dir web install
pnpm --dir web build
go build ./cmd/craftlane
```

all work from a clean clone.

A production binary built after the web build must serve its embedded frontend successfully.

---

## Milestone 1 — Local Runtime and Protocol

Feature:

```text
001-foundation-runtime
```

Deliver:

- native Go CLI;
- project detection;
- local Go HTTP server;
- localhost binding;
- runtime token;
- WebSocket;
- versioned event protocol;
- protocol schemas;
- file-backed storage;
- embedded web shell;
- health endpoint;
- graceful shutdown through contexts/signals;
- `craftlane start`;
- `craftlane stop`;
- `craftlane open`;
- `craftlane status`;
- `craftlane doctor`.

UI can be minimal.

Exit criteria:

After installing or building the `craftlane` executable:

```bash
cd any-git-project
craftlane
```

opens a local page showing detected repository information.

The released executable must not require Node.js, npm, or pnpm at runtime.

No cloud dependency.

---

## Milestone 2 — Claude Observability

Feature:

```text
002-claude-observability
```

Deliver:

- Claude detection;
- hook setup;
- hook receiver;
- normalized agent state;
- subagent relationships;
- tool activity;
- permission waiting;
- session stop/failure;
- agent event stream;
- redaction;
- fixture tests.

Exit criteria:

Running Claude Code and subagents causes Craftlane state to update without refreshing the browser.

Must support:

```text
working
waiting permission
idle/inferred idle
completed
failed
```

with documented confidence/source.

---

## Milestone 3 — Command Dashboard

Feature:

```text
003-command-dashboard
```

Deliver:

- real dashboard;
- agent list;
- status summary;
- attention queue;
- detail panel;
- command palette;
- keyboard shortcuts;
- connection/reconnect behavior;
- accessible UI;
- responsive layout.

Exit criteria:

A developer can monitor at least 10 agents without needing Three.js.

This milestone must already be useful in daily work.

---

## Milestone 4 — Spec Kit Integration

Feature:

```text
004-speckit-integration
```

Deliver:

- Spec Kit detection;
- artifact watcher;
- spec/plan/tasks/checklist model;
- task parser;
- phase detection;
- task completion state;
- dependency model;
- parallelizable markers;
- Spec Kit UI;
- recommended next step.

Exit criteria:

Opening Craftlane in a Spec Kit project shows the current SDD state and task graph accurately.

No Spec Kit files are mutated automatically.

---

## Milestone 5 — Task Graph and Worktrees

Feature:

```text
005-task-graph-and-worktrees
```

Deliver:

- Git worktree discovery;
- branch state;
- dirty/conflict indicators;
- task/worktree associations;
- task readiness calculation;
- parallel opportunity detection;
- worktree UI;
- safe worktree create flow;
- task assignment metadata.

Exit criteria:

Craftlane can answer:

```text
Which tasks can run now?
Which worktree belongs to which task?
Which tasks are blocked?
Which agents are free?
```

---

## Milestone 6 — Three.js Workspace

Feature:

```text
006-threejs-workspace
```

Deliver:

- world projection package;
- Three.js scene;
- robot agents;
- workbenches/worktrees;
- stations;
- task routes;
- status animations;
- camera modes;
- selection;
- detail panels;
- performance controls;
- reduced-motion mode.

Exit criteria:

3D state always matches the normal dashboard.

Turning World Mode off loses no functionality.

---

## Milestone 7 — Agent Control and Dispatch

Feature:

```text
007-agent-control-and-dispatch
```

Deliver:

- provider capability model;
- agent stop where supported;
- send message where supported;
- task dispatch;
- isolated worktree assignment;
- concurrency limits;
- dispatch groups;
- parallel dispatch recommendations;
- action audit events;
- failure recovery.

Exit criteria:

A user can select independent ready tasks and safely dispatch them to available agents.

No automatic merge.

---

## Milestone 8 — Hardening and Public Release

Feature:

```text
008-hardening-packaging-release
```

Deliver:

- security audit;
- cross-platform testing;
- package smoke tests;
- native cross-platform binary packaging;
- clean first-run UX;
- installation docs;
- CLI help;
- error handling;
- performance budgets;
- public roadmap;
- GitHub templates;
- release automation.

Exit criteria:

A new developer can install the platform binary, run:

```bash
craftlane
```

inside a supported project and reach a working local dashboard with clear setup guidance.

Target public release:

```text
v0.1.0
```

---

## Milestone 9 — Provider SDK

Feature:

```text
009-provider-sdk
```

Deliver:

- documented provider interface;
- provider fixture kit;
- capability negotiation;
- reference second adapter;
- provider authoring docs.

Exit criteria:

A contributor can add another agent provider without changing Craftlane Core.

---

## Milestone 10 — Insights and Optimization

Feature:

```text
010-insights-and-optimization
```

Deliver:

- local execution history;
- wait-time analysis;
- retry detection;
- idle-agent opportunities;
- dependency bottleneck detection;
- rule-based workflow suggestions;
- insights dashboard.

Do not introduce opaque AI scoring at first.

Exit criteria:

Craftlane can identify actionable workflow inefficiencies from real local execution history.

---

# 57. v0.1 Scope

Public v0.1 should be intentionally smaller than the entire roadmap.

Required:

- local-first runtime;
- CLI;
- browser dashboard;
- Claude Code observability;
- parent/subagent display;
- task/current-tool visibility;
- attention states;
- Git/worktree detection;
- Spec Kit detection;
- basic task graph;
- basic Three.js world;
- command palette;
- no Craftlane cloud account.

Optional if stable:

- basic stop/message controls.

Not required for v0.1:

- multi-provider support;
- automatic task dispatch;
- remote access;
- mobile app;
- plugin marketplace;
- cloud sync;
- AI-based routing;
- automated merging;
- team collaboration;
- complex 3D exploration.

---

# 58. v1.0 Definition

Craftlane v1.0 should mean:

- stable local runtime;
- stable event protocol;
- reliable Claude integration;
- at least one additional provider;
- stable Spec Kit integration;
- safe task dispatch;
- safe worktree isolation;
- command dashboard;
- optional Three.js world;
- provider SDK documentation;
- robust installation;
- cross-platform support;
- clear security model;
- backwards-compatible config migration.

---

# 59. Explicit Non-Goals

Do not turn Craftlane into:

- a full code editor;
- a replacement for VS Code/IntelliJ;
- a Git hosting service;
- a CI/CD platform;
- a cloud IDE;
- a project management SaaS;
- a social game;
- a cryptocurrency/token system;
- a productivity surveillance product.

Avoid scope creep.

---

# 60. Decision Rules for Future Features

Before implementing a feature, ask:

## Does it save time?

If no, do not prioritize it.

## Does it improve understanding?

If no, do not add it merely because it looks impressive.

## Does it require cloud infrastructure?

If yes, confirm it cannot reasonably remain local-first.

## Does it duplicate an existing IDE/terminal feature?

If yes, integrate or link instead of rebuilding unless Craftlane gains meaningful orchestration value.

## Can it be done from keyboard?

Frequently used actions should have keyboard access.

## Is the 3D version slower?

Keep a faster standard UI path.

---

# 61. Architecture Constraints for Spec Kit

When generating implementation tasks, preserve these constraints.

1. Core state must not live in React.
2. Three.js must not become the source of truth.
3. Provider-specific event formats must terminate at adapters.
4. Browser must never receive secrets.
5. UI must never execute arbitrary shell strings.
6. Spec Kit integration begins read-only.
7. localhost binding is the default.
8. state-changing API endpoints require runtime protection.
9. write-heavy parallel tasks should use isolated worktrees.
10. no auto-merge in early versions.
11. no telemetry by default.
12. every 3D state must have a 2D/accessibility equivalent.
13. provider capabilities must be explicit.
14. unsupported provider features must degrade gracefully.
15. protocol schemas must be versioned.
16. Go Core must remain independent of React/Three.js implementation details.
17. released binaries must not require Node.js at runtime.
18. frontend assets must be embedded or otherwise shipped with the native executable.
19. long-running goroutines must be cancellable through explicit lifecycle ownership.
20. provider adapters must not mutate canonical state directly.

---

# 62. Suggested Domain Types

Use this only as an architectural starting point; implementation may refine names.

Canonical domain models live in Go.

```go
type ProjectID string
type AgentID string
type TaskID string
type WorktreeID string

type Project struct {
    ID             ProjectID
    Name           string
    Root           string
    RepositoryRoot string
}

type Agent struct {
    ID        AgentID
    ProjectID ProjectID
    Provider  string

    ParentAgentID AgentID

    DisplayName string
    Role        string

    Status       AgentStatus
    StatusSource StatusSource

    TaskID     TaskID
    WorktreeID WorktreeID

    CurrentTool *ToolActivity

    StartedAt      *time.Time
    LastActivityAt *time.Time
}

type Task struct {
    ID    TaskID
    Title string

    State TaskState

    Dependencies   []TaskID
    Parallelizable bool

    AssignedAgentIDs []AgentID
    WorktreeID       WorktreeID

    Source *TaskSource
}

type Worktree struct {
    ID        WorktreeID
    ProjectID ProjectID

    Path   string
    Branch string

    Dirty      bool
    Conflicted bool

    TaskID TaskID
}

type ToolActivity struct {
    Name      string
    StartedAt time.Time
    Summary   string
}
```

Do not expose these Go structs directly as a promise to the browser.

The browser consumes versioned wire DTOs defined by the protocol schema.

Frontend TypeScript types should be generated from, or validated against, that wire contract.

---

# 63. Suggested State Machines

## Agent

```text
DISCOVERED
   ↓
STARTING
   ↓
WORKING
   ├──────────────→ WAITING_PERMISSION
   │                      │
   │                      ▼
   │                  WORKING
   │
   ├──────────────→ WAITING
   │                      │
   │                      ▼
   │                  WORKING
   │
   ├──────────────→ IDLE
   │                      │
   │                      ▼
   │                  WORKING
   │
   ├──────────────→ BLOCKED
   │                      │
   │                      ▼
   │                  WORKING
   │
   ├──────────────→ COMPLETED
   ├──────────────→ FAILED
   └──────────────→ STOPPED
```

Do not force every provider through every transition.

## Task

```text
NOT_READY
   ↓ dependencies satisfied
READY
   ↓ assign
ASSIGNED
   ↓ agent begins
RUNNING
   ├────→ BLOCKED
   │        │
   │        └────→ RUNNING
   ├────→ COMPLETED
   └────→ FAILED
```

---

# 64. Suggested UI Information Hierarchy

Always show highest-value information first.

Priority:

```text
1. Attention required
2. Blocked work
3. Active agents
4. Ready parallel work
5. Spec Kit stage
6. Worktrees
7. Recent history
8. Raw logs
```

Do not make users scan animations to find a problem.

---

# 65. First-Run Experience

Target flow:

```text
$ craftlane

Craftlane detected:
✓ Git repository
✓ Claude Code
✓ Spec Kit

Claude integration is not installed for this project.

Craftlane can add local lifecycle hooks so it can show
agent and subagent activity.

Install integration? [Y/n]

✓ Integration installed

Craftlane running:
http://127.0.0.1:7331
```

Browser opens to Command Mode.

First screen:

```text
No agents running.

Start Claude Code normally in this repository.
Craftlane will detect activity automatically.
```

Avoid forcing users to learn Craftlane-specific agent commands before they see value.

---

# 66. Daily User Experience

Ideal workflow:

```text
Open repository
↓
craftlane
↓
Dashboard appears
↓
Developer starts normal Claude workflow
↓
Craftlane detects main agent and subagents
↓
Agent statuses update
↓
Spec tasks show readiness/dependencies
↓
Craftlane highlights one permission request
↓
Developer resolves it
↓
Craftlane detects two independent READY tasks
↓
Developer chooses parallel dispatch
↓
Agents work in isolated worktrees
↓
Tests/review complete
↓
Spec Kit convergence is recommended
↓
Developer reviews and merges
```

The user should spend less time managing terminal windows.

---

# 67. Reference Projects and Boundaries

## Three.js

Use Three.js documentation as the rendering reference for:

- renderer;
- scene/camera;
- raycasting;
- animation;
- instancing;
- level of detail;
- asset loading.

Reference:

https://threejs.org/

## ORCA

ORCA is an architectural/product reference for ideas such as:

- parallel agent visibility;
- worktree-oriented workflows;
- local developer experience;
- canonical agent state;
- CLI + UI coexistence.

Reference:

https://github.com/stablyai/orca

ORCA is currently MIT-licensed, but Craftlane should remain an independent codebase and brand.

Do not copy ORCA source or assets unnecessarily.

If code is ever reused, review the applicable license and preserve required notices.

## GitHub Spec Kit

Use official Spec Kit behavior as the source of truth for the SDD lifecycle.

References:

https://github.com/github/spec-kit
https://github.com/github/spec-kit/blob/main/docs/reference/agentic-sdd.md

Current conceptual workflow:

```text
constitution
→ specify
→ clarify
→ plan
→ checklist
→ tasks
→ analyze
→ implement
→ converge
```

## Claude Code

Use official Claude Code documentation for:

- hooks;
- subagents;
- agent teams;
- provider lifecycle;
- CLI behavior.

References:

https://code.claude.com/docs/
https://docs.anthropic.com/en/docs/claude-code/

Do not depend on undocumented Claude terminal UI strings when structured lifecycle data is available.

---

# 68. Recommended Spec Kit Execution Strategy

Do not send this entire project into one `/speckit.implement`.

Instead:

## Step 1 — Project constitution

Create Craftlane constitution around:

```text
local-first
productivity-first
security
accessibility
testability
provider abstraction
backwards-compatible protocol
open-source maintainability
```

## Step 2 — Bootstrap Feature 001

Run:

```text
/speckit.specify
```

for:

```text
Craftlane local runtime, CLI, local server, web shell,
event protocol, storage abstraction, project detection,
localhost security, and diagnostics.
```

Then:

```text
/speckit.clarify
/speckit.plan
/speckit.checklist
/speckit.tasks
/speckit.analyze
```

Fix all high-severity artifact problems before implementation.

Then implement in stages:

```text
Setup + foundational
Core runtime
CLI
Web shell
Security
Tests
```

Finally:

```text
/speckit.converge
```

Repeat implementation/convergence until clean.

## Step 3 — Repeat per milestone

Create one feature spec per roadmap milestone.

This keeps:

- context manageable;
- reviews smaller;
- regressions easier to isolate;
- parallel work safer;
- Spec Kit artifacts useful.

---

# 69. Recommended Feature Prompts

These can be used as inputs to `/speckit.specify`.

## 001 Foundation

```text
Build the local-first Craftlane runtime.

A user inside a Git repository can run Craftlane from the command line,
start a localhost-only service, open a browser dashboard, see the detected
project, receive real-time events through a versioned protocol, persist
local runtime state, and diagnose installation problems.

The system must require no Craftlane account or cloud backend.
```

## 002 Claude Observability

```text
Add Claude Code observability to Craftlane.

Craftlane should detect Claude sessions and subagents using supported
Claude lifecycle integrations, normalize their activity into provider-
agnostic agent states, display current tool activity and permission waits,
and stream updates to the local UI without exposing secrets.
```

## 003 Command Dashboard

```text
Build the productivity-first Craftlane command dashboard.

A developer should be able to understand active agents, blocked work,
attention requests, idle agents, current tasks, and recent completions
without opening multiple terminals. Frequently used actions must be
keyboard-accessible.
```

## 004 Spec Kit

```text
Integrate GitHub Spec Kit into Craftlane as a read-only first-class workflow.

Craftlane should detect Spec Kit artifacts, show the current SDD lifecycle
stage, parse tasks and completion state, understand dependency and parallel
markers, and recommend the correct next workflow action without mutating
Spec Kit artifacts automatically.
```

## 005 Worktrees

```text
Add Git worktree and task graph coordination to Craftlane.

Craftlane should detect worktrees, branches, dirty/conflict state, map
tasks to worktrees, determine which tasks are READY or BLOCKED, and identify
parallel work opportunities safely.
```

## 006 World

```text
Add an optional Three.js spatial workspace to Craftlane.

The world must project existing Craftlane state into a miniature software
workshop with agents, worktrees, tasks, development/testing/review stations,
status indicators, dependency routes, selection, and camera controls.

All functionality must remain available outside the 3D view.
```

## 007 Dispatch

```text
Add controlled task dispatch to Craftlane.

Users should be able to assign dependency-ready tasks to supported agents,
isolate write work using Git worktrees, respect concurrency limits, monitor
dispatch groups, stop or message agents where supported, and explicitly
launch safe parallel batches.

Craftlane must not automatically merge code.
```

---

# 70. Acceptance Criteria for the Overall Product

Craftlane is successful when:

## Ease of use

A new developer can go from repository to dashboard with one platform installation step followed by:

```bash
craftlane
```

plus only necessary provider integration confirmation.

## Local-first

Normal operation requires no Craftlane server outside the machine.

## Visibility

Developer can identify within seconds:

- working agents;
- blocked agents;
- waiting agents;
- idle agents;
- current tasks;
- available parallel work.

## Spec Kit

Developer can understand the project's SDD stage without manually searching multiple files.

## Parallelism

Craftlane can explain which tasks are safe to run concurrently.

## Safety

Write tasks can be isolated through worktrees.

## Productivity

Common operations are accessible via command palette/keyboard.

## 3D

World Mode improves understanding but is never mandatory.

## Open source

A contributor can clone, build, test, and understand the architecture from documented commands.

---

# 71. Product Success Metrics

Do not optimize for vanity metrics.

Useful local/user-facing measurements:

```text
time to identify blocked agent
time to resolve permission wait
number of terminal switches avoided
parallel work utilization
dependency waiting time
failed/retried task count
setup success rate
startup time
UI responsiveness
```

For public project health:

```text
installation success
GitHub issues resolution
contributor setup success
provider compatibility
release regression rate
```

Do not require telemetry to measure local metrics.

---

# 72. Major Risks

## Risk: 3D becomes distracting

Mitigation:

- Command Mode default;
- reduced motion;
- World Mode optional;
- strict performance budgets.

## Risk: Claude integration changes

Mitigation:

- provider adapter;
- capability detection;
- fixture tests;
- version diagnostics.

## Risk: Spec Kit changes format

Mitigation:

- isolated parser package;
- defensive parsing;
- fixtures from multiple versions;
- no destructive rewriting.

## Risk: task dispatch causes code conflicts

Mitigation:

- dependency validation;
- worktree isolation;
- write concurrency limits;
- no auto merge.

## Risk: localhost attack surface

Mitigation:

- loopback binding;
- strict origin/host rules;
- runtime token;
- typed action API;
- path restrictions.

## Risk: cross-platform native distribution becomes difficult

Mitigation:

- prefer pure-Go runtime dependencies in v0.1;
- avoid CGO by default;
- keep storage file-backed behind an interface;
- embed static frontend assets;
- cross-compile in CI;
- smoke-test real binaries across macOS, Linux, and Windows.

## Risk: agent status inference is inaccurate

Mitigation:

- mark explicit vs inferred state;
- surface confidence/source internally;
- do not call inactivity "blocked."

---

# 73. Deferred Ideas

Only consider after core product is proven.

```text
Desktop wrapper
Remote access
Mobile companion
Multi-user team mode
Cloud sync
Provider marketplace
Agent templates
Custom world themes
Plugin marketplace
Advanced analytics
AI task routing
Automatic worktree lifecycle
Automated PR creation
CI integrations
GitHub issue synchronization
```

---

# 74. Final Architecture Principle

Craftlane should remain understandable as:

```text
Agents
   ↓
Provider Adapters
   ↓
Normalized Craftlane State
   ↓
Task + Spec + Git Intelligence
   ↓
Local API / Events
   ↓
CLI + Dashboard + Three.js World
```

The 3D world is a projection.

The provider adapters are replaceable.

Spec Kit is a first-class development model.

Git worktrees provide safe parallelism.

The user remains in control.

---

# 75. Definition of Done for This Master Plan

Before implementation begins, the first Spec Kit feature must preserve the following decisions:

- Name: Craftlane
- Open source
- Local-first
- GitHub published
- Installable locally
- Browser-based localhost UI
- CLI
- Claude-first but provider-agnostic architecture
- Spec Kit first-class integration
- Git worktree-aware
- Command Mode before World Mode
- Three.js optional
- Productivity before entertainment
- No mandatory account
- No cloud backend required
- No telemetry by default
- No auto merge
- Security around localhost actions
- Cross-platform intent
- Feature-by-feature Spec Kit execution

Any generated implementation plan that contradicts one of these decisions should be corrected before `/speckit.implement`.

---

# 76. Suggested Immediate Next Step

Use this master plan as context, then create only:

```text
001-foundation-runtime
```

first.

Do **not** start Three.js implementation yet.

The first engineering objective is:

> A developer can enter any Git repository, run the native Craftlane executable locally, open the localhost dashboard served by the Go runtime, see the detected project, and receive normalized real-time events through a secure versioned local protocol.

Once that foundation is stable, add Claude observability.

Once Claude observability is useful in a normal 2D dashboard, add Spec Kit.

Only after those contracts are reliable should Craftlane render them as a 3D world.

That order protects Craftlane from becoming a visually impressive demo with an unreliable developer workflow.
