## ticloud serverless authorized-network delete

Delete an authorized network

```
ticloud serverless authorized-network delete [flags]
```

### Examples

```
  Delete an authorized network in interactive mode:
  $ ticloud serverless authorized-network delete

  Delete an authorized network in non-interactive mode:
  $ ticloud serverless authorized-network delete -c <cluster-id> --start-ip-address <start-ip-address> --end-ip-address <end-ip-address>
```

### Options

```
  -c, --cluster-id string         The ID of the cluster.
      --end-ip-address string     The end IP address of the authorized network.
      --force                     Delete an authorized network without confirmation.
  -h, --help                      help for delete
      --start-ip-address string   The start IP address of the authorized network.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless authorized-network](ticloud_serverless_authorized-network.md)	 - Manage TiDB Cloud Serverless cluster authorized networks

