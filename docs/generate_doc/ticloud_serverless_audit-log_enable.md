## ticloud serverless audit-log enable

Enable the database audit logging

```
ticloud serverless audit-log enable [flags]
```

### Examples

```
  enable the database audit logging in interactive mode:
  $ ticloud serverless audit-log enable

  enable the database audit logging in non-interactive mode:
  $ ticloud serverless audit-log enable -c <cluster-id> 
```

### Options

```
  -c, --cluster-id string   The cluster ID.
  -h, --help                help for enable
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless audit-log](ticloud_serverless_audit-log.md)	 - Manage TiDB Cloud Serverless database audit logging

