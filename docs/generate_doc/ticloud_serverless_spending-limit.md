## ticloud serverless spending-limit

Set spending limit for a TiDB Cloud Serverless cluster

```
ticloud serverless spending-limit [flags]
```

### Examples

```
  Set spending limit for a TiDB Cloud Serverless cluster in interactive mode:
  $ ticloud serverless spending-limit

  Set spending limit for a TiDB Cloud Serverless cluster in non-interactive mode:
  $ ticloud serverless spending-limit -c <cluster-id> --monthly <spending-limit-monthly>
```

### Options

```
  -c, --cluster-id string   The ID of the cluster.
  -h, --help                help for spending-limit
      --monthly int32       Maximum monthly spending limit in USD cents.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless](ticloud_serverless.md)	 - Manage TiDB Cloud Serverless clusters

