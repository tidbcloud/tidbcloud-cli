## ticloud serverless audit-log config update

Update the database audit logging configuration

```
ticloud serverless audit-log config update [flags]
```

### Examples

```
  Conigure database audit logging in interactive mode:
  $ ticloud serverless audit-log config update

  Unredact the database audit log in non-interactive mode:
  $ ticloud serverless audit-log config update -c <cluster-id> --unredacted

  Enable database audit logging in non-interactive mode:
  $ ticloud serverless audit-log config update -c <cluster-id> --enabled

  Disable database audit logging in non-interactive mode:
  $ ticloud serverless audit-log config update -c <cluster-id> --enabled=false
```

### Options

```
      --azblob.sas-token string           The SAS token of Azure Blob.
      --azblob.uri string                 The Azure Blob URI in azure://<account>.blob.core.windows.net/<container>/<path> format. Required when cloud storage is AZURE_BLOB.
      --cloud-storage string              The cloud storage. One of ["TIDB_CLOUD" "S3" "GCS" "AZURE_BLOB" "OSS"].
  -c, --cluster-id string                 The ID of the cluster to be updated.
      --enabled                           enable or disable database audit logging.
      --gcs.service-account-key string    The base64 encoded service account key of GCS.
      --gcs.uri string                    The GCS URI in gs://<bucket>/<path> format. Required when cloud storage is GCS.
  -h, --help                              help for update
      --oss.access-key-id string          The access key ID of the OSS.
      --oss.access-key-secret string      The access key secret of the OSS.
      --oss.uri string                    The OSS URI in oss://<bucket>/<path> format. Required when cloud storage is OSS.
      --rotation-interval-minutes int32   The rotation interval in minutes, range [10, 1440].
      --rotation-size-mib int32           The rotation size in MiB, range [1, 1024].
      --s3.access-key-id string           The access key ID of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].
      --s3.role-arn string                The role arn of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].
      --s3.secret-access-key string       The secret access key of the S3. You only need to set one of the s3.role-arn and [s3.access-key-id, s3.secret-access-key].
      --s3.uri string                     The S3 URI in s3://<bucket>/<path> format. Required when cloud storage is S3.
      --unredacted                        unredact or redact the database audit log.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless audit-log config](ticloud_serverless_audit-log_config.md)	 - Manage TiDB Cloud Serverless database audit logging configuration

