## ticloud serverless audit-log download

Download the database audit log files

```
ticloud serverless audit-log download [flags]
```

### Examples

```
  Download the database audit logs in interactive mode:
  $ ticloud serverless audit-log download

  Download the database audit logs in non-interactive mode:
  $ ticloud serverless audit-log download -c <cluster-id> --start-date <start-date> --end-date <end-date>
```

### Options

```
  -c, --cluster-id string    Cluster ID.
      --concurrency int      Download concurrency. (default 3)
      --end-date string      The end date of the audit log you want to download in the format of 'YYYY-MM-DD', e.g. '2025-01-01'.
      --force                Download without confirmation.
  -h, --help                 help for download
      --output-path string   The path where you want to download to. If not specified, download to the current directory.
      --start-date string    The start date of the audit log you want to download in the format of 'YYYY-MM-DD', e.g. '2025-01-01'.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless audit-log](ticloud_serverless_audit-log.md)	 - Manage TiDB Cloud Serverless database audit logging

