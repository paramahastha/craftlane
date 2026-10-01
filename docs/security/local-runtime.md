# Local runtime security

The runtime accepts only its exact `127.0.0.1:<port>` Host. State-changing
requests require `X-Craftlane-Token` from the current runtime and reject a
different browser Origin. Secret-like keys and credential-shaped values are
redacted recursively before event storage or presentation. No generic shell
or arbitrary filesystem endpoint is exposed.
