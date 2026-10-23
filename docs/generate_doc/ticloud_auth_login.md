## ticloud auth login

Authenticate with TiDB Cloud using OAuth

### Synopsis

Authenticate with TiDB Cloud using OAuth and save the token to the selected profile.
This command does not validate API keys. A complete public/private API key pair
in the selected profile takes precedence over the OAuth token for API requests.
-P selects a profile for this command only. Use 'ticloud config use <profile-name>'
to change the default profile for subsequent commands.

```
ticloud auth login [flags]
```

### Examples

```
  To log into TiDB Cloud:
  $ ticloud auth login
```

### Options

```
  -h, --help               help for login
      --insecure-storage   Save authentication credentials in plain text instead of credential store.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud auth](ticloud_auth.md)	 - Login and logout via TiDB Cloud API

