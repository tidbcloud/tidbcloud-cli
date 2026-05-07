# TiDB Cloud FS Commands

The `ticloud fs` command group provides filesystem-like operations for managing files and directories in TiDB Cloud FS.

## Overview

TiDB Cloud FS is a managed file storage service that allows you to store and retrieve files using familiar filesystem semantics. The CLI provides both direct API operations and FUSE mount capabilities.

## Prerequisites

- TiDB Cloud CLI (`ticloud`) installed
- TiDB Cloud authentication configured (OAuth login or API key pair)

## Configuration

The `fs` command reuses your existing TiDB Cloud CLI authentication:

```bash
# OAuth login
ticloud auth login

# Or configure API keys
ticloud config set public-key <your-public-key>
ticloud config set private-key <your-private-key>
```

### Database Association

FS operations can be associated with a specific database:

```bash
# Associate with a TiDB Cloud serverless cluster
ticloud config set fs.cluster-id <cluster-id>

# Associate with a TiDB Zero instance
ticloud config set fs.zero-instance-id <instance-id>
```

These are sent as `X-TIDBCLOUD-CLUSTER-ID` and `X-TIDBCLOUD-ZERO-INSTANCE-ID` headers respectively. If both are set, `cluster-id` takes precedence.

### Server URL

The default FS server URL is `https://fs.tidbapi.com/`. You can override it:

```bash
ticloud config set fs-endpoint <url>
# or
export TICLOUD_FS_SERVER=<url>
```

## Path Format

Remote paths use the `:/` prefix:

- `:/` - Root directory
- `:/path/to/file.txt` - File path

## Commands

### Basic Operations

#### List Directory

```bash
# List root directory
ticloud fs ls :/

# List with details
ticloud fs ls -l :/myfolder

# List subdirectory
ticloud fs ls :/docs
```

#### Display File Contents

```bash
ticloud fs cat :/readme.txt
ticloud fs cat :/docs/guide.md
```

#### Copy Files

Upload local to remote:
```bash
ticloud fs cp local.txt :/remote/
ticloud fs cp ./data.json :/backup/
```

Download remote to local:
```bash
ticloud fs cp :/readme.txt ./
ticloud fs cp :/docs/guide.md ./local-docs/
```

Server-side copy:
```bash
ticloud fs cp :/old/path :/new/path
```

From stdin:
```bash
echo "content" | ticloud fs cp - :/file.txt
```

Resume interrupted upload:
```bash
ticloud fs cp --resume large.zip :/remote/
```

#### Create Directory

```bash
ticloud fs mkdir :/mydir
ticloud fs mkdir :/parent/child
```

#### Move/Rename

```bash
ticloud fs mv :/oldname :/newname
ticloud fs mv :/folder/file.txt :/other/file.txt
```

#### Remove

```bash
ticloud fs rm :/file.txt
ticloud fs rm :/empty_folder/
```

#### File Metadata

```bash
ticloud fs stat :/file.txt
```

### SQL Execution

Execute SQL queries against the FS backend:

```bash
ticloud fs sql -q "SELECT 1"
ticloud fs sql -f query.sql
```

### Secret Management

Manage secrets in the TiDB Cloud FS vault:

```bash
# Set a secret
ticloud fs secret set myapp DATABASE_URL=mysql://...

# Get a secret
ticloud fs secret get myapp
ticloud fs secret get myapp/password

# List secrets
ticloud fs secret ls

# Delete a secret
ticloud fs secret rm myapp

# Run a command with secret env vars
ticloud fs secret exec myapp -- ./run.sh

# Issue a scoped capability token
ticloud fs secret grant --agent myagent --ttl 1h myapp/password

# Revoke a token
ticloud fs secret revoke tok_abc123

# Query audit events
ticloud fs secret audit --limit 50
```

### Search Operations

#### Search File Contents (Grep)

```bash
# Search for pattern
ticloud fs grep "function main" :/

# Search with limit
ticloud fs grep "TODO" :/myproject --limit 20
```

#### Find Files

```bash
# Find by name pattern
ticloud fs find :/ -name "*.go"

# Find by size (larger than 1MB)
ticloud fs find :/ -size +1M

# Find modified after date
ticloud fs find :/ -newer 2024-01-01

# Combined filters
ticloud fs find :/data -name "*.json" -newer 2024-01-01
```

### Interactive Shell

Start an interactive shell for filesystem operations:

```bash
ticloud fs shell
```

Shell commands:
- `cd <path>` - Change directory
- `pwd` - Print current directory
- `ls [path]` - List directory
- `cat <path>` - Display file
- `cp <src> <dst>` - Copy files
- `mkdir <path>` - Create directory
- `mv <old> <new>` - Move/rename
- `rm <path>` - Remove
- `sql <query>` - Execute SQL query
- `stat <path>` - Show metadata
- `help` - Show help
- `exit` - Exit shell

Prompt: `ticloud:fs>`

### FUSE Mount

Mount the remote filesystem as a local directory using FUSE.

#### Prerequisites

**Linux:**
- Install FUSE: `sudo apt-get install fuse3` (Debian/Ubuntu) or `sudo yum install fuse3` (RHEL/CentOS)

**macOS:**
- Install macFUSE from https://osxfuse.github.io/

#### Mount

```bash
# Create mount point
mkdir -p /mnt/tidbcloud

# Mount (read-only, default for MVP)
ticloud fs mount /mnt/tidbcloud

# Mount with debug logging
ticloud fs mount /mnt/tidbcloud --debug

# Mount allowing other users
ticloud fs mount /mnt/tidbcloud --allow-other
```

#### Use Mounted Filesystem

Once mounted, use standard filesystem commands:

```bash
# List files
ls /mnt/tidbcloud

# Read files
cat /mnt/tidbcloud/readme.txt

# Copy from mount
cp /mnt/tidbcloud/file.txt ./
```

#### Unmount

```bash
ticloud fs umount /mnt/tidbcloud
```

Or use system commands:
```bash
# Linux
fusermount3 -u /mnt/tidbcloud

# macOS
umount /mnt/tidbcloud
```

## Examples

### Backup Workflow

```bash
# Create backup directory
ticloud fs ls :/backups 2>/dev/null || ticloud fs mkdir :/backups

# Upload backup
tar czf - /data | ticloud fs cp - :/backups/data-$(date +%Y%m%d).tar.gz

# List backups
ticloud fs ls :/backups
```

### Development Workflow

```bash
# Start interactive shell
ticloud fs shell

# In shell:
ticloud:fs> cd /myproject
ticloud:fs> ls
ticloud:fs> cat config.json
ticloud:fs> exit
```

### Search and Filter

```bash
# Find all Go files modified recently
ticloud fs find :/ -name "*.go" -newer $(date -v-7d +%Y-%m-%d)

# Search for function definitions
ticloud fs grep "^func " :/src --limit 50
```

## Limitations

### MVP (Current Version)

- **FUSE mount is read-only** - Write operations via FUSE will be added in a future release
- **No symlink support** - Symlinks are not supported
- **Limited concurrent uploads** - Large file uploads may be slower

### Planned Features

- Full read-write FUSE mount
- Streaming large file operations
- Progress bars for transfers
- Batch operations
- File versioning

## Troubleshooting

### FUSE Mount Issues

**Error: "fuse: device not found"**
- Ensure FUSE is installed: `modprobe fuse` (Linux) or install macFUSE (macOS)

**Error: "permission denied"**
- Check mount point permissions
- Try with `--allow-other` flag
- On Linux, add user to `fuse` group: `sudo usermod -a -G fuse $USER`

**Error: "transport endpoint is not connected"**
- The mount is stale, try unmounting and remounting

### Connection Issues

**Error: "cannot reach FS server"**
- Check network connectivity
- Verify `fs-endpoint` config: `ticloud config describe`
- Check TiDB Cloud auth status: `ticloud auth status`

## See Also

- `ticloud config` - Manage CLI configuration
- `ticloud auth` - Authentication commands
