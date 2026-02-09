# Unit tests

This document describes how to write TiDB Cloud CLI unit tests.

Follow the patterns under `internal/cli/serverless/branch`:

- Use generated mocks for service calls.
- Use `suite.Suite` with `SetupTest` to set `NO_COLOR`, create `iostream.Test()`, and inject a mock client.
- Use table-driven tests for flag combinations and error cases.
- Validate both `stdout` and `stderr` contents exactly.
- Assert mock expectations only on success paths.
- Cover:
  - Required flags errors.
  - Shorthand flags (`-c`, `-b`, `-o`, etc.).
  - Output formats (`json` vs human).
  - Multi-page list behavior when applicable.