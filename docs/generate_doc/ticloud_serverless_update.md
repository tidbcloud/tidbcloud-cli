## ticloud serverless update

Update a TiDB Serverless cluster

```
ticloud serverless update [flags]
```

### Examples

```
  Update a TiDB Serverless cluster in interactive mode:
  $ ticloud serverless update

  Update displayName of a TiDB Serverless cluster in non-interactive mode:
  $ ticloud serverless update -c <cluster-id> --display-name <new-cluster-name>
 
  Update labels of a TiDB Serverless cluster in non-interactive mode:
  $ ticloud serverless update -c <cluster-id> --labels "{\"label1\":\"value1\"}"
```

### Options

```
      --annotations string        The annotations of the cluster to be added or updated.
                                  Interactive example: {"annotation1":"value1","annotation2":"value2"}.
                                  NonInteractive example: "{\"annotation1\":\"value1\",\"annotation2\":\"value2\"}".
  -c, --cluster-id string         The ID of the cluster to be updated.
      --disable-public-endpoint   Disable the public endpoint of the cluster.
  -n, --display-name string       The new displayName of the cluster to be updated.
  -h, --help                      help for update
      --labels string             The labels of the cluster to be added or updated.
                                  Interactive example: {"label1":"value1","label2":"value2"}.
                                  NonInteractive example: "{\"label1\":\"value1\",\"label2\":\"value2\"}".
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless](ticloud_serverless.md)	 - Manage TiDB Serverless clusters

