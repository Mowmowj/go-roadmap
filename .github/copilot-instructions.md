# Project Guidelines

## Overview

Go 语言完整学习路线图 — 18 个章节的交互式学习项目，包含可运行的 Go 代码示例和一个 Python 静态网站生成器。

## Architecture

```
<NN>-<kebab-title>/main.go   # Go 教学章节（01–18）
scripts/build-site.py         # 静态网站生成器（Python 3，无外部依赖）
scripts/assets/theme.js       # 客户端主题切换（暗/亮）
site/                         # 生成输出（gitignored）
.github/workflows/            # GitHub Pages 部署
```

**站点生成器核心数据结构**（`build-site.py`）：

- `CHAPTERS` — 章节元数据列表（dir, title, emoji, level, tag）
- `PHASES` — 学习阶段分组（name, range, color, level），`range` 为 **1-indexed 闭区间**

## Build and Test

```bash
# 运行单个章节
cd <NN>-<name> && go run main.go

# 运行测试（第14章）
cd 14-testing && go test -v ./...

# 生成网站
python3 scripts/build-site.py

# 本地预览
cd site && python3 -m http.server 8080
```

Go module: `go 1.19`，唯一依赖 `github.com/lib/pq`（第 18 章 PostgreSQL）。

## Conventions

### 添加新章节（必须同时完成以下步骤）

1. 创建 `<NN>-<kebab-name>/main.go`
2. 在 `CHAPTERS` 列表中添加条目（按序号位置）
3. 在 `PHASES` 列表中更新或新增阶段的 `range`

遗漏任一步骤会导致新章节在网站中缺失。

### Go 文件格式

```go
// ============================================================================
// 第XX章：<标题> <emoji> <等级> | <标签>
// ============================================================================
// 本章内容：
// 1. ...
// ============================================================================

package main
```

- 节标记：`// N. 标题` + `=` 分隔线
- 概念块：`// 📚 概念：<名称>`
- 注释全部使用中文

### 网站生成器

- 纯 Python 标准库，无外部依赖
- CSS 暗色优先：`:root` 为暗色，`[data-theme="light"]` 覆盖
- Go Playground API（`play.golang.org/compile`）仅在 HTTPS 下可用；localhost 会触发 CORS，已有 fallback UI

## Pitfalls

- **CORS on localhost**：Go Playground API 在 `http://` 下不可用，部署到 GitHub Pages 后正常
- **PHASES range**：1-indexed 闭区间，添加章节后必须核查

详细说明见 [README.md](../README.md)。
