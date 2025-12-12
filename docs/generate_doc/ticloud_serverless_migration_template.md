## ticloud serverless migration template

Show migration JSON templates

```
ticloud serverless migration template [flags]
```

### Examples

```
  Show the ALL mode migration template:
  $ ticloud serverless migration template --mode all

  Show the INCREMENTAL migration template:
  $ ticloud serverless migration template --mode incremental
```

### Options

```
  -h, --help          help for template
      --mode string   Migration mode template to show, one of [all, incremental].
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud serverless migration](ticloud_serverless_migration.md)	 - Manage TiDB Cloud Serverless migrations

