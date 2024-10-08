## ticloud serverless import cancel

Cancel a data import task

```
ticloud serverless import cancel [flags]
```

### Examples

```
  Cancel an import task in interactive mode:
  $ ticloud serverless import cancel

  Cancel an import task in non-interactive mode:
  $ ticloud serverless import cancel --cluster-id <cluster-id> --import-id <import-id>
```

### Options

```
  -c, --cluster-id string   Cluster ID.
      --force               Cancel an import task without confirmation.
  -h, --help                help for cancel
      --import-id string    The ID of import task.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless import](ticloud_serverless_import.md)	 - Manage TiDB Cloud Serverless data imports

