## ticloud serverless private-link-connection zones

Get account and availability zones

```
ticloud serverless private-link-connection zones [flags]
```

### Examples

```
  Get availability zones (interactive):
  $ ticloud serverless private-link-connection get-zones

  Get availability zones (non-interactive):
  $ ticloud serverless private-link-connection zones -c <cluster-id>
```

### Options

```
  -c, --cluster-id string   The cluster ID.
  -h, --help                help for zones
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless private-link-connection](ticloud_serverless_private-link-connection.md)	 - Manage private link connections for dataflow

