## ticloud serverless audit-log filter-rule create

create an audit log filter rule

```
ticloud serverless audit-log filter-rule create [flags]
```

### Examples

```
  Create a filter rule in interactive mode:
  $ ticloud serverless audit-log filter-rule create

  Create a filter rule which filter all audit logs in non-interactive mode:
  $ ticloud serverless audit-log filter-rule create --cluster-id <cluster-id> --rule-name <rule-name> --users %@% --filters {}

    Create a filter rule which filter QUERY and EXECUTE for test.t1 and filter QUERY for all tables in non-interactive mode:
  $ ticloud serverless audit-log filter-rule create --cluster-id <cluster-id> --rule-name <rule-name> --users user1,user2 --filters '{"classes": ["QUERY", "EXECUTE"], "tables": ["test.t1"]}' --filters '{"classes": ["QUERY"]}'

```

### Options

```
  -c, --cluster-id string     The ID of the cluster.
      --filters stringArray   Filter expressions. e.g. '{"classes": ["QUERY"]' or '{}' to filter all audit logs.
  -h, --help                  help for create
      --rule-name string      The name of the filter rule.
      --users strings         Users to apply the rule to. e,g. %@%.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless audit-log filter-rule](ticloud_serverless_audit-log_filter-rule.md)	 - Manage TiDB Cloud Serverless database audit logging filter rules

