## ticloud serverless export download

Download the exported data

```
ticloud serverless export download [flags]
```

### Examples

```
  Download the local type export in interactive mode:
  $ ticloud serverless export download

  Download the local type export in non-interactive mode:
  $ ticloud serverless export download -c <cluster-id> -e <export-id>
```

### Options

```
  -c, --cluster-id string    The cluster ID of the export.
      --concurrency int      Download concurrency. (default 3)
  -e, --export-id string     The ID of the export.
      --force                Download without confirmation.
  -h, --help                 help for download
      --output-path string   Where you want to download to. If not specified, download to the current directory.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless export](ticloud_serverless_export.md)	 - Manage TiDB Cloud Serverless exports

