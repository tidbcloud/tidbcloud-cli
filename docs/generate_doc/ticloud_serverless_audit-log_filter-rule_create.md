## ticloud serverless audit-log filter-rule create

Create an audit log filter rule

```
ticloud serverless audit-log filter-rule create [flags]
```

### Examples

```
  Create a filter rule in interactive mode:
  $ ticloud serverless audit-log filter create

  Create a filter rule which filters all audit logs in non-interactive mode:
  $ ticloud serverless audit-log filter create --cluster-id <cluster-id> --display-name <rule-name> --rule '{"users":["%@%"],"filters":[{}]}'

  Create a filter rule which filters QUERY and EXECUTE for test.t and filter QUERY for all tables in non-interactive mode:
  $ ticloud serverless audit-log filter create --cluster-id <cluster-id> --display-name <rule-name> --rule '{"users":["%@%"],"filters":[{"classes":["QUERY","EXECUTE"],"tables":["test.t"]},{"classes":["QUERY"]}]}'

```

### Options

```
  -c, --cluster-id string     The ID of the cluster.
      --display-name string   The display name of the filter rule.
  -h, --help                  help for create
      --rule string           Filter rule expressions, use "ticloud serverless audit-log filter template" to see filter templates.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless audit-log filter-rule](ticloud_serverless_audit-log_filter-rule.md)	 - Manage TiDB Cloud Serverless database audit logging filter rules

