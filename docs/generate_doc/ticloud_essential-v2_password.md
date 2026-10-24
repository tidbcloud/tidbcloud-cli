## ticloud essential-v2 password

Reset the root password of a TiDB Cloud Essential V2 instance

```
ticloud essential-v2 password [flags]
```

### Examples

```
  Reset the root password interactively:
  $ ticloud essential-v2 password
  $ ticloud essential-v2 password -c <instance-id>

  Reset the root password from a protected file (or pipe from a secret manager):
  $ ticloud essential-v2 password -c <instance-id> --password-stdin < /path/to/password-file
```

### Options

```
  -c, --cluster-id string   The ID of the instance. If omitted, select an instance interactively.
  -h, --help                help for password
      --password-stdin      Read the new root password from a pipe or file on standard input until EOF; terminal input is rejected. Requires --cluster-id; removes one trailing LF or CRLF. If omitted, prompt with hidden input and confirmation.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud essential-v2](ticloud_essential-v2.md)	 - Manage TiDB Cloud Essential V2 instances

