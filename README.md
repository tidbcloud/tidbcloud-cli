# TiDB Cloud CLI

The `ticloud` command line tool brings deploy cluster requests, and other TiDB Cloud concepts to your fingertips.

## Table of Contents
* [TiDB Cloud CLI](#tidb-cloud-cli)
    * [Installation](#installation)
        * [macOS and Linux](#macos-and-linux)
            * [Installing via script](#installing-via-script)
            * [Installing via TiUP](#installing-via-tiup)
        * [Manually](#manually)
        * [GitHub Action](#github-action)
    * [Quick Start](#quick-start)
        * [Config a profile](#config-a-profile)
        * [Create a cluster](#create-a-cluster)
    * [Documentation](#documentation)
        * [Set up TiDB Cloud API host](#set-up-tidb-cloud-api-host)
    * [Roadmap](#roadmap)

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

Config a profile with your TiDB Cloud [API key](https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management).

```shell
ticloud config create

# via TiUP
tiup cloud config create
```

> :information_source: The config name **MUST NOT** contain '.'

### Create a cluster

```shell
ticloud cluster create

# via TiUP
tiup cloud cluster create
```

## Documentation

Please check the CLI help for more information.

### Set up TiDB Cloud API host

Usually you don't need to set up the TiDB Cloud API url, the default value is `https://api.tidbcloud.com`.

```shell
ticloud config set api-url https://api.tidbcloud.com

# via TiUP
tiup cloud config set api-url https://api.tidbcloud.com
```

## Roadmap

There are many features we want to add to the CLI.
- CLI supports auth login.
- CLI supports connecting to the TiDB Cloud cluster.
- CLI supports generating connect-string for multiple frameworks and MySQL clients.
- CLI supports generating demo code for using TiDB Cloud cluster.
- CLI supports creating the dedicated cluster(on condition that dedicated support flavor).
- CLI supports import data to TiDB Cloud cluster.
- CLI supports backing up and restoring the TiDB Cloud cluster.
- CLI supports telemetry.
