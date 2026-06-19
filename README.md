# Go 项目布局指南 — 完整 v2 参考实现

> **分支说明**：这是 **完整四域示例**（user / order / payment / notification）分支。
> 日常开新项目请使用默认分支 `main`（精简 starter）→ GitHub **Use this template**。

本仓库是一套 Go 后端项目的目录结构参考，随业务规模分阶段演进。

## 本分支包含什么

- `internal/user/`、`internal/order/`、`internal/payment/`、`internal/notification/` 四个业务域
- `internal/shared/` 共享内核 + `internal/router/` 路由聚合
- 可编译运行的 Gin HTTP 服务骨架

## 文档索引

| 文档 | 适用阶段 | 说明 |
|------|---------|------|
| [docs/version-1.0.md](docs/version-1.0.md) | 新项目起步 | 标准企业级结构，按技术分层 |
| [docs/version-2.0.md](docs/version-2.0.md) | 业务域增多 | 按业务域拆分 + 共享内核 |
| [docs/migration-guide.md](docs/migration-guide.md) | v1 → v2 迁移 | 分 5 阶段逐步重构 |

## 快速运行

```bash
make build && make run
# GET http://localhost:8080/api/v1/users/1
# POST http://localhost:8080/api/v1/orders
```

## 回到 Template Starter

```bash
git checkout main
```

或在 GitHub 上对 `main` 分支点击 **Use this template** 创建新项目。
