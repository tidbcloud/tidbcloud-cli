## ticloud serverless private-link-connection create

Create a private link connection

```
ticloud serverless private-link-connection create [flags]
```

### Examples

```
  Create a private link connection (interactive):
  $ ticloud serverless private-link-connection create

  Create a private link connection which connect to alicloud endpoint service (non-interactive):
  $ ticloud serverless private-link-connection create -c <cluster-id> --display-name <name> --type ALICLOUD_ENDPOINT_SERVICE --alicloud.endpoint-service.name <name>

  Create a private link connection which connect to aws endpoint service (non-interactive):
  $ ticloud serverless private-link-connection create -c <cluster-id> --display-name <name> --type AWS_ENDPOINT_SERVICE --aws.endpoint-service.name <name>
```

### Options

```
      --alicloud.endpoint-service.name string   Alicloud endpoint service name
      --aws.endpoint-service.name string        AWS endpoint service name
      --aws.endpoint-service.region string      AWS endpoint service region
  -c, --cluster-id string                       The cluster ID.
      --display-name string                     Display name for the private link connection.
  -h, --help                                    help for create
      --type string                             Type of the private link connection, one of ["AWS_ENDPOINT_SERVICE" "ALICLOUD_ENDPOINT_SERVICE"]
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless private-link-connection](ticloud_serverless_private-link-connection.md)	 - Manage private link connections

