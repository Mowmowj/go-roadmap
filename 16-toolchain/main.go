// ============================================================================
// 第16章：Go 工具链 🟡 中级 | ⭐ 必须学习
// ============================================================================
// 本章内容：
// 1. 核心命令 (go run/build/install/fmt/vet/test)
// 2. 代码格式化与质量
// 3. 构建与交叉编译
// 4. 性能分析 (pprof)
// 5. 代码生成 (go generate)
// 6. Build Tags / 构建约束
// 7. 实用工具
// ============================================================================

package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	fmt.Println("=== 第16章：Go 工具链 ===")

	// ========================================================================
	// 1. 核心命令
	// ========================================================================
	fmt.Println("\n=== 1. 核心命令 ===")
	fmt.Println(`
  go run main.go         — 编译并运行（不生成二进制）
  go build               — 编译生成二进制文件
  go install             — 编译并安装到 $GOPATH/bin
  go fmt ./...           — 格式化代码（等价于 gofmt -w）
  go vet ./...           — 静态分析，检测常见错误
  go test ./...          — 运行测试
  go test -race ./...    — 竞态检测
  go test -cover ./...   — 测试覆盖率
  go clean                — 清理编译缓存
  go doc fmt.Println     — 查看文档
  go version             — 当前 Go 版本
  go env                 — 查看所有环境变量`)

	fmt.Printf("\n  当前环境:\n")
	fmt.Printf("  Go 版本: %s\n", runtime.Version())
	fmt.Printf("  操作系统: %s\n", runtime.GOOS)
	fmt.Printf("  架构:    %s\n", runtime.GOARCH)
	fmt.Printf("  CPU 核心: %d\n", runtime.NumCPU())
	fmt.Printf("  Goroutines: %d\n", runtime.NumGoroutine())

	// ========================================================================
	// 2. 代码格式化与质量
	// ========================================================================
	fmt.Println("\n=== 2. 代码格式化与质量 ===")
	fmt.Println(`
  go fmt ./...              — 标准格式化（必做）
  goimports -w .            — 格式化 + 自动管理 import
  go vet ./...              — 官方静态分析
  golangci-lint run         — 聚合多个 linter（推荐）
  staticcheck ./...         — 高质量静态分析

  安装:
  go install golang.org/x/tools/cmd/goimports@latest
  go install honnef.co/go/tools/cmd/staticcheck@latest
  
  golangci-lint (推荐):
  brew install golangci-lint  # macOS
  # 或者
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

  golangci-lint 集成的常用 linter:
  - errcheck:    检查未处理的 error
  - gosimple:    简化代码建议
  - govet:       go vet 的增强
  - ineffassign: 检测无效赋值
  - staticcheck: 高级静态分析
  - unused:      未使用的代码
  - gosec:       安全检查
  - misspell:    拼写检查`)

	// ========================================================================
	// 3. 构建与交叉编译
	// ========================================================================
	fmt.Println("\n=== 3. 构建与交叉编译 ===")
	os.Stdout.WriteString(`
  # 基本构建
  go build -o myapp          # 指定输出名
  go build -v                # 显示编译的包
  go build -ldflags "-s -w"  # 去除调试信息,减小体积

  # 注入版本信息
  go build -ldflags "-X main.version=1.0.0 -X main.buildTime=$(date -u +%Y%m%d%H%M%S)"

  # 交叉编译 — Go 的杀手级特性！
  GOOS=linux   GOARCH=amd64 go build -o myapp-linux
  GOOS=windows GOARCH=amd64 go build -o myapp.exe
  GOOS=darwin  GOARCH=arm64 go build -o myapp-mac
  GOOS=linux   GOARCH=arm64 go build -o myapp-arm

  # 支持的平台列表
  go tool dist list

  # 静态链接（完全独立的二进制）
  CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' -o myapp

  # 减小二进制体积
  go build -ldflags "-s -w"   # 去除符号表和调试信息
  # 进一步压缩
  upx --best myapp            # 需要安装 upx
`)

	// ========================================================================
	// 4. 性能分析 (Profiling)
	// ========================================================================
	fmt.Println("\n=== 4. 性能分析 ===")
	fmt.Println(`
  📚 Go 内置强大的性能分析工具

  # CPU 分析
  import "runtime/pprof"
  f, _ := os.Create("cpu.prof")
  pprof.StartCPUProfile(f)
  defer pprof.StopCPUProfile()
  // ... 运行代码 ...
  go tool pprof cpu.prof       # 交互式分析

  # 内存分析
  f, _ := os.Create("mem.prof")
  pprof.WriteHeapProfile(f)
  go tool pprof mem.prof

  # HTTP pprof（给服务器用）
  import _ "net/http/pprof"
  go func() { http.ListenAndServe(":6060", nil) }()
  
  # 然后访问:
  go tool pprof http://localhost:6060/debug/pprof/heap
  go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
  go tool pprof http://localhost:6060/debug/pprof/goroutine

  # 基准测试分析
  go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof
  go tool pprof cpu.prof

  # 常用 pprof 命令:
  top10        — 显示最耗时的 10 个函数
  list funcName — 显示函数的详细分析
  web          — 在浏览器中查看调用图
  svg          — 生成 SVG 图

  # 执行追踪
  go test -trace=trace.out
  go tool trace trace.out  # 浏览器可视化

  # 逃逸分析
  go build -gcflags="-m" main.go
  go build -gcflags="-m -m" main.go  # 更详细`)

	// ========================================================================
	// 5. 代码生成
	// ========================================================================
	fmt.Println("\n=== 5. 代码生成 ===")
	fmt.Println(`
  # go generate — 在源码中嵌入生成命令

  //go:generate stringer -type=Color
  type Color int
  const (
      Red Color = iota
      Green
      Blue
  )
  // 运行 go generate ./... 会生成 color_string.go

  # 常用生成工具:
  - stringer:     为 iota 生成 String() 方法
  - mockgen:      生成接口 mock
  - protoc-gen-go: 从 .proto 生成 Go 代码
  - sqlc:         从 SQL 生成类型安全代码
  - ent:          从 schema 生成 ORM 代码
  - wire:         Google 的依赖注入代码生成

  安装 stringer:
  go install golang.org/x/tools/cmd/stringer@latest`)

	// ========================================================================
	// 6. Build Tags / 构建约束
	// ========================================================================
	fmt.Println("\n=== 6. Build Tags ===")
	fmt.Println(`
  # 文件头部的构建约束（Go 1.17+ 语法）
  //go:build linux
  //go:build !windows
  //go:build linux && amd64
  //go:build linux || darwin
  //go:build ignore        // never build this file

  # 按文件名自动匹配（约定）
  config_linux.go     — 只在 Linux 编译
  config_darwin.go    — 只在 macOS 编译
  config_windows.go   — 只在 Windows 编译

  # 自定义标签
  //go:build integration
  // 运行: go test -tags=integration ./...

  # 实际用途：
  1. 平台特定代码（syscall 等）
  2. 集成测试 vs 单元测试
  3. 开发环境 vs 生产环境
  4. 有/无 CGO`)

	// ========================================================================
	// 7. 其他实用工具
	// ========================================================================
	fmt.Println("\n=== 7. 实用工具 ===")
	fmt.Println(`
  # 文档
  go doc fmt                 — 包文档
  go doc fmt.Println         — 函数文档
  godoc -http=:6060          — 本地文档服务器
  pkgsite                    — 新版文档工具

  # 依赖分析
  go list -m all             — 列出所有依赖
  go list -m -json all       — JSON 格式
  go mod graph               — 依赖图
  go mod why <pkg>           — 为什么需要某包

  # 漏洞检查 (Go 1.18+)
  go install golang.org/x/vuln/cmd/govulncheck@latest
  govulncheck ./...

  # 调试
  go install github.com/go-delve/delve/cmd/dlv@latest
  dlv debug main.go          — 启动调试器
  dlv test                   — 调试测试

  # Delve 常用命令:
  break main.main            — 设置断点
  continue                   — 继续执行
  next                       — 下一步
  step                       — 进入函数
  print <var>                — 打印变量
  goroutines                 — 列出 goroutines`)

	// ========================================================================
	// 8. 版本信息注入示例
	// ========================================================================
	fmt.Println("\n=== 8. 版本信息注入 ===")
	fmt.Printf("  version: %s\n", version)
	fmt.Printf("  buildTime: %s\n", buildTime)
	fmt.Println(`
  编译时注入:
  go build -ldflags "-X main.version=1.0.0 -X main.buildTime=2024-01-01" -o myapp
  ./myapp`)

	fmt.Println(`
📚 工具链最佳实践：
1. 每次提交前: go fmt + go vet + go test
2. CI/CD 中: golangci-lint + go test -race -cover
3. 交叉编译: GOOS + GOARCH 覆盖目标平台
4. 性能优化: 先 benchmark，再 pprof，后优化
5. 二进制优化: -ldflags "-s -w" 减小体积
6. 安全检查: govulncheck 检测已知漏洞
7. 调试: delve 是 Go 最好的调试器`)
}

// 通过 -ldflags 注入的变量
var (
	version   = "dev"
	buildTime = "unknown"
)

// ============================================================================
// 💻 运行方式：
//   cd 16-toolchain && go run main.go
//
// 📝 练习：
// 1. 用 -ldflags 注入版本信息并构建
// 2. 交叉编译到 Linux/Windows
// 3. 用 pprof 分析一段代码的性能
// 4. 用 golangci-lint 分析你的项目
// 5. 用 go build -gcflags="-m" 查看逃逸分析结果
// ============================================================================
