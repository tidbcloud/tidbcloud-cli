## ticloud premium public-endpoint enable

Enable the public endpoint of a TiDB Cloud Premium instance

### Synopsis

Enable the public endpoint of a TiDB Cloud Premium instance, preserving its existing IP access list. No IP access list entries are added. The request is asynchronous; use --describe to check endpoint readiness after it is accepted.

```
ticloud premium public-endpoint enable [flags]
```

### Examples

```
  $ ticloud premium public-endpoint enable
  $ ticloud premium public-endpoint enable -c <instance-id>
```

### Options

```
  -c, --cluster-id string   The ID of the instance. If omitted, select an instance interactively.
  -h, --help                help for enable
  -o, --output string       Output format, one of ["human" "json"]. For the complete result, please use json format. (default "human")
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud premium public-endpoint](ticloud_premium_public-endpoint.md)	 - Manage the public endpoint of a TiDB Cloud Premium instance

