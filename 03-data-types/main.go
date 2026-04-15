// ============================================================================
// 第03章：数据类型 🟢 初级 | ⭐ 必须学习
// ============================================================================
// 本章内容：
// 1. 布尔类型 (bool)
// 2. 数值类型 - 整数(有符号/无符号)、浮点数、复数
// 3. 字符串 - 原始字符串、解释字符串
// 4. Rune (Unicode 字符)
// 5. 类型转换
// ============================================================================

package main

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

func main() {
	// ========================================================================
	// 1. 布尔类型 bool
	// ========================================================================
	fmt.Println("=== 1. 布尔类型 bool ===")

	var isActive bool = true
	var isDeleted bool // 零值: false
	fmt.Printf("isActive=%t, isDeleted=%t\n", isActive, isDeleted)

	// 布尔运算
	a, b := true, false
	fmt.Printf("AND: %t && %t = %t\n", a, b, a && b)
	fmt.Printf("OR:  %t || %t = %t\n", a, b, a || b)
	fmt.Printf("NOT: !%t = %t\n", a, !a)

	// 比较运算返回 bool
	x, y := 10, 20
	fmt.Printf("%d > %d = %t\n", x, y, x > y)
	fmt.Printf("%d == %d = %t\n", x, y, x == y)
	fmt.Printf("%d != %d = %t\n", x, y, x != y)

	// ⚠️ Go 中不能将整数当布尔用
	// if 1 { } // ❌ 编译错误！

	// ========================================================================
	// 2. 整数类型
	// ========================================================================
	fmt.Println("\n=== 2. 整数类型 ===")

	// 📚 概念：有符号整数 (Signed)
	// int8:  -128 ~ 127
	// int16: -32768 ~ 32767
	// int32: -2^31 ~ 2^31-1
	// int64: -2^63 ~ 2^63-1
	// int:   平台相关（32位系统=int32，64位系统=int64）

	var i8 int8 = 127 // 最大值
	var i16 int16 = 32767
	var i32 int32 = 2147483647
	var i64 int64 = 9223372036854775807
	fmt.Printf("int8:  %d (max)\n", i8)
	fmt.Printf("int16: %d (max)\n", i16)
	fmt.Printf("int32: %d (max)\n", i32)
	fmt.Printf("int64: %d (max)\n", i64)

	// 📚 概念：无符号整数 (Unsigned)
	// uint8 (byte): 0 ~ 255
	// uint16: 0 ~ 65535
	// uint32: 0 ~ 4294967295
	// uint64: 0 ~ 18446744073709551615
	// uint: 不能表示负数，范围是 0 和正整数，因此它是“无符号整数”，不是严格数学意义上的“正整数”。

	var u8 uint8 = 255
	var u16 uint16 = 65535
	fmt.Printf("uint8: %d (max), uint16: %d (max)\n", u8, u16)

	// byte 是 uint8 的别名
	var myByte byte = 'A'
	fmt.Printf("byte: %d = '%c'\n", myByte, myByte)

	// 整数运算
	fmt.Println("\n整数运算:")
	fmt.Printf("10 + 3 = %d\n", 10+3)
	fmt.Printf("10 - 3 = %d\n", 10-3)
	fmt.Printf("10 * 3 = %d\n", 10*3)
	fmt.Printf("10 / 3 = %d (整数除法，截断小数)\n", 10/3)
	fmt.Printf("10 %% 3 = %d (取模/余数)\n", 10%3)

	// 位运算
	fmt.Println("\n位运算:")
	fmt.Printf("5 & 3  = %d  (AND)\n", 5&3)
	fmt.Printf("5 | 3  = %d  (OR)\n", 5|3)
	fmt.Printf("5 ^ 3  = %d  (XOR)\n", 5^3)
	fmt.Printf("5 << 1 = %d (左移)\n", 5<<1)
	fmt.Printf("5 >> 1 = %d  (右移)\n", 5>>1)

	// ========================================================================
	// 3. 浮点数
	// ========================================================================
	fmt.Println("\n=== 3. 浮点数 ===")

	// float32: ~7 位有效数字
	// float64: ~15 位有效数字（推荐使用）
	var f32 float32 = 3.14
	var f64 float64 = 3.141592653589793
	fmt.Printf("float32: %.10f (精度有限)\n", f32)
	fmt.Printf("float64: %.15f (精度更高)\n", f64)

	// ⚠️ 浮点数精度陷阱
	fmt.Println("\n⚠️  浮点数精度陷阱:")
	result := 0.1 + 0.2
	fmt.Printf("0.1 + 0.2 = %.20f (不是精确的 0.3!)\n", result)
	fmt.Printf("0.1 + 0.2 == 0.3? %t\n", result == 0.3)

	// 📚 注意：这里在 Go 中常常会输出 true，
	// 因为 0.1、0.2、0.3 是无类型常量，编译器会以较高精度计算。
	// 但在真实业务里，如果参与运算的是 float32/float64 变量，
	// 就更容易出现精度误差，因此不要过度依赖 == 直接比较。
	var a64 float64 = 0.1
	var b64 float64 = 0.2
	var c64f float64 = 0.3
	fmt.Printf("变量比较: a64+b64 == c64f ? %t\n", a64+b64 == c64f)

	// 正确比较浮点数的方式（工程实践）
	const epsilon = 1e-9 // 允许的微小误差范围
	fmt.Printf("差值 < epsilon? %t (正确比较方式)\n",
		math.Abs((a64+b64)-c64f) < epsilon)

	// 💰 最佳实践：金额不要用 float 直接计算
	// 金额建议用最小货币单位的整数表示，例如“分”。
	priceFen := 1999 // 19.99 元
	discountFen := 200
	finalFen := priceFen - discountFen
	fmt.Printf("金额计算(分): %d - %d = %d (即 %.2f 元)\n",
		priceFen, discountFen, finalFen, float64(finalFen)/100)

	// 特殊浮点值
	fmt.Println("\n特殊浮点值:")
	fmt.Printf("MaxFloat64: %e\n", math.MaxFloat64)
	fmt.Printf("+Inf: %f\n", math.Inf(1))
	fmt.Printf("-Inf: %f\n", math.Inf(-1))
	fmt.Printf("NaN: %f\n", math.NaN())
	fmt.Printf("NaN == NaN? %t (NaN 不等于自身!)\n", math.NaN() == math.NaN())

	// ========================================================================
	// 4. 复数（💡 可选学习，科学计算用）
	// ========================================================================
	fmt.Println("\n=== 4. 复数 (可选) ===")

	var c64 complex64 = 3 + 4i
	var c128 complex128 = complex(3.0, 4.0) // 等价写法
	fmt.Printf("c64:  %v\n", c64)
	fmt.Printf("c128: %v\n", c128)
	fmt.Printf("实部: %f, 虚部: %f\n", real(c128), imag(c128))
	fmt.Printf("模长: %f\n", math.Sqrt(real(c128)*real(c128)+imag(c128)*imag(c128)))

	// ========================================================================
	// 5. Rune (Unicode 字符)
	// ========================================================================
	fmt.Println("\n=== 5. Rune ===")

	// 📚 概念：Rune
	// rune 是 int32 的别名，表示一个 Unicode 码点
	// Go 的字符串是 UTF-8 编码的字节序列
	// 一个中文字符 = 1 个 rune = 3 个 byte

	var r1 rune = 'A' // ASCII 字符
	var r2 rune = '中' // 中文字符
	var r3 rune = '🚀' // Emoji
	fmt.Printf("'A': rune=%d, byte size=1\n", r1)
	fmt.Printf("'中': rune=%d, byte size=3\n", r2)
	fmt.Printf("'🚀': rune=%d, byte size=4\n", r3)

	// 字符串中的 rune vs byte
	s := "Hello, 世界!"
	fmt.Printf("\n字符串: %s\n", s)
	fmt.Printf("字节数 len(): %d\n", len(s)) // 字节长度
	fmt.Printf("字符数 RuneCountInString(): %d\n",
		utf8.RuneCountInString(s)) // rune 长度

	// 遍历 rune
	fmt.Println("\n遍历 rune:")
	for i, r := range s {
		fmt.Printf("  索引=%d, rune='%c', 码点=%d\n", i, r, r)
	}

	// ========================================================================
	// 6. 字符串
	// ========================================================================
	fmt.Println("\n=== 6. 字符串 ===")

	// 📚 概念：Go 字符串是不可变的字节序列
	// 两种字面量：
	// - 解释字符串 "..." （支持转义字符 \n \t 等）
	// - 原始字符串 `...` （所见即所得，不转义）

	// 解释字符串 (Interpreted String Literals)
	interpreted := "Hello\tWorld\n换行了"
	fmt.Printf("解释字符串: %s\n", interpreted)

	// 原始字符串 (Raw String Literals) - 非常适合正则、JSON、多行文本
	raw := `这是原始字符串
可以直接换行
\n 不会被转义
"双引号" 也不用转义`
	fmt.Printf("原始字符串:\n%s\n", raw)

	// 字符串操作
	fmt.Println("\n字符串常用操作:")
	str := "Hello, Go World"
	fmt.Printf("长度: %d\n", len(str))
	fmt.Printf("包含: %t\n", strings.Contains(str, "Go"))
	fmt.Printf("前缀: %t\n", strings.HasPrefix(str, "Hello"))
	fmt.Printf("后缀: %t\n", strings.HasSuffix(str, "World"))
	fmt.Printf("索引: %d\n", strings.Index(str, "Go"))
	fmt.Printf("大写: %s\n", strings.ToUpper(str))
	fmt.Printf("小写: %s\n", strings.ToLower(str))
	fmt.Printf("替换: %s\n", strings.Replace(str, "World", "语言", 1))
	fmt.Printf("分割: %v\n", strings.Split(str, ", "))
	fmt.Printf("连接: %s\n", strings.Join([]string{"a", "b", "c"}, "-"))
	fmt.Printf("去空格: '%s'\n", strings.TrimSpace("  hello  "))
	fmt.Printf("重复: %s\n", strings.Repeat("Go ", 3))

	// 字符串拼接
	fmt.Println("\n字符串拼接方式:")
	// 方式1: + 运算符（简单但大量拼接效率低）
	s1 := "Hello" + " " + "World"
	fmt.Printf("+ 运算符: %s\n", s1)

	// 方式2: fmt.Sprintf（格式化拼接）
	s2 := fmt.Sprintf("%s is %d years old", "Go", 15)
	fmt.Printf("Sprintf: %s\n", s2)

	// 方式3: strings.Builder（高效拼接，推荐用于循环中）
	var builder strings.Builder
	for i := 0; i < 5; i++ {
		builder.WriteString(fmt.Sprintf("item%d ", i))
	}
	fmt.Printf("Builder: %s\n", builder.String())

	// ========================================================================
	// 7. 类型转换
	// ========================================================================
	fmt.Println("\n=== 7. 类型转换 ===")

	// 📚 概念：Go 没有隐式类型转换，必须显式转换
	// 语法：T(v) 将值 v 转换为类型 T

	// 数值类型之间的转换
	var intVal int = 42
	var floatVal float64 = float64(intVal)
	var int32Val int32 = int32(intVal)
	fmt.Printf("int->float64: %d -> %f\n", intVal, floatVal)
	fmt.Printf("int->int32: %d -> %d\n", intVal, int32Val)

	// ⚠️ 浮点转整数会截断小数
	var pi float64 = 3.99
	var piInt int = int(pi) // 截断，不是四舍五入！
	fmt.Printf("float64->int: %f -> %d (截断!)\n", pi, piInt)

	// ⚠️ 大范围转小范围可能溢出
	var bigInt int64 = 256
	var smallInt int8 = int8(bigInt) // 溢出！256 mod 256 = 0
	fmt.Printf("int64->int8: %d -> %d (溢出!)\n", bigInt, smallInt)

	// 字符串与数值的转换需要用 strconv 包（详见后续章节）
	// string(65)  得到 "A" 不是 "65"！
	fmt.Printf("string(65) = '%s' (是字符'A'，不是\"65\"!)\n", string(rune(65)))

	// string 和 []byte 互转
	str2 := "Hello"
	bytes := []byte(str2) // string -> []byte
	str3 := string(bytes) // []byte -> string
	fmt.Printf("string->[]byte: %v\n", bytes)
	fmt.Printf("[]byte->string: %s\n", str3)

	// string 和 []rune 互转
	chinese := "你好世界"
	runes := []rune(chinese) // string -> []rune
	str4 := string(runes)    // []rune -> string
	fmt.Printf("string->[]rune: %v\n", runes)
	fmt.Printf("[]rune->string: %s\n", str4)
}

// ============================================================================
// 💻 运行方式：
//   cd 03-data-types && go run main.go
//
// 📝 练习：
// 1. 计算 int8 能表示的最大值和最小值
// 2. 写一个函数统计一个中文字符串的字符数（不是字节数）
// 3. 实现一个安全的类型转换函数，检查是否溢出
// ============================================================================
