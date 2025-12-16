## ticloud serverless private-link-connection detach-domains

Detach domains from a private link connection

```
ticloud serverless private-link-connection detach-domains [flags]
```

### Examples

```
  Detach domains (interactive):
  $ ticloud serverless private-link-connection detach

  Detach domains (non-interactive):
  $ ticloud serverless private-link-connection detach-domains -c <cluster-id> --private-link-connection-id <plc-id> --attach-domain-id <attach-id>
```

### Options

```
      --attach-domain-id string             The private link connection attach domain ID.
  -c, --cluster-id string                   The cluster ID.
  -h, --help                                help for detach-domains
      --private-link-connection-id string   The private link connection ID.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless private-link-connection](ticloud_serverless_private-link-connection.md)	 - Manage private link connections for dataflow

