## ticloud serverless migration resume

Resume a paused migration

```
ticloud serverless migration resume [flags]
```

### Examples

```
  Resume a migration in interactive mode:
  $ ticloud serverless migration resume

  Resume a migration in non-interactive mode:
  $ ticloud serverless migration resume -c <cluster-id> --migration-id <migration-id>
```

### Options

```
  -c, --cluster-id string     Cluster ID that owns the migration.
  -h, --help                  help for resume
  -m, --migration-id string   ID of the migration to resume.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless migration](ticloud_serverless_migration.md)	 - Manage TiDB Cloud Serverless migrations

