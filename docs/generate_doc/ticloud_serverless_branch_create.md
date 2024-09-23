## ticloud serverless branch create

Create a branch

```
ticloud serverless branch create [flags]
```

### Examples

```
  Create a branch in interactive mode:
  $ ticloud serverless branch create

  Create a branch in non-interactive mode:
  $ ticloud serverless branch create --cluster-id <cluster-id> --display-name <branch-name> --parent-id <parent-id>
```

### Options

```
  -c, --cluster-id string         The ID of the cluster, in which the branch will be created.
  -n, --display-name string       The displayName of the branch to be created.
  -h, --help                      help for create
      --parent-id string          The ID of the branch parent, default is cluster id.
      --parent-timestamp string   The timestamp of the parent branch, default is current time. (RFC3339 format, e.g., 2024-01-01T00:00:00Z)
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless branch](ticloud_serverless_branch.md)	 - Manage TiDB Cloud Serverless branches

