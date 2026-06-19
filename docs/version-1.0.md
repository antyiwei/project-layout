# 方案一：标准企业级项目结构（推荐）

下面是目前 Go 社区最主流、最成熟的结构，参考了 golang-standards/project-layout，适合中大型后端服务。

```text

your-project/
├── cmd/                      # 程序入口
│   └── server/
│       └── main.go           # 极薄的一层，只做启动和配置组装
├── internal/                 # 私有业务代码（外部项目不可导入）
│   ├── handler/             # 接口层（HTTP/RPC）：接收请求，参数校验
│   │   └── user_handler.go
│   ├── service/             # 业务逻辑层：核心业务编排
│   │   └── user_service.go
│   ├── repository/          # 数据访问层：数据库/Cache操作
│   │   └── user_repo.go
│   ├── model/               # 数据结构定义（DO/PO）
│   │   └── user.go
│   ├── dto/                 # 数据传输对象（Request/Response）
│   │   └── user_dto.go
│   └── pkg/                 # 内部共享的工具包（非业务核心）
│       ├── middleware/
│       └── utils/
├── pkg/                     # 可公开引用的库代码（可选）
│   └── client/              # 提供给外部调用的SDK
├── api/                     # API 契约
│   ├── openapi/             # Swagger/YAPI 文档
│   └── proto/               # gRPC Protobuf 文件
├── configs/                 # 配置文件
│   ├── config.yaml
│   └── dev.yaml
├── scripts/                 # 脚本
│   ├── build.sh
│   └── deploy.sh
├── test/                    # 集成测试、E2E 测试
├── docs/                    # 设计文档
├── .gitignore
├── go.mod
├── go.sum
└── Makefile                 # 项目管理命令（build, run, test）

```

#### 各层职责说明（Clean Architecture 思想）

1. **Handler**：只关心协议（HTTP/gRPC），解析参数，调用 Service，返回结果。**不包含业务逻辑**。
2. **Service**：处理核心业务逻辑，编排 Repo 层，做事务控制。
3. **Repository**：抽象数据来源，屏蔽 MySQL/Redis 的具体实现细节。
4. **Model**：纯数据结构，不依赖外部库（除了 ORM 标签）。

# 方案二：单体模块化结构（微内核）

如果项目初期不想分得太细，或者这是一个工具类项目，可以使用这种扁平化结构。

```text
your-project/
├── main.go
├── router/                  # 路由定义
│   └── router.go
├── controller/              # 控制器
├── logic/                   # 业务逻辑（核心）
├── dao/                     # 数据访问
├── model/
├── middleware/
├── config/
└── test/
```

**特点**：上手快，心智负担低。随着业务增长，可以将 `logic`拆分成不同的业务域（Domain）。

# 方案三：基于业务领域的 DDD 风格

当业务极其复杂，需要多人协作时，按“业务边界”划分比按“技术层级”划分更有效。

```text

internal/
├── user/                    # 用户域
│   ├── domain/              # 领域实体、聚合根
│   ├── application/         # 应用服务（对外接口）
│   ├── infrastructure/      # 基础设施实现（DB, Cache）
│   └── interfaces/          # 对外接口适配（HTTP/RPC）
├── order/                   # 订单域
│   ├── domain/
│   ├── application/
│   └── ...
└── common/                  # 全局通用代码

```

---

### 关键设计决策与最佳实践

1. **为什么要有** `internal`**？**
  - 利用 Go 编译器的强制限制，防止项目中的公共代码被外部项目错误引用，保证封装性。
2. **依赖注入（Dependency Injection）**

- 不要在 `init()`函数中隐式初始化全局变量（尤其是 DB、MQ 连接）。在 `main.go` 用**结构体字面量**显式组装依赖，方便测试和替换实现。
- 推荐包级函数风格，避免 `NewXxx()` 构造函数模式。

```go
// cmd/server/main.go — 在入口组装
userRepo := &userrepo.Repo{DB: db}
userSvc := &userservice.Service{Repo: userRepo}
userH := &userhandler.Handler{Service: userSvc}

// internal/user/service/user_service.go — 域内保持薄结构体
type Service struct {
    Repo *repository.Repo
}

func (s *Service) GetByID(ctx context.Context, id int64) (*model.User, error) {
    return s.Repo.FindByID(ctx, id)
}
```

1. **配置文件管理**

- 使用 `viper`读取配置，支持本地文件、环境变量和远程配置中心。
- 配置结构体要清晰，避免使用 `map[string]interface{}`。

1. **错误处理**

- 定义统一的错误码和错误结构，方便追踪和国际化。
- 在 Repository 层记录原始错误，在 Service 层包装业务语义，在 Handler 层转换为用户友好的信息。



1. **Makefile 的使用**

- 将常用的命令固化下来，降低新人上手成本。

```
.PHONY: build
build:
    go build -o bin/server ./cmd/server

.PHONY: test
test:
    go test -v ./...
```

### 我的建议

- **起步阶段**：直接从**方案一**开始。虽然看起来文件多，但它能迫使你一开始就养成“关注点分离”的好习惯，避免后期重构的痛苦。
- **不要做的事**：不要引入过多的第三方框架去强行改变 Go 的语法风格（比如试图模仿 Java Spring 那种重度注解的模式）。拥抱 Go 的简单性。

### 下一步

当业务域增多（3 个以上模块、多人并行冲突频繁）时，演进为按业务域拆分：

- [version-2.0.md](version-2.0.md) — 目标结构与代码示例
- [migration-guide.md](migration-guide.md) — 从本结构迁移的分步指南

