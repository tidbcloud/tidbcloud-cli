# TiDB Cloud CLI

The `ticloud` command line tool brings deploy cluster requests, and other TiDB Cloud concepts to your fingertips.

## Installation

#### macOS and Linux

```
curl https://raw.githubusercontent.com/tidbcloud/tidbcloud-cli/main/install.sh | sh
```

#### Manually

Download the pre-compiled binaries from the [releases](https://github.com/tidbcloud/tidbcloud-cli/releases/latest) page and copy to the desired location.

## Quick Start

In order to use the `ticloud` CLI, you need to have a TiDB Cloud account. If you don't have one, you can sign up for a free trial [here](https://tidbcloud.com/).

#### Config a profile

Config a profile with your TiDB Cloud [API key](https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management).

```
ticloud config init
```

#### Create a cluster

```
ticloud cluster create
```

## Documentation

Please check the CLI help for more information.

Documentation page is on the way.

## Roadmap

There are many features we want to add to the CLI.
- CLI Supports auth login.
- CLI Supports connect to TiDB Cloud cluster.
- CLI Supports generate connect-string for multiple frameworks and mysql clients.
- CLI Supports generate demo code for using TiDB Cloud cluster.
- CLI Supports create the dedicated cluster(on condition that dedicated support flavor).
- CLI Supports import data to TiDB Cloud cluster.
- CLI Supports backup and restore TiDB Cloud cluster.
