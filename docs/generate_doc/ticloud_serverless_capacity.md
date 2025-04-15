## ticloud serverless capacity

Set capacity for a TiDB Cloud Serverless cluster

```
ticloud serverless capacity [flags]
```

### Examples

```
  Set capacity for a TiDB Cloud Serverless cluster in interactive mode:
  $ ticloud serverless capacity

  Set capacity for a TiDB Cloud Serverless cluster in non-interactive mode:
  $ ticloud serverless capacity -c <cluster-id> --max-rcu <maximum-rcu> --min-rcu <minimum-rcu>
```

### Options

```
  -c, --cluster-id string   The ID of the cluster.
  -h, --help                help for capacity
      --max-rcu int32       Maximum RCU for the cluster, at most 100000.
      --min-rcu int32       Minimum RCU for the cluster, at least 2000.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless](ticloud_serverless.md)	 - Manage TiDB Cloud Serverless clusters

