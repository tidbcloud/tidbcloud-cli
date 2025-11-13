## ticloud serverless changefeed template

Show changefeed Kafka, MySQL and Filter JSON templates

```
ticloud serverless changefeed template [flags]
```

### Examples

```
  Show all changefeed templates:
  $ ticloud serverless changefeed template

  Show Kafka JSON template:
  $ ticloud serverless changefeed template --type kafka
```

### Options

```
      --explain       Show template with explanations.
  -h, --help          help for template
      --type string   The type of changefeed template to show, one of ["kafka", "mysql", "filter"].
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless changefeed](ticloud_serverless_changefeed.md)	 - Manage TiDB Cloud Serverless changefeeds

