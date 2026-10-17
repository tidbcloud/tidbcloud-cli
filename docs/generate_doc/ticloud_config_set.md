## ticloud config set

Configure specific properties of the active profile

### Synopsis

Configure specific properties of the active profile.
Available properties : [public-key private-key serverless-endpoint nextgen-endpoint nextgen-ca-cert-path iam-endpoint oauth-endpoint oauth-client-id oauth-client-secret telemetry-enabled].

If using -P flag, the config in the specific profile will be set.
If not, the config in the active profile will be set

```
ticloud config set <property-name> <value> [flags]
```

### Examples

```
  Set the value of the public-key in active profile:
  $ ticloud config set public-key <public-key>

  Set the value of the public-key in the specific profile "test":
  $ ticloud config set public-key <public-key> -P test

  Configure a custom NextGen endpoint and CA certificate for profile "test":
  $ ticloud config set nextgen-endpoint https://gateway.example.com -P test
  $ ticloud config set nextgen-ca-cert-path /absolute/path/to/ca.pem -P test
```

### Options

```
  -h, --help   help for set
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud config](ticloud_config.md)	 - Configure and manage your user profiles

