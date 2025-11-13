## ticloud serverless changefeed pause

Pause a changefeed

```
ticloud serverless changefeed pause [flags]
```

### Examples

```
  Pause a changefeed in interactive mode:
  $ ticloud serverless changefeed pause

  Pause a changefeed in non-interactive mode:
  $ ticloud serverless changefeed pause -c <cluster-id> --changefeed-id <changefeed-id>
```

### Options

```
  -f, --changefeed-id string   The changefeed ID.
  -c, --cluster-id string      The cluster ID.
  -h, --help                   help for pause
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless changefeed](ticloud_serverless_changefeed.md)	 - Manage TiDB Cloud Serverless changefeeds

