## ticloud premium password

Reset the root password of a TiDB Cloud Premium instance

```
ticloud premium password [flags]
```

### Examples

```
  Reset the root password interactively:
  $ ticloud premium password
  $ ticloud premium password -c <instance-id>

  Reset the root password non-interactively:
  $ ticloud premium password -c <instance-id> --password <password>
```

### Options

```
  -c, --cluster-id string   The ID of the instance. If omitted, select an instance interactively.
  -h, --help                help for password
      --password string     The new root password. Prefer interactive input to avoid shell history exposure.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud premium](ticloud_premium.md)	 - Manage TiDB Cloud Premium instances

