## ticloud serverless audit-log filter-rule list

List audit log filter rules

```
ticloud serverless audit-log filter-rule list [flags]
```

### Examples

```
  List all audit log filter rules in interactive mode:
  $ ticloud serverless auditlog filter list

  List all audit log filter rules in non-interactive mode:
  $ ticloud serverless auditlog filter list -c <cluster-id>

  List all audit log filter rules with json format in non-interactive mode:
  $ ticloud serverless auditlog filter list -c <cluster-id> -o json
```

### Options

```
  -c, --cluster-id string   The cluster ID of the audit log filter rules to be listed.
  -h, --help                help for list
  -o, --output string       Output format, one of ["human" "json"]. For the complete result, please use json format. (default "human")
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless audit-log filter-rule](ticloud_serverless_audit-log_filter-rule.md)	 - Manage TiDB Cloud Serverless database audit logging filter rules

