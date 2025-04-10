## ticloud serverless authorized-network create

Create an authorized network

```
ticloud serverless authorized-network create [flags]
```

### Examples

```
  Create an authorized network in interactive mode:
  $ ticloud serverless authorized-network create

  Create an authorized network in non-interactive mode:
  $ ticloud serverless authorized-network create -c <cluster-id> --display-name <display-name> --start-ip-address <start-ip-address> --end-ip-address <end-ip-address>
```

### Options

```
  -c, --cluster-id string         The ID of the cluster.
  -n, --display-name string       The name of the authorized network.
      --end-ip-address string     The end IP address of the authorized network.
  -h, --help                      help for create
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

