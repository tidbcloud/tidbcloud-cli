# TiDB Cloud FS (ticloud fs) User Guide

**TiDB Cloud FS** is a managed file storage service provided by TiDB Cloud that offers a unified filesystem abstraction for TiDB instances. You can store, read, and manage files in TiDB Cloud just like a local filesystem, with advanced capabilities such as SQL queries, a secret vault, and FUSE mount support.

**`ticloud fs`** is the command group in the TiDB Cloud CLI for operating TiDB Cloud FS, supporting upload/download, directory browsing, SQL execution, secret management, and FUSE mount.

---

## Prerequisites: Associate a TiDB Instance

FS data storage depends on a specific TiDB instance. Depending on the type of instance you have, you need to bind its ID to the CLI configuration and complete a one-time FS initialization (`init`).

### Using Serverless TiDB

If you already have a **TiDB Cloud Serverless** cluster (created via the [TiDB Cloud Web Console](https://tidbcloud.com/) or the `ticloud serverless create` command), configure it as follows:

1. **Configure Authentication**

   Serverless clusters are managed by the TiDB Cloud platform, so you must provide valid platform credentials to access their FS service. Choose one of the following methods:

   - **OAuth Login (Recommended)**

     ```bash
     ticloud auth login
     ```

   - **API Key Authentication**

     ```bash
     ticloud config set public-key <your-public-key>
     ticloud config set private-key <your-private-key>
     ```

2. **Configure cluster-id**

   ```bash
   ticloud config set fs.cluster-id <your-serverless-cluster-id>
   ```

3. **Initialize FS**

   ```bash
   ticloud fs init --user admin --password <your-password>
   ```

   > Every new Serverless cluster must run `init` before using FS for the first time to provision the FS tenant.

### Using Zero TiDB

If you are using a **TiDB Zero** instance, configure it as follows. If you are not familiar with TiDB Zero, visit [https://zero.tidbcloud.com/](https://zero.tidbcloud.com/) to learn more.

1. **Configure zero-instance-id**

   ```bash
   ticloud config set fs.zero-instance-id <your-zero-instance-id>
   ```

   TiDB Zero uses an independent deployment and access model. You **do not** need to run `ticloud auth login` or configure `public-key` / `private-key`.

2. **Initialize FS**

   ```bash
   ticloud fs init --user admin --password <your-password>
   ```

   > Like Serverless, every new Zero instance must run `init` once to provision the FS tenant.

---

## Path Format

FS uses the `:/` prefix to denote remote paths:

- `:/` — Remote root directory
- `:/path/to/file.txt` — Remote file path

---

## Subcommand Manual

### Initialization and Configuration

Before performing any file operations, you must create an FS tenant for the associated database. This step only needs to be performed once.

#### `ticloud fs init`

Initialize the FS tenant for the currently associated database. This must be run before using FS for the first time on a new database.

```bash
ticloud fs init --user admin --password secret
```

### Basic File Operations

`ticloud fs` provides a set of Unix-like filesystem commands for managing remote files and directories. You can upload, download, browse, and delete files as if they were local.

#### `ticloud fs ls`

List remote directory contents. Supports `-l` for detailed view.

```bash
# List root directory
ticloud fs ls :/

# List a folder with size details
ticloud fs ls -l :/myfolder
```

#### `ticloud fs cat`

Display remote file contents, similar to the Unix `cat` command. Useful for quickly viewing text files.

```bash
ticloud fs cat :/readme.txt
```

#### `ticloud fs cp`

Copy files. Supports local↔remote, remote↔remote, stdin upload, and resuming large file transfers.

```bash
# Upload a local file to remote
ticloud fs cp local.txt :/remote/

# Download a remote file to current directory
ticloud fs cp :/remote/file.txt .

# Upload from stdin via pipe
echo "hello" | ticloud fs cp - :/file.txt

# Resume a large file upload (useful for unstable networks)
ticloud fs cp --resume large.zip :/remote/
```

#### `ticloud fs mkdir`

Create a remote directory. Supports nested paths.

```bash
ticloud fs mkdir :/mydir
```

#### `ticloud fs mv`

Move or rename remote files/directories.

```bash
ticloud fs mv :/oldname :/newname
```

#### `ticloud fs rm`

Delete remote files or empty directories.

```bash
ticloud fs rm :/file.txt
```

#### `ticloud fs stat`

View remote file or directory metadata, such as size and modification time.

```bash
ticloud fs stat :/file.txt
```

### Search and Query

In addition to basic file operations, `ticloud fs` supports searching file contents and executing SQL queries directly on remote storage, making it easy to locate information or perform data analysis.

#### `ticloud fs grep`

Search for matching content in remote files. Supports `--limit` to restrict the number of results.

```bash
ticloud fs grep "TODO" :/project --limit 20
```

#### `ticloud fs find`

Find remote files by criteria. Supports `-name`, `-size`, `-newer` filters, similar to the Unix `find` command.

```bash
# Find by file extension
ticloud fs find :/ -name "*.go"

# Find files larger than 1MB
ticloud fs find :/ -size +1M
```

#### `ticloud fs sql`

Execute SQL queries against the FS backend for structured queries on stored data.

```bash
ticloud fs sql -q "SELECT 1"
```

### Secret Management

TiDB Cloud FS includes a built-in secret vault for securely storing sensitive information such as database connection strings and API keys. You can manage these secrets via subcommands, and issue scoped capability tokens for third-party agents.

#### `ticloud fs secret set`

Create or update a secret. Field values can be set directly, read from a file (`@file`), or read from stdin (`-`).

```bash
# Set fields directly
ticloud fs secret set myapp DATABASE_URL=mysql://...

# Read a field from file and password from stdin
ticloud fs secret set myapp key=@secret.txt password=-
```

#### `ticloud fs secret get`

Read a secret or a specific field. Supports `--json` and `--env` output formats.

```bash
# Read the entire secret
ticloud fs secret get myapp

# Read a specific field and output as JSON
ticloud fs secret get myapp/password --json
```

#### `ticloud fs secret exec`

Execute a command with secret fields injected as environment variables. Commonly used in local scripts or CI pipelines to securely use secrets.

```bash
ticloud fs secret exec myapp -- ./run.sh
```

#### `ticloud fs secret ls`

List all secrets. Supports `--json` output.

```bash
ticloud fs secret ls
```

#### `ticloud fs secret rm`

Delete a specific secret.

```bash
ticloud fs secret rm myapp
```

#### `ticloud fs secret grant`

Issue a scoped capability token for a specific agent, restricting access to certain secrets. Requires `--agent` and `--ttl`.

```bash
ticloud fs secret grant --agent myagent --ttl 1h myapp/password
```

#### `ticloud fs secret revoke`

Revoke an issued capability token, invalidating it immediately.

```bash
ticloud fs secret revoke tok_abc123
```

#### `ticloud fs secret audit`

Query the secret audit log to track who accessed which secrets and when. Supports filtering by `--secret`, `--agent`, `--since`, and limiting results with `--limit`.

```bash
# View the latest 50 audit events
ticloud fs secret audit --limit 50

# View audit events for a secret in the last 24 hours
ticloud fs secret audit --secret myapp --since 24h
```

### Interactive Shell

If you need to execute multiple FS commands frequently, you can use the interactive shell to avoid typing the full command prefix each time.

#### `ticloud fs shell`

Launch the interactive FS Shell. Supports commands such as `cd`, `pwd`, `ls`, `cat`, `cp`, `mkdir`, `mv`, `rm`, `sql`, `stat`, `help`, `exit`. The prompt is `ticloud:fs>`.

```bash
ticloud fs shell
```

### FUSE Mount

Via FUSE mount, you can map the remote FS to a local directory and access it using the system file manager or standard command-line tools.

#### `ticloud fs mount`

Mount the remote FS to a local directory. The current version defaults to read-only. Supports `--debug` (enable debug logging) and `--allow-other` (allow other users to access).

```bash
ticloud fs mount /mnt/tidbcloud
```

#### `ticloud fs umount`

Unmount the local mount point, safely disconnecting from the remote FS.

```bash
ticloud fs umount /mnt/tidbcloud
```

---

## Configuration Quick Reference

| Property | Description | Default |
|----------|-------------|---------|
| `fs.cluster-id` | Associated Serverless cluster ID | — |
| `fs.zero-instance-id` | Associated Zero instance ID | — |
| `fs-endpoint` | FS server endpoint | `https://fs.tidbapi.com/` |
| `public-key` / `private-key` | API Key authentication for Serverless | — |

You can also override these via environment variables:

- `TICLOUD_FS_ENDPOINT` — overrides `fs-endpoint`
- `TICLOUD_FS_CLUSTER_ID` — overrides `fs.cluster-id`
- `TICLOUD_FS_ZERO_INSTANCE_ID` — overrides `fs.zero-instance-id`

---

## FAQ

**Q: Running `ticloud fs ls :/` returns `tenant not found`?**

A: The associated database has not been initialized for FS yet. Please run `ticloud fs init --user <user> --password <password>` to complete initialization.

**Q: Serverless cluster returns 401 Unauthorized?**

A: Please confirm you have completed `ticloud auth login` or correctly configured `public-key` / `private-key`.
