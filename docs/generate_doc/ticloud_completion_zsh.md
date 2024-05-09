## ticloud completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(ticloud completion zsh); compdef _ticloud ticloud

To load completions for every new session, execute once:

#### Linux:

	ticloud completion zsh > "${fpath[1]}/_ticloud"

#### macOS:

	ticloud completion zsh > $(brew --prefix)/share/zsh/site-functions/_ticloud

You will need to start a new shell for this setup to take effect.


```
ticloud completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
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

