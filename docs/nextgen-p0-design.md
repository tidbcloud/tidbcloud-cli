# Premium / Essential V2 P0 design

Source: [TiDB Cloud CLI for Premium/Essential v2](https://pingcap.feishu.cn/wiki/KTAkw40RmiCpy5kQnKDcfPtcnlg), read through Feishu MCP on 2026-09-24. The P0 command table is provided in the task screenshot; the document's embedded sheet is not expanded by the document reader.

## Scope and decisions

The original P0 surface is create/list/update/delete/describe, region, and shell with interactive selection, an explicit instance ID, default credentials, or an explicit username/password. A subsequent CLI extension adds public-endpoint enable/disable for both plans. Premium and Essential V2 retain plan guards. The current published contract governs flags and request fields: max RCU is user supplied; minimum RCU is backend controlled. Import/export is P1. Shell selects PUBLIC endpoints by default and retains TLS verification; the private-connection extension below selects existing PRIVATE_ENDPOINT addresses without provisioning networking.

The CMEK increment closes the create gap using existing mgmt handlers, changing the public OpenAPI contract, generated CLI client, CLI commands, tests and documentation. The subsequent private-connection correction also fills existing endpoint fields in NextGen Global's cluster read path. Global/Regional/infra encryption behavior is reused. Existing root-password and public-connection OpenAPIs supply shell prerequisites. The CLI public-endpoint extension reuses the published public-connection PATCH without changing the API or backend behavior.

## OpenAPI

- Publish `GET /v1beta2/cmeks:principal` and `POST /v1beta2/cmeks:verifyPrincipal` by adding the `public` marker, preserving the existing `CustomerManagedEncryptionKeyService` tag, paths, operation IDs and payload shapes.
- Declare the principal request's `regionId` as required only in the newly published public Swagger projection. Keep the existing shared proto annotation and full/Console Swagger requiredness unchanged so existing SDK consumers are not tightened. The runtime already rejects an omitted or invalid region.
- Keep `servicePlan` and the shared CMEK model unchanged. Document that the helper APIs support Premium on AWS or AliCloud. Keep returned principal fields, including AWS external ID, intact.
- Reuse existing portal resource mappings. For CMEK verification without `tidbId`, keep organization-scoped pre-create behavior. With `tidbId`, use the instance permission checker before its resource pool is resolved: enforce organization ownership, token project scope and the existing instance-read privilege. Keep the original CMEK RBAC metadata intact. No Global/Regional permission or encryption changes are needed.
- Preserve the first API tag for existing full/Console SDK consumers (`CustomerManagedEncryptionKeyService`, `PublicConnectionSettingService`); append the publication marker without renaming the original groups. Verify Digest/Bearer transport behavior. Runtime Gateway/IAM enablement remains a deployment acceptance check.
- Reuse existing `POST /tidbs` and `GET /tidbs/{id}` for creation/readback. Do not publish CMEK list/get or encryption-update endpoints for this increment.

## CLI

- Preserve `--encryption none` (default) and `--encryption default-key`, their JSON payloads and request count.
- Add `--encryption cmek --cmek-key-arn <key ARN>` to Premium create. Register the ARN flag and `premium cmek principal --region <region>` only for Premium. Essential V2 rejects CMEK before obtaining a client.
- Keep `premium --create`, existing hidden action subcommands, `essential-v2` and `essential` alias routing intact. Preserve optional project selection, output formats and credential precedence.
- Validate output format, existing create fields, encryption/ARN combination and supported provider before obtaining a client. Avoid implementing cloud-provider key-policy validation in the CLI; the server owns key validity and access rules.
- Only the CMEK create branch fetches the principal, builds the provider-specific key object, verifies access, and then creates the instance. HTTP errors or `valid=false` stop creation and preserve the cause. Verification itself follows the existing backend behavior and is not advertised as a side-effect-free dry run.
- The principal helper allows users to configure their KMS policy before creation. CLI neither changes IAM policy nor handles raw encryption key material.
- Add `public-endpoint enable|disable` to both plan commands. An omitted `-c` selects an instance interactively; an explicit ID skips selection, and an explicitly empty ID is rejected. Read the selected instance again and validate its current plan before mutation, send only the nullable `enabled` field, and leave the IP access list unchanged. Disabling requires confirmation naming the target instance or `--force`; both operations report request acceptance rather than endpoint convergence. Enabling does not add IP access list entries or guarantee connectivity.
- No automatic downgrade, create retry, capability cache or endpoint-version inference. A 404 can also be routing/configuration failure, so preserve the API error instead of asserting that the server is old.

## Compatibility and validation

The CMEK increment does not add helper requests to existing commands: old CLI/new Gateway and new CLI/old Gateway keep their existing wire contracts. New CMEK commands require the two newly published endpoints; new OAuth/PUBLIC-setting capabilities also require their corresponding server rollout. Deploy the permission-aware mgmt-service before enabling OAuth in Portal, then publish Gateway support before CLI. Preserve prior public operation IDs, shared schemas and API group names. Pin OpenAPI Generator 7.12.0 and explicitly map the new principal query's deduplicated service-plan enum to the existing `TidbServiceListTidbsServicePlanParameter`; otherwise the generator renames List/Region's existing parameter type. Generate public Swagger before UI Swagger so a single `make gen` produces consistent artifacts; check idempotence.

Test AWS/AliCloud payloads; helper output; missing/conflicting flags; unsupported plan/provider; invalid output before requests; principal/verify failure before create; `valid=false`; existing none/default-key payloads and one-request behavior; action routing and aliases; Digest/Bearer requests; the public, full and UI Swagger contract. Run existing P0 lifecycle, pagination, plan-guard, password and shell tests, CLI build/vet/race checks, and affected mgmt module generation/lint/tests. Inspect both complete working diffs, including untracked files, against their recorded baselines.

Real CMEK acceptance requires a target Gateway exposing these APIs, an accessible KMS key/policy, a successful create followed by CMEK readback, and explicit test-resource cleanup. Local transport tests do not prove deployed IAM/KMS or staging readiness.

## P0 requirement mapping

| Requirement | Implementation and boundary |
|---|---|
| Create | `premium --create` / `essential-v2 --create`; explicit service plan, optional project label, backend-controlled minimum RCU; no automatic public endpoint or IAM-policy creation |
| List | Plan-scoped `GET /v1beta2/tidbs`; follows pagination and rejects repeated tokens |
| Describe | `GET /v1beta2/tidbs/{id}`; rejects a different plan before returning the resource |
| Update | `PATCH /v1beta2/tidbs/{id}`; sends only changed display name/max RCU after plan check; preserves capacity cooldown/backend validation |
| Delete | `DELETE /v1beta2/tidbs/{id}` after plan check and explicit confirmation or `--force`; reports acceptance, not completed deletion |
| Region | Plan-scoped `GET /v1beta2/regions`; paginated with token-loop detection |
| Interactive shell | Without `-c`, `-u` or `--password`, select a plan-scoped instance and SQL user interactively; with `-c`, skip instance selection and prompt for a password if absent; network flags alone preserve interactive selection |
| Explicit shell credentials | `-c ... --password ...` and `-c ... -u ... --password ...`; same meaning as Starter: no credential prompt, still a terminal-based SQL shell, not a new batch SQL mode |
| Shell prerequisites | Existing root-password API wrapped by `password`; public settings available through OpenAPI; fetch the instance CA and enforce TLS hostname/chain verification |
| Public endpoint | `public-endpoint enable|disable`; interactive or explicit instance selection, plan guard, PATCH only `enabled`, confirmation/`--force` for disable; no IP allowlist mutation or convergence claim |
| Authentication | Existing API-key Digest and CLI OAuth Bearer; preserve credential precedence and legacy Starter transport; do not persist Bearer secrets in workflow JSON |
| CMEK | Premium create with AWS/Alibaba Cloud key ARN; principal retrieval, verification and existing create body; Essential V2 remains rejected by the unchanged backend capability |

The draft spec's camelCase flags and `minRcu` are not a new contract: the user selected the current published contract. The initial PUBLIC-only shell is extended below at the user's request. CMEK updates/rotation, private-network setup, SQL-user administration and import/export are not added by this increment.

## Existing private-endpoint connections

The existing OpenAPI schema and converter support `PRIVATE_ENDPOINT` host/port entries, but runtime review found that NextGen Global did not populate them. Its cluster Get/List read path now reads existing private services/connections: active AWS services supply their service-wide DNS name; GCP/Azure/AliCloud supply active, registered connection hostnames. Empty hosts, invalid ports and deleted/unregistered connections are not published. Private metadata lookup failures are logged without failing the existing cluster query. PUBLIC/VPC_PEERING projection is preserved. Backend `reachable=true` is not proof that the client can resolve or reach the private address. The CLI therefore selects a published address and lets the actual TLS/SQL connection establish connectivity; it does not add a separate network probe.

- Add `--connection-type public|private-endpoint`, defaulting to `public`, to both plan commands. PUBLIC endpoint selection and Starter behavior remain unchanged. `VPC_PEERING` is a different connection type and is not supported by this extension.
- Add `--endpoint <host:port>` only for private connections. Validate its syntax before obtaining an API client and require an exact match against this instance's API-returned PRIVATE_ENDPOINT addresses (normalizing the numeric port). This is not an arbitrary host override. IPv6 addresses use `[address]:port`.
- Preserve instance selection and credential semantics. Network flags alone keep interactive selection. A single usable private address is selected automatically; multiple addresses prompt only in the existing interactive path. With `-c`, require `--endpoint` for multiple addresses instead of silently choosing one or adding a prompt.
- Missing private endpoints, addresses not yet returned, an unmatched selector, and explicitly unreachable endpoints return errors. Never fall back to PUBLIC or another endpoint after a failed connection. Keep plan and ACTIVE-state checks before connecting.
- Use the selected host for both the SQL destination and TLS ServerName, with the existing instance CA download and strict verification. Do not add skip-verify, custom host/SNI overrides or retries.
- Bound private initial connection setup to 30 seconds, starting after password entry. Check the child context immediately after `Open`, before canceling it: usql can swallow the deadline error from its initial version query. Only enter the SQL shell if setup succeeded before the deadline, then cancel the child context so it does not impose a SQL-session/query timeout. PUBLIC and Starter retain their existing context behavior.

No new API/schema/SDK or resource-provisioning change is required, but the NextGen Global read-path fix must be deployed before the CLI can receive private addresses. Users remain responsible for endpoint creation/authorization, DNS, routing and a client with private-network access.

Local regression coverage includes both plans, default-public compatibility, plan rejection, invalid flags before requests, one/multiple/missing addresses, explicit selection, cancellation, no public/peering fallback, and timeout/session separation. Real private TLS acceptance still requires an existing private endpoint and a network-reachable runner: test interactive and `-c` selection, execute `SELECT CURRENT_USER()` and `SHOW STATUS LIKE 'Ssl_cipher'`, verify certificate/hostname failures remain failures, and verify an inaccessible endpoint exits within the initial connection timeout. Local tests do not prove deployed CA/SAN correctness or private connectivity.

## Review outcome and non-goals

The compatibility review resulted in three targeted corrections: optional-instance authorization for verification, public-only requiredness for `regionId`, and stable SDK group/enum names. Existing none/default-key create paths still issue one create call. Existing Starter commands and v1beta1 generated clients are not extended or rewritten. No retry engine, capability cache, encryption downgrade, TLS bypass, KMS-policy editor, database migration or internal encryption workflow change is introduced.

Tests and the final cross-repository review are recorded in `docs/nextgen-p0-review.md`. Source completion and deployment acceptance are separate: staging capacity or CA-chain failures must not be converted into code workarounds or reported as passing E2E.
