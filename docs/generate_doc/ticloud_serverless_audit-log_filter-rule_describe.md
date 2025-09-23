## ticloud serverless audit-log filter-rule describe

Describe an audit log filter rule

```
ticloud serverless audit-log filter-rule describe [flags]
```

### Examples

```
  Describe an audit log filter rule in interactive mode:
  $ ticloud serverless audit-log filter describe

  Describe an audit log filter rule in non-interactive mode:
  $ ticloud serverless audit-log filter describe --cluster-id <cluster-id> --filter-rule-id <rule-id>

```

### Options

```
  -c, --cluster-id string       The ID of the cluster.
      --filter-rule-id string   The ID of the filter rule.
  -h, --help                    help for describe
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless audit-log filter-rule](ticloud_serverless_audit-log_filter-rule.md)	 - Manage TiDB Cloud Serverless database audit logging filter rules

