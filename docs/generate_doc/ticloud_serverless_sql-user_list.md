## ticloud serverless sql-user list

List all accessible SQL users

```
ticloud serverless sql-user list [flags]
```

### Examples

```
  List all SQL users in interactive mode:
  $ ticloud serverless sql-user list

  List all SQL users in non-interactive mode:
  $ ticloud serverless sql-user list -c <cluster-id>

  List all SQL users with json format:
  $ ticloud serverless sql-user list -o json
```

### Options

```
  -c, --cluster-id string   The ID of the cluster.
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

* [ticloud serverless sql-user](ticloud_serverless_sql-user.md)	 - Manage cluster SQL users

