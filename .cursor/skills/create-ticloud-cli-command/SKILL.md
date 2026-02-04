---
name: create-ticloud-cli-command
description: Builds TiDB Cloud CLI commands. Use when creating new commands.
license: Complete terms in LICENSE.txt
metadata:
  author: shiyuhang0
---

# Create TiCloud CLI Command

This skill builds TiDB Cloud CLI commands, helping users create new commands with production-ready code.

## When to use

- The user asks to create new CLI commands or subcommands.
- Do not use when the agent does not support plan mode.

## About TiDB Cloud CLI

TiDB Cloud CLI (TiCloud CLI) is a command-line interface for interacting with TiDB Cloud, built on the Cobra library.

Key design of TiDB Cloud CLI:
- Built on the Cobra library.
- Uses the TiDB Cloud Open API as the client and keeps the SDK inside the project.
- Every command supports both interactive and non-interactive modes.

- More about TiDB Cloud CLI: `references/tidbcloud-cli.md`.
- More about Cobra: `references/cobra.md`.

## Core Principles

### Plan first

The AI agent must switch to plan mode to collect information from the user and generate a plan first.

Then use agent mode to generate the code.

### Test first

If the user requires testing, testing takes precedence over implementation. Write tests first, then implement, and finally validate against the tests.

### Declarative

The user must declare the command format; the implementation must comply with that declared format.

## Workflow

Must follow the workflow below:

### Generate SDK phase

Detect whether there are new or modified Swagger files under `/pkg/tidbcloud`. Skip this phase if none are found.

If detection fails, prompt the user: "Do you need to update swagger and generate SDK?" Skip this phase if the user does not need it.

Once in this phase, follow the guide in `references/sdk.md` to generate the SDK.

### Plan phase

The agent must switch to plan mode in this phase.

The agent must ask the user for the following information during the plan phase:
1. **Command format**: If the user does not provide the command format, ask them to provide it first. See `assets/command.md` for the command format template.
2. **Tests**: Ask whether the user needs tests. Tests are recommended.
3. **Other necessary information**

Generate the plan after the user confirms all the information.

### Agent phase

The agent must switch to agent mode in this phase. Follow the workflow below:

1. **Write unit tests**
2. **Implement command**
3. **Run tests**

#### Write unit tests

Skip this step if the user does not need tests.

Write unit tests following `references/ut.md`. If the tests need to invoke the command, create the empty command first. See `assets/empty.go`.

#### Implement command

Implement the command following the patterns in `assets/example`.

The implementation must meet the following requirements:

1. **Command definition**: Include name, aliases, example, and other necessary attributes.

2. **Flags definition**: Include flag type, description, default value, etc. Infer flag information from the Swagger spec and SDK client parameters. See `references/flag.md` for more details.

3. **Flags retrieval**: Implement retrieval logic for both interactive and non-interactive modes. For interactive mode, UI is required. The UI is based on the Bubbletea library. See `references/ui.md` for more details.

4. **Dual mode support**: Support both interactive and non-interactive modes. Rules:
   1. **Detect mode** with `MarkInteractive`:
      - Default to interactive.
      - If any non-interactive flag is provided, switch to non-interactive and mark required flags.
   2. **Interactive constraints**:
      - If `!h.IOStreams.CanPrompt`, return an error instructing the user to use non-interactive flags.

5. **SDK calls**: Use the SDK client in `api_client.go` to call the corresponding method. Note that `assets/example` does not include this part (it is only a template); it must be included when implementing the actual command.

#### Run the tests

Skip this step if the user does not need tests.

Run tests with: `go test -race -cover -count=1 path -v`

For example, to run tests under `internal/cli/serverless/branch`: `go test -race -cover -count=1 ./internal/cli/serverless/branch -v`

Ensure all tests pass. If they do not pass, fix the implementation or the tests.

## Final Code Structure

The final code structure must follow the structure in `assets/example`:

- Command Go files.
- UI Go files.
- Test files if needed.

Command output should follow the patterns in `assets/example`:

- **Human output**: Print a concise, user-friendly message; use `color.GreenString` for success. Avoid raw JSON in human mode.
- **JSON output**: Use `output.PrintJson` with a stable, predictable payload (e.g. keys like `message`, resource IDs, `items`).
- **No TTY**: Default to JSON output for list/describe-style commands (see `assets/example/list.go`).
- **Examples**: Include `Example` blocks for both interactive and non-interactive usage (see `assets/example/create.go` and `assets/example/list.go`).