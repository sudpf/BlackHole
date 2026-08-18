# 项目架构概览

BlackHole 是一个基于 Go 的网络运营相关项目，包含两个独立运行的组件：

- **VoidEngine**：HTTP API 服务，提供用户管理、网络流量查询、统一错误响应、请求上下文与 API 文档能力。
- **Stash**：配置驱动的数据采集与处理管道，从 syslog、Kafka 等输入读取事件，经过过滤转换后写入 ClickHouse、Elasticsearch、syslog 或标准输出。

## 总体设计原则

- 明确的配置驱动，服务启动所需参数全部通过配置文件加载。
- 统一的请求上下文、错误码、日志与 trace ID。
- 清晰的分层：入口层、业务层、数据访问层、基础设施层。
- 可测试：service、filter、handler 等核心逻辑优先通过接口或纯函数组织。

## 目录结构

```text
cmd/                         程序入口
  voidengine/                VoidEngine HTTP 服务入口
  stash/                     Stash 管道服务入口
  swagger-generator/         Swagger 文档生成工具

api/                         HTTP 接入层
  middleware/                通用中间件：错误处理、恢复、请求上下文、访问日志
  router/                    路由抽象
  wrapper/                   Handler 返回值适配
  swagger/                   Swagger UI 挂载
  voidengine/openapi/        VoidEngine OpenAPI 服务装配
    v1/handler/              v1 接口处理器
    v1/router/               v1 路由注册

internal/                    业务内部实现
  runtime/                   通用运行时能力（优雅退出、配置路径等）
  voidengine/                VoidEngine 业务
    config/                  配置加载与校验
    contract/                请求 / 响应契约
    service/                 业务逻辑
    model/                   数据模型与 DAO
    errorcode/               业务错误码与多语言文案定义
  stash/                     Stash 管道业务
    config/                  配置加载与校验
    app/                     应用启动装配
    service/                 管道服务
      input/                 输入源：syslog、Kafka 等
      filter/                过滤器：drop / remove_field / transfer
      handler/               消息处理入口
      output/                输出目标：ClickHouse、Elasticsearch、syslog 等

pkg/                         可复用基础包
  apperror/                  业务错误与多语言错误目录
  auth/                      认证抽象
  config/                    配置通用工具
  constant/                  常量
  db/                        数据库连接与基础封装
  env/                       请求环境上下文（语言、IP、Principal 等）
  localize/                  本地化相关工具
  logger/                    结构化日志
  requestctx/                请求上下文（trace ID 等）
  units/                     单位解析

conf/                        配置样例
docs/                        文档与生成的 Swagger
```

## VoidEngine 架构

### 调用链路

```text
HTTP Request
  -> Gin Router
  -> 中间件链
       RequestContext   // 生成 trace id、超时、语言
       Recovery         // panic 恢复
       ErrorHandler     // 统一错误响应
       ApiLog           // 访问日志
  -> Handler
       参数绑定 + 校验
       调用 service
  -> Service
       业务逻辑
       调用 DAO
  -> DAO / Model
       GORM 访问 MySQL / ClickHouse
```

### 关键能力

- **统一响应**：`api/common/response` 定义 `ApiResponse`，成功与错误统一结构。
- **错误目录**：`pkg/apperror.Catalog` 管理错误码、HTTP 状态码和中英文消息。
- **校验错误翻译**：`api/validation` 将 validator 错误翻译为字段级中文/英文信息。
- **请求上下文**：`pkg/requestctx` 贯穿 trace ID，日志和 SQL 都可携带。
- **优雅退出**：`internal/runtime.Run` 统一处理信号、超时与 shutdown 回调。

## Stash 架构

### 调用链路

```text
Input (Syslog / Kafka)
  -> MessageHandler.Consume
       JSON 解析
       Filter 链
       Output 写入
  -> Writer
       ClickHouse / Elasticsearch / Syslog / stdout
```

### 核心抽象

- **Input**：事件来源。目前支持 syslog server、Kafka consumer 等。
- **Filter**：对事件进行过滤或转换。
  - `drop`：按条件丢弃。
  - `remove_field`：删除指定字段。
  - `transfer`：解析 JSON 字段并展开或重命名。
- **Handler**：串联 filter 与 output，是输入和输出之间的处理入口。
- **Output / Writer**：写入目标。
  - ClickHouse：批量插入，使用 chunk executor 按大小与间隔刷盘。
  - Elasticsearch：Bulk 写入，支持按时间切索引。
  - Syslog：按配置列与条件转发。
  - stdout：开发调试用。

### 配置模型

Stash 以 cluster 为单位组织 pipeline：

```yaml
clusters:
- input:
    syslogs: [...]
    kafka: ...
  filters:
  - action: drop
    conditions: [...]
  - action: remove_field
    fields: [...]
  - action: transfer
    field: content
  output:
    clickhouse: ...
    elasticsearch: ...
    syslogs: [...]
```

每个 cluster 内部：

- Input 可以有多个来源；
- Filters 按顺序执行；
- Output 可以配置多个目标。

## 上下文与日志

- HTTP 请求会生成或透传 `X-Trace-ID`。
- 业务日志通过 `pkg/logger.FromContext(ctx)` 自动携带 trace ID。
- `pkg/env.Env` 保存语言、客户端 IP、Principal 等请求级环境信息。
- 错误目录按 `Accept-Language` 做中英文本地化。

## 错误处理模型

1. Handler / service 中产生 `apperror.Error`。
2. 通过 `c.Error(...)` 传入 Gin 错误链。
3. `ErrorHandler` 取出最后一个错误：
   - 已知业务错误：查目录、本地化、返回统一响应。
   - 未知错误：包装为系统错误、打 error 日志。
   - 校验错误：翻译为字段级详情。
4. 响应结构为 `{ code, message, data, details }`。

## 扩展点

### VoidEngine

- 新增资源：在 `contract`、`service`、`model`、`handler`、`router` 中对应扩展。
- 新增认证方式：实现 `pkg/auth.Authenticator`。
- 新增错误码：在 `internal/voidengine/errorcode` 中注册。

### Stash

- 新增输入：在 `internal/stash/service/input` 中实现并接入 `service.New`。
- 新增过滤器：在 `internal/stash/service/filter` 中新增 `FilterFunc` 并注册到 `CreateFilters`。
- 新增输出：在 `internal/stash/service/output` 中实现 `Writer` 接口并注册到 `NewWriters`。

## 测试策略

- 纯逻辑函数优先单元测试，如 filter、handler、service。
- service 层通过 stub DAO 进行测试，避免依赖真实数据库。
- HTTP 层通过 `httptest` 测试路由、错误处理、鉴权、响应结构。

## 安全与运行注意

- 生产环境必须配置用户接口鉴权，避免默认开放。
- 用户密码应使用 hash 存储，不得明文入库。
- 启动日志建议对敏感字段脱敏。
- Stash 输出失败的语义应根据业务容忍度明确（允许丢 / 重试 / 死信队列）。
