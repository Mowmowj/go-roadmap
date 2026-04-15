# 🚀 Go 语言完整学习路线图 (基于 roadmap.sh/golang)

> 从零到精通 Go 语言的系统化学习项目，所有代码可直接运行

## 📋 学习等级说明

| 标记 | 含义 | 说明 |
|------|------|------|
| 🟢 初级 | Beginner | 零基础必学，Go 语言入门核心 |
| 🟡 中级 | Intermediate | 有基础后必学，日常开发必备技能 |
| 🔴 高级 | Advanced | 深入理解，高级开发和性能优化 |
| ⭐ 必须学习 | Required | 不学无法进行后续内容 |
| 📌 推荐学习 | Recommended | 强烈建议掌握 |
| 💡 可选学习 | Optional | 根据工作方向选学 |

## 📂 项目结构

### 第一阶段：语言基础 (🟢 初级 ⭐ 必须)
```
01-getting-started/     Go 简介、环境搭建、Hello World、go 命令
02-variables-constants/ 变量声明、var vs :=、零值、const 与 iota、作用域
03-data-types/          布尔、数值类型、字符串、Rune、类型转换
04-composite-types/     数组、切片、Map、结构体、JSON 标签
05-control-flow/        if/else、switch、for 循环、range、break/continue
06-functions/           函数基础、可变参数、多返回值、匿名函数、闭包
```

### 第二阶段：核心进阶 (🟡 中级 ⭐ 必须)
```
07-pointers/            指针基础、结构体指针、Map/Slice 与指针、内存管理
08-methods-interfaces/  方法、接口、空接口、类型断言、类型 Switch
09-generics/            泛型函数、泛型类型、类型约束、类型推断
10-error-handling/      error 接口、自定义错误、Wrap/Unwrap、panic/recover
11-modules-packages/    模块管理、包、导入规则、第三方包
```

### 第三阶段：并发编程 (🔴 高级 ⭐ 必须)
```
12-concurrency/         Goroutine、Channel、Select、sync 包、context、并发模式
```

### 第四阶段：标准库与工程实践 (🟡 中级 ⭐ 必须)
```
13-standard-library/    fmt、strings、os、io、time、encoding/json、net/http、sort、regexp
14-testing/             单元测试、表驱动测试、Mock、基准测试、示例测试
```

### 第五阶段：生态与工具 (🟡 中级 📌 推荐)
```
15-ecosystem/           Web 框架(Gin/Echo)、数据库(GORM/sqlx)、CLI(Cobra)、REST API 实战
16-toolchain/           go 命令、代码质量、Lint、交叉编译、pprof 性能分析
```

### 第六阶段：高级主题 (🔴 高级 💡 可选)
```
17-advanced/            反射、unsafe 包、编译器指令、go:embed、内存模型、插件系统
```

## 🏃 如何运行

每个章节目录下的代码都可以独立运行：

```bash
# 运行某个章节
cd 01-getting-started && go run main.go

# 运行测试
cd 14-testing && go test -v ./...

# 运行 REST API 示例（在 15-ecosystem 中）
# 取消 startServer() 调用的注释后运行
cd 15-ecosystem && go run main.go
```

## 📖 学习建议

1. **按顺序学习**：章节之间有依赖关系，建议从 01 开始
2. **动手实践**：每个文件中的代码都可以运行，建议修改后再运行观察变化
3. **先必须后可选**：优先完成 ⭐ 标记内容，再学习 💡 可选内容
4. **做练习**：每个章节末尾都有练习建议，动手完成它们
5. **参考文档**：遇到不清楚的地方，查阅 [Go 官方文档](https://go.dev/doc/)

## 🎯 学习完成后你将掌握

- Go 语言所有核心语法和数据结构
- 并发编程模型（Goroutine、Channel、sync）
- 测试驱动开发（TDD）和基准测试
- Web API 开发和数据库操作
- Go 工具链和代码质量管理
- 性能分析和优化技巧
- 反射、unsafe 等高级主题
