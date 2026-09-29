# Premium / Essential V2 P0 最终 review

日期：2026-09-24。最终设计见 [nextgen-p0-design.md](nextgen-p0-design.md)。

## 结论和范围

按用户确认的范围，Premium / Essential V2 的 P0 命令均已有代码实现和对应本地测试；本轮补齐 Premium CMEK 创建所需的公开 API、CLI 和兼容性修复。整体 review 未发现需要继续修改的确定性代码问题，但这不等于已通过部署后的全部 E2E，也不承诺不存在未覆盖问题。

范围按当前发布契约执行，不按草案中的字段名重新设计 API：project 可选、min RCU 由服务端管理、shell 初期只连接 PUBLIC endpoint（后续私网扩展见 2026-09-26 记录）；`--password` 与 Starter 一致，仅免去凭据提示，不新增批处理 SQL 模式。Import/Export 为 P1。Essential V2 的 CMEK 不在此次能力承诺内：现有内部实现不支持，本次不修改内部能力。

本轮没有执行生产或 staging 资源操作，没有提交、推送、切换分支或同步远端。

## Review 基线

| 仓库 | 检查时本地状态 | 比较范围 |
|---|---|---|
| CLI | detached HEAD `19aea6adfceb24cd5fdb716de6369f6ca83ff8a0`；本地 `origin/main` 相同 | 已跟踪未暂存修改 + 全部未跟踪实现/测试/生成文件；无暂存修改 |
| mgmt | `codex/premium-openapi-oauth`，HEAD `9c9fa4ef680340083909e8d76a9dcf782bcb1b47`；本地 `origin/main` 为 `7b8db0a6dd5913313089c6e6ab2ffc2387e5ca34`，0 ahead / 27 behind | 当前 HEAD 上的工作区修改和未跟踪文件；另核对与本地 `origin/main` 的基线差异，未把落后提交造成的差异当作本次删除 |

以上是本地 ref，不代表实时远端状态。mgmt 合入较新基线前需要同步并重新验证；本轮未自行 rebase。原有未提交的 P0 改动被保留，不能只看 `git diff --stat` 而漏掉未跟踪目录。

## 兼容性修正

1. **保持现有 API/SDK 形状。** 公开 CMEK principal/verify 时保留原 path、operationId、请求 body 和 `CustomerManagedEncryptionKeyService` 首个 tag；PUBLIC 设置保留 `PublicConnectionSettingService` 首个 tag。避免 Console/full SDK 的服务类改名。
2. **只修新公开文档的 requiredness。** `regionId` 必填仅作用于新发布的 public Swagger；共享 proto 和既有 Console/full Swagger 注解保持原样。服务端原本就拒绝缺失 region。
3. **保持已有生成类型。** 固定生成器 7.12.0，显式映射复用的 service-plan enum，避免增加 principal API 后导致现有 List/Region 参数类型改名。旧 SDK Go 文件主体比较中，除 `client.go` 的新增服务注册外，既有主体未变。
4. **不改变旧创建路径。** `none/default-key` 在两个 plan 下保持原 payload，均只调用一次 create；不会额外访问 CMEK helper。只有 CMEK 分支执行 principal → verify → create，验证失败停止创建。
5. **补齐公开接口权限边界。** `verifyPrincipal` 未指定 `tidbId` 时保持组织级预创建行为；指定时检查实例归属和 token project scope，并使用已有实例读取权限。没有新增 IAM privilege 或修改 Global/Regional KMS 处理。
6. **发布顺序有约束。** 先部署 permission-aware mgmt-service，再启用 Portal OAuth，再开放 Gateway API，最后发布 CLI。旧 Gateway 不支持新能力时保留 API 错误，不自动降级、不重试创建。

## P0 覆盖核对

| 功能 | CLI / OpenAPI 实现 | 本地验证 |
|---|---|---|
| Create | 两个 plan 的 `--create`；可选 project、max RCU、加密设置 | plan/project/payload；不指定 project；非法 output 在网络请求前拒绝 |
| List | plan 过滤、分页 | 过滤参数、翻页、重复 token 中止 |
| Describe | 获取实例并校验 plan | Premium ↔ Essential V2 双向 guard |
| Update | 只提交改变的字段；保留服务端容量限制 | changed fields、非法输入前置、跨 plan 拒绝 |
| Delete | plan guard + 确认或 `--force`；输出请求已接受 | 两个 plan 的确认要求、force、跨 plan 无删除调用 |
| Region | plan 过滤和分页 | 请求契约、输出校验、重复 token 中止 |
| Shell 交互选择 | 不带连接参数时选择本 plan 实例；`-c` 跳过选择 | 列表选择、显式 ID、空列表、非交互缺少 ID |
| Shell 显式凭据 | `-c --password`、`-c -u --password` | 默认/指定用户名路径、DSN 特殊字符、保持终端 SQL shell 语义 |
| Shell 前置条件 | root-password API/CLI；PUBLIC 设置 OpenAPI；实例 CA + 严格 TLS | 密码输入、plan guard、reachability、CA 下载协议、TLS 错误链与诊断 |
| Public endpoint | 两个 plan 的 `public-endpoint enable|disable`；交互选择或 `-c`；disable 确认/`--force` | PATCH 仅含 `enabled`、选中后重新读取并校验 plan、非法 output/空 ID/非交互参数前置拒绝、取消/读取失败不发 PATCH、更新错误透传、Digest/Bearer 的 true/false 请求契约 |
| API key / OAuth | Digest / Bearer，保留既有凭据优先级 | transport 测试、Portal 认证/RBAC、实例权限、metadata 传递和 token 非持久化 |
| Premium CMEK | AWS/AliCloud principal、verify、create；`premium cmek principal` | 两种云与两种认证的 HTTP 契约；输入冲突；helper 失败/valid=false 不 create；Essential 拒绝 |

关键测试在 `internal/cli/nextgen/*_test.go`、`internal/service/cloud/nextgen_client_test.go`，以及 mgmt 的 API contract、Portal auth/RBAC、service permission/accountcli、utilities metadata 测试中。

## 整体差异 review

- **冗余和过度防御：** 未加入 capability cache、版本探测、自动重试、IAM policy 编辑器或本地云端 key-policy 校验。生成 SDK 包含契约中的关联模型，不手工裁剪生成物。主要防御对应明确边界：非法写请求前校验、分页 token 循环、plan/权限隔离和 TLS 验证。
- **非本次功能语义：** Starter/v1beta1 的请求契约、Digest 传输和 shell 入口保留。共享 SQL helper 提取没有改动 Starter DSN 构造；错误由 `%s` 改为 `%w`，保留显示文本并暴露原始错误链。
- **有意的共享行为变化：** debug 输出隐藏 Authorization/重置密码内容；Bearer token 不进入 workflow JSON；NextGen OAuth 的实际实例 project 权限在 mgmt 校验。这些是有测试的安全边界变化，不能概括为所有共享行为完全不变。
- **语法和类型：** CLI build、vet、全量单测、lint 通过；mgmt 受影响模块 generation、fmt/lint、单测通过。
- **假设与实现：** 未假定 Essential 支持 CMEK，未把 API 接受创建/删除等同于完成，未把 PUBLIC host 存在等同于 DNS/TLS 可用，未把本地 Bearer 测试等同于已部署 Gateway/IAM 接受 token。
- **内部业务边界：** CMEK 阶段没有 Global/Regional/infra 业务源码、内部加密工作流或数据库迁移改动；后续私网地址读取修复涉及 NextGen Global，见 2026-09-27 记录。现有 KMS 校验可能有内部动作，不将 verify 宣称为无副作用 dry run。

## 2026-09-25 review 修复

- 密码校验改为与 mgmt/Global 一致的 8–64 字节，补充 ASCII 边界及多字节密码用例。
- OAuth 实例权限保留授权项目内的实例级角色：通过本地集群元数据按组织、token 项目和角色实例 ID 求交；包括备份/恢复使用的已删除实例，查询失败时不放行。
- 旧 Starter 独立鉴权链路恢复原有项目加载、空项目拒绝及默认项目 metadata 行为；日志仍不记录原始 token。NextGen/BFF 共用的 token 范围检查及两层鉴权仍保留，不宣称全部 OAuth 路径与旧实现完全相同。
- 删除 debug request 中重复的 Header 拷贝；保留 Authorization/密码脱敏及 workflow token 非持久化。
- 新增角色撤权、项目范围过滤、查询错误、旧 Starter 兼容性及密码长度回归测试。未部署 Gateway/Account，也未执行真实云资源 E2E。
- 本次验证通过：CLI 全量测试及相关 race；Portal 全量测试及新增兼容性回归；mgmt Account、permission、NextGen、相关 BFF/Internal API 和原有 Dedicated/Starter 创建模块测试；mgmt 权限 race；两个服务模块的 generate/fmt/lint/vet。CLI 全量测试和 mgmt vet 遇到沙箱限制后重跑成功。

## 2026-09-26 已有 Private Endpoint 连接扩展

- 方案 review 后明确：仅消费当前实例 API 返回的 `PRIVATE_ENDPOINT` 地址，不创建、授权或配置网络；mgmt 的 `reachable=true` 不等于客户端私网可达。
- Premium / Essential V2 Shell 新增 `--connection-type private-endpoint`，默认仍为 `public`。`--endpoint host:port` 只能匹配当前实例返回的私网地址；无地址报错，多地址在交互模式中选择，显式 `-c` 时要求指定地址。不支持 `VPC_PEERING`，不自动回退公网。
- 继续使用实例 CA 和所选 host 做严格 TLS 校验。私网初始连接设置 30 秒超时，密码输入不占用该时间，成功连接后取消超时上下文；不改变正常 SQL 会话、PUBLIC 或 Starter 的超时行为。
- 补充非法参数前置校验、plan guard、单/多地址选择、候选地址未就绪、禁止任意地址和公网回退、交互语义、TLS 错误保留、连接超时与会话分离的回归测试。
- 本地通过：`go test -count=1 ./...`、`go vet ./...`、NextGen/util/cloud 三包 race、CLI 构建、v1.64.7 源码构建的 golangci-lint 全量检查、帮助/非法参数 smoke、生成文档幂等和 `git diff --check`。lint 按仓库规则补齐候选列表预分配后通过；测试和 lint 的缓存/回环访问在沙箱外完成。
- 本轮未改 mgmt、SDK 或云端资源；保留既有未提交改动，未提交或推送。真实私网 TLS/SQL E2E **未执行**，仍需已有私网端点、可达客户端和有效证书链，不能用本地测试替代部署验收。

## 2026-09-27 Private Endpoint review 修复

- 修正“契约里有字段就代表运行时会返回”的假设：NextGen Global 的 Get/List 原来只组装 PUBLIC/VPC_PEERING，现在补读已有 PrivateLink 服务及连接，将地址交给既有 OpenAPI converter。未新增 API、SDK 字段、资源创建流程或权限逻辑。
- AWS 使用 ACTIVE 服务的 DNS；GCP/Azure/AliCloud 使用 ACTIVE 且仍在仓库登记的连接域名。私网元数据读取失败时保留原有实例查询；不把服务端 `reachable` 当作客户端网络探测结果。
- 修复 usql 初始版本查询吞掉超时错误的问题：有连接超时的路径在 `Open` 返回后、主动取消 context 前检查 context 错误，超时不能继续进入 Shell；PUBLIC/Starter 不增加这一检查。
- 新增云厂商地址映射、状态/登记过滤、读取失败、PUBLIC/VPC_PEERING 兼容性、两种 plan 的 OpenAPI 转换，以及被吞掉超时的回归用例。
- 本轮通过：CLI 全量单测、vet、NextGen/util/cloud race、固定版本 v1.64.7 lint、构建和帮助 smoke；SQL dialog 定向测试重复 20 次通过。NextGen Global、mgmt-service 两个模块的全量单测、vet、generate/fmt/lint 通过，Global cluster / OpenAPI tidb 定向 race 通过。macOS race 链接出现 `LC_DYSYMTAB` warning，但测试退出码为 0。
- 核对生成/格式化前快照后，只有本次新测试的格式被自动调整，没有额外生成物或无关源码变化。两仓库 `git diff --check` 通过。全量后端单测日志：`/private/tmp/private-shell-fix.7Cqtu1/global-tests.log`、`/private/tmp/private-shell-fix.7Cqtu1/mgmt-tests.log`（临时文件，不是永久 CI artifact）。
- 必须先部署 NextGen Global 读取修复再做真实私网验收；本轮不创建、授权或配置私网资源，不放宽 TLS 验证。真实私网 E2E 未执行。

## 2026-09-27 Public Endpoint enable/disable

- Premium 与 Essential V2 新增 `public-endpoint enable|disable`，均支持不带 `-c` 的交互实例选择和带 `-c` 的显式选择。
- 两个操作先读取实例并执行 plan guard，再调用既有 `PATCH /v1beta2/tidbs/{id}/publicConnectionSetting`。请求只携带 `enabled`，不会读取、替换或清空现有 IP access list。
- disable 属于连接中断操作，交互模式要求输入 `yes`；自动化必须显式传入 `--force`。human 输出只声明请求已接受，JSON 输出返回服务端 setting，不把 PATCH 成功表述为 endpoint 已完成收敛。
- 本地覆盖两种 plan、true/false 请求体、交互选择、跨 plan 拒绝、非交互参数/确认和非法输出在网络请求前拒绝，并验证生成文档和命令帮助。

### Public Endpoint review 优化

- 对照 mgmt public-connection PATCH 和 NextGen Global 实现：`enabled: false` 是显式布尔更新；不传 `ipAccessList` / `clearIpAccessList` 会保留原 IP access list。复用原 API/SDK，不增加后端、权限或 IP 放行规则改动。
- 交互选择后统一执行 `GetTiDB` + plan guard，与 password/shell 一致，避免只用列表快照判断。显式传入空 `-c` 在客户端初始化前拒绝，不意外转入交互选择。
- disable 确认提示包含目标实例及连接中断提醒；补充确认、取消、`--force`、读取/更新失败和空响应的回归测试。成功响应与公网实际可用仍分开判断。
- 删除只包一层的 `publicEndpointCmdWithSelector`，合并参数校验和执行，去掉重复读取 ID 和跨回调保存 output 状态。确认逻辑独立为可测试依赖，不增加全局 hook 或通用命令框架。
- 请求测试覆盖 Digest challenge 重试和 OAuth Bearer，并断言 true/false payload 都只包含 `enabled`；两种 plan 和 `essential` 别名的命令分发均覆盖。本地请求测试不证明已部署 Gateway 接受 OAuth，也不证明公网已收敛；本轮没有切换真实集群的公网状态。
- 本轮验证通过：CLI 全量单测、全量 vet、NextGen/cloud race、NextGen/cloud 固定版本 v1.64.7 lint、构建、Premium/Essential V2/essential alias 帮助 smoke。构建产物为 `/private/tmp/ticloud-nextgen-public-endpoint`，未覆盖原 `/private/tmp/ticloud-nextgen`。

## 此前 CMEK 实现阶段的本地验证记录

| 检查 | 结果 |
|---|---|
| CLI `go test -count=1 ./...` | 通过 |
| CLI `go vet ./...`、`go build -o /private/tmp/ticloud-nextgen-cmek ./cmd/ticloud` | 通过 |
| CLI `go test -race -count=1 ./internal/cli/nextgen ./internal/service/cloud ./internal/cli/config ./internal/util` | 通过 |
| CLI 固定版本 golangci-lint 1.64.7 | 通过 |
| CLI Premium create / CMEK principal / Essential V2 / essential alias help | 通过 |
| API `make gen`、Buf/AIP lint、契约测试 | 通过 |
| mgmt-service、Portal、utilities 各模块 `make generate`、`make fmt`（含配置的 lint） | 通过 |
| mgmt-service、Portal、utilities 各模块 `go test -count=1 ./...` | 通过 |
| mgmt-service 权限/accountcli race、全模块 vet | 通过 |
| API、CLI SDK、命令文档重复生成 | 通过；SDK 在干净目录重生成，比较包含生成器 manifest |
| 两仓库 `git diff --check` | 通过 |
| Buf breaking 相对本轮修改前 API 快照 | 通过 |
| API `make lint-buf-breaking`（Makefile 基线 `release-v0.79`） | **失败：历史 `EstimatePriceRequest/Response` 删除/重命名/类型变化**；相关 `tidb.proto` 相对本地 HEAD 和 origin/main 均无本次 diff，不扩展范围修复 |

中间执行遇到过 sandbox Go cache 权限限制、并行 lint 锁，以及生成期间 Portal vet 的短暂导入失败；在稳定生成结果上重跑成功，上表记录最终结果。不能把 Buf 全量历史失败隐藏成“所有检查通过”。

本地日志目录：`/private/tmp/ticloud-cmek.Upys1Z/`。主要日志：`cli-tests-final.log`、`mgmt-service-tests.log`、`mgmt-portal-tests-final.log`、`utilities-tests.log`、`mgmt-permission-race.log`、`api-gen-idempotence.log`、`cli-sdk-clean-repeat.log`、`api-breaking-final.log`。日志位于临时目录，不是永久 CI artifact。

## 未执行的部署验收

本轮是源码实现和兼容性 review，没有部署服务，也没有创建真实 CMEK 实例。以下需要目标环境补证：

1. Gateway 路由发布，以及真实 API key / OAuth token 的 IAM 权限检查（包括跨 org/project 拒绝）。
2. AWS/AliCloud key policy 按 principal 设置；verify → create → ACTIVE → CMEK readback，完成后按明确授权清理资源。
3. Essential V2 有可用容量时完成其生命周期 E2E。历史 staging 容量不足不能算通过，也不能通过 CLI 绕过。
4. PUBLIC endpoint、密码、可信 CA 就绪后验证两种 plan 的实际 shell；不使用跳过证书校验代替验收。
5. mgmt 同步较新目标基线后重跑以上本地门禁及相应集成检查。

因此：当前结果是约定范围的 **P0 命令源码实现完成 + 本地验证通过（全量历史 Buf 差异单列）**，不是全部线上功能已经验收通过。
