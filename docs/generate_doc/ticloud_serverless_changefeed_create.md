## ticloud serverless changefeed create

Create a changefeed

```
ticloud serverless changefeed create [flags]
```

### Examples

```
  Create a changefeed in interactive mode:
  $ ticloud serverless changefeed create

  Create a changefeed in non-interactive mode:
  $ ticloud serverless changefeed create -c <cluster-id> --type KAFKA --kafka '{"network_info":{"network_type":"PUBLIC"},"broker":{"kafka_version":"VERSION_2XX","broker_endpoints":"52.34.156.155:9092","compression":"NONE"},"authentication":{"auth_type":"DISABLE"},"topic_partition_config":{"dispatch_type":"ONE_TOPIC","default_topic":"default-topic","replication_factor":1,"partition_num":1,"partition_dispatchers":[{"partition_type":"TABLE","matcher":["*.*"]}]},"data_format":{"protocol":"CANAL_JSON"}}' --filter '{"filterRule":["test.*"], "mode": "IGNORE_NOT_SUPPORT_TABLE"}'

  Create a changefeed named "myfeed" with specified start TSO in non-interactive mode:
  $ ticloud serverless changefeed create -c <cluster-id> --name myfeed --type KAFKA --kafka <kafka-json> --filter <filter-json> --start-tso 458996254096228352

```

### Options

```
  -c, --cluster-id string     The ID of the cluster.
  -n, --display-name string   The name of the changefeed.
      --filter string         Filter in JSON format, use "ticloud serverless changefeed template" to see templates.
  -h, --help                  help for create
      --kafka string          Kafka information in JSON format, use "ticloud serverless changefeed template" to see templates.
      --mysql string          MySQL information in JSON format, use "ticloud serverless changefeed template" to see templates.
      --start-time string     Start Time for the changefeed (RFC3339 format, e.g., 2024-01-01T00:00:00Z). If both start-tso and start-time are provided, start-tso will be used.
      --start-tso uint        Start TSO for the changefeed, default to current TSO. See https://docs.pingcap.com/tidb/stable/tso/ for more information about TSO.
      --type string           The type of the changefeed, one of ["KAFKA" "MYSQL"]
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless changefeed](ticloud_serverless_changefeed.md)	 - Manage TiDB Cloud Serverless changefeeds

