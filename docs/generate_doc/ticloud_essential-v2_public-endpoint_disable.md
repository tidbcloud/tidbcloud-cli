## ticloud essential-v2 public-endpoint disable

Disable the public endpoint of a TiDB Cloud Essential V2 instance

### Synopsis

Disable the public endpoint of a TiDB Cloud Essential V2 instance, preserving its existing IP access list. No IP access list entries are added. The request is asynchronous; use --describe to check endpoint readiness after it is accepted. Disabling the public endpoint interrupts public connections.

```
ticloud essential-v2 public-endpoint disable [flags]
```

### Examples

```
  $ ticloud essential-v2 public-endpoint disable
  $ ticloud essential-v2 public-endpoint disable -c <instance-id>
  $ ticloud essential-v2 public-endpoint disable -c <instance-id> --force
```

### Options

```
  -c, --cluster-id string   The ID of the instance. If omitted, select an instance interactively.
      --force               Disable the public endpoint without confirmation.
  -h, --help                help for disable
  -o, --output string       Output format, one of ["human" "json"]. For the complete result, please use json format. (default "human")
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud essential-v2 public-endpoint](ticloud_essential-v2_public-endpoint.md)	 - Manage the public endpoint of a TiDB Cloud Essential V2 instance

