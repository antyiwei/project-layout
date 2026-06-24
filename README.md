# project-layout

Go 后端项目布局 **GitHub Template Repository** — 精简 starter，开箱即用。

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## 快速开始

> **重要**：无论哪种方式，都必须运行 `init-project.sh`，将默认 module 名 `project-layout` 改为你的仓库路径。

### 方式 A：GitHub 网页创建（Use this template）

1. 在 GitHub 点击 **Use this template** → 输入新项目名（如 `my-api`）→ Create repository
2. 克隆并初始化：

```bash
git clone git@github.com:antyiwei/my-api.git
cd my-api
./scripts/init-project.sh github.com/antyiwei/my-api --yes
```

### 方式 B：本地创建（任意 Git 平台）

不在网页先建仓库，本地起步；代码可推到 **Gitee / GitHub / GitLab / 公司自建 Git** 等任意平台：

```bash
git clone git@github.com:antyiwei/project-layout.git my-api && cd my-api
rm -rf .git && git init && git branch -M main
./scripts/init-project.sh <你的-module-路径> --yes
```

`<你的-module-路径>` 由目标 Git 地址推导，例如：

| 计划使用的仓库 | module 路径 |
|---------------|-------------|
| `git@gitee.com:antyiwei/my-api.git` | `gitee.com/antyiwei/my-api` |
| `git@github.com:acme/my-api.git` | `github.com/acme/my-api` |
| `git@git.company.com:backend/my-api.git` | `git.company.com/backend/my-api` |

详细说明见 **[docs/local-setup.md](docs/local-setup.md)**。

### 构建并运行

```bash
make build
make run
```

验证：

```bash
curl http://localhost:8080/api/v1/users/1
```

## 这个模板包含什么

| 内容 | 说明 |
|------|------|
| **main 分支（本分支）** | 精简 starter：`user` 单域 + `shared` + `router` |
| **[example/v2-full 分支](https://github.com/antyiwei/project-layout/tree/example/v2-full)** | 完整四域示例（user / order / payment / notification） |
| **docs/** | v1.0 / v2.0 布局文档 + 迁移指南 |

## 目录结构

```text
cmd/server/          # 入口，依赖注入组装
internal/user/       # 示例业务域
internal/shared/     # 共享内核（db / cache / middleware / errors ...）
internal/router/     # 路由聚合
configs/             # 配置文件
scripts/             # build.sh / init-project.sh
```

## 文档索引

| 文档 | 适用阶段 |
|------|---------|
| [docs/version-1.0.md](docs/version-1.0.md) | 新项目起步，技术分层 |
| [docs/version-2.0.md](docs/version-2.0.md) | 业务增多，按域拆分 |
| [docs/migration-guide.md](docs/migration-guide.md) | 从 v1 迁移到 v2 |
| [docs/TEMPLATE.md](docs/TEMPLATE.md) | 团队模板使用与维护手册 |
| [docs/local-setup.md](docs/local-setup.md) | 本地创建项目（任意 Git 平台） |

## 演进到多业务域

业务域增多后，参考 [docs/version-2.0.md](docs/version-2.0.md) 和 [docs/migration-guide.md](docs/migration-guide.md)。

完整可运行四域示例：

```bash
git clone -b example/v2-full git@github.com:antyiwei/project-layout.git
```

## License

MIT — 详见 [LICENSE](LICENSE)

Issues 和 PR 欢迎。
