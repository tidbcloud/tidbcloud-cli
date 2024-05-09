## ticloud serverless export list

List export tasks

```
ticloud serverless export list [flags]
```

### Examples

```
  List all exports in interactive mode:
  $ ticloud serverless export list

  List all exports in non-interactive mode:
  $ ticloud serverless export list -c <cluster-id> 

  List all exports with json format in non-interactive mode:
  $ ticloud serverless export list -c <cluster-id> -o json
```

### Options

```
  -c, --cluster-id string   The cluster ID of the exports to be listed.
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

* [ticloud serverless export](ticloud_serverless_export.md)	 - Manage TiDB Serverless exports

