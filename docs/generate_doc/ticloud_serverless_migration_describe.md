## ticloud serverless migration describe

Describe a migration

```
ticloud serverless migration describe [flags]
```

### Examples

```
  Describe a migration in interactive mode:
  $ ticloud serverless migration describe

  Describe a migration in non-interactive mode:
  $ ticloud serverless migration describe -c <cluster-id> --migration-id <migration-id>
```

### Options

```
  -c, --cluster-id string     Cluster ID that owns the migration.
  -h, --help                  help for describe
  -m, --migration-id string   ID of the migration to describe.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless migration](ticloud_serverless_migration.md)	 - Manage TiDB Cloud Serverless migrations

