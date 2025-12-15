## ticloud serverless migration create

Create a migration

```
ticloud serverless migration create [flags]
```

### Examples

```
  Create a migration:
  $ ticloud serverless migration create -c <cluster-id> --display-name <name> --config-file <file-path> --dry-run
  $ ticloud serverless migration create -c <cluster-id> --display-name <name> --config-file <file-path>

```

### Options

```
  -c, --cluster-id string     The ID of the target cluster.
      --config-file string    Path to a migration config JSON file. Use "ticloud serverless migration template --mode <mode>" to print templates.
  -n, --display-name string   Display name for the migration.
      --dry-run               Run a migration precheck (dry run) with the provided inputs without creating a migration.
  -h, --help                  help for create
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless migration](ticloud_serverless_migration.md)	 - Manage TiDB Cloud Serverless migrations

