## ticloud serverless authorized-network update

Update an authorized network

```
ticloud serverless authorized-network update [flags]
```

### Examples

```
  Update an authorized network in interactive mode:
  $ ticloud serverless authorized-network update

  Update an authorized network in non-interactive mode:
  $ ticloud serverless authorized-network update -c <cluster-id> --start-ip-address <start-ip-address> --end-ip-address <end-ip-address> --new-start-ip-address <new-start-ip-address> --new-end-ip-address <new-end-ip-address>
```

### Options

```
  -c, --cluster-id string             The ID of the cluster.
      --end-ip-address string         The end IP address of the authorized network.
  -h, --help                          help for update
      --new-display-name string       The new display name of the authorized network.
      --new-end-ip-address string     The new end IP address of the authorized network.
      --new-start-ip-address string   The new start IP address of the authorized network.
      --start-ip-address string       The start IP address of the authorized network.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless authorized-network](ticloud_serverless_authorized-network.md)	 - Manage TiDB Cloud Serverless cluster authorized networks

