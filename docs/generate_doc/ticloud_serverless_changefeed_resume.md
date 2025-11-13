## ticloud serverless changefeed resume

Resume a paused changefeed

```
ticloud serverless changefeed resume [flags]
```

### Examples

```
  Resume a changefeed in interactive mode:
  $ ticloud serverless changefeed resume

  Resume a changefeed in non-interactive mode:
  $ ticloud serverless changefeed resume -c <cluster-id> --changefeed-id <changefeed-id>
```

### Options

```
  -f, --changefeed-id string   The changefeed ID.
  -c, --cluster-id string      The cluster ID.
  -h, --help                   help for resume
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless changefeed](ticloud_serverless_changefeed.md)	 - Manage TiDB Cloud Serverless changefeeds

