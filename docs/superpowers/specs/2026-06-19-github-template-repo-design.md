# GitHub Template Repo 设计规格

**日期**：2026-06-19  
**仓库**：`git@github.com:antyiwei/project-layout.git`  
**状态**：待用户审阅

---

## 1. 背景与目标

### 1.1 背景

`project-layout` 是一套 Go 后端项目布局参考，含 v1.0 / v2.0 文档及可编译的 v2.0 四域示例代码。用户希望将其发布为 **GitHub Template Repository**，供自己、团队日常使用，并开源给外部使用者。

### 1.2 目标用户

| 用户 | 需求 |
|------|------|
| 本人 | 快速从模板创建新项目，统一目录规范 |
| 团队 | 共享同一套布局与文档，降低 onboarding 成本 |
| 外部使用者 | README 能看懂、能跑起来、License 清晰 |

### 1.3 成功标准

- [ ] GitHub 仓库勾选 Template repository，「Use this template」可用
- [ ] 从模板创建的新项目，3 步内可 `make build && make run`
- [ ] `init-project.sh` 一键替换 module 路径并通过编译
- [ ] CI 在 `main` 与 `example/v2-full` 分支均通过
- [ ] 文档说明 starter 与完整 v2 示例的区别及获取方式

### 1.4 非目标（YAGNI）

- 不做 Docker / K8s 部署模板
- 不做 golangci-lint CI（后续按需加）
- 不维护第二个 GitHub 仓库
- 不在 starter 中保留 order / payment / notification 示例域

---

## 2. 架构决策

### 2.1 双模板策略：单仓库 + 双分支

GitHub「Use this template」仅复制**默认分支**内容，因此：

| 分支 | 角色 | 获取方式 |
|------|------|----------|
| `main` | 精简 starter（Template 默认产出） | GitHub → Use this template |
| `example/v2-full` | 完整四域 v2.0 参考实现 | `git clone -b example/v2-full` 或文档链接 |

**不采用** `examples/v2-full/` 子目录（同一 go.mod 下 import 易混乱）。  
**不采用** 两个 GitHub 仓库（维护成本翻倍）。

### 2.2 `main` 分支目录（精简 starter）

```text
project-layout/
├── cmd/server/main.go              # 只 wire user 域
├── internal/
│   ├── user/                       # 唯一示例业务域
│   │   ├── handler/
│   │   ├── service/
│   │   ├── repository/
│   │   ├── model/
│   │   └── dto/
│   ├── shared/                     # config db cache logger event middleware utils errors
│   └── router/                     # 只注册 user 路由
├── pkg/client/
├── api/openapi/  api/proto/
├── configs/config.yaml  dev.yaml
├── migrations/
├── scripts/
│   ├── build.sh
│   ├── deploy.sh
│   └── init-project.sh             # 新增
├── test/integration/  test/e2e/
├── docs/                           # version-1.0 / 2.0 / migration / TEMPLATE
├── .github/workflows/ci.yml        # 新增
├── LICENSE                         # MIT，新增
├── Makefile
├── go.mod / go.sum
└── README.md                       # 重写
```

### 2.3 `example/v2-full` 分支

- 保留当前完整四域代码（user / order / payment / notification）
- README 顶部注明：「完整 v2 参考实现，日常开新项目请用 main 模板」
- 与 main 共享同一套 docs/ 与 CI workflow 文件结构

### 2.4 License

**MIT License** — Go 社区模板最常见，个人 / 团队 / 商业友好。

---

## 3. `scripts/init-project.sh`

### 3.1 用法

```bash
./scripts/init-project.sh github.com/antyiwei/my-api
./scripts/init-project.sh github.com/antyiwei/my-api --yes   # 跳过确认
```

### 3.2 行为

1. 校验参数为合法 Go module 路径（含至少一段 `/`）
2. 若 git remote 指向 `antyiwei/project-layout` 且未传 `--force`，提示警告并退出
3. 打印将要替换的内容预览（`project-layout` → 新 module）
4. 无 `--yes` 时等待用户确认
5. 修改 `go.mod` 第一行
6. 在 `*.go`、`go.mod`、`Makefile`（如有硬编码）中替换 import 路径 `project-layout`
7. **不修改** `docs/` 内作为文档示例的字符串
8. 执行 `go mod tidy && go build ./...`
9. 成功时打印后续步骤提示

### 3.3 错误处理

| 场景 | 行为 |
|------|------|
| 参数缺失 / 格式非法 | exit 1，打印用法 |
| `go build` 失败 | exit 1，提示手动检查 |
| 在非 git 目录运行 | 允许，跳过 remote 检测 |

---

## 4. CI：`.github/workflows/ci.yml`

### 4.1 触发条件

- `push` 到 `main`、`example/v2-full`
- `pull_request` 目标为 `main`

### 4.2 Job 步骤

```yaml
- checkout
- setup-go: "1.22"
- go vet ./...
- go test ./...
- go build -o /dev/null ./cmd/server
```

### 4.3 不做

- Docker build、deploy、golangci-lint、release automation

---

## 5. 文档

### 5.1 README.md（重写）

结构优先级：**用法 > 说明 > 文档索引**

1. 一句话介绍 + Template Repo 徽章说明
2. **快速开始（3 步）**：Use template → clone → init-project.sh → make run
3. **这个模板包含什么**：starter vs 完整示例分支
4. 文档索引表（保留现有）
5. 推荐阅读顺序
6. 从 starter 演进到多域（链到 version-2.0 / migration-guide）
7. License + Contributing 简述

### 5.2 docs/TEMPLATE.md（新增）

面向本人与团队的运维手册：

1. 首次发布：git push、GitHub Settings 勾选 Template repository、添加 topics
2. 从模板创建新项目的标准流程
3. `example/v2-full` 分支维护：何时更新、如何从 main cherry-pick shared 改动
4. 团队 module 命名约定（建议 `github.com/antyiwei/<项目名>`）

### 5.3 现有文档

- `docs/version-1.0.md`、`version-2.0.md`、`migration-guide.md` 保留，不改动核心内容
- version-2.0 文首可加一句：完整可运行示例见 `example/v2-full` 分支

---

## 6. GitHub 发布与配置

### 6.1 首次推送

```bash
git init
git add .
git commit -m "feat: initial project-layout template with v2 starter scaffold"
git remote add origin git@github.com:antyiwei/project-layout.git
git branch -M main
git push -u origin main
```

### 6.2 创建完整示例分支

在 main 精简完成**之前**，先基于当前完整代码创建分支：

```bash
# 在当前完整代码状态下
git checkout -b example/v2-full
git push -u origin example/v2-full
git checkout main
# 然后在 main 上执行精简
```

> **顺序很重要**：先保存完整版到分支，再精简 main，避免丢失四域示例。

### 6.3 GitHub Settings

| 设置项 | 值 |
|--------|-----|
| Template repository | ✅ 勾选 |
| Default branch | `main` |
| Topics | `go`, `golang`, `template`, `project-layout`, `gin`, `clean-architecture` |
| Description | Go backend project layout template with domain-driven structure |

---

## 7. 分支维护流程

### 7.1 日常开发（main）

- 新功能 / 文档改进在 `main` 进行
- starter 保持「一个 user 域 + shared + router」最小可运行状态

### 7.2 同步到 example/v2-full

当 `shared/`、`router/` 模式、CI、init 脚本等**公共部分**在 main 更新时：

```bash
git checkout example/v2-full
git cherry-pick <commit>   # 或 merge main（需手动解决 order 等域冲突）
git push
```

**原则**：main 的 shared/router 改动应定期同步到 example/v2-full；业务域示例代码仅在 example 分支维护。

### 7.3 版本标签（可选）

- `v1.0.0` — starter 首个稳定 Template 版本
- 不在 example 分支单独打 tag，除非有外部引用需求

---

## 8. 实施顺序

| 步骤 | 内容 | 验证 |
|------|------|------|
| 1 | 在当前完整代码上创建并 push `example/v2-full` 分支 | 分支可 clone、可 build |
| 2 | 精简 `main`：删 order/payment/notification，简化 router/main | `go build ./...` |
| 3 | 新增 `scripts/init-project.sh` | 本地模拟替换 module 后 build 通过 |
| 4 | 新增 `.github/workflows/ci.yml` | push 后 Actions 绿 |
| 5 | 新增 `LICENSE`（MIT） | — |
| 6 | 重写 `README.md`、新增 `docs/TEMPLATE.md` | 人工读一遍 |
| 7 | `git init` + push 到 GitHub（若尚未） | remote 可访问 |
| 8 | GitHub Settings 勾选 Template repository | Use this template 按钮出现 |
| 9 | 端到端验证：Use template 创建测试仓库 → init → build → run | 全流程 < 5 分钟 |

---

## 9. 风险与缓解

| 风险 | 缓解 |
|------|------|
| 先精简 main 导致丢失完整示例 | **步骤 1 必须先创建 example/v2-full** |
| init 脚本误改 docs 示例 | 脚本排除 `docs/` 目录 |
| 两分支 shared 代码漂移 | docs/TEMPLATE.md 写清 cherry-pick 流程 |
| 用户忘记跑 init-project.sh | README 加粗为必填步骤；init 脚本检测默认 module 名并警告 |

---

## 10. 审阅清单

- [ ] 双分支策略（main starter + example/v2-full）是否符合预期
- [ ] MIT License 是否 OK
- [ ] init-project.sh 行为（含 --yes / --force）是否合理
- [ ] CI 范围是否足够
- [ ] 实施顺序是否清晰

**审阅通过后**：进入 implementation plan 阶段（writing-plans）。
