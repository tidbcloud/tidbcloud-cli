---
name: update-ticloud-cli-command
description: Update TiDB Cloud CLI commands. Use when adding new flags or updating existing commands.
license: Complete terms in LICENSE.txt
metadata:
  author: shiyuhang0
---

# Update TiCloud CLI Command

This skill updates TiDB Cloud CLI existing commands, helping users update existing commands with production-ready code.

## When to use

- The user asks to add flags to existing commands.
- The user asks to update logic of existing commands.

## About TiDB Cloud CLI

TiDB Cloud CLI (TiCloud CLI) is a command-line interface for interacting with TiDB Cloud, built on the Cobra library.

Key design of TiDB Cloud CLI:
- Built on the Cobra library.
- Uses the TiDB Cloud Open API as the client and keeps the SDK inside the project.
- Every command supports both interactive and non-interactive modes.

- More about TiDB Cloud CLI: `references/tidbcloud-cli.md`.
- More about Cobra: `references/cobra.md`.

## Workflow

Must follow the workflow below:

### Generate SDK phase

Always prompt the user: "Do you need to add or update swagger? Please provide the swagger path if you need."

Skip this phase if user does not need.

Once in this phase, follow the guide in `references/sdk.md` to generate the SDK.

After SDK is generated, ask user to use go>=1.24 to run `make generate-mocks` manually!

### Plan phase

The agent needs to switch to plan mode if supported. This phase can be skipped if user already provide enough informations.

The agent must ask the user for the following information during the plan phase:
1. **Updating Information**: Ask the user for details of the update, including whether new flags need to be added (if yes, request the user to provide them), and whether existing logic needs to be modified (if yes, ask the user to provide a specific description of the logic).
2. **Other necessary information**

Generate the plan after the user confirms all the information.

### Agent phase

The agent must switch to agent mode in this phase. Follow the workflow below:

1. **Update command**
2. **Write and run tests**

#### Update command

Update the command according to user-provided information and current command code style.

If new flags are added:
- Refer to `references/flag.md` for the flag definition.
- Refer to `references/ui.md` for the UI design.

#### Write and run the tests

Skip this step if the test file does not exist. Otherwise, write unit tests following `references/ut.md`. Then run tests with: `go test -race -cover -count=1 path -v`

For example, to run tests under `internal/cli/serverless/branch`: `go test -race -cover -count=1 ./internal/cli/serverless/branch -v`

Ensure all tests pass. If they do not pass, fix the implementation or the tests.