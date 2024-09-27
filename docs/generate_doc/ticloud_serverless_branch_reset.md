## ticloud serverless branch reset

Reset a branch to its parent's latest state

```
ticloud serverless branch reset [flags]
```

### Examples

```
  Reset a branch in interactive mode:
  $ ticloud serverless branch reset

  Reset a branch in non-interactive mode:
  $ ticloud serverless branch reset -c <cluster-id> -b <branch-id>
```

### Options

```
  -b, --branch-id string    The ID of the branch to be reset.
  -c, --cluster-id string   The cluster ID of the branch to be reset.
      --force               Reset a branch without confirmation.
  -h, --help                help for reset
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless branch](ticloud_serverless_branch.md)	 - Manage TiDB Cloud Serverless branches

