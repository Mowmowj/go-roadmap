---
name: add-chapter
description: "添加新的 Go 学习章节。Use when: 新增章节、创建 chapter、add chapter、添加教程。自动完成目录创建、main.go 生成、CHAPTERS/PHASES 注册、站点构建验证。"
argument-hint: "章节编号 标题 等级 标签（如：19 HTTP客户端 中级 推荐）"
---

# 添加新章节

为 Go 学习路线图添加新章节的完整流程。

## 输入参数

从用户输入中提取以下信息（缺失时询问）：

| 参数     | 说明             | 示例                           | 可选值                 |
| -------- | ---------------- | ------------------------------ | ---------------------- |
| 编号     | 两位数字序号     | `19`                           | 下一个可用编号         |
| 标题     | 中文章节标题     | `HTTP 客户端`                  | —                      |
| kebab    | 英文目录名       | `http-client`                  | —                      |
| emoji    | 章节图标         | `🌐`                           | —                      |
| 等级     | 难度等级         | `中级`                         | 初级 / 中级 / 高级     |
| 标签     | 重要性标签       | `推荐`                         | 必须 / 推荐 / 可选     |
| 阶段     | 所属学习阶段名   | `网络编程`                     | 现有阶段名 或 新阶段名 |
| 主题列表 | 章节涵盖的知识点 | `1. GET 请求 2. POST 请求 ...` | —                      |

## 操作步骤

**严格按顺序执行，不可跳过任何步骤。**

### Step 1 — 创建目录和 main.go

创建 `<NN>-<kebab>/main.go`，使用 [main.go 模板](./assets/main.go.tmpl) 的格式。

等级对应 emoji 映射：

- 初级 → `🟢`，标签前缀：必须 →`⭐`，推荐 →`📌`，可选 →`💡`
- 中级 → `🟡`，同上
- 高级 → `🔴`，同上

### Step 2 — 注册到 CHAPTERS

在 `scripts/build-site.py` 的 `CHAPTERS` 列表中，按章节编号顺序插入：

```python
{"dir": "<NN>-<kebab>", "title": "<标题>", "emoji": "<emoji>", "level": "<等级>", "tag": "<标签>"},
```

对齐已有条目的格式。

### Step 3 — 更新 PHASES

在 `scripts/build-site.py` 的 `PHASES` 列表中：

- **归入现有阶段**：扩展该阶段的 `range` 上界包含新编号
- **新建阶段**：追加新条目，`range` 为 `(<NN>, <NN>)`

等级到颜色映射：

```
初级 → color: "#10b981", level: "🟢 初级"
中级 → color: "#f59e0b", level: "🟡 中级"
高级 → color: "#ef4444", level: "🔴 高级"
```

`range` 是 **1-indexed 闭区间**。

### Step 4 — 验证构建

```bash
cd <NN>-<kebab> && go build -o /dev/null . && go vet .
python3 scripts/build-site.py
```

确认输出中包含 `✅ <NN>-<kebab>.html`。

### Step 5 — 填充章节内容

根据用户提供的主题列表，为每个主题编写：

1. `// ========` 节分隔线 + `// N. 标题`
2. `// 📚 概念：<名称>` 概念说明块（中文注释）
3. 可运行的示例代码（在 `func main()` 中）

所有注释使用中文。参考已有章节的风格。

## 检查清单

完成后逐项确认：

- [ ] `<NN>-<kebab>/main.go` 存在且 `go vet` 通过
- [ ] `CHAPTERS` 包含新条目，位置正确
- [ ] `PHASES` 的 range 覆盖新章节编号
- [ ] `python3 scripts/build-site.py` 成功生成对应 HTML
