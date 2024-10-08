## ticloud serverless shell

Connect to a TiDB Cloud Serverless cluster

### Synopsis

Connect to a TiDB Cloud Serverless cluster.
The connection forces the [ANSI SQL mode](https://dev.mysql.com/doc/refman/8.0/en/sql-mode.html#sqlmode_ansi) for the session.

```
ticloud serverless shell [flags]
```

### Examples

```
  Connect to a TiDB Cloud Serverless cluster in interactive mode:
  $ ticloud serverless shell

  Connect to a TiDB Cloud Serverless cluster with default user in non-interactive mode:
  $ ticloud serverless shell -c <cluster-id>

  Connect to a serverless cluster with default user and password in non-interactive mode:
  $ ticloud serverless shell -c <cluster-id> --password <password>

  Connect to a serverless cluster with specific user and password in non-interactive mode:
  $ ticloud serverless shell -c <cluster-id> -u <user-name> --password <password>
```

### Options

```
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

* [ticloud serverless](ticloud_serverless.md)	 - Manage TiDB Cloud Serverless clusters

