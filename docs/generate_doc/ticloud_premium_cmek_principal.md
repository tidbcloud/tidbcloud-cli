## ticloud premium cmek principal

Get the IAM principal for your KMS key policy

### Synopsis

Get the IAM principal required to grant TiDB Cloud access to your AWS or Alibaba Cloud KMS key. Configure the key policy before creating an instance with --encryption cmek.

```
ticloud premium cmek principal [flags]
```

### Options

```
  -h, --help            help for principal
  -o, --output string   Output format, one of ["human" "json"]. For the complete result, please use json format. (default "human")
  -r, --region string   The region ID, for example aws-us-west-2.
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud premium cmek](ticloud_premium_cmek.md)	 - Configure customer-managed encryption for Premium instances

