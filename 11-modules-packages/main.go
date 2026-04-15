// ============================================================================
// 第11章：模块与包 🟡 中级 | ⭐ 必须学习
// ============================================================================
// 本章内容：
// 1. Go 模块系统 (go mod)
// 2. 包 (package) 基础
// 3. 导出规则（大小写）
// 4. 包的组织与设计
// 5. init() 函数
// 6. internal 包
// 7. go.mod 和 go.sum
// 8. 常用 go mod 命令
// ============================================================================
// 📝 注意：本文件主要是概念说明和演示
// 由于模块系统涉及文件组织，大部分内容以注释形式呈现
// ============================================================================

package main

import (
	"fmt"
)

// ============================================================================
// 1. Go 模块系统概述
// ============================================================================

// 📚 概念：Go Modules (Go 1.11+)
//
// 模块是相关 Go 包的集合，以 go.mod 文件定义
// 模块是依赖管理和版本控制的基本单元
//
// 核心命令：
//   go mod init <module-path>  — 初始化新模块
//   go mod tidy               — 添加缺少的依赖，删除未用的
//   go mod download           — 下载依赖到本地缓存
//   go mod vendor             — 将依赖复制到 vendor 目录
//   go mod graph              — 打印依赖图
//   go mod verify             — 验证依赖的完整性
//   go mod edit               — 编辑 go.mod (如替换依赖)
//   go mod why <pkg>          — 解释为什么需要某个依赖
//
// go.mod 文件示例：
//   module github.com/user/myproject
//
//   go 1.21
//
//   require (
//       github.com/gin-gonic/gin v1.9.1
//       go.uber.org/zap v1.26.0
//   )
//
// go.sum 文件：
//   - 记录依赖的加密哈希
//   - 确保可重现构建
//   - 应该提交到版本控制

// ============================================================================
// 2. 包 (Package) 基础
// ============================================================================

// 📚 概念：包
//
// 1. 每个 .go 文件必须声明所属包
// 2. 同一目录下所有文件必须属于同一个包
// 3. main 包是可执行程序的入口
// 4. 包名通常和目录名相同（但不是强制的）
// 5. 一个包可以有多个 .go 文件
//
// 包的导入：
//   import "fmt"                    // 标准库
//   import "github.com/user/pkg"   // 第三方
//   import "./mypackage"           // 相对路径（不推荐）
//
// 导入别名：
//   import (
//       myfmt "fmt"                // 别名导入
//       _ "database/sql"           // 空白导入（只执行 init）
//       . "math"                   // 点导入（直接使用 Sqrt）不推荐
//   )

// ============================================================================
// 3. 导出规则
// ============================================================================

// 📚 Go 的可见性规则非常简单：
// 首字母大写 = 导出（公开的，其他包可访问）
// 首字母小写 = 未导出（私有的，仅包内访问）
//
// 适用于：类型、函数、方法、变量、常量、结构体字段

// 导出的类型和函数
type PublicStruct struct {
	PublicField  string // 导出的字段
	privateField string // 未导出的字段
}

func PublicFunction() string {
	return "I'm accessible from other packages"
}

// 未导出的
type privateStruct struct {
	field string
}

func privateFunction() string {
	return "I'm only accessible within this package"
}

// ============================================================================
// 4. init() 函数
// ============================================================================

// 📚 概念：init() 函数
// - 每个包可以有一个或多个 init() 函数
// - init() 在 main() 之前自动执行
// - 不能被直接调用
// - 执行顺序：导入包的 init() → 当前包的 init() → main()
// - 常用于初始化全局变量、注册驱动、验证配置等

var appConfig map[string]string

func init() {
	// init 在 main 之前运行
	appConfig = map[string]string{
		"env":     "development",
		"version": "1.0.0",
	}
	fmt.Println("[init] 配置已初始化")
}

// 可以有多个 init 函数
func init() {
	fmt.Println("[init] 第二个 init 函数执行")
}

// ============================================================================
// 5. 包的设计原则
// ============================================================================

// 📚 包设计最佳实践：
//
// 1. 命名：
//    - 简短、有意义、全小写
//    - 好: http, json, sync, errors
//    - 坏: utilities, helpers, common, misc
//
// 2. 单一职责：
//    每个包应该有一个清晰的目的
//    如果包名需要 "and" 来描述，可能需要拆分
//
// 3. 项目结构示例：
//    myproject/
//    ├── go.mod
//    ├── go.sum
//    ├── main.go              // package main
//    ├── cmd/                  // 多个可执行程序
//    │   ├── server/main.go
//    │   └── cli/main.go
//    ├── internal/             // 私有包（外部不可导入）
//    │   ├── auth/
//    │   ├── database/
//    │   └── middleware/
//    ├── pkg/                  // 公共库（外部可导入）
//    │   ├── models/
//    │   └── utils/
//    ├── api/                  // API 定义 (protobuf, OpenAPI)
//    ├── configs/              // 配置文件
//    ├── docs/                 // 文档
//    └── tests/                // 集成测试
//
// 4. internal 包：
//    - internal/ 目录下的包只能被其父目录下的包导入
//    - 用于防止外部依赖你的内部实现

// ============================================================================
// 6. 演示代码
// ============================================================================

// 模拟一个简单的包结构
type Logger struct {
	prefix string
	level  string
}

func NewLogger(prefix, level string) *Logger {
	return &Logger{prefix: prefix, level: level}
}

func (l *Logger) Info(msg string) {
	fmt.Printf("[%s] %s: %s\n", l.level, l.prefix, msg)
}

func (l *Logger) Error(msg string) {
	fmt.Printf("[%s] %s: ERROR: %s\n", l.level, l.prefix, msg)
}

// 模拟一个 Config 包
type Config struct {
	Host string
	Port int
	DB   DatabaseConfig
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

func DefaultConfig() Config {
	return Config{
		Host: "localhost",
		Port: 8080,
		DB: DatabaseConfig{
			Driver: "postgres",
			DSN:    "postgres://localhost:5432/mydb",
		},
	}
}

func main() {
	fmt.Println("\n=== Go 模块与包 ===")

	// init() 已经在 main 之前执行
	fmt.Println("\n=== 1. init() 函数 ===")
	fmt.Println("App config:", appConfig)

	// 2. 导出规则演示
	fmt.Println("\n=== 2. 导出规则 ===")
	pub := PublicStruct{PublicField: "visible"}
	fmt.Printf("PublicStruct: %+v\n", pub)
	fmt.Println("PublicFunction:", PublicFunction())
	fmt.Println("privateFunction:", privateFunction())

	priv := privateStruct{field: "hidden"}
	fmt.Printf("privateStruct: %+v\n", priv)

	// 3. 包的使用演示
	fmt.Println("\n=== 3. 包的使用 ===")
	logger := NewLogger("app", "INFO")
	logger.Info("应用启动")
	logger.Error("数据库连接失败")

	config := DefaultConfig()
	fmt.Printf("Config: %+v\n", config)

	// 4. 常用 go mod 命令
	fmt.Println("\n=== 4. go mod 命令速查 ===")
	commands := []struct {
		Cmd  string
		Desc string
	}{
		{"go mod init <path>", "初始化模块"},
		{"go mod tidy", "整理依赖（添加缺少的，删除多余的）"},
		{"go mod download", "下载依赖到缓存"},
		{"go mod vendor", "复制依赖到 vendor/"},
		{"go mod graph", "打印依赖关系图"},
		{"go mod verify", "验证依赖完整性"},
		{"go mod edit -replace=old=new", "替换依赖（本地开发常用）"},
		{"go mod why <pkg>", "解释为什么需要某依赖"},
		{"go get <pkg>@version", "添加/更新依赖"},
		{"go get <pkg>@none", "删除依赖"},
		{"go list -m all", "列出所有依赖"},
		{"go list -m -versions <pkg>", "列出可用版本"},
	}

	for _, cmd := range commands {
		fmt.Printf("  %-40s → %s\n", cmd.Cmd, cmd.Desc)
	}

	// 5. 版本规则
	fmt.Println("\n=== 5. 语义化版本 (SemVer) ===")
	fmt.Println(`
  版本格式: vMAJOR.MINOR.PATCH
  - MAJOR: 不兼容的 API 变更 (v2.0.0)
  - MINOR: 向后兼容的功能添加 (v1.1.0)
  - PATCH: 向后兼容的 bug 修复 (v1.0.1)

  Go 的特殊规则：
  - v0.x.x 视为不稳定，可以随时变更
  - v2+ 需要在模块路径中加版本号
    例: github.com/user/pkg/v2

  版本选择策略：
  go get pkg@latest       — 最新稳定版
  go get pkg@v1.2.3       — 指定版本
  go get pkg@master       — 指定分支
  go get pkg@abc123       — 指定 commit`)

	// 6. 工作区模式
	fmt.Println("\n=== 6. Go 工作区 (Go 1.18+) ===")
	fmt.Println(`
  多模块开发时使用 go.work 文件：
  
  go work init ./module1 ./module2
  
  go.work 文件示例:
    go 1.21
    use (
        ./frontend
        ./backend
        ./shared
    )
  
  好处：
  - 同时开发多个模块
  - 不需要 replace 指令
  - 类似 monorepo 的体验`)

	fmt.Println(`
📚 关键总结：
1. go mod 管理依赖，go.mod 定义模块
2. 包名小写，简短有意义
3. 大写开头=导出，小写开头=私有
4. init() 在 main 前执行，慎用
5. internal/ 目录限制包的可见性
6. "Accept interfaces, return structs"
7. 优先使用标准库和成熟的第三方库
8. go mod tidy 是最常用的命令`)
}

// ============================================================================
// 💻 运行方式：
//   cd 11-modules-packages && go run main.go
//
// 📝 练习：
// 1. 创建一个多包项目：main + helper + model
// 2. 用 go mod init 创建模块，添加第三方依赖
// 3. 理解 go.sum 的作用
// 4. 用 go mod graph 查看依赖关系
// ============================================================================
