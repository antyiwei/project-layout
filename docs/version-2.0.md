# 方案一演进：按业务域拆分（Domain Layout）

> **前置阅读**：[version-1.0.md](version-1.0.md) 中的「方案一：标准企业级项目结构」
>
> **适用场景**：业务域增多、多人并行开发、改 A 域容易影响 B 域时，从「技术分层」演进为「按业务域拆分 + 共享内核」。

---

## 何时演进

出现以下信号时，建议从 v1.0 结构迁移到本结构：

1. **业务域 ≥ 3**：除 User 外还有 Order、Payment 等独立模块，且各自持续迭代。
2. **横向目录冲突多**：`internal/handler/` 下多人改不同业务，PR 冲突频繁。
3. **职责边界模糊**：Service 层出现大量 `if biz == "order"` 分支，或 Repo 被跨业务直接引用。

若项目仍只有 1～2 个简单 CRUD 模块，继续用 v1.0 即可，不必过早拆分。

---

## 与 v1.0 的对应关系

| v1.0（技术分层） | v2.0（按域拆分） | 说明 |
|-----------------|-----------------|------|
| `internal/handler/` | `internal/<domain>/handler/` | 各域独立 HTTP 层 |
| `internal/service/` | `internal/<domain>/service/` | 各域独立业务逻辑 |
| `internal/repository/` | `internal/<domain>/repository/` | 各域独立数据访问 |
| `internal/model/`、`internal/dto/` | `internal/<domain>/model/`、`dto/` | 数据结构随域走 |
| `internal/pkg/` | `internal/shared/` | 改名强调「共享内核」，禁止依赖业务域 |
| `main.go` 内注册路由 | `internal/router/` | 路由聚合，main 只做组装 |
| `configs/` | `configs/` + `shared/config/` | 配置文件仍放顶层；加载逻辑放 shared |
| `test/` | `test/integration/`、`test/e2e/` + 各域 `*_test.go` | 单测在域内；集成/E2E 放顶层 test |
| — | `internal/shared/event/` | v2 新增（可选）：跨域异步解耦 |
| — | `internal/<domain>/service/deps.go` | v2 新增：跨域依赖的窄接口 |
| — | `migrations/` | v2 新增（可选）：数据库 schema 版本管理 |

详细迁移步骤见 [migration-guide.md](migration-guide.md)。

---

## 业务增长后的典型目录结构

假设项目增加了 Order（订单）、Payment（支付）和 Notification（通知）等业务，结构如下。

> 以下在 [v1.0 标准结构](version-1.0.md) 基础上展示**分域后的变化**；根目录工程文件（`go.mod`、`Makefile` 等）与 v1.0 一致，此处一并列出以便直接参照搭建。

```text
your-project/
├── cmd/
│   └── server/
│       └── main.go           # 负责初始化所有模块的依赖注入
├── internal/
│   ├── user/                 # 用户域（独立模块）
│   │   ├── handler/
│   │   │   └── user_handler.go
│   │   ├── service/
│   │   │   └── user_service.go
│   │   ├── repository/
│   │   │   └── user_repo.go
│   │   ├── model/
│   │   └── dto/
│   │
│   ├── order/                # 订单域（独立模块）
│   │   ├── handler/
│   │   ├── service/
│   │   │   ├── order_service.go
│   │   │   └── deps.go       # 跨域依赖接口（如 UserLookup）
│   │   ├── repository/
│   │   ├── model/
│   │   └── dto/
│   │
│   ├── payment/              # 支付域（独立模块，无 HTTP 入口）
│   │   ├── service/
│   │   └── client/           # 对接第三方支付平台（支付宝/微信）
│   │
│   ├── notification/         # 通知域（独立模块，无 HTTP 入口）
│   │   ├── service/
│   │   └── template/         # 消息模板管理
│   │
│   ├── shared/               # 共享内核（Common Layer）
│   │   ├── config/           # 配置加载
│   │   ├── db/               # 数据库连接池
│   │   ├── cache/            # Redis 客户端
│   │   ├── logger/           # 日志
│   │   ├── event/            # 事件发布/订阅（跨域异步解耦，可选）
│   │   ├── middleware/       # Gin 中间件（JWT、限流）
│   │   ├── utils/            # 通用工具函数
│   │   └── errors/           # 统一错误定义
│   │
│   └── router/               # 路由聚合层
│       ├── router.go         # 总路由
│       ├── user_router.go
│       └── order_router.go
│
├── pkg/                      # 可导出的公共库（供其他微服务使用）
│   └── client/               # 对外 SDK（可选）
├── api/                      # API 契约
│   ├── openapi/              # Swagger / YAPI 文档
│   └── proto/                # gRPC Protobuf 文件（可选）
├── configs/                  # 配置文件
│   ├── config.yaml
│   └── dev.yaml
├── migrations/               # 数据库迁移脚本（可选，亦可用 goose/atlas 等工具）
├── scripts/                  # 运维脚本
│   ├── build.sh
│   └── deploy.sh
├── test/                     # 跨域 / 端到端测试
│   ├── integration/          # 集成测试（需真实 DB/Redis）
│   └── e2e/                  # E2E 测试（从 HTTP 入口跑完整链路）
├── docs/                     # 设计文档
├── .gitignore
├── go.mod
├── go.sum
└── Makefile                  # 项目管理命令（build, run, test）
```

### 目录说明

| 路径 | 说明 |
|------|------|
| `internal/<domain>/service/deps.go` | 跨域依赖的窄接口定义，由调用方（如 order）声明，提供方（如 user）隐式实现 |
| `internal/shared/event/` | 跨域异步解耦时使用；同步调用仍走 Service 接口，不必引入 |
| `migrations/` | 数据库 schema 版本管理；小项目也可把迁移逻辑放在 `scripts/` |
| `test/integration/`、`test/e2e/` | 单测放在各域内 `*_test.go`；需多域协作的测试放此处 |
| `pkg/client/` | 仅当其他服务需要 import 你的公共库时才需要 |

### 可选扩展（按项目需要添加）

以下不属于 v2 核心结构，但中大型项目常见：

```text
your-project/
├── cmd/
│   ├── server/               # HTTP 服务（必有）
│   ├── worker/               # 异步任务 / 队列消费（可选）
│   └── migrate/              # 独立迁移命令（可选）
├── deploy/                   # K8s / Docker Compose 部署清单（可选）
├── Dockerfile                # 容器镜像（可选）
└── internal/shared/
    ├── tracing/              # 链路追踪（OpenTelemetry 等，可选）
    └── metrics/              # Prometheus 指标（可选）
```

> **注意**：并非每个域都需要完整的 handler/service/repository/dto 四层。例如 `payment` 可能只有 service + client（对接第三方），`notification` 可能只有 service + template。按实际职责组织，不要机械复制目录。

---

## 核心变化解析

### 1. 业务隔离（Business Isolation）

- **以前**：所有 `handler` 都在一个文件夹里，改用户代码容易误伤订单代码。
- **现在**：`internal/user` 和 `internal/order` 物理隔离。开发者只需关注自己负责的 Domain，编译范围更小，代码冲突更少。

### 2. 共享内核（Shared Kernel）

随着业务增加，**基础设施（Infra）和业务无关的代码必须抽离**。

- `shared/db`：统一管理 GORM/SQLX 实例，避免每个模块都 `Open` 一次数据库。
- `shared/middleware`：认证（Auth）、链路追踪（Tracing）是所有接口通用的。
- `shared/errors`：定义统一的错误码（如 `ErrUnauthorized`），避免各模块自定义混乱。

### 3. 路由聚合（Router Aggregation）

不再在 `main.go` 里注册路由，而是由各业务模块定义自己的路由组，在 `router` 层进行组装。

```go
// internal/router/router.go
package router

import (
    "github.com/gin-gonic/gin"

    orderhandler "your-project/internal/order/handler"
    userhandler "your-project/internal/user/handler"
)

type Handlers struct {
    User  *userhandler.Handler
    Order *orderhandler.Handler
}

func Register(r *gin.Engine, h Handlers) {
    v1 := r.Group("/api/v1")
    registerUserRoutes(v1, h.User)
    registerOrderRoutes(v1, h.Order)
}
```

```go
// internal/router/user_router.go
package router

import (
    "github.com/gin-gonic/gin"

    userhandler "your-project/internal/user/handler"
)

func registerUserRoutes(v1 *gin.RouterGroup, h *userhandler.Handler) {
    users := v1.Group("/users")
    users.GET("/:id", h.GetByID)
}
```

### 4. 依赖流向（Dependency Flow）

必须遵守：**业务域依赖 shared，shared 不依赖任何业务域**。

```text
main.go（组装依赖）
   ↓
router（HTTP 分发）
   ↓
handler → service → repository
   ↓         ↓          ↓
   └─────────┴──────────┴──→ shared（db / logger / errors / middleware）
```

- `repository` 使用 `shared/db` 获取连接，不自己 `Open`。
- `handler` 使用 `shared/middleware`，不在域内重复实现 JWT。
- **禁止** `shared` import `internal/user` 等业务包。

---

## 依赖注入：main.go 组装示例

推荐在 `main.go` 用**结构体字面量**显式组装依赖，避免 `init()` 隐式初始化和 `NewXxx()` 构造函数模式。

```go
// cmd/server/main.go
package main

import (
    "log"

    "github.com/gin-gonic/gin"

    orderhandler "your-project/internal/order/handler"
    orderrepo "your-project/internal/order/repository"
    orderservice "your-project/internal/order/service"
    "your-project/internal/router"
    "your-project/internal/shared/config"
    "your-project/internal/shared/db"
    userhandler "your-project/internal/user/handler"
    userrepo "your-project/internal/user/repository"
    userservice "your-project/internal/user/service"
)

func main() {
    cfg := config.Load("configs/config.yaml")
    database := db.Open(cfg.Database)

    // 按域组装：repo → service → handler
    userRepo := &userrepo.Repo{DB: database}
    userSvc := &userservice.Service{Repo: userRepo}
    userH := &userhandler.Handler{Service: userSvc}

    orderRepo := &orderrepo.Repo{DB: database}
    orderSvc := &orderservice.Service{
        Repo:  orderRepo,
        Users: userSvc, // 跨域依赖：注入 user 域的 service
    }
    orderH := &orderhandler.Handler{Service: orderSvc}

    r := gin.Default()
    router.Register(r, router.Handlers{
        User:  userH,
        Order: orderH,
    })

    if err := r.Run(cfg.Server.Addr); err != nil {
        log.Fatal(err)
    }
}
```

各域内部保持薄结构体 + 方法，依赖通过字段注入：

```go
// internal/user/service/user_service.go
package service

import (
    "context"

    "your-project/internal/user/model"
    "your-project/internal/user/repository"
)

type Service struct {
    Repo *repository.Repo
}

func (s *Service) GetByID(ctx context.Context, id int64) (*model.User, error) {
    return s.Repo.FindByID(ctx, id)
}
```

```go
// internal/user/handler/user_handler.go
package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

    "your-project/internal/user/service"
)

type Handler struct {
    Service *service.Service
}

func (h *Handler) GetByID(c *gin.Context) {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    user, err := h.Service.GetByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, user)
}
```

---

## 跨域调用

### 原则

如果 `order` 需要 `user` 的数据，**通过 user 域的 Service 接口调用**，禁止直接 import `user/repository`。

### 推荐做法：在消费方定义接口

接口定义在**调用方**（order），由**提供方**（user）隐式实现，避免 order import user 的内部实现细节。

```go
// internal/order/service/deps.go
package service

import (
    "context"

    usermodel "your-project/internal/user/model"
)

// UserLookup 由 user/service.Service 满足，order 只依赖这个窄接口。
type UserLookup interface {
    GetByID(ctx context.Context, id int64) (*usermodel.User, error)
}
```

```go
// internal/order/service/order_service.go
package service

import (
    "context"
    "fmt"

    "your-project/internal/order/repository"
)

type Service struct {
    Repo  *repository.Repo
    Users UserLookup
}

func (s *Service) Create(ctx context.Context, userID int64, amount int64) error {
    user, err := s.Users.GetByID(ctx, userID)
    if err != nil {
        return err
    }
    if !user.Active {
        return fmt.Errorf("user %d is inactive", userID)
    }
    return s.Repo.Insert(ctx, userID, amount)
}
```

### 同步调用 vs 事件解耦

| 方式 | 适用场景 | 示例 |
|------|---------|------|
| **同步 Service 调用** | 强一致性、需立即返回结果 | 下单时校验用户状态 |
| **事件/MQ 解耦** | 最终一致、可异步处理 | 订单创建后发通知、写审计日志 |

事件方式下，各域通过 `shared/event/` 发布/订阅，仍不直接访问对方 Repository：

```go
// order 域：创建订单后发布事件
event.Publish("order.created", OrderCreated{OrderID: id, UserID: userID})

// notification 域：订阅事件，独立消费
event.Subscribe("order.created", notification.SendOrderConfirm)
```

---

## 测试与配置如何组织

### 配置

- **配置文件**：继续放在顶层 `configs/`（与 v1.0 一致）。
- **加载逻辑**：放在 `internal/shared/config/`，各域通过注入的配置结构体读取，不在域内直接读文件。

### 测试

| 类型 | 位置 | 说明 |
|------|------|------|
| 单元测试 | `internal/<domain>/**/*_test.go` | 与源码同包或 `_test` 包，mock 接口 |
| 集成测试 | `test/integration/` | 需要真实 DB/Redis 的跨域流程 |
| E2E 测试 | `test/e2e/` | 从 HTTP 入口跑完整链路 |

单测示例（mock 跨域依赖）：

```go
// internal/order/service/order_service_test.go
package service_test

import (
    "context"
    "testing"

    "your-project/internal/order/service"
    usermodel "your-project/internal/user/model"
)

type fakeUserLookup struct{}

func (f fakeUserLookup) GetByID(_ context.Context, id int64) (*usermodel.User, error) {
    return &usermodel.User{ID: id, Active: true}, nil
}

func TestCreate_RejectsInactiveUser(t *testing.T) {
    svc := &service.Service{
        Repo:  &fakeRepo{},
        Users: fakeUserLookup{},
    }
    // ...
}
```

---

## 什么时候需要进一步拆分？

当单个项目过大（如超过 10 万行代码，或 50+ 开发人员）时，可以考虑**物理拆分**（微服务化）：

| 阶段 | 状态 | 建议 |
|------|------|------|
| 初期 | 单仓单模块 | 使用上述结构，按文件夹隔离 |
| 中期 | 单仓多模块 | 使用 **Multi-Module**（一个 Git 仓库包含多个 `go.mod`），例如 `payment` 模块单独发版 |
| 后期 | 多仓多服务 | 彻底拆分为微服务，`pkg/` 中的代码可能独立成 Git 仓库作为 SDK |

---

## 避坑指南

1. **禁止 `shared` 依赖业务模块**：`shared` 包应该是纯净的，不能 `import internal/user`。如果出现这种情况，说明抽象有问题，应把共享逻辑下沉或把耦合逻辑上移到 main/router。
2. **跨域走 Service 接口，不碰 Repository**：`order` 不应 `import user/repository`，只依赖 `user/service` 暴露的窄接口。
3. **循环依赖**：Go 不允许循环引用。如果 A 调 B、B 又调 A，说明边界划分有问题——合并为一个域，或提取更底层的 shared 契约，或改用事件解耦。
4. **不要机械复制四层**：没有 HTTP 入口的域（如纯后台 payment）不必强行建 handler/；按职责裁剪目录。
5. **router 只做注册，不写业务逻辑**：路由文件里不应出现 SQL 或业务判断，保持薄聚合层。

---

## 相关文档

- [version-1.0.md](version-1.0.md) — 起步阶段的标准结构
- [migration-guide.md](migration-guide.md) — 从 v1.0 迁移到 v2.0 的分步指南
