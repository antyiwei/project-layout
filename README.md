# Go 项目布局指南

本仓库是一套 Go 后端项目的目录结构参考，随业务规模分阶段演进。

## 文档索引

| 文档 | 适用阶段 | 说明 |
|------|---------|------|
| [docs/version-1.0.md](docs/version-1.0.md) | 新项目起步 | 标准企业级结构，按技术分层（handler / service / repository） |
| [docs/version-2.0.md](docs/version-2.0.md) | 业务域增多 | 从方案一演进为按业务域拆分 + 共享内核 |
| [docs/migration-guide.md](docs/migration-guide.md) | v1 → v2 迁移 | 分 5 个阶段逐步重构，无需大爆炸 |

## 推荐阅读顺序

1. 新项目 → 直接读 **version-1.0**（方案一）
2. 业务变复杂 → 读 **version-2.0**，了解目标结构
3. 已有 v1.0 项目要改造 → 按 **migration-guide** 逐步执行

## GitHub Template Repo

本仓库也可作为 GitHub Template Repo 使用，创建属于你自己的项目模板。
