## ticloud serverless changefeed edit

Edit a changefeed

```
ticloud serverless changefeed edit [flags]
```

### Examples

```
  Update a changefeed in interactive mode:
  $ ticloud serverless changefeed edit

  Update the name, kafka, and filter of a changefeed in non-interactive mode:
  $ ticloud serverless changefeed edit -c <cluster-id> --changefeed-id <changefeed-id> --name newname --kafka <full-specified-kafka> --filter <full-specified-filter>

```

### Options

```
      --changefeed-id string   The ID of the changefeed to be updated.
  -c, --cluster-id string      The ID of the cluster.
  -n, --display-name string    The name of the changefeed.
      --filter string          Complete filter in JSON format, use "ticloud serverless changefeed template" to see templates.
  -h, --help                   help for edit
      --kafka string           Complete Kafka information in JSON format, use "ticloud serverless changefeed template" to see templates.
      --mysql string           Complete MySQL information in JSON format, use "ticloud serverless changefeed template" to see templates.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless changefeed](ticloud_serverless_changefeed.md)	 - Manage TiDB Cloud Serverless changefeeds

