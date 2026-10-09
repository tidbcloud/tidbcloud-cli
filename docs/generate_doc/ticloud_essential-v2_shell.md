## ticloud essential-v2 shell

Connect to a TiDB Cloud Essential V2 instance through a public or private endpoint

### Synopsis

Connect through a public endpoint by default, or an existing private endpoint with --connection-type private-endpoint. With --endpoint host:port, connect directly without discovering or checking the address against API results. Otherwise, discovering private addresses requires private endpoint read permissions and has a 30-second total timeout. Private networking and DNS must already be configured. This command does not create network resources. Private connections have a 30-second initial connection timeout and retain TLS certificate verification; there is no fallback to public connections. For a known public DNS probe failure, the CLI checks that public access is enabled, then warns and attempts a TLS-verified connection with a 30-second initial timeout after password entry.

```
ticloud essential-v2 shell [flags]
```

### Examples

```
  $ ticloud essential-v2 shell
  $ ticloud essential-v2 shell -c <instance-id>
  $ ticloud essential-v2 shell -c <instance-id> --password <password>
  $ ticloud essential-v2 shell -c <instance-id> -u <user-name> --password <password>
  $ ticloud essential-v2 shell --connection-type private-endpoint
  $ ticloud essential-v2 shell -c <instance-id> --connection-type private-endpoint
  $ ticloud essential-v2 shell -c <instance-id> --connection-type private-endpoint --endpoint <host>:4000
```

### Options

```
  -c, --cluster-id string        The ID of the instance.
      --connection-type string   The connection type: public or private-endpoint. Does not configure networking. (default "public")
      --endpoint string          Connect directly to host:port without private endpoint discovery or existence checks. Requires --connection-type private-endpoint; TLS certificate verification remains enabled.
  -h, --help                     help for shell
      --password string          The password of the SQL user.
  -u, --user string              The SQL user name. The default is root.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud essential-v2](ticloud_essential-v2.md)	 - Manage TiDB Cloud Essential V2 instances

