// ============================================================================
// 第02章：变量与常量 🟢 初级 | ⭐ 必须学习
// ============================================================================
// 本章内容：
// 1. var 声明 vs := 短声明
// 2. 零值 (Zero Values)
// 3. const 常量与 iota
// 4. 作用域与变量遮蔽 (Scope & Shadowing)
// ============================================================================

package main

import "fmt"

// ============================================================================
// 📚 概念：Go 的变量声明
// ============================================================================
// Go 是静态类型语言，变量一旦声明类型就不能改变。
// 有两种主要声明方式：
//   var name type = value   （显式声明）
//   name := value           （短声明，自动推断类型）
// ============================================================================

// 包级变量 - 只能使用 var，不能使用 :=
var globalVar = "我是包级变量"

// 批量声明
var (
	packageName    = "go-roadmap"
	packageVersion = "1.0.0"
	isStable       = true
)

// ============================================================================
// 常量 - 编译时确定，不能修改
// ============================================================================
const Pi = 3.14159265358979

const (
	StatusOK       = 200
	StatusNotFound = 404
	StatusError    = 500
)

// ============================================================================
// iota - Go 的常量计数器，在 const 块中自动递增
// ============================================================================
// 📚 概念：iota
// iota 在每个 const 块中从 0 开始，每行自动 +1
// 非常适合定义枚举类型
// ============================================================================

// 基础枚举
const (
	Sunday    = iota // 0
	Monday           // 1
	Tuesday          // 2
	Wednesday        // 3
	Thursday         // 4
	Friday           // 5
	Saturday         // 6
)

// 使用 iota 做位运算（权限系统常用）
const (
	ReadPermission    = 1 << iota // 1  (001)
	WritePermission               // 2  (010)
	ExecutePermission             // 4  (100)
)

// 跳过某些值
const (
	_  = iota             // 跳过 0
	KB = 1 << (10 * iota) // 1 << 10 = 1024
	MB                    // 1 << 20
	GB                    // 1 << 30
	TB                    // 1 << 40
)

func main() {
	// ========================================================================
	// 1. var 声明
	// ========================================================================
	fmt.Println("=== 1. var 声明 ===")

	// 声明并赋值
	var age int = 25
	var name string = "Alice"
	fmt.Printf("name=%s, age=%d\n", name, age)

	// 类型推断（省略类型）
	var score = 95.5  // 自动推断为 float64
	var active = true // 自动推断为 bool
	fmt.Printf("score=%v(%T), active=%v(%T)\n", score, score, active, active)

	// 只声明不赋值（使用零值）
	var count int
	var message string
	var flag bool
	fmt.Printf("零值: count=%d, message='%s', flag=%t\n", count, message, flag)

	// ========================================================================
	// 2. := 短声明（只能在函数内部使用）
	// ========================================================================
	fmt.Println("\n=== 2. := 短声明 ===")

	// 短声明 - 自动推断类型
	city := "Beijing"
	temperature := 28.5
	isRaining := false
	fmt.Printf("city=%s, temp=%.1f, rain=%t\n", city, temperature, isRaining)

	// 多变量短声明
	x, y, z := 1, 2.0, "three"
	fmt.Printf("x=%v(%T), y=%v(%T), z=%v(%T)\n", x, x, y, y, z, z)

	// ⚠️ 注意：:= 不能用于已存在的变量（除非至少有一个新变量）
	// x := 10  // ❌ 编译错误！x 已经声明过
	x, w := 10, 20 // ✅ 因为 w 是新变量
	fmt.Printf("x=%d, w=%d\n", x, w)

	// ========================================================================
	// 3. 零值 (Zero Values) - Go 的重要特性
	// ========================================================================
	fmt.Println("\n=== 3. 零值 Zero Values ===")

	// 📚 概念：零值
	// Go 中所有变量声明后都有默认零值，不存在"未初始化"的变量
	// 这是 Go 安全性的重要保证

	var zeroInt int               // 0
	var zeroFloat float64         // 0.0
	var zeroBool bool             // false
	var zeroString string         // ""（空字符串）
	var zeroPointer *int          // nil
	var zeroSlice []int           // nil
	var zeroMap map[string]int    // nil
	var zeroFunc func()           // nil
	var zeroInterface interface{} // nil

	fmt.Printf("int:       %v\n", zeroInt)
	fmt.Printf("float64:   %v\n", zeroFloat)
	fmt.Printf("bool:      %v\n", zeroBool)
	fmt.Printf("string:    '%v' (空字符串)\n", zeroString)
	fmt.Printf("pointer:   %v\n", zeroPointer)
	fmt.Printf("slice:     %v (nil=%t)\n", zeroSlice, zeroSlice == nil)
	fmt.Printf("map:       %v (nil=%t)\n", zeroMap, zeroMap == nil)
	fmt.Printf("func:      %v\n", zeroFunc == nil)
	fmt.Printf("interface: %v\n", zeroInterface)

	// ========================================================================
	// 4. 常量
	// ========================================================================
	fmt.Println("\n=== 4. 常量 const ===")

	fmt.Printf("Pi = %.10f\n", Pi)
	fmt.Printf("HTTP Status: OK=%d, NotFound=%d, Error=%d\n",
		StatusOK, StatusNotFound, StatusError)

	// 常量不能修改
	// Pi = 3.0  // ❌ 编译错误！

	// 无类型常量 - Go 的特殊特性
	const untypedInt = 42         // 无类型整数常量
	const untypedFloat = 3.14     // 无类型浮点常量
	const untypedString = "hello" // 无类型字符串常量

	// 无类型常量可以赋给兼容的类型
	var i int = untypedInt
	var f float64 = untypedInt // 42 可以作为 float64
	var f2 float32 = untypedFloat
	fmt.Printf("i=%d, f=%f, f2=%f\n", i, f, f2)

	// ========================================================================
	// 5. iota 示例
	// ========================================================================
	fmt.Println("\n=== 5. iota 枚举 ===")

	fmt.Printf("星期: Sun=%d, Mon=%d, Fri=%d, Sat=%d\n",
		Sunday, Monday, Friday, Saturday)

	// 权限位运算
	fmt.Printf("\n权限位: Read=%d(%b), Write=%d(%b), Exec=%d(%b)\n",
		ReadPermission, ReadPermission,
		WritePermission, WritePermission,
		ExecutePermission, ExecutePermission)

	// 组合权限
	myPermission := ReadPermission | WritePermission
	fmt.Printf("我的权限: %d (%b)\n", myPermission, myPermission)
	fmt.Printf("有读权限? %t\n", myPermission&ReadPermission != 0)
	fmt.Printf("有写权限? %t\n", myPermission&WritePermission != 0)
	fmt.Printf("有执行权限? %t\n", myPermission&ExecutePermission != 0)

	// 存储单位
	fmt.Printf("\n存储单位: KB=%d, MB=%d, GB=%d\n", KB, MB, GB)

	// ========================================================================
	// 6. 作用域与变量遮蔽 (Shadowing)
	// ========================================================================
	fmt.Println("\n=== 6. 作用域与 Shadowing ===")

	// 📚 概念：作用域
	// Go 使用词法作用域（lexical scoping）
	// 内层作用域可以访问外层变量，也可以声明同名变量"遮蔽"外层变量

	outer := "外层"
	fmt.Printf("外层: outer = %s\n", outer)

	{
		// 新的作用域块
		outer := "内层" // 遮蔽了外层的 outer
		inner := "只在内层可见"
		fmt.Printf("内层: outer = %s\n", outer)
		fmt.Printf("内层: inner = %s\n", inner)
	}

	// 回到外层，outer 还是原来的值
	fmt.Printf("外层: outer = %s (没有被修改)\n", outer)
	// fmt.Println(inner) // ❌ 编译错误！inner 不在此作用域

	// ⚠️ Shadowing 的常见陷阱
	fmt.Println("\n⚠️  Shadowing 常见陷阱:")
	err := doSomething()
	fmt.Printf("外层 err = %v\n", err)

	if true {
		// 这里的 err 是新变量，不会修改外层的 err！
		err := doSomethingElse()
		fmt.Printf("内层 err = %v (新变量)\n", err)
	}
	fmt.Printf("外层 err 没变 = %v\n", err)

	// 正确做法：使用 = 而不是 :=
	if true {
		err = doSomethingElse() // 用 = 修改外层的 err
	}
	fmt.Printf("外层 err 被修改 = %v\n", err)
}

func doSomething() error {
	return nil
}

func doSomethingElse() error {
	return fmt.Errorf("something else error")
}

// ============================================================================
// 💻 运行方式：
//   cd 02-variables-constants && go run main.go
//
// 📝 练习：
// 1. 使用 iota 定义一个季节枚举（Spring, Summer, Autumn, Winter）
// 2. 创建一个权限系统：Admin/Editor/Viewer 三种角色
// 3. 故意制造一个 shadowing 陷阱，然后修复它
// ============================================================================
