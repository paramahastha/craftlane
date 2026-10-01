# Runtime architecture

The command detects one canonical Git root, creates a project-scoped runtime
record, binds an IPv4 loopback listener, and serves the embedded dashboard and
typed API. Durable records are written with temp-file replacement and control
credentials are stored separately with owner-only permissions.
