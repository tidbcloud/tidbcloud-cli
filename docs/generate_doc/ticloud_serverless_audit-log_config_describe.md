## ticloud serverless audit-log config describe

Describe the database audit logging configuration

```
ticloud serverless audit-log config describe [flags]
```

### Examples

```
  Get the database audit logging configuration in interactive mode:
  $ ticloud serverless audit-log config describe

  Get the database audit logging configuration in non-interactive mode:
  $ ticloud serverless audit-log config describe -c <cluster-id> 
```

### Options

```
  -c, --cluster-id string   The cluster ID.
  -h, --help                help for describe
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless audit-log config](ticloud_serverless_audit-log_config.md)	 - Manage TiDB Cloud Serverless database audit logging configuration

