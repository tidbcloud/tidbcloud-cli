## ticloud config create

Configure a user profile to store settings

### Synopsis

Configure a user profile to store settings, where profile names are case-insensitive and do not contain periods.

```
ticloud config create [flags]
```

### Examples

```
  To configure a new user profile in interactive mode:
  $ ticloud config create

  To configure a new user profile in non-interactive mode:
  $ ticloud config create --profile-name <profile-name>

  To configure a new user profile in non-interactive mode with api keys:
  $ ticloud config create --profile-name <profile-name> --public-key <public-key> --private-key <private-key>
```

### Options

```
  -h, --help                  help for create
      --private-key string    The private key of the TiDB Cloud API. (optional)
      --profile-name string   The name of the profile, must not contain periods.
      --public-key string     The public key of the TiDB Cloud API. (optional)
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud config](ticloud_config.md)	 - Configure and manage your user profiles

