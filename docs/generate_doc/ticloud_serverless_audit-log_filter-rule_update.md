## ticloud serverless audit-log filter-rule update

Update an audit log filter rule

```
ticloud serverless audit-log filter-rule update [flags]
```

### Examples

```
  Update an audit log filter rule in interactive mode:
  $ ticloud serverless audit-log filter update

  Enable audit log filter rule in non-interactive mode:
  $ ticloud serverless audit-log filter update --cluster-id <cluster-id> --filter-rule-id <rule-id> --enabled

  Disable audit log filter rule in non-interactive mode:
  $ ticloud serverless audit-log filter update --cluster-id <cluster-id> --filter-rule-id <rule-id> --enabled=false

  Update filters of an audit log filter rule in non-interactive mode:
  $ ticloud serverless audit-log filter update --cluster-id <cluster-id> --filter-rule-id <rule-id> --rule '{"users":["%@%"],"filters":[{"classes":["QUERY"],"tables":["test.t"]}]}'

```

### Options

```
  -c, --cluster-id string       The ID of the cluster.
      --display-name string     The display name of the filter rule.
      --enabled                 Enable or disable the filter rule.
      --filter-rule-id string   The ID of the filter rule.
  -h, --help                    help for update
      --rule string             Complete filter rule expressions, use "ticloud serverless audit-log filter template" to see filter templates.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless audit-log filter-rule](ticloud_serverless_audit-log_filter-rule.md)	 - Manage TiDB Cloud Serverless database audit logging filter rules

