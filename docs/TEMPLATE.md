# Template 使用与维护手册

面向本人与团队的操作指南。

## 首次发布到 GitHub

```bash
git remote add origin git@github.com:antyiwei/project-layout.git
git branch -M main
git push -u origin main
git push -u origin example/v2-full
```

### GitHub Settings

| 设置项 | 值 |
|--------|-----|
| **Template repository** | ✅ 勾选 |
| **Default branch** | `main` |
| **Topics** | `go`, `golang`, `template`, `project-layout`, `gin`, `clean-architecture` |
| **Description** | Go backend project layout template with domain-driven structure |

勾选 Template repository 后，仓库页会出现 **Use this template** 按钮。

## 从模板创建新项目（团队标准流程）

1. GitHub → **Use this template** → 创建 `my-api`
2. `git clone git@github.com:antyiwei/my-api.git`
3. `./scripts/init-project.sh github.com/antyiwei/my-api --yes`
4. `make build && make run`
5. 删除或改写 `internal/user/` 为你的第一个业务域
6. 业务变复杂时读 `docs/version-2.0.md`

### Module 命名约定

建议统一使用：

```text
github.com/antyiwei/<项目名>
```

## 双分支说明

| 分支 | 用途 |
|------|------|
| `main` | Template 默认产出，精简 starter |
| `example/v2-full` | 完整四域参考实现，不给 Template 用 |

GitHub「Use this template」**只复制 main 分支**。

## 维护 example/v2-full 分支

当 `main` 更新了 `shared/`、`router/` 模式、CI、init 脚本等公共部分时，同步到 example 分支：

```bash
git checkout example/v2-full
git cherry-pick <commit-hash>   # 来自 main 的公共改动
# 如有冲突，保留 example 分支的四域业务代码
git push
```

**原则**：

- `main`：只维护 starter（单 user 域）
- `example/v2-full`：维护完整四域示例 + 接收公共改动

## init-project.sh 说明

```bash
./scripts/init-project.sh github.com/antyiwei/my-api [--yes] [--force]
```

| 参数 | 作用 |
|------|------|
| `--yes` | 跳过确认 |
| `--force` | 在模板源仓库中强制运行（一般不需要） |

脚本会替换 `project-layout` module 路径（**不修改 docs/**），然后 `go mod tidy && go build ./...`。

## 发布版本（可选）

```bash
git tag v1.0.0
git push origin v1.0.0
```

仅在 main starter 稳定时打 tag。
