# 本地创建新项目（任意 Git 平台）

面向习惯**先在本地开发、之后再推到任意 Git 托管平台**的流程——Gitee、GitHub、GitLab、公司自建 Git 均可。

模板源在 [GitHub](https://github.com/antyiwei/project-layout)，**新项目的 module 路径与 remote 由你自己决定**。

---

## 核心原则

只需记住两件事：

| 概念 | 说明 |
|------|------|
| **Go module 路径** | `init-project.sh` 的第一个参数，写入 `go.mod` 和各文件 import |
| **Git remote** | `git remote add origin ...`，决定代码推到哪里 |

**推荐**：module 路径与将来仓库的「主机 + 路径」一致，便于 `go get`、团队协作和 CI。

```text
Git remote URL          →    Go module 路径
git@gitee.com:antyiwei/my-api.git     →    gitee.com/antyiwei/my-api
git@github.com:acme/pay-svc.git       →    github.com/acme/pay-svc
git@git.company.com:backend/api.git   →    git.company.com/backend/api
```

---

## 如何确定 module 路径

从计划使用的 **Git 仓库地址** 推导，去掉协议和 `.git` 后缀：

| Git remote（示例） | module 路径 |
|-------------------|-------------|
| `git@github.com:antyiwei/my-api.git` | `github.com/antyiwei/my-api` |
| `https://github.com/acme/pay-svc.git` | `github.com/acme/pay-svc` |
| `git@gitee.com:antyiwei/my-api.git` | `gitee.com/antyiwei/my-api` |
| `git@gitlab.com:team/backend.git` | `gitlab.com/team/backend` |
| `git@git.company.com:org/my-api.git` | `git.company.com/org/my-api` |
| `https://codeup.aliyun.com/org/repo.git` | `codeup.aliyun.com/org/repo` |

**公司自建 Git**：用你们实际 clone 地址的主机名 + 仓库路径即可，无需是 GitHub / Gitee。

**自定义域名**（较少见）：若团队用 `go.company.com/my-api` 做 module proxy，则 module 填该域名路径。

> 仓库尚未创建时，也可以先按「将来会用的地址」填写 module；之后在对应平台建**同名空仓库**再 push。

---

## 标准流程

### 1. 克隆模板

```bash
git clone git@github.com:antyiwei/project-layout.git my-api
cd my-api
```

完整四域参考（可选）：

```bash
git clone -b example/v2-full git@github.com:antyiwei/project-layout.git my-api-full
```

### 2. 断开与模板仓库的关联

```bash
rm -rf .git
git init
git branch -M main
```

### 3. 初始化 module（必做）

将 `<你的-module-路径>` 换成上表推导出的路径：

```bash
./scripts/init-project.sh <你的-module-路径> --yes
```

示例：

```bash
# 个人 Gitee
./scripts/init-project.sh gitee.com/antyiwei/my-api --yes

# 公司 GitLab
./scripts/init-project.sh gitlab.com/acme-corp/order-api --yes

# 公司自建
./scripts/init-project.sh git.corp.example.com/backend/my-api --yes
```

### 4. 构建运行

```bash
make build && make run
curl http://localhost:8080/api/v1/users/1
```

---

## 之后推到任意平台

在 **Gitee / GitHub / GitLab / 公司 Git** 上新建**空仓库**（不要初始化 README），然后：

```bash
git add -A
git commit -m "chore: init from project-layout template"

# remote 与 module 对应即可，示例：
git remote add origin git@gitee.com:antyiwei/my-api.git
# git remote add origin git@github.com:acme/my-api.git
# git remote add origin git@git.company.com:backend/my-api.git

git push -u origin main
```

**检查一致性**（可选）：

```text
module 路径：  gitee.com/antyiwei/my-api
remote 主机：  gitee.com
remote 路径：  antyiwei/my-api
```

三者逻辑上应对得上。

---

## 一句话备忘

```bash
git clone git@github.com:antyiwei/project-layout.git <目录名> && cd <目录名>
rm -rf .git && git init && git branch -M main
./scripts/init-project.sh <主机>/<组织或用户>/<仓库名> --yes
make build && make run
# 在任意 Git 平台建空仓库后：
git remote add origin <你的-clone-地址> && git push -u origin main
```

---

## 备选：保留 Git 历史

```bash
git clone git@github.com:antyiwei/project-layout.git my-api
cd my-api
git remote remove origin
./scripts/init-project.sh <你的-module-路径> --yes
```

---

## 常见问题

### Q：平台还没定，能先开发吗？

可以。先选一个**预计不会改**的路径跑 `init-project.sh`。若后来换平台，需再跑一次脚本或手动改 `go.mod` 和 import（成本较高，尽量事先想好）。

### Q：module 路径和 remote 必须完全一致吗？

- **应用项目**（不对外 `go get）：路径与 remote 一致即可，团队内统一最重要。
- **要发布的库**：建议与真实托管地址一致，或配置 [GOPROXY / vanity import](https://go.dev/ref/mod#vcs-find)。

### Q：`init-project.sh` 报 template source repo

当前仍关联 `antyiwei/project-layout`。执行 `rm -rf .git && git init`，或 `git remote remove origin` 后再运行。

### Q：换公司 / 换平台怎么办？

重新执行 `./scripts/init-project.sh <新-module-路径> --yes`，并更新 `git remote set-url origin <新地址>`。

### Q：业务变复杂后怎么演进？

见 [version-2.0.md](version-2.0.md)、[migration-guide.md](migration-guide.md)。完整示例：`example/v2-full` 分支。
