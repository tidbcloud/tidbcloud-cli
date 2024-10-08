## ticloud serverless export describe

Describe an export task

```
ticloud serverless export describe [flags]
```

### Examples

```
  Get an export in interactive mode:
  $ ticloud serverless export describe

  Get an export in non-interactive mode:
  $ ticloud serverless export describe -c <cluster-id> -e <export-id>
```

### Options

```
  -c, --cluster-id string   The cluster ID of the export.
  -e, --export-id string    The ID of the export.
  -h, --help                help for describe
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless export](ticloud_serverless_export.md)	 - Manage TiDB Cloud Serverless exports

