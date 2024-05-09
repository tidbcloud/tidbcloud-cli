## ticloud serverless branch describe

Describe a branch

```
ticloud serverless branch describe [flags]
```

### Examples

```
  Get a branch in interactive mode:
  $ ticloud serverless branch describe

  Get a branch in non-interactive mode:
  $ ticloud serverless branch describe -c <cluster-id> -b <branch-id>
```

### Options

```
  -b, --branch-id string    The ID of the branch to be described.
  -c, --cluster-id string   The cluster ID of the branch to be described.
  -h, --help                help for describe
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless branch](ticloud_serverless_branch.md)	 - Manage TiDB Serverless branches

