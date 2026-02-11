# TiDB Cloud CLI Knowledge

## TiDB Cloud

TiDB Cloud is a fully-managed Database-as-a-Service (DBaaS) that brings TiDB, an open-source Hybrid Transactional and Analytical Processing (HTAP) database, to your cloud. TiDB Cloud offers an easy way to deploy and manage databases to let you focus on your applications, not the complexities of the databases. You can create TiDB Cloud clusters to quickly build mission-critical applications on Amazon Web Services (AWS), Google Cloud, Microsoft Azure, and Alibaba Cloud.

## TiDB Cloud CLI

TiDB Cloud provides a command-line interface (CLI) `ticloud` for you to interact with TiDB Cloud from your terminal with a few lines of commands. For example, you can easily perform the following operations using `ticloud`:

- Create, delete, and list your clusters.
- Import data to your clusters.
- Export data from your clusters.

### TiDB Cloud CLI Auth

- Create a user profile with your TiDB Cloud API key

    ```shell
    ticloud config create
    ```

    > **Warning:**
    >
    > The profile name **MUST NOT** contain `.`.

- Log into TiDB Cloud with authentication:

    ```shell
    ticloud auth login
    ```

    After successful login, an OAuth token will be assigned to the current profile. If no profiles exist, the token will be assigned to a profile named `default`.

> **Note:**
>
> In the preceding two methods, the TiDB Cloud API key takes precedence over the OAuth token. If both are available, the API key will be used.

### Use the TiDB Cloud CLI

View all commands available:

```shell
ticloud --help
```

Verify that you are using the latest version:

```shell
ticloud version
```

If not, update to the latest version:

```shell
ticloud update
```

Create a TiDB Cloud cluster:

```shell
ticloud serverless create
```