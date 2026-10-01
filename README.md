# Craftlane

Craftlane is a local-first project runtime. The foundation executable detects the
current Git repository, serves a semantic dashboard on IPv4 loopback, and keeps
runtime state on the developer’s machine. It requires no account, cloud service,
or network access at runtime.

## Development

Install Go 1.25+ and, when available, Node.js/pnpm for the contributor web
source. `make test` builds/verifies the embedded dashboard before running Go
tests; in an offline checkout it uses the checked-in fallback artifact.

```sh
make test
make build
./bin/craftlane start --no-open
```

Run commands from a Git repository. The service binds to `127.0.0.1`, validates
the exact Host header, and never exposes its runtime control token in status,
health, events, or diagnostics.

Set `CRAFTLANE_EVENT_RETENTION` locally to increase event history; values below
1000 are clamped to the documented minimum.
