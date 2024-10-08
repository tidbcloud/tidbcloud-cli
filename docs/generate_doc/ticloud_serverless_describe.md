## ticloud serverless describe

Describe a TiDB Cloud Serverless cluster

```
ticloud serverless describe [flags]
```

### Examples

```
  Get a TiDB Cloud Serverless cluster in interactive mode:
 $ ticloud serverless describe

 Get a TiDB Cloud Serverless cluster in non-interactive mode:
 $ ticloud serverless describe -c <cluster-id>
```

### Options

```
  -c, --cluster-id string   The ID of the cluster.
  -h, --help                help for describe
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless](ticloud_serverless.md)	 - Manage TiDB Cloud Serverless clusters

