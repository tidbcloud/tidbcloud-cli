# DM (Data Migration) CLI Commands

这个目录包含了 TiDB Cloud Serverless 数据迁移 (DM) 的 CLI 命令实现。

## 命令概览

### 1. 基础命令

#### `dm list` - 列出 DM 任务
```bash
# 交互模式
ticloud serverless dm list

# 非交互模式
ticloud serverless dm list --cluster-id <cluster-id>

# JSON 输出
ticloud serverless dm list --cluster-id <cluster-id> --output json
```

#### `dm describe` - 描述 DM 任务详情
```bash
# 交互模式
ticloud serverless dm describe

# 非交互模式
ticloud serverless dm describe --cluster-id <cluster-id> --task-id <task-id>
```

#### `dm delete` - 删除 DM 任务
```bash
# 交互模式
ticloud serverless dm delete

# 非交互模式
ticloud serverless dm delete --cluster-id <cluster-id> --task-id <task-id>

# 强制删除（无确认）
ticloud serverless dm delete --cluster-id <cluster-id> --task-id <task-id> --force
```

### 2. 任务管理

#### `dm create` - 创建 DM 任务
```bash
# 使用配置文件创建
ticloud serverless dm create --cluster-id <cluster-id> --config-file <config-file>

# 从标准输入读取配置
ticloud serverless dm create --cluster-id <cluster-id> --config-file -
```

#### `dm operate` - 操作 DM 任务
```bash
# 暂停任务
ticloud serverless dm operate --cluster-id <cluster-id> --task-id <task-id> --operation pause

# 恢复任务
ticloud serverless dm operate --cluster-id <cluster-id> --task-id <task-id> --operation resume
```

### 3. 预检查管理

#### `dm precheck` - 运行预检查
```bash
# 使用配置文件运行预检查
ticloud serverless dm precheck --cluster-id <cluster-id> --config-file <config-file>

# 从标准输入读取配置
ticloud serverless dm precheck --cluster-id <cluster-id> --config-file -
```

#### `dm get-precheck` - 获取预检查结果
```bash
# 获取预检查结果
ticloud serverless dm get-precheck --cluster-id <cluster-id> --precheck-id <precheck-id>
```

#### `dm cancel-precheck` - 取消预检查
```bash
# 取消预检查
ticloud serverless dm cancel-precheck --cluster-id <cluster-id> --precheck-id <precheck-id>

# 强制取消（无确认）
ticloud serverless dm cancel-precheck --cluster-id <cluster-id> --precheck-id <precheck-id> --force
```

## 配置文件格式

### DM 任务创建配置示例
```json
{
  "name": "my-dm-task",
  "mode": "incremental",
  "sourceConfig": {
    "host": "source.example.com",
    "port": 3306,
    "user": "root",
    "password": "password",
    "database": "source_db"
  },
  "targetConfig": {
    "database": "target_db"
  },
  "filterRules": [
    {
      "sourceSchema": "source_db",
      "sourceTable": "users"
    }
  ]
}
```

### DM 预检查配置示例
```json
{
  "name": "my-dm-task",
  "sourceConfig": {
    "host": "source.example.com",
    "port": 3306,
    "user": "root",
    "password": "password",
    "database": "source_db"
  },
  "targetConfig": {
    "database": "target_db"
  }
}
```

## 通用选项

- `--cluster-id`: TiDB Cloud Serverless 集群 ID
- `--output`: 输出格式 (human, json)
- `--force`: 跳过确认提示

## 注意事项

1. **配置文件格式**: 所有配置文件必须是有效的 JSON 格式
2. **交互模式**: 在支持交互的终端中，可以使用交互模式进行操作选择
3. **操作限制**: 目前支持的操作仅限于 `pause` 和 `resume`
4. **权限要求**: 需要有相应集群的 DM 操作权限

## 示例配置文件

项目中包含了一个示例配置文件 `example-dm-config.json`，可以作为参考使用。
