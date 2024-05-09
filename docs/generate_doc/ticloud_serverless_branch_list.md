## ticloud serverless branch list

List branches

```
ticloud serverless branch list [flags]
```

### Examples

```
  List all branches in interactive mode:
  $ ticloud serverless branch list

  List all branches in non-interactive mode:
  $ ticloud serverless branch list -c <cluster-id> 

  List all branches with json format in non-interactive mode:
  $ ticloud serverless branch list -c <cluster-id> -o json
```

### Options

```
  -c, --cluster-id string   The cluster ID of the branch to be listed.
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

* [ticloud serverless branch](ticloud_serverless_branch.md)	 - Manage TiDB Serverless branches

