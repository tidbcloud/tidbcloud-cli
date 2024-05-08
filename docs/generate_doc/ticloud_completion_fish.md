## ticloud completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	ticloud completion fish | source

To load completions for every new session, execute once:

	ticloud completion fish > ~/.config/fish/completions/ticloud.fish

You will need to start a new shell for this setup to take effect.


```
ticloud completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -D, --debug            Enable debug mode
      --no-color         Disable color output
  -P, --profile string   Profile to use from your configuration file
```

### SEE ALSO

* [ticloud completion](ticloud_completion.md)	 - Generate the autocompletion script for the specified shell

