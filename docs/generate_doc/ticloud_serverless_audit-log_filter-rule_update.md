## ticloud serverless audit-log filter-rule update

Update an audit log filter rule

```
ticloud serverless audit-log filter-rule update [flags]
```

### Examples

```
  Update an audit log filter rule in interactive mode:
  $ ticloud serverless auditlog filter-rule update

  Update users of an audit log filter rule in non-interactive mode:
  $ ticloud serverless auditlog filter-rule update --cluster-id <cluster-id> --rule-name <rule-name> --users user1,user2

  Update filters of an audit log filter rule in non-interactive mode:
  $ ticloud serverless auditlog filter-rule update --cluster-id <cluster-id> --rule-name <rule-name> --filters '{"classes": ["QUERY", "EXECUTE"], "tables": ["test.t1"]}' --filters '{"classes": ["QUERY"]}'

```

### Options

```
  -c, --cluster-id string     The ID of the cluster.
      --filters stringArray   Filter expressions. e.g. '{"classes": ["QUERY"]' or '{}' to filter all audit logs.
  -h, --help                  help for update
      --rule-name string      The name of the filter rule to update.
      --users strings         Users to apply the rule to. e.g. %@%.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless audit-log filter-rule](ticloud_serverless_audit-log_filter-rule.md)	 - Manage TiDB Cloud Serverless database audit logging filter rules

