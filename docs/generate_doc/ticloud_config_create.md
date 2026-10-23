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

  To configure a staging profile for Premium / Essential V2 and OAuth:
  $ ticloud config create --profile-name staging --environment staging

  To configure a new user profile in non-interactive mode with api keys:
  $ ticloud config create --profile-name <profile-name> --public-key <public-key> --private-key <private-key>
```

Interactive creation prompts for `prod` or `staging` unless `--environment` is supplied. Non-interactive creation defaults to `prod`. The choice sets `cloud-api-endpoint`, `iam-endpoint`, and `oauth-endpoint`; the OAuth client ID is selected automatically. Legacy `serverless-endpoint` is not configured by this option.

### Options

```
      --environment string    Environment for Premium / Essential V2 and OAuth: prod or staging. (default "prod")
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

