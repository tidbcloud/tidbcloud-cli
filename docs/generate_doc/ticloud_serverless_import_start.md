## ticloud serverless import start

Start a data import task

```
ticloud serverless import start [flags]
```

### Examples

```
  Start an import task in interactive mode:
  $ ticloud serverless import start

  Start a local import task in non-interactive mode:
  $ ticloud serverless import start --local.file-path <file-path> --cluster-id <cluster-id> --file-type <file-type> --local.target-database <target-database> --local.target-table <target-table>

  Start a local import task with custom upload concurrency:
  $ ticloud serverless import start --local.file-path <file-path> --cluster-id <cluster-id> --file-type <file-type> --local.target-database <target-database> --local.target-table <target-table> --local.concurrency 10
	
  Start a local import task with custom CSV format:
  $ ticloud serverless import start --local.file-path <file-path> --cluster-id <cluster-id> --file-type CSV --local.target-database <target-database> --local.target-table <target-table> --csv.separator \" --csv.delimiter \' --csv.backslash-escape=false --csv.trim-last-separator=true

```

### Options

```
      --azure.blob-sas-url string      The SAS URL for Azure Blob.
  -c, --cluster-id string              Cluster ID.
      --csv.backslash-escape           Specifies whether to parse backslash inside fields as escape characters in the CSV file. (default true)
      --csv.delimiter string           The delimiter used for quoting of CSV file. (default "\"")
      --csv.not-null                   Specifies whether a CSV file can contain any NULL values.
      --csv.null-value string          The representation of NULL values in the CSV file. (default "\\N")
      --csv.separator string           The field separator of CSV file. (default ",")
      --csv.skip-header                Specifies whether the CSV file contains a header line.
      --csv.trim-last-separator        Specifies whether to treat separator as the line terminator and trim all trailing separators in the CSV file.
      --file-type string               The import file type, one of ["CSV"].
      --gcs.credentials-path string    The local path of GCS credentials.
      --gcs.uri string                 The GCS folder URI for import.
  -h, --help                           help for start
      --local.concurrency int          The concurrency for uploading file. (default 5)
      --local.file-path string         The local file path to import.
      --local.target-database string   Target database to which import data.
      --local.target-table string      Target table to which import data.
      --s3.access-key-id string        The access key ID for S3.
      --s3.role-arn string             The role ARN for S3.
      --s3.secret-access-key string    The secret access key for S3.
      --s3.uri string                  The S3 folder URI for import.
      --source-type string             The import source type, one of ["S3" "LOCAL" "GCS" "AzBlob"]. (default "LOCAL")
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless import](ticloud_serverless_import.md)	 - Manage TiDB Serverless data imports

