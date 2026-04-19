# TiDB Cloud FS (ticloud fs) 使用指南

**TiDB Cloud FS** 是 TiDB Cloud 提供的托管文件存储服务，为 TiDB 实例提供统一的文件系统抽象。你可以像操作本地文件系统一样在 TiDB Cloud 中存储、读取和管理文件，同时支持 SQL 查询、密钥保险箱（Vault）以及 FUSE 挂载等高级能力。

**`ticloud fs`** 是 TiDB Cloud CLI 中用于操作 TiDB Cloud FS 的命令组，支持上传下载、目录浏览、SQL 执行、密钥管理以及 FUSE 挂载等功能。

---

## 前置依赖：关联一个 TiDB 实例

FS 的数据存储依赖于一个具体的 TiDB 实例。根据你拥有的实例类型，需要将其 ID 绑定到 CLI 配置中，并完成一次性的 FS 初始化（`init`）。

### 使用 Serverless TiDB

如果你已经有一个 **TiDB Cloud Serverless** 集群（可以通过 [TiDB Cloud Web Console](https://tidbcloud.com/) 创建，也可以使用 `ticloud serverless create` 命令创建），请按以下步骤配置：

1. **配置鉴权**

   Serverless 集群受 TiDB Cloud 平台统一管控，访问其 FS 服务时必须提供有效的平台鉴权信息。你可以选择以下任一方式：

   - **OAuth 登录（推荐）**

     ```bash
     ticloud auth login
     ```

   - **API Key 鉴权**

     ```bash
     ticloud config set public-key <your-public-key>
     ticloud config set private-key <your-private-key>
     ```

2. **配置 cluster-id**

   ```bash
   ticloud config set fs.cluster-id <your-serverless-cluster-id>
   ```

3. **初始化 FS**

   ```bash
   ticloud fs init --user admin --password <your-password>
   ```

   > 每个新的 Serverless 集群在首次使用 FS 前，都必须执行 `init` 来为其开通 FS 租户。

### 使用 Zero TiDB

如果你使用的是 **TiDB Zero** 实例，配置方式如下。如果你还不了解 TiDB Zero，可以访问 [https://zero.tidbcloud.com/](https://zero.tidbcloud.com/) 了解更多。

1. **配置 zero-instance-id**

   ```bash
   ticloud config set fs.zero-instance-id <your-zero-instance-id>
   ```

   TiDB Zero 采用独立的部署与访问模型，**不需要**执行 `ticloud auth login`，也**不需要**配置 `public-key` / `private-key`。

2. **初始化 FS**

   ```bash
   ticloud fs init --user admin --password <your-password>
   ```

   > 与 Serverless 相同，新的 Zero 实例也需要执行一次 `init` 来开通 FS 租户。

---

## 路径格式说明

FS 使用 `:/` 前缀表示远程路径：

- `:/` — 远程根目录
- `:/path/to/file.txt` — 远程文件路径

---

## 子命令手册

### 初始化与配置

在使用任何文件操作之前，必须先为关联的数据库创建 FS 租户。这一步只需要执行一次。

#### `ticloud fs init`

为当前关联的数据库初始化 FS 租户，每个新数据库在首次使用 FS 前必须执行一次。

```bash
ticloud fs init --user admin --password secret
```

### 基础文件操作

`ticloud fs` 提供了一组类似 Unix 文件系统的命令，用于管理远程文件和目录。你可以像操作本地文件一样进行上传、下载、浏览和删除。

#### `ticloud fs ls`

列出远程目录内容。支持 `-l` 查看详情。

```bash
# 查看根目录
ticloud fs ls :/

# 查看某个文件夹，并显示文件大小等详细信息
ticloud fs ls -l :/myfolder
```

#### `ticloud fs cat`

显示远程文件内容，类似于 Unix 的 `cat` 命令，适合快速查看文本文件。

```bash
ticloud fs cat :/readme.txt
```

#### `ticloud fs cp`

复制文件。支持本地↔远程、远程↔远程、从标准输入上传，以及断点续传大文件。

```bash
# 上传本地文件到远程目录
ticloud fs cp local.txt :/remote/

# 从远程下载文件到当前目录
ticloud fs cp :/remote/file.txt .

# 通过管道将标准输入内容写入远程文件
echo "hello" | ticloud fs cp - :/file.txt

# 断点续传大文件，适合网络不稳定场景
ticloud fs cp --resume large.zip :/remote/
```

#### `ticloud fs mkdir`

创建远程目录，支持多级路径。

```bash
ticloud fs mkdir :/mydir
```

#### `ticloud fs mv`

移动或重命名远程文件/目录。

```bash
ticloud fs mv :/oldname :/newname
```

#### `ticloud fs rm`

删除远程文件或空目录。

```bash
ticloud fs rm :/file.txt
```

#### `ticloud fs stat`

查看远程文件或目录的元数据，例如大小、修改时间等。

```bash
ticloud fs stat :/file.txt
```

### 搜索与查询

除了基础文件操作，`ticloud fs` 还支持直接在远程存储中搜索文件内容和执行 SQL 查询，方便你快速定位信息或进行数据分析。

#### `ticloud fs grep`

在远程文件中搜索匹配内容。支持 `--limit` 限制结果数量。

```bash
ticloud fs grep "TODO" :/project --limit 20
```

#### `ticloud fs find`

按条件查找远程文件。支持 `-name`、`-size`、`-newer` 等过滤条件，类似 Unix 的 `find` 命令。

```bash
# 按文件名后缀查找
ticloud fs find :/ -name "*.go"

# 查找大于 1MB 的文件
ticloud fs find :/ -size +1M
```

#### `ticloud fs sql`

在 FS 后端执行 SQL 查询，可以直接对存储的数据进行结构化查询。

```bash
ticloud fs sql -q "SELECT 1"
```

### Secret 管理

TiDB Cloud FS 内置了一个密钥保险箱（Vault），用于安全地存储敏感信息，例如数据库连接串、API Key 等。你可以通过子命令对这些密钥进行增删改查，还可以为第三方 Agent 颁发受限访问令牌。

#### `ticloud fs secret set`

创建或更新密钥。字段值可以直接赋值、从文件读取（`@file`）或从标准输入读取（`-`）。

```bash
# 直接设置字段
ticloud fs secret set myapp DATABASE_URL=mysql://...

# 从文件读取字段值，并从标准输入读取密码
ticloud fs secret set myapp key=@secret.txt password=-
```

#### `ticloud fs secret get`

读取密钥或指定字段。支持 `--json` 和 `--env` 输出格式。

```bash
# 读取整个密钥
ticloud fs secret get myapp

# 读取指定字段，并以 JSON 格式输出
ticloud fs secret get myapp/password --json
```

#### `ticloud fs secret exec`

将密钥字段作为环境变量注入后执行指定命令，常用于在本地脚本或 CI 流程中安全地使用密钥。

```bash
ticloud fs secret exec myapp -- ./run.sh
```

#### `ticloud fs secret ls`

列出所有密钥。支持 `--json` 输出。

```bash
ticloud fs secret ls
```

#### `ticloud fs secret rm`

删除指定密钥。

```bash
ticloud fs secret rm myapp
```

#### `ticloud fs secret grant`

为指定 Agent 颁发受限能力令牌，限制其只能访问特定的密钥范围。需要指定 `--agent` 和 `--ttl`。

```bash
ticloud fs secret grant --agent myagent --ttl 1h myapp/password
```

#### `ticloud fs secret revoke`

撤销一个已颁发的能力令牌，使其立即失效。

```bash
ticloud fs secret revoke tok_abc123
```

#### `ticloud fs secret audit`

查询密钥审计日志，追踪谁在什么时间访问了哪些密钥。支持按 `--secret`、`--agent`、`--since` 过滤，以及 `--limit` 限制条数。

```bash
# 查看最近的 50 条审计记录
ticloud fs secret audit --limit 50

# 查看某个密钥最近 24 小时内的访问记录
ticloud fs secret audit --secret myapp --since 24h
```

### 交互式 Shell

如果你需要频繁执行多条 FS 命令，可以使用交互式 Shell，避免每次都要输入完整命令前缀。

#### `ticloud fs shell`

启动交互式 FS Shell。支持 `cd`、`pwd`、`ls`、`cat`、`cp`、`mkdir`、`mv`、`rm`、`sql`、`stat`、`help`、`exit` 等命令。提示符为 `ticloud:fs>`。

```bash
ticloud fs shell
```

### FUSE 挂载

通过 FUSE 挂载，你可以将远程 FS 映射为本地目录，直接使用系统自带的文件管理器或命令行工具访问。

#### `ticloud fs mount`

将远程 FS 挂载到本地目录。当前版本默认只读。支持 `--debug`（开启调试日志）和 `--allow-other`（允许其他用户访问）。

```bash
ticloud fs mount /mnt/tidbcloud
```

#### `ticloud fs umount`

卸载本地挂载点，安全断开与远程 FS 的连接。

```bash
ticloud fs umount /mnt/tidbcloud
```

---

## 配置项速查

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `fs.cluster-id` | 关联的 Serverless 集群 ID | — |
| `fs.zero-instance-id` | 关联的 Zero 实例 ID | — |
| `fs-endpoint` | FS 服务端点地址 | `https://fs.tidbapi.com/` |
| `public-key` / `private-key` | Serverless 场景下的 API Key 鉴权 | — |

也支持通过环境变量覆盖：

- `TICLOUD_FS_ENDPOINT` — 覆盖 `fs-endpoint`
- `TICLOUD_FS_CLUSTER_ID` — 覆盖 `fs.cluster-id`
- `TICLOUD_FS_ZERO_INSTANCE_ID` — 覆盖 `fs.zero-instance-id`

---

## 常见问题

**Q: 执行 `ticloud fs ls :/` 时提示 `tenant not found`？**

A: 说明当前关联的数据库尚未初始化 FS。请运行 `ticloud fs init --user <user> --password <password>` 完成初始化。

**Q: Serverless 集群报 401 Unauthorized？**

A: 请确认已完成 `ticloud auth login` 或正确配置了 `public-key` / `private-key`。
