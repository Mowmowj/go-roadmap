// ============================================================================
// 第17章：高级主题 🔴 高级 | 📌 可选学习
// ============================================================================
// 本章内容：
// 1. 反射 (reflect)
// 2. unsafe 包
// 3. CGO 基础
// 4. 编译器指令 (compiler directives)
// 5. 内存模型深入
// 6. Plugin 系统
// ============================================================================
// ⚠️ 高级内容，日常开发较少使用
// 理解这些概念有助于阅读标准库源码和框架代码
// ============================================================================

package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// ============================================================================
// 1. 反射 (Reflection)
// ============================================================================

// 📚 概念：反射
// 反射允许在运行时检查和操作类型信息
// 核心类型：reflect.Type 和 reflect.Value
// 三大法则（Rob Pike）：
//   1. 从接口值到反射对象
//   2. 从反射对象到接口值
//   3. 要修改反射对象，值必须可设置（addressable）

type Config struct {
	Host    string `json:"host" default:"localhost" required:"true"`
	Port    int    `json:"port" default:"8080"`
	Debug   bool   `json:"debug" default:"false"`
	Timeout int    `json:"timeout" default:"30"`
}

// 用反射打印结构体的所有字段和标签
func inspectStruct(v interface{}) {
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	// 如果是指针，获取底层类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		val = val.Elem()
	}

	fmt.Printf("  类型: %s (Kind: %s)\n", t.Name(), t.Kind())
	fmt.Printf("  字段数: %d\n", t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := val.Field(i)
		fmt.Printf("  - %s (%s) = %v", field.Name, field.Type, value.Interface())

		// 读取 tag
		if tag := field.Tag.Get("json"); tag != "" {
			fmt.Printf("  json:%q", tag)
		}
		if tag := field.Tag.Get("default"); tag != "" {
			fmt.Printf("  default:%q", tag)
		}
		if tag := field.Tag.Get("required"); tag != "" {
			fmt.Printf("  required:%q", tag)
		}
		fmt.Println()
	}
}

// 用反射设置字段值
func setFieldByName(obj interface{}, fieldName string, value interface{}) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("must pass a pointer")
	}
	v = v.Elem()
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field %q not found", fieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("field %q cannot be set", fieldName)
	}
	newVal := reflect.ValueOf(value)
	if field.Type() != newVal.Type() {
		return fmt.Errorf("type mismatch: %s vs %s", field.Type(), newVal.Type())
	}
	field.Set(newVal)
	return nil
}

// 用反射调用方法
type Calculator struct{}

func (c Calculator) Add(a, b int) int { return a + b }
func (c Calculator) Mul(a, b int) int { return a * b }

func callMethod(obj interface{}, methodName string, args ...interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(obj)
	m := v.MethodByName(methodName)
	if !m.IsValid() {
		return nil, fmt.Errorf("method %q not found", methodName)
	}

	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	results := m.Call(in)
	out := make([]interface{}, len(results))
	for i, r := range results {
		out[i] = r.Interface()
	}
	return out, nil
}

// ============================================================================
// 2. unsafe 包
// ============================================================================

// 📚 概念：unsafe
// 提供绕过 Go 类型安全系统的操作
// 主要用途：
// - 获取类型大小和对齐方式
// - 指针算术
// - 低级内存操作
// - 与 C 代码交互
//
// ⚠️ 使用 unsafe 的代码不保证跨版本兼容

type Example struct {
	A bool  // 1 byte
	B int64 // 8 bytes
	C int32 // 4 bytes
}

type OptimizedExample struct {
	B int64 // 8 bytes
	C int32 // 4 bytes
	A bool  // 1 byte
}

// ============================================================================
// 3. 编译器指令
// ============================================================================

// 📚 常用编译器指令：
// go:noinline  — 禁止内联
// go:nosplit   — 禁止栈分裂检查
// go:linkname  — 链接到其他包的私有函数（谨慎使用）
// go:embed     — 嵌入文件到二进制（Go 1.16+）
// go:generate  — 代码生成命令
// go:build     — 构建约束
// (注意：实际使用时 // 和 go: 之间没有空格)

//go:noinline
func noInlineFunc() int {
	return 42
}

// ============================================================================
func main() {
	// 1. 反射
	fmt.Println("=== 1. 反射 (reflect) ===")

	cfg := Config{
		Host:    "localhost",
		Port:    8080,
		Debug:   true,
		Timeout: 30,
	}

	fmt.Println("\n--- 检查结构体 ---")
	inspectStruct(cfg)

	// 通过反射修改值
	fmt.Println("\n--- 反射修改值 ---")
	fmt.Printf("  修改前: Port=%d\n", cfg.Port)
	err := setFieldByName(&cfg, "Port", 3000)
	if err != nil {
		fmt.Println("  错误:", err)
	}
	fmt.Printf("  修改后: Port=%d\n", cfg.Port)

	// 反射调用方法
	fmt.Println("\n--- 反射调用方法 ---")
	calc := Calculator{}
	result, err := callMethod(calc, "Add", 3, 5)
	if err == nil {
		fmt.Printf("  Calculator.Add(3, 5) = %v\n", result[0])
	}
	result, err = callMethod(calc, "Mul", 4, 6)
	if err == nil {
		fmt.Printf("  Calculator.Mul(4, 6) = %v\n", result[0])
	}

	// 反射检查类型
	fmt.Println("\n--- reflect.TypeOf / ValueOf ---")
	values := []interface{}{42, "hello", true, 3.14, []int{1, 2, 3}}
	for _, v := range values {
		fmt.Printf("  值: %-12v 类型: %-12s Kind: %s\n",
			v, reflect.TypeOf(v), reflect.TypeOf(v).Kind())
	}

	// DeepEqual
	fmt.Println("\n--- reflect.DeepEqual ---")
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	s3 := []int{1, 2, 4}
	fmt.Printf("  [1,2,3] == [1,2,3]: %v\n", reflect.DeepEqual(s1, s2))
	fmt.Printf("  [1,2,3] == [1,2,4]: %v\n", reflect.DeepEqual(s1, s3))

	// 2. unsafe 包
	fmt.Println("\n=== 2. unsafe 包 ===")

	// Sizeof — 类型大小
	fmt.Println("\n--- Sizeof ---")
	fmt.Printf("  bool:    %d bytes\n", unsafe.Sizeof(bool(false)))
	fmt.Printf("  int:     %d bytes\n", unsafe.Sizeof(int(0)))
	fmt.Printf("  int32:   %d bytes\n", unsafe.Sizeof(int32(0)))
	fmt.Printf("  int64:   %d bytes\n", unsafe.Sizeof(int64(0)))
	fmt.Printf("  float64: %d bytes\n", unsafe.Sizeof(float64(0)))
	fmt.Printf("  string:  %d bytes\n", unsafe.Sizeof(string("")))
	fmt.Printf("  slice:   %d bytes\n", unsafe.Sizeof([]int{}))

	// 结构体内存布局与对齐
	fmt.Println("\n--- 结构体内存对齐 ---")
	fmt.Printf("  Example 大小:          %d bytes (有填充)\n", unsafe.Sizeof(Example{}))
	fmt.Printf("  OptimizedExample 大小: %d bytes (优化后)\n", unsafe.Sizeof(OptimizedExample{}))

	// 字段偏移
	fmt.Println("\n--- Offsetof ---")
	var e Example
	fmt.Printf("  Example.A offset: %d\n", unsafe.Offsetof(e.A))
	fmt.Printf("  Example.B offset: %d\n", unsafe.Offsetof(e.B))
	fmt.Printf("  Example.C offset: %d\n", unsafe.Offsetof(e.C))

	fmt.Println(`
  📚 结构体内存优化：
  按字段大小降序排列可以减少内存填充
  Example:   bool(1) + padding(7) + int64(8) + int32(4) + padding(4) = 24
  Optimized: int64(8) + int32(4) + bool(1) + padding(3) = 16`)

	// 3. 编译器指令
	fmt.Println("\n=== 3. 编译器指令 ===")
	fmt.Println(`
  //go:noinline          — 禁止函数内联
  //go:nosplit           — 禁止栈分裂检查
  //go:norace            — 跳过竞态检测
  //go:linkname          — 链接私有符号
  //go:embed file.txt    — 嵌入文件
  //go:generate cmd      — 代码生成
  //go:build tag         — 构建约束`)

	// 4. go:embed 示例
	fmt.Println("\n=== 4. go:embed (Go 1.16+) ===")
	fmt.Println(`
  import "embed"

  //go:embed config.json
  var configData []byte

  //go:embed templates/*.html
  var templates embed.FS

  //go:embed version.txt
  var version string

  用途：将静态文件嵌入二进制
  - 配置文件
  - HTML 模板
  - 静态资源
  - SQL 迁移文件`)

	// 5. 内存模型
	fmt.Println("\n=== 5. Go 内存模型 ===")
	fmt.Println(`
  📚 Go 内存模型定义了多条件下内存读写的可见性规则

  核心概念：happens-before 关系
  如果事件 A happens-before 事件 B，则 A 的写入对 B 可见

  保证的 happens-before:
  1. 包的 init() happens-before main()
  2. go f() happens-before f() 开始执行
  3. ch <- v (发送) happens-before <-ch (接收) 完成
  4. mu.Unlock() happens-before 下一个 mu.Lock()
  5. sync.Once.Do(f) 中 f 的完成 happens-before Do 返回

  ⚠️ 不保证的：
  - 不同 goroutine 中对同一变量的无同步读写
  - 编译器和 CPU 可能重排指令
  
  → 始终使用 channel 或 sync 包来同步！`)

	// 6. Plugin 系统
	fmt.Println("\n=== 6. Plugin (Go 插件) ===")
	fmt.Println(`
  # 编译插件
  go build -buildmode=plugin -o plugin.so plugin.go
  
  # 加载插件
  p, err := plugin.Open("plugin.so")
  sym, err := p.Lookup("Hello")
  hello := sym.(func() string)
  fmt.Println(hello())
  
  限制：
  - 仅 Linux 和 macOS （不支持 Windows）
  - 编译器版本必须完全一致
  - 不能卸载
  - 实际项目中很少使用`)

	fmt.Println(`
📚 高级主题使用建议：
1. 反射：JSON 库、ORM、依赖注入框架使用
   日常代码尽量避免，性能差且失去类型安全
2. unsafe：极少需要，标准库和驱动可能用到
3. go:embed：推荐使用，嵌入静态文件很方便
4. 内存对齐：高性能场景下注意结构体字段排列
5. CGO：需要调用 C 库时使用，但会增加复杂度
6. "清晰胜于巧妙" — Go 格言`)
}

// ============================================================================
// 💻 运行方式：
//   cd 17-advanced && go run main.go
//
// 📝 练习：
// 1. 用反射实现一个简单的结构体验证器
// 2. 比较 Example 和 OptimizedExample 的内存大小
// 3. 用 go:embed 嵌入一个配置文件
// 4. 用 go build -gcflags="-m" 分析逃逸
// ============================================================================
