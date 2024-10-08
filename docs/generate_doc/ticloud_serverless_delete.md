## ticloud serverless delete

Delete a TiDB Cloud Serverless cluster

```
ticloud serverless delete [flags]
```

### Examples

```
  Delete a TiDB Cloud Serverless cluster in interactive mode:
 $ ticloud serverless delete

 Delete a TiDB Cloud Serverless cluster in non-interactive mode:
 $ ticloud serverless delete -c <cluster-id>
```

### Options

```
  -c, --cluster-id string   The ID of the cluster to be deleted.
      --force               Delete a cluster without confirmation.
  -h, --help                help for delete
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless](ticloud_serverless.md)	 - Manage TiDB Cloud Serverless clusters

