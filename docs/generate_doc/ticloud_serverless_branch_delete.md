## ticloud serverless branch delete

Delete a branch

```
ticloud serverless branch delete [flags]
```

### Examples

```
  Delete a branch in interactive mode:
  $ ticloud serverless branch delete

  Delete a branch in non-interactive mode:
  $ ticloud serverless branch delete -c <cluster-id> -b <branch-id>
```

### Options

```
  -b, --branch-id string    The ID of the branch to be deleted.
  -c, --cluster-id string   The cluster ID of the branch to be deleted.
      --force               Delete a branch without confirmation.
  -h, --help                help for delete
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless branch](ticloud_serverless_branch.md)	 - Manage TiDB Cloud Serverless branches

