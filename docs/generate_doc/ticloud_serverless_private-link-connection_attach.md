## ticloud serverless private-link-connection attach

Attach domains to a private link connection

```
ticloud serverless private-link-connection attach [flags]
```

### Examples

```
  Attach domain (interactive):
  $ ticloud serverless private-link-connection attach

  Attach domain (non-interactive):
  $ ticloud serverless private-link-connection attach -c <cluster-id> --private-link-connection-id <plc-id> --type <type> --unique-name <unique-name>
```

### Options

```
  -c, --cluster-id string                   The cluster ID.
      --dry-run                             set dry run mode to only show generated domains without attaching them.
  -h, --help                                help for attach
      --private-link-connection-id string   The private link connection ID.
      --type string                         The type of domain to attach, one of: [TIDBCLOUD_MANAGED CONFLUENT]
      --unique-name string                  The unique name of the domain to attach, you can use --dry-run to generate the unique name when attaching a TiDB Cloud managed domain.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless private-link-connection](ticloud_serverless_private-link-connection.md)	 - Manage private link connections for dataflow

