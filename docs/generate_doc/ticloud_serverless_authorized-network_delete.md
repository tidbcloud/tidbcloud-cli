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
  $ ticloud serverless authorized-network delete -c <cluster-id> --ip-range <ip-range>
```

### Options

```
  -c, --cluster-id string   The ID of the cluster.
      --force               Delete an authorized network without confirmation.
  -h, --help                help for delete
      --ip-range string     The IP range of the authorized network.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless authorized-network](ticloud_serverless_authorized-network.md)	 - Manage TiDB Cloud Serverless cluster authorized networks

