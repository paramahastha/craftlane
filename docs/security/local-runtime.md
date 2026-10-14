# Local observation security

Observation uses a distinct owner-only token stored with the project runtime state. Callback requests must use the exact runtime Host, no Origin or the exact runtime Origin, a supported provider path, a bounded JSON body, and the observation token. Project identity is checked before the reducer can mutate state.

Provider setup is explicit and reversible. Claude settings and Codex hooks are structurally merged around Craftlane markers; OpenCode uses one owned generated plugin. Codex and OpenCode artifacts may be tracked, so setup warns without changing ignore rules. Normalized state is redacted and bounded; raw provider payloads are never persisted or sent to the browser.
