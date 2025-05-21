## ticloud serverless audit-log filter-rule enable

Enable an audit log filter rule

```
ticloud serverless audit-log filter-rule enable [flags]
```

### Examples

```
  Enable an audit log filter rule in interactive mode:
  $ ticloud serverless auditlog filter-rule enable

  Enable an audit log filter rule in non-interactive mode:
  $ ticloud serverless auditlog filter-rule enable --cluster-id <cluster-id> --rule-name <rule-name>

```

### Options

```
  -c, --cluster-id string   The ID of the cluster.
  -h, --help                help for enable
      --rule-name string    The name of the filter rule.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless audit-log filter-rule](ticloud_serverless_audit-log_filter-rule.md)	 - Manage TiDB Cloud Serverless database audit logging filter rules

