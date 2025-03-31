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

  Start an S3 import task in non-interactive mode:
  $ ticloud serverless import start --source-type S3 --s3.uri <s3-uri> --cluster-id <cluster-id> --file-type <file-type> --s3.role-arn <role-arn>

  Start a GCS import task in non-interactive mode:
  $ ticloud serverless import start --source-type GCS --gcs.uri <gcs-uri> --cluster-id <cluster-id> --file-type <file-type> --gcs.service-account-key <service-account-key>

  Start an Azure Blob import task in non-interactive mode:
  $ ticloud serverless import start --source-type AZURE_BLOB --azblob.uri <azure-blob-uri> --cluster-id <cluster-id> --file-type <file-type> --azblob.sas-token <sas-token>

```

### Options

```
      --azblob.sas-token string          The SAS token of Azure Blob.
      --azblob.uri string                The Azure Blob URI in azure://<account>.blob.core.windows.net/<container>/<path> format.
  -c, --cluster-id string                Cluster ID.
      --csv.backslash-escape             Specifies whether to interpret backslash escapes inside fields in the CSV file. (default true)
      --csv.delimiter string             The delimiter used for quoting of CSV file. (default """)
      --csv.not-null                     Specifies whether a CSV file can contain any NULL values.
      --csv.null-value string            The representation of NULL values in the CSV file. (default "\N")
      --csv.separator string             The field separator of CSV file. (default ",")
      --csv.skip-header                  Specifies whether the CSV file contains a header line.
      --csv.trim-last-separator          Specifies whether to treat separator as the line terminator and trim all trailing separators in the CSV file.
      --file-type string                 The import file type, one of ["CSV" "SQL" "AURORA_SNAPSHOT" "PARQUET"].
      --gcs.service-account-key string   The base64 encoded service account key of GCS.
      --gcs.uri string                   The GCS URI in gs://<bucket>/<path> format. Required when source type is GCS.
  -h, --help                             help for start
      --local.concurrency int            The concurrency for uploading file. (default 5)
      --local.file-path string           The local file path to import.
      --local.target-database string     Target database to which import data.
      --local.target-table string        Target table to which import data.
      --oss.access-key-id string         The AccessKey ID of the Alibaba Cloud OSS.
      --oss.access-key-secret string     The AccessKey Secret of the Alibaba Cloud OSS.
      --oss.uri string                   The OSS URI in oss://<bucket>/<path> format. Required when source type is Alibaba Cloud OSS.
      --s3.access-key-id string          The access key ID of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].
      --s3.role-arn string               The role arn of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].
      --s3.secret-access-key string      The secret access key of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].
      --s3.uri string                    The S3 URI in s3://<bucket>/<path> format. Required when source type is S3.
      --source-type string               The import source type, one of ["LOCAL" "S3" "GCS" "AZURE_BLOB" "OSS"]. (default "LOCAL")
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless import](ticloud_serverless_import.md)	 - Manage TiDB Cloud Serverless data imports

