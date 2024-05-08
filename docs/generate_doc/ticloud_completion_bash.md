## ticloud completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(ticloud completion bash)

To load completions for every new session, execute once:

#### Linux:

	ticloud completion bash > /etc/bash_completion.d/ticloud

#### macOS:

	ticloud completion bash > $(brew --prefix)/etc/bash_completion.d/ticloud

You will need to start a new shell for this setup to take effect.


```
ticloud completion bash
```

### Options

```
  -h, --help              help for bash
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

