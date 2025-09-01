## ticloud serverless sql-user update

Update a SQL user

```
ticloud serverless sql-user update [flags]
```

### Examples

```
  Update a SQL user in interactive mode:
  $ ticloud serverless sql-user update

  Update a SQL user in non-interactive mode:
  $ ticloud serverless sql-user update -c <cluster-id> --user <user-name> --password <password> --role <role>
```

### Options

```
      --add-role strings      The role(s) to be added to the SQL user.
  -c, --cluster-id string     The cluster ID of the SQL user to be updated.
      --delete-role strings   The role(s) to be deleted from the SQL user.
  -h, --help                  help for update
      --password string       The new password of the SQL user.
      --role strings          The new role(s) of the SQL user. Passing this flag replaces preexisting data, supported roles ["role_admin" "role_readwrite" "role_readonly"]
  -u, --user string           The name of the SQL user to be updated.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless sql-user](ticloud_serverless_sql-user.md)	 - Manage cluster SQL users

