## ticloud serverless list

List all TiDB Serverless clusters

```
ticloud serverless list [flags]
```

### Examples

```
  List all TiDB Serverless clusters in interactive mode):
  $ ticloud serverless list

  List all TiDB Serverless clusters in non-interactive mode:
  $ ticloud serverless list -p <project-id>

  List all TiDB Serverless clusters in non-interactive mode:
  $ ticloud serverless list -p <project-id> -o json
```

### Options

```
  -h, --help                help for list
  -o, --output string       Output format, one of ["human" "json"]. For the complete result, please use json format. (default "human")
  -p, --project-id string   The ID of the project.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless](ticloud_serverless.md)	 - Manage TiDB Serverless clusters

