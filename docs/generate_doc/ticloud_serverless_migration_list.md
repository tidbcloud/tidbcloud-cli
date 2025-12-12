## ticloud serverless migration list

List migrations

```
ticloud serverless migration list [flags]
```

### Examples

```
  List migrations in interactive mode:
  $ ticloud serverless migration list

  List migrations in non-interactive mode with JSON output:
  $ ticloud serverless migration list -c <cluster-id> -o json
```

### Options

```
  -c, --cluster-id string   The cluster ID of the migration tasks to list.
  -h, --help                help for list
  -o, --output string       Output format, one of ["human" "json"]. For the complete result, please use json format. (default "human")
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless migration](ticloud_serverless_migration.md)	 - Manage TiDB Cloud Serverless migrations

