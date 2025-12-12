## ticloud serverless migration pause

Pause a migration

```
ticloud serverless migration pause [flags]
```

### Examples

```
  Pause a migration in interactive mode:
  $ ticloud serverless migration pause

  Pause a migration in non-interactive mode:
  $ ticloud serverless migration pause -c <cluster-id> --migration-id <migration-id>
```

### Options

```
  -c, --cluster-id string     Cluster ID that owns the migration.
  -h, --help                  help for pause
  -m, --migration-id string   ID of the migration to pause.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless migration](ticloud_serverless_migration.md)	 - Manage TiDB Cloud Serverless migrations

