# Quickstart: Validate Local Runtime Foundation

This guide validates the completed Foundation runtime from a contributor checkout. See [data-model.md](./data-model.md) for state ownership and [contracts](./contracts/) for public behavior.

## Prerequisites

- A supported Go toolchain and the pinned Node.js/pnpm contributor toolchain.
- Git available on `PATH`.
- A modern desktop browser for the dashboard check.

## Build the contributor artifact

```bash
pnpm --dir web install --frozen-lockfile
make test
make build
```

Expected outcome: the Make targets build the web assets before compiling any Go package that embeds them, and `./bin/craftlane` is a self-contained executable with embedded dashboard assets. The final executable must not need Node.js or pnpm to run.

## Validate start, dashboard, status, and stop

Create a disposable Git project outside the Craftlane checkout, then run:

```bash
git init /tmp/craftlane-foundation-fixture
cd /tmp/craftlane-foundation-fixture
<craftlane-checkout>/bin/craftlane start --no-open
<craftlane-checkout>/bin/craftlane status
<craftlane-checkout>/bin/craftlane doctor
```

Expected outcome:

1. `start --no-open` reports a `http://127.0.0.1:<port>` origin and does not launch a browser.
2. Opening that origin shows the project name, root, and runtime health in semantic visible content.
3. `status` reports the same project and a running state.
4. `doctor` reports required Git/runtime/dashboard checks and marks unavailable optional integrations as warnings rather than blocking failures.

Then stop and recheck:

```bash
<craftlane-checkout>/bin/craftlane stop
<craftlane-checkout>/bin/craftlane status
```

Expected outcome: stop completes gracefully and status reports no matching live runtime (or an explicit stopped state) without presenting stale data as live.

## Validate local security

1. Confirm the reported origin uses `127.0.0.1`, never a LAN address.
2. Send requests with an invalid Host and a WebSocket connection with a missing or different Origin; both must be rejected according to the [HTTP](./contracts/http-api.yaml) and [WebSocket](./contracts/websocket.md) contracts.
3. Send a shutdown request without `X-Craftlane-Token` and with an invalid token; both must be rejected. Confirm a secret from a stopped, interrupted, stale, or restarted runtime is also rejected. Confirm the valid control secret is absent from CLI output, health, event responses, and diagnostics.
4. Publish test values shaped like credentials and verify their raw values do not occur in persisted event segments, dashboard data, or doctor output.

## Validate continuity and recovery

1. Connect a test client and record its latest event sequence.
2. Disconnect it, trigger project lifecycle events, and reconnect with that sequence.
3. Confirm each retained missed event is delivered once in ascending sequence.
4. Confirm the default retained window contains the latest 10,000 events and cannot be configured below 1,000. Repeat with a sequence older than retention; confirm `resync_required` precedes a fresh snapshot.
5. Leave a stale runtime record or stop the runtime unexpectedly, then run `start`; confirm it identifies stale state, recovers safely, and reports the fresh origin.

## Automated acceptance coverage

Run the feature's Go tests with the race detector, web tests, and Playwright suite. The browser suite must launch the built executable against a temporary Git repository—not only a Vite development server. Verify dashboard readiness within 10 seconds across 20 clean starts with at least 19 successes, 100 reconnect cycles with no missed or duplicate retained events, and a hermetic offline startup. Finally, run a five-developer task-completion study and require at least four developers to start, find status, and stop a workspace without help. Record results before release.

## Validation record (2026-10-01)

- `make test`, `make build`, `make lint`, `go test -race ./...`, and the built-binary Playwright dashboard test passed on macOS arm64.
- The automated start exercise completed 20/20 clean detached starts within 10 seconds, with loopback-only origins and a successful authenticated stop after each start.
- The automated WebSocket exercise completed 100/100 reconnects with the retained event delivered once per connection using the same identity and sequence.
- A release binary started with unreachable `HTTP_PROXY` and `HTTPS_PROXY` values, reached ready state in 0.63 seconds, remained live after the invoking command exited, passed `status` and `doctor`, and stopped cleanly. The normal startup path therefore has no network dependency.
- Web unit tests cover labelled project/root/origin/health output plus cursor persistence, duplicate suppression, and resynchronization. The browser acceptance test loads the embedded dashboard from the built binary rather than Vite.
- macOS, Linux, and Windows clean-install smoke jobs are defined in `.github/workflows/release-smoke.yml`; local cross-compilation supplements, but does not replace, those hosted platform runs.
- The five-developer unaided usability study remains release-owner work. No participant result has been fabricated; SC-007 and the final release-study task remain open until at least five real participants are observed.

## Automated validation rerun (2026-10-01)

- `go test ./...` passed with the runtime integration suite allowed to bind loopback ports.
- `make test` passed: production Vite build, web unit tests (3 tests), and all Go tests/integration tests.
- Event storage coverage now verifies bounded segment rollover and rollback of the latest journal append when snapshot persistence fails.

## Convergence validation rerun (2026-10-02)

- The full `make test` target passed: production web build, 3 web unit tests, and all Go packages/integration tests.
- Startup regression coverage confirms a pre-ready failure releases its loopback listener and does not leave usable control state behind.
- Active-runtime coverage confirms a new Git revision is detected by the bounded polling monitor and persisted as a redacted `project.changed` event.
- The five-developer usability study remains outstanding and must be run with five real Git-familiar participants; no usability result is fabricated here.
