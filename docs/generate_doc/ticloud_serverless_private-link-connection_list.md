## ticloud serverless private-link-connection list

List private link connections

```
ticloud serverless private-link-connection list [flags]
```

### Examples

```
 List private link connections (interactive):
  $ ticloud serverless private-link-connection list

  Describe a private link connection (non-interactive):
  $ ticloud serverless private-link-connection list -c <cluster-id>
```

### Options

```
  -c, --cluster-id string   The cluster ID.
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

* [ticloud serverless private-link-connection](ticloud_serverless_private-link-connection.md)	 - Manage private link connections

