## ticloud serverless changefeed delete

Delete a changefeed

```
ticloud serverless changefeed delete [flags]
```

### Examples

```
  Delete a changefeed in interactive mode:
  $ ticloud serverless changefeed delete

  Delete a changefeed in non-interactive mode:
  $ ticloud serverless changefeed delete -c <cluster-id> --changefeed-id <changefeed-id>
```

### Options

```
  -f, --changefeed-id string   The changefeed ID.
  -c, --cluster-id string      The cluster ID.
      --force                  Delete a changefeed without confirmation.
  -h, --help                   help for delete
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless changefeed](ticloud_serverless_changefeed.md)	 - Manage TiDB Cloud Serverless changefeeds

