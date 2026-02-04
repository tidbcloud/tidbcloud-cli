# TiDB Cloud UI

## Overview

The basic UI lives under `internal/ui`; command-specific UI is under each command package (e.g. `assets/example/ui.go`).

Live example: `internal/cli/serverless/changefeed/ui.go`

## Forbidden

github.com/AlecAivazis/survey/v2 is no longer maintained.

Use it only for deletion confirmation. Do not use it for any other case.

Deletion confirmation example:
```
if err := survey.AskOne(&survey.Input{Message: DeleteConfirmPrompt}, &confirm); err != nil {
					return err
				}
```

## Selection UI (IDs, enums, predefined options)

Use selection UI for IDs, enum fields, and any fixed option set.

- Use `ui.InitialSelectModel(items, prompt)` with `bubbletea`.
- Enable pagination and filter:
  - `model.EnablePagination(6)`
  - `model.EnableFilter()`
- Handle interrupt:
  - If model reports `Interrupted`, return `util.InterruptError`.
- Example pattern: see `assets/example/ui.go` (`GetSelectedExample`).

Live Example: see `internal/service/cloud/logic.go` (`GetSelectedChangefeed`,`GetSelectedCluster`)

## Simple input UI (single field)

Use `ui.InitialOneInputModel` for any single field.

- Example: `GetDisplayNameInput` in `assets/example/ui.go`.

## Complex composite input (multiple related fields)

For grouped or multi-field inputs, use `ui.InitialInputModel`.

- Provide a list of keys and a description map for each field.
- Read values from `textInput.Inputs[i].Value()`.
- Validate each required field and return explicit errors.
- Example pattern: `GetS3Inputs` in `assets/example/ui.go`.

## Others

More UI components are in `internal/ui`. If what you need is not there, create a basic UI under `internal/ui` based on Bubbletea. See `references/bubbletea.md` for more information.
