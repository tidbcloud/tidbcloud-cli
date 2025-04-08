## ticloud serverless create

Create a TiDB Cloud Serverless cluster

```
ticloud serverless create [flags]
```

### Examples

```
  Create a TiDB Cloud Serverless cluster in interactive mode:
  $ ticloud serverless create

  Create a TiDB Cloud Serverless cluster of the default project in non-interactive mode:
  $ ticloud serverless create --display-name <cluster-name> --region <region>

  Create a TiDB Cloud Serverless cluster in non-interactive mode:
  $ ticloud serverless create --project-id <project-id> --display-name <cluster-name> --region <region>
```

### Options

```
      --authorized-networks strings    The authorized networks of the public endpoint.
      --disable-public-endpoint        Whether the public endpoint is disabled.
  -n, --display-name string            Display name of the cluster to de created.
      --encryption                     Whether Enhanced Encryption at Rest is enabled.
  -h, --help                           help for create
      --max-rcu int32                  Maximum RCU for the cluster, at most 100000.
      --min-rcu int32                  Minimum RCU for the cluster, at least 2000.
  -p, --project-id string              The ID of the project, in which the cluster will be created. (default "default project")
  -r, --region string                  The name of cloud region. You can use "ticloud serverless region" to see all regions.
      --spending-limit-monthly int32   Maximum monthly spending limit in USD cents.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless](ticloud_serverless.md)	 - Manage TiDB Cloud Serverless clusters

