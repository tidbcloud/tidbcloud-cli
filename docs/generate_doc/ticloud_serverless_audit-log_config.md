## ticloud serverless audit-log config

Configure the database audit logging

```
ticloud serverless audit-log config [flags]
```

### Examples

```
  Conigure the database audit logging in interactive mode:
  $ ticloud serverless audit-log config

  Unredacted the database audit logging in non-interactive mode:
  $ ticloud serverless audit-log config -c <cluster-id> --unredacted

  Enable the database audit logging in non-interactive mode:
  $ ticloud serverless audit-log config -c <cluster-id> --enabled

  Disable the database audit logging in non-interactive mode:
  $ ticloud serverless audit-log config -c <cluster-id> --enabled=false
```

### Options

```
  -c, --cluster-id string   The ID of the cluster to be updated.
      --enabled             enable or disable the database audit logging.
  -h, --help                help for config
      --unredacted          unredacted or redacted the database audit logging.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless audit-log](ticloud_serverless_audit-log.md)	 - Manage TiDB Cloud Serverless database audit logging

