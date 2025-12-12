## ticloud serverless migration delete

Delete a migration

```
ticloud serverless migration delete [flags]
```

### Examples

```
  Delete a migration in interactive mode:
  $ ticloud serverless migration delete

  Delete a migration in non-interactive mode:
  $ ticloud serverless migration delete -c <cluster-id> --migration-id <migration-id>
```

### Options

```
  -c, --cluster-id string     Cluster ID that owns the migration.
      --force                 Delete without confirmation.
  -h, --help                  help for delete
  -m, --migration-id string   ID of the migration to delete.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless migration](ticloud_serverless_migration.md)	 - Manage TiDB Cloud Serverless migrations

