## ticloud serverless import list

List data import tasks

```
ticloud serverless import list [flags]
```

### Examples

```
  List import tasks in interactive mode:
  $ ticloud serverless import list

  List import tasks in non-interactive mode:
  $ ticloud serverless import list --cluster-id <cluster-id>
  
  List the clusters in the project with json format:
  $ ticloud serverless import list --cluster-id <cluster-id> --output json
```

### Options

```
  -c, --cluster-id string   Cluster ID.
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

* [ticloud serverless import](ticloud_serverless_import.md)	 - Manage TiDB Serverless data imports

