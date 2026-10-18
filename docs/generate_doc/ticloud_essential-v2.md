## ticloud essential-v2

Manage TiDB Cloud Essential V2 instances

### Synopsis

Manage TiDB Cloud Essential V2 instances.

Action-specific flags:

--create

      -n, --display-name string   The display name of the instance.
          --encryption string     Dual-layer data encryption, one of ["none" "default-key"]. (default "none")
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
  $ ticloud essential-v2 --create --display-name <name> --region <region-id> --max-rcu <max-rcu>
  $ ticloud essential-v2 --create --project-id <project-id> --display-name <name> --region <region-id> --max-rcu <max-rcu>
  $ ticloud essential-v2 --list
  $ ticloud essential-v2 --describe -c <instance-id>
  $ ticloud essential-v2 --update -c <instance-id> --max-rcu <max-rcu>
  $ ticloud essential-v2 --delete -c <instance-id>
  $ ticloud essential-v2 public-endpoint enable -c <instance-id>
  $ ticloud essential-v2 public-endpoint disable -c <instance-id>
  $ ticloud essential-v2 password -c <instance-id>
```

### Options

```
      --create     Create a TiDB Cloud Essential V2 instance.
      --delete     Delete a TiDB Cloud Essential V2 instance.
      --describe   Describe a TiDB Cloud Essential V2 instance.
  -h, --help       help for essential-v2
      --list       List TiDB Cloud Essential V2 instances.
      --update     Update a TiDB Cloud Essential V2 instance.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud](ticloud.md)	 - CLI tool to manage TiDB Cloud
* [ticloud essential-v2 password](ticloud_essential-v2_password.md)	 - Reset the root password of a TiDB Cloud Essential V2 instance
* [ticloud essential-v2 public-endpoint](ticloud_essential-v2_public-endpoint.md)	 - Manage the public endpoint of a TiDB Cloud Essential V2 instance
* [ticloud essential-v2 region](ticloud_essential-v2_region.md)	 - List regions available for TiDB Cloud Essential V2
* [ticloud essential-v2 shell](ticloud_essential-v2_shell.md)	 - Connect to a TiDB Cloud Essential V2 instance through a public or private endpoint

