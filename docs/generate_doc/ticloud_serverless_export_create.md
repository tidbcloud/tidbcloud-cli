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

  Export all data with s3 type in non-interactive mode:
  $ ticloud serverless export create -c <cluster-id> --target-type S3 --s3.uri <s3-uri> --s3.access-key-id <access-key-id> --s3.secret-access-key <secret-access-key>

  Export all data and customize csv format in non-interactive mode:
  $ ticloud serverless export create -c <cluster-id> --file-type CSV --csv.separator ";" --csv.delimiter "\"" --csv.null-value 'NaN' --csv.skip-header

  Export test.t1 and test.t2 in non-interactive mode:
  $ ticloud serverless export create -c <cluster-id> --filter 'test.t1,test.t2'

  Export tables with special characters, for example, if you want to export `test,`.`t1` and `"test`.`t1`:
  $ ticloud serverless export create -c <cluster-id> --filter '"`test1,`.t1","`""test`.t1"'
```

### Options

```
  -c, --cluster-id string             The ID of the cluster, in which the export will be created.
      --compression string            The compression algorithm of the export file. One of ["GZIP" "SNAPPY" "ZSTD" "NONE"]. (default "GZIP")
      --csv.delimiter string          Delimiter of string type variables in CSV files. (default "\"")
      --csv.null-value string         Representation of null values in CSV files. (default "\\N")
      --csv.separator string          Separator of each value in CSV files. (default ",")
      --csv.skip-header               Export CSV files of the tables without header.
      --file-type string              The export file type. One of ["CSV" "SQL"]. (default "SQL")
      --filter strings                Specify the exported table(s) with table filter patterns. See https://docs.pingcap.com/tidb/stable/table-filter to learn table filter.
      --force                         Create without confirmation. You need to confirm when you want to export the whole cluster in non-interactive mode.
  -h, --help                          help for create
      --s3.access-key-id string       The access key ID of the S3. Required when target type is S3.
      --s3.secret-access-key string   The secret access key of the S3. Required when target type is S3.
      --s3.uri string                 The s3 uri in s3://<bucket>/<path> format. Required when target type is S3.
      --sql string                    Filter the exported data with SQL SELECT statement.
      --target-type string            The export target. One of ["LOCAL" "S3"]. (default "LOCAL")
      --where string                  Filter the exported table(s) with the where condition.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless export](ticloud_serverless_export.md)	 - Manage TiDB Serverless exports

