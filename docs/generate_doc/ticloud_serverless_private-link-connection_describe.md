## ticloud serverless private-link-connection describe

Describe a private link connection

```
ticloud serverless private-link-connection describe [flags]
```

### Examples

```
  Describe a private link connection (interactive):
  $ ticloud serverless private-link-connection describe

  Describe a private link connection (non-interactive):
  $ ticloud serverless private-link-connection describe -c <cluster-id> -p <private-link-connection-id>
```

### Options

```
  -c, --cluster-id string                   The cluster ID.
  -h, --help                                help for describe
  -p, --private-link-connection-id string   The private link connection ID.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless private-link-connection](ticloud_serverless_private-link-connection.md)	 - Manage private link connections

