# GitHub Template Repo Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 project-layout 发布为 GitHub Template Repo（main=精简 starter，example/v2-full=完整四域示例）。

**Architecture:** 单仓库双分支；main 仅保留 user 域 + shared + router；init-project.sh 替换 module 路径；GitHub Actions CI。

**Tech Stack:** Go 1.22+、Gin、bash、GitHub Actions

---

### Task 1: 保存完整示例到 example/v2-full 分支

- [ ] git init + 首次 commit（完整四域代码）
- [ ] 创建 example/v2-full 分支并添加分支说明 README

### Task 2: 精简 main 分支

- [ ] 删除 order/payment/notification 域及 order_router.go
- [ ] 更新 router.go、main.go 仅 wire user
- [ ] go build ./... 验证

### Task 3: 模板工具链

- [ ] scripts/init-project.sh
- [ ] .github/workflows/ci.yml
- [ ] LICENSE (MIT)

### Task 4: 文档

- [ ] README.md 重写
- [ ] docs/TEMPLATE.md
- [ ] docs/version-2.0.md 补充分支说明

### Task 5: Git 推送

- [ ] push main + example/v2-full 到 origin
