# TiDB Cloud CLI


The `ticloud` command line tool brings deploy cluster requests, and other TiDB Cloud concepts to your fingertips.

## Table of Contents
- [TiDB Cloud CLI](#tidb-cloud-cli)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
    - [macOS and Linux](#macos-and-linux)
      - [Installing via script](#installing-via-script)
      - [Installing via TiUP](#installing-via-tiup)
    - [Manually](#manually)
    - [GitHub Action](#github-action)
  - [Quick Start](#quick-start)
    - [Config a profile](#config-a-profile)
    - [Create a cluster](#create-a-cluster)
  - [Documentation](#documentation)
  - [Telemetry](#telemetry)
  - [Roadmap](#roadmap)

## Installation

### macOS and Linux

#### Installing via script

```shell
curl https://raw.githubusercontent.com/tidbcloud/tidbcloud-cli/main/install.sh | sh
```

#### Installing via [TiUP](https://tiup.io/)

```shell
tiup install cloud
```

### Manually

Download the pre-compiled binaries from the [releases](https://github.com/tidbcloud/tidbcloud-cli/releases/latest) page and copy to the desired location.

### GitHub Action

To set up `ticloud` in GitHub Action, use [`setup-tidbcloud-cli`](https://github.com/tidbcloud/setup-tidbcloud-cli).

## Quick Start

In order to use the `ticloud` CLI, you need to have a TiDB Cloud account. If you don't have one, you can sign up for a free trial [here](https://tidbcloud.com/).

### Config a profile

Config a profile with your TiDB Cloud [API key](https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management) or log in via OAuth.

```shell
ticloud config create

# via TiUP
tiup cloud config create

# via OAuth
ticloud auth login
```

> :information_source: The config name **MUST NOT** contain '.'

### Create a cluster

```shell
ticloud serverless create

# via TiUP
tiup cloud serverless create
```

## Documentation

Please check the CLI help for more information.

## Telemetry

See [here](/docs/telemetry.md)

## Roadmap

There are many features we want to add to the CLI.
- CLI supports generating demo code for using TiDB Cloud cluster.
- CLI supports managing the dedicated cluster.
- CLI supports backing up and restoring the TiDB Cloud cluster.
- CLI supports connecting in shell command without SQL user.
