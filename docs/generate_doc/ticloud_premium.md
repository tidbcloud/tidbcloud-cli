## ticloud premium

Manage TiDB Cloud Premium instances

### Synopsis

Manage TiDB Cloud Premium instances.

Action-specific flags:

--create

          --cmek-key-arn string   The AWS or Alibaba Cloud KMS key ARN. Required with --encryption cmek.
      -n, --display-name string   The display name of the instance.
          --encryption string     Dual-layer data encryption, one of ["none" "default-key" "cmek"]. (default "none")
          --max-rcu int           The maximum number of Request Capacity Units (RCUs).
      -o, --output string         Output format, one of ["human" "json"]. For the complete result, please use json format. (default "human")
      -p, --project-id string     The ID of the project in which the instance will be created. If omitted, the default TiDB X project is used.
      -r, --region string         The region ID, for example aws-us-west-2.

--list

      -o, --output string   Output format, one of ["human" "json"]. For the complete result, please use json format. (default "human")

--describe

      -c, --cluster-id string   The ID of the instance.

--update

      -c, --cluster-id string     The ID of the instance.
      -n, --display-name string   The new display name of the instance.
          --max-rcu int           The new maximum number of Request Capacity Units (RCUs).
      -o, --output string         Output format, one of ["human" "json"]. For the complete result, please use json format. (default "human")

--delete

      -c, --cluster-id string   The ID of the instance to delete.
          --force               Delete the instance without confirmation.

### Examples

```
  $ ticloud premium --create --display-name <name> --region <region-id> --max-rcu <max-rcu>
  $ ticloud premium --create --project-id <project-id> --display-name <name> --region <region-id> --max-rcu <max-rcu>
  $ ticloud premium --list
  $ ticloud premium --describe -c <instance-id>
  $ ticloud premium --update -c <instance-id> --max-rcu <max-rcu>
  $ ticloud premium --delete -c <instance-id>
  $ ticloud premium public-endpoint enable -c <instance-id>
  $ ticloud premium public-endpoint disable -c <instance-id>
  $ ticloud premium password -c <instance-id>
```

### Options

```
      --create     Create a TiDB Cloud Premium instance.
      --delete     Delete a TiDB Cloud Premium instance.
      --describe   Describe a TiDB Cloud Premium instance.
  -h, --help       help for premium
      --list       List TiDB Cloud Premium instances.
      --update     Update a TiDB Cloud Premium instance.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud](ticloud.md)	 - CLI tool to manage TiDB Cloud
* [ticloud premium cmek](ticloud_premium_cmek.md)	 - Configure customer-managed encryption for Premium instances
* [ticloud premium password](ticloud_premium_password.md)	 - Reset the root password of a TiDB Cloud Premium instance
* [ticloud premium public-endpoint](ticloud_premium_public-endpoint.md)	 - Manage the public endpoint of a TiDB Cloud Premium instance
* [ticloud premium region](ticloud_premium_region.md)	 - List regions available for TiDB Cloud Premium
* [ticloud premium shell](ticloud_premium_shell.md)	 - Connect to a TiDB Cloud Premium instance through a public or private endpoint

