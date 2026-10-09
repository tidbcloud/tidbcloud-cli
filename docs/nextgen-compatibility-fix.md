# NextGen CLI response compatibility fix

> Status (2026-10-08): The capacity-field correction below remains current.
> The proposed Management/private-reachability changes and their validation
> are historical and have been superseded by [Existing private-endpoint
> connections](nextgen-p0-design.md#existing-private-endpoint-connections).
> Current private discovery uses the existing networking GET APIs and does
> not require Global enrichment or a Management reachability change.

## Problem and acceptance criteria

The current SDK requires `maxRcu` in every TiDB response. The management service
omits it for Elastic instances, so a successful HTTP GET or LIST fails during SDK
decoding. Separately, Global publishes private endpoint addresses without a
reachability assessment. Management currently converts that missing assessment
to `reachable: false`, preventing the CLI from attempting a connection.

Both failures were reproduced locally. Existing tests do not cover these response
boundaries. Acceptance requires successful decoding and JSON output for both
capacity modes, and preservation of unknown versus explicitly false reachability.
Live DNS, TLS and SQL connectivity remain separate deployment acceptance checks.

## CLI change

Update only the TiDB instance capacity fields in `nextgen.swagger.json` from the
management service contract: optional/nullable `minRcu` and `maxRcu`, plus
`baselineRcu`, `capacityMode` and its enum definition. Preserve existing paths,
including CMEK and Public Connection Setting. Leave the PATCH model unchanged.

Regenerate using the repository's pinned OpenAPI generator and custom templates.
Adapt constructor calls and optional-field access. Create continues to require a
positive `--max-rcu` and explicitly sets it on the request. Update continues to
send only fields selected by the user. No new capacity flags are introduced.

JSON output preserves omitted, null and populated capacity fields, including
unknown capacity-mode strings. Human-readable lists show `-` when Max RCU is
absent or null and retain the existing columns.

## Management change (superseded)

Preserve nil reachability only when converting a `PRIVATE_ENDPOINT` address.
Leave the shared `cvtConnectionReachability` helper unchanged: PUBLIC and
VPC_PEERING retain their existing nil-to-false behavior. Keep conversion of
every non-nil object unchanged, including an empty object whose `reachable`
value is false. Use the gateway's `EmitDefaultValues: true` JSON settings to
verify that unknown private-endpoint reachability is omitted while explicit
false remains false.

The CLI retains private-endpoint selection, explicit-false rejection, connection
timeout and TLS verification. Global publishes these private addresses without
inventing a client-side reachability result. The separate PUBLIC DNS-probe
recovery path requires public access to be enabled and retains strict TLS.

## Validation

- Exercise the real generated HTTP client with successful GET/LIST JSON for
  Premium and Essential V2: Elastic, Max RCU, mixed lists, omitted/null fields,
  and `CAPACITY_MODE_UNSPECIFIED` or a future mode. Verify serialization too.
- Verify create still sends the requested `maxRcu`, required CLI validation is
  retained, and a rename sends no capacity fields. Check human-readable output.
- For private endpoints, verify that nil reachability is omitted and non-nil
  false/true assessments survive the converter and gateway JSON format.
  Decode equivalent wire responses in CLI endpoint-selection tests. For PUBLIC
  and VPC_PEERING, verify that nil still produces the existing false assessment.
- Run focused tests, CLI/package checks and the affected management module's
  required generation, formatting and lint checks. Inspect generated changes and
  preserve unrelated existing worktree edits.
- Live acceptance uses an existing reachable private endpoint to run SQL with
  TLS, and an unreachable endpoint to verify the existing connection timeout.
  Local test success does not establish that live acceptance.

## Historical verification (2026-09-29)

Passed CLI main-module `go test ./...`, SDK-module `go test ./...` (compilation),
and golangci-lint v1.64.7. Management passed the complete
`app/openapi/nextgen/tidb` package tests and `make generate`, `make fmt`, and
`make lint`. Swagger paths and the PATCH schema are unchanged. Existing
management worktree files were compared against a pre-edit snapshot; only the
intended converter changed, with one new regression-test file. Live private
network acceptance has not been run.
