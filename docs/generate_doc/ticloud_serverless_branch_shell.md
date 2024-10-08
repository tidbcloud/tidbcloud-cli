## ticloud serverless branch shell

Connect to a branch

### Synopsis

Connect to a branch. 
The connection forces the [ANSI SQL mode](https://dev.mysql.com/doc/refman/8.0/en/sql-mode.html#sqlmode_ansi) for the session.

```
ticloud serverless branch shell [flags]
```

### Examples

```
  Connect to a branch in interactive mode:
  $ ticloud serverless branch shell

  Connect to a branch with default user in non-interactive mode:
  $ ticloud serverless branch shell -c <cluster-id> -b <branch-id>

  Connect to a branch with default user and password in non-interactive mode:
  $ ticloud serverless branch shell -c <cluster-id> -b <branch-id> --password <password>

  Connect to a branch with specific user and password in non-interactive mode:
  $ ticloud serverless branch shell -c <cluster-id> -b <branch-id> -u <user-name> --password <password>
```

### Options

```
  -b, --branch-id string    The ID of the branch.
  -c, --cluster-id string   The ID of the cluster.
  -h, --help                help for shell
      --password string     The password of the user.
  -u, --user string         A specific user for login if not using the default user.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless branch](ticloud_serverless_branch.md)	 - Manage TiDB Cloud Serverless branches

