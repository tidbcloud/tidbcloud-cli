## ticloud serverless authorized-network list

List all authorized networks

```
ticloud serverless authorized-network list [flags]
```

### Examples

```
  List all authorized networks in interactive mode:
  $ ticloud serverless authorized-network list

  List all authorized networks in non-interactive mode:
  $ ticloud serverless authorized-network list -c <cluster-id>

  List all authorized networks with json format:
  $ ticloud serverless authorized-network list -o json
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

* [ticloud serverless authorized-network](ticloud_serverless_authorized-network.md)	 - Manage TiDB Cloud Serverless cluster authorized networks

