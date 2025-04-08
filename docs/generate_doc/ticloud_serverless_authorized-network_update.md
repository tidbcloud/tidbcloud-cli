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
  $ ticloud serverless authorized-network update -c <cluster-id> --ip-range <ip-range> --display-name <display-name>
```

### Options

```
  -c, --cluster-id string        The ID of the cluster.
  -n, --display-name string      The name of the authorized network.
  -h, --help                     help for update
      --ip-range string          The new IP range of the authorized network.
      --target-ip-range string   The IP range of the authorized network to be updated.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless authorized-network](ticloud_serverless_authorized-network.md)	 - Manage TiDB Cloud Serverless cluster authorized networks

