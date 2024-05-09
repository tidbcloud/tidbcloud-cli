## ticloud serverless export cancel

Cancel an export task

```
ticloud serverless export cancel [flags]
```

### Examples

```
  Cancel an export in interactive mode:
  $ ticloud serverless export cancel

  Cancel an export in non-interactive mode:
  $ ticloud serverless export cancel -c <cluster-id> -e <export-id>
```

### Options

```
  -c, --cluster-id string   The cluster ID of the export to be canceled.
  -e, --export-id string    The ID of the export to be canceled.
      --force               Cancel an export without confirmation.
  -h, --help                help for cancel
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless export](ticloud_serverless_export.md)	 - Manage TiDB Serverless exports

