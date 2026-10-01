# CLI Contract: Local Runtime Foundation

## Commands

| Command | Success behavior | Failure behavior |
|---|---|---|
| `craftlane` | Equivalent to `craftlane start`; detects the current Git project, reports local origin, then opens the browser best-effort. | Nonzero exit with a corrective message. |
| `craftlane start` | Starts or reuses the verified runtime for the current Git project. | Nonzero exit for invalid project, unrecoverable stale state, or explicit occupied port. |
| `craftlane start --no-open` | Starts/reuses runtime and prints local origin without launching a browser. | Same as `start`. |
| `craftlane start --headless` | Starts/reuses runtime with no browser action and terminal-safe output. | Same as `start`. |
| `craftlane start --port <port>` | Requests that exact valid loopback port. | Fails rather than silently choosing another port when the requested port is unavailable. |
| `craftlane stop` | Authenticates a typed graceful shutdown request for the current project's verified runtime. | Nonzero exit when no matching live runtime exists. |
| `craftlane restart` | Gracefully stops a verified runtime then starts a fresh runtime for the current project. | Nonzero exit if either stage fails, with current state reported. |
| `craftlane open` | Opens the verified running runtime origin. | Nonzero exit and next action when no live runtime exists. |
| `craftlane status` | Prints runtime state, project identity, origin, and safe capability health. | Nonzero exit only for invalid invocation or unreadable/corrupt state; a stopped runtime is a valid reported state. |
| `craftlane doctor` | Prints typed diagnostic checks and corrective actions. | Uses nonzero exit when a required check fails; optional integration absence is a warning. |

## General rules

- Commands resolve the project from the current working directory using Git; non-Git directories are rejected in this milestone.
- Human-readable output is mandatory. Exit status is script-friendly: `0` for successful command completion, nonzero for an error or failed required diagnostic.
- Browser opening never determines startup success.
- Commands do not display control secrets, credential-like values, or raw provider environment values.
- No Foundation command auto-installs integrations, mutates Spec Kit files, deletes worktrees, merges code, or executes a caller-supplied shell command.

