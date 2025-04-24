## ticloud serverless audit-log disable

Disable the database audit logging

```
ticloud serverless audit-log disable [flags]
```

### Examples

```
  disable the database audit logging in interactive mode:
  $ ticloud serverless audit-log disable

  disable the database audit logging in non-interactive mode:
  $ ticloud serverless audit-log disable -c <cluster-id> 
```

### Options

```
  -c, --cluster-id string   The cluster ID.
  -h, --help                help for disable
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless audit-log](ticloud_serverless_audit-log.md)	 - Manage TiDB Cloud Serverless database audit logging

