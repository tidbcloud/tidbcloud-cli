# TiDB Cloud CLI

The `ticloud` command line tool brings deploy cluster requests, and other TiDB Cloud concepts to your fingertips.

## Installation

#### macOS

- amd64
```
curl -L https://github.com/tidbcloud/tidbcloud-cli/releases/download/v0.1.0-rc1/ticloud_0.1.0-rc1_macos_x86_64.tar.gz | tar -xz 
```

- arm64
```
curl -L https://github.com/tidbcloud/tidbcloud-cli/releases/download/v0.1.0-rc1/ticloud_0.1.0-rc1_macos_arm64.tar.gz | tar -xz
```

#### Linux

- amd64
```
curl -L https://github.com/tidbcloud/tidbcloud-cli/releases/download/v0.1.0-rc1/ticloud_0.1.0-rc1_linux_x86_64.tar.gz | tar -xz
```

- arm64
```
curl -L https://github.com/tidbcloud/tidbcloud-cli/releases/download/v0.1.0-rc1/ticloud_0.1.0-rc1_linux_arm64.tar.gz | tar -xz
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
