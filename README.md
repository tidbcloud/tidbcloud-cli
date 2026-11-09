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
  - [Credential storage](#credential-storage)
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

## Credential storage

API private keys are stored in `~/.ticloud/config` by default for compatibility
with existing scripts and environments without an OS credential store. On
Unix-like systems, the CLI restricts the directory to `0700` and the config
file to `0600`. These permissions do not encrypt the file.

To store a new profile's API private key in the OS credential store, use:

```shell
ticloud config create --secure-storage
```

To move an existing profile's private key to the credential store, set it again:

```shell
ticloud config set private-key <private-key> --secure-storage
```

Add `-P <profile-name>` to select a different profile. A successful secure write
removes the plaintext private key from the config. The public key and a
`private-key-storage = "keyring"` marker remain in the config. Existing profiles
are not migrated automatically. Later private-key updates continue to use the
credential store even without `--secure-storage`, and deleting the profile
also removes its stored private key. Secure writes and reads report credential
store errors rather than silently falling back to plaintext or OAuth.

Secure storage requires a working OS credential store; copying the config alone
does not transfer the private key. To clear an API private key, use
`ticloud config set private-key ''`.

OAuth tokens use the OS credential store by default. `auth login --insecure-storage`
explicitly enables plaintext token storage. When both stores
contain a token, the credential store takes precedence; plaintext remains usable
when no token can be loaded from the credential store.

## Telemetry

See [here](/docs/telemetry.md)

## Roadmap

There are many features we want to add to the CLI.
- CLI supports generating demo code for using TiDB Cloud cluster.
- CLI supports managing the dedicated cluster.
- CLI supports backing up and restoring the TiDB Cloud cluster.
- CLI supports connecting in shell command without SQL user.
