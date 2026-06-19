# 从 v1.0 迁移到 v2.0

> **前置阅读**：[version-1.0.md](version-1.0.md) → [version-2.0.md](version-2.0.md)

本文描述如何将「技术分层」结构（v1.0 方案一）逐步演进为「按业务域拆分」结构（v2.0），**无需一次性大爆炸式重构**。

---

## 迁移概览

```text
阶段 1：抽 shared          → 基础设施与业务解耦
阶段 2：试点单域（user）    → 验证目录与依赖规范
阶段 3：抽 router           → main.go 瘦身
阶段 4：迁移第二域（order） → 验证跨域调用
阶段 5：清理旧目录          → 删除空的 handler/service/repository
```

每完成一个阶段，运行 `go build ./...` 和 `go test ./...` 确认无回归。

---

## 阶段 1：抽取 shared

**目标**：把 v1.0 中 `internal/pkg/` 及散落的基础设施代码迁到 `internal/shared/`。

### 操作步骤

1. 创建 `internal/shared/`，按职责分子目录：

   ```text
   internal/shared/
   ├── config/
   ├── db/
   ├── cache/
   ├── logger/
   ├── middleware/
   ├── utils/
   └── errors/
   ```

2. 将 `internal/pkg/middleware/` → `internal/shared/middleware/`
3. 将 `internal/pkg/utils/` → `internal/shared/utils/`
4. 若 DB 初始化在 `main.go` 中，提取到 `internal/shared/db/`
5. 全局搜索 `internal/pkg`，批量替换 import 路径为 `internal/shared/...`
6. 删除空的 `internal/pkg/` 目录

### 验证

- [ ] `go build ./...` 通过
- [ ] `shared/` 下没有任何 `import ".../internal/user"` 等业务包
- [ ] 配置文件仍位于顶层 `configs/`，加载逻辑在 `shared/config/`

---

## 阶段 2：试点迁移 user 域

**目标**：以 user 为第一个业务域，建立 `internal/user/` 的标准四层结构。

### 操作步骤

1. 创建目录：

   ```text
   internal/user/
   ├── handler/
   ├── service/
   ├── repository/
   ├── model/
   └── dto/
   ```

2. 移动文件（包名保持不变或按新路径调整）：

   | 原路径 | 新路径 |
   |--------|--------|
   | `internal/handler/user_handler.go` | `internal/user/handler/user_handler.go` |
   | `internal/service/user_service.go` | `internal/user/service/user_service.go` |
   | `internal/repository/user_repo.go` | `internal/user/repository/user_repo.go` |
   | `internal/model/user.go` | `internal/user/model/user.go` |
   | `internal/dto/user_dto.go` | `internal/user/dto/user_dto.go` |

3. 更新 import 路径（IDE 重构或 `gofmt` + 手动修正）
4. **暂时保留**旧目录中其他业务的文件（order 等尚未迁移）

### 验证

- [ ] user 相关接口功能正常
- [ ] `internal/user/` 不 import `internal/order/` 等其他域
- [ ] user 域内依赖方向：handler → service → repository → shared

---

## 阶段 3：抽取 router

**目标**：路由注册从 `main.go` 迁到 `internal/router/`。

### 操作步骤

1. 创建 `internal/router/router.go` 和 `internal/router/user_router.go`
2. 将 `main.go` 中的 user 路由注册逻辑移到 `registerUserRoutes`
3. `main.go` 只保留：配置加载 → 依赖组装 → 调用 `router.Register`

参考 [version-2.0.md 路由聚合示例](version-2.0.md#3-路由聚合router-aggregation)。

### 验证

- [ ] `main.go` 不含具体路由路径字符串（或仅保留 health check 等全局路由）
- [ ] HTTP 接口与迁移前行为一致

---

## 阶段 4：迁移 order 域 + 跨域依赖

**目标**：迁移第二个域，并验证跨域调用规范。

### 操作步骤

1. 按阶段 2 同样方式创建 `internal/order/` 并移动文件
2. 若 order 需要 user 数据，在 `order/service/deps.go` 定义窄接口：

   ```go
   type UserLookup interface {
       GetByID(ctx context.Context, id int64) (*usermodel.User, error)
   }
   ```

3. 在 `main.go` 组装时注入：

   ```go
   orderSvc := &orderservice.Service{
       Repo:  orderRepo,
       Users: userSvc, // user/service.Service 满足 UserLookup
   }
   ```

4. 添加 `internal/router/order_router.go`

### 验证

- [ ] `order` 不 import `user/repository`
- [ ] `go vet ./...` 无 import cycle 报错
- [ ] 跨域流程（如「下单校验用户」）测试通过

---

## 阶段 5：迁移剩余域 + 清理

**目标**：迁移 payment、notification 等域，删除旧顶层目录。

### 操作步骤

1. 逐个迁移剩余业务域（无 HTTP 入口的域可省略 handler/）
2. 确认 `internal/handler/`、`internal/service/`、`internal/repository/`、`internal/model/`、`internal/dto/` 已空
3. 删除上述空目录
4. 更新 `Makefile`、`README`、CI 路径（如有硬编码）

### 验证

- [ ] 项目目录与 [version-2.0.md 目标结构](version-2.0.md#业务增长后的典型目录结构) 一致
- [ ] `go test ./...` 全部通过
- [ ] 无废弃 import 或 dead code

---

## 常见问题

### Q：迁移期间新旧结构并存，可以接受吗？

可以。阶段 2～4 期间 `internal/handler/`（旧）和 `internal/user/handler/`（新）会短暂共存。关键是**新代码只写在新结构里**，旧代码按域逐步搬完。

### Q：一个函数被多个域共用，放哪里？

- 纯工具函数（无业务语义）→ `shared/utils/`
- 有业务语义但跨域 → 提取到 owning 域的 Service 接口，或下沉到 shared 的**无业务耦合**抽象（如 `shared/validator/`）
- 若 shared 需要知道业务概念 → 说明不应放 shared，应通过 Service 接口或事件解耦

### Q：payment 没有 HTTP 入口，怎么组织？

只建 `internal/payment/service/` 和 `internal/payment/client/`，在 `main.go` 注入给 order 使用即可，不必强行加 handler/。

### Q：要不要上 Wire / Fx 等 DI 框架？

小中型项目：**不必**。`main.go` 结构体字面量组装足够清晰。团队规模大、依赖图复杂时再考虑。

---

## 相关文档

- [version-1.0.md](version-1.0.md) — 起步阶段结构
- [version-2.0.md](version-2.0.md) — 目标结构与代码示例
