## ticloud serverless changefeed list

List changefeeds

```
ticloud serverless changefeed list [flags]
```

### Examples

```
  List all changefeeds in interactive mode:
  $ ticloud serverless changefeed list

  List all changefeeds in non-interactive mode:
  $ ticloud serverless changefeed list -c <cluster-id>

  List all changefeeds with json format in non-interactive mode:
  $ ticloud serverless changefeed list -c <cluster-id> -o json
```

### Options

```
  -c, --cluster-id string   The cluster ID of the changefeeds to be listed.
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

* [ticloud serverless changefeed](ticloud_serverless_changefeed.md)	 - Manage TiDB Cloud Serverless changefeeds

