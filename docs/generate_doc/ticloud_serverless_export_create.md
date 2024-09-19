## ticloud serverless export create

Export data from a TiDB Serverless cluster

```
ticloud serverless export create [flags]
```

### Examples

```
  Create an export in interactive mode:
  $ ticloud serverless export create

  Export all data with local type in non-interactive mode:
  $ ticloud serverless export create -c <cluster-id>

  Export all data with S3 type in non-interactive mode:
  $ ticloud serverless export create -c <cluster-id> --target-type S3 --s3.uri <s3-uri> --s3.access-key-id <access-key-id> --s3.secret-access-key <secret-access-key>

  Export all data and customize CSV format in non-interactive mode:
  $ ticloud serverless export create -c <cluster-id> --file-type CSV --csv.separator ";" --csv.delimiter "\"" --csv.null-value 'NaN' --csv.skip-header

  Export test.t1 and test.t2 in non-interactive mode:
  $ ticloud serverless export create -c <cluster-id> --filter 'test.t1,test.t2'

  Export tables with special characters, for example, if you want to export `test,`.`t1` and `"test`.`t1`:
  $ ticloud serverless export create -c <cluster-id> --filter '"`test1,`.t1","`""test`.t1"'
```

### Options

```
      --azblob.sas-token string          The SAS token of Azure Blob.
      --azblob.uri string                The Azure Blob URI in azure://<account>.blob.core.windows.net/<container>/<path> format. Required when target type is AZURE_BLOB.
  -c, --cluster-id string                The ID of the cluster, in which the export will be created.
      --compression string               The compression algorithm of the export file. One of ["GZIP" "SNAPPY" "ZSTD" "NONE"].
      --csv.delimiter string             Delimiter of string type variables in CSV files. (default "\"")
      --csv.null-value string            Representation of null values in CSV files. (default "\\N")
      --csv.separator string             Separator of each value in CSV files. (default ",")
      --csv.skip-header                  Export CSV files of the tables without header.
      --display-name string              The display name of the export. default is SNAPSHOT_{snapshot_time}.
      --file-type string                 The export file type. One of ["SQL" "CSV" "PARQUET"]. (default "CSV")
      --filter strings                   Specify the exported table(s) with table filter patterns. See https://docs.pingcap.com/tidb/stable/table-filter to learn table filter.
      --force                            Create without confirmation. You need to confirm when you want to export the whole cluster in non-interactive mode.
      --gcs.service-account-key string   The base64 encoded service account key of GCS.
      --gcs.uri string                   The GCS URI in gcs://<bucket>/<path> format. Required when target type is GCS.
  -h, --help                             help for create
      --parquet.compression string       The parquet compression algorithm. One of ["GZIP" "SNAPPY" "ZSTD" "NONE"]. (default "ZSTD")
      --s3.access-key-id string          The access key ID of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].
      --s3.role-arn string               The role arn of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].
      --s3.secret-access-key string      The secret access key of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].
      --s3.uri string                    The S3 URI in s3://<bucket>/<path> format. Required when target type is S3.
      --sql string                       Filter the exported data with SQL SELECT statement.
      --target-type string               The export target. One of ["LOCAL" "S3" "GCS" "AZURE_BLOB"]. (default "LOCAL")
      --where string                     Filter the exported table(s) with the where condition.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless export](ticloud_serverless_export.md)	 - Manage TiDB Serverless exports

