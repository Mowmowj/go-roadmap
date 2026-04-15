// ============================================================================
// 第06章：函数 🟢 初级 | ⭐ 必须学习
// ============================================================================
// 本章内容：
// 1. 函数基础
// 2. 可变参数函数 (Variadic)
// 3. 多返回值
// 4. 匿名函数与闭包
// 5. 命名返回值
// 6. 值传递 (Call by Value)
// ============================================================================

package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// ============================================================================
// 📚 概念：Go 的函数
// ============================================================================
// Go 中函数是一等公民（first-class citizens），可以：
// - 赋值给变量
// - 作为参数传递
// - 作为返回值
// - 匿名定义
// ============================================================================

// ========================================================================
// 1. 函数基础
// ========================================================================

// 无参数无返回值
func sayHello() {
	fmt.Println("Hello!")
}

// 有参数有返回值
func add(a int, b int) int {
	return a + b
}

// 相同类型参数可以合并声明
func multiply(a, b float64) float64 {
	return a * b
}

// ========================================================================
// 2. 多返回值
// ========================================================================

// 返回两个值
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("除数不能为零")
	}
	return a / b, nil
}

// 返回多个值（交换）
func swap(a, b string) (string, string) {
	return b, a
}

// ========================================================================
// 3. 命名返回值 (Named Return Values)
// ========================================================================
// 📚 概念：命名返回值
// 返回值可以命名，相当于在函数顶部声明了变量
// 可以使用 "裸 return"（bare return），但不推荐在长函数中使用

func rectProperties(width, height float64) (area, perimeter float64) {
	area = width * height
	perimeter = 2 * (width + height)
	return // 裸 return，返回 area 和 perimeter
}

// 命名返回值在错误处理中很有用
func parseConfig(data string) (config map[string]string, err error) {
	config = make(map[string]string)
	if data == "" {
		err = fmt.Errorf("空数据")
		return
	}
	for _, line := range strings.Split(data, "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return
}

// ========================================================================
// 4. 可变参数函数 (Variadic Functions)
// ========================================================================

// ... 表示可变参数，在函数内部是一个切片
func sum(numbers ...int) int {
	total := 0
	for _, n := range numbers {
		total += n
	}
	return total
}

// 可变参数与固定参数混用（可变参数必须在最后）
func printf(format string, args ...interface{}) {
	fmt.Printf("[LOG] "+format+"\n", args...)
}

// ========================================================================
// 5. 值传递 (Call by Value)
// ========================================================================

// 📚 概念：Go 所有函数参数都是值传递
// 传入的是值的副本，修改不影响原值
// 如果想修改原值，传指针

func modifyValue(n int) {
	n = 100 // 只修改了副本
}

func modifySlice(s []int) {
	// 切片是引用类型，传递的是切片头的副本
	// 但底层数组共享，所以修改元素会影响原切片
	s[0] = 999
	// 但 append 可能不影响原切片（如果触发扩容）
}

func modifyPointer(p *int) {
	*p = 100 // 通过指针修改原值
}

func main() {
	// ========================================================================
	// 1. 函数基础
	// ========================================================================
	fmt.Println("=== 1. 函数基础 ===")
	sayHello()
	fmt.Printf("add(3, 5) = %d\n", add(3, 5))
	fmt.Printf("multiply(3.5, 2.0) = %.1f\n", multiply(3.5, 2.0))

	// ========================================================================
	// 2. 多返回值
	// ========================================================================
	fmt.Println("\n=== 2. 多返回值 ===")

	result, err := divide(10, 3)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("10 / 3 = %.4f\n", result)
	}

	// 处理除零
	result, err = divide(10, 0)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	}

	// 忽略不需要的返回值
	a, b := swap("hello", "world")
	fmt.Printf("swap: %s, %s\n", a, b)

	_, onlySecond := swap("first", "second")
	fmt.Printf("只取第二个: %s\n", onlySecond)

	// ========================================================================
	// 3. 命名返回值
	// ========================================================================
	fmt.Println("\n=== 3. 命名返回值 ===")

	area, perimeter := rectProperties(5, 3)
	fmt.Printf("面积=%.1f, 周长=%.1f\n", area, perimeter)

	config, err := parseConfig("host = localhost\nport = 8080\nmode = debug")
	if err == nil {
		fmt.Printf("配置: %v\n", config)
	}

	// ========================================================================
	// 4. 可变参数
	// ========================================================================
	fmt.Println("\n=== 4. 可变参数 ===")

	fmt.Printf("sum() = %d\n", sum())
	fmt.Printf("sum(1) = %d\n", sum(1))
	fmt.Printf("sum(1,2,3) = %d\n", sum(1, 2, 3))
	fmt.Printf("sum(1,2,3,4,5) = %d\n", sum(1, 2, 3, 4, 5))

	// 传入切片：用 ... 展开
	nums := []int{10, 20, 30}
	fmt.Printf("sum(slice...) = %d\n", sum(nums...))

	printf("用户 %s 登录，IP: %s", "Alice", "192.168.1.1")

	// ========================================================================
	// 5. 值传递
	// ========================================================================
	fmt.Println("\n=== 5. 值传递 ===")

	n := 42
	modifyValue(n)
	fmt.Printf("modifyValue 后: n = %d (没变)\n", n)

	modifyPointer(&n)
	fmt.Printf("modifyPointer 后: n = %d (被修改)\n", n)

	slice := []int{1, 2, 3}
	modifySlice(slice)
	fmt.Printf("modifySlice 后: %v (元素被修改)\n", slice)

	// ========================================================================
	// 6. 匿名函数
	// ========================================================================
	fmt.Println("\n=== 6. 匿名函数 ===")

	// 赋值给变量
	greet := func(name string) string {
		return "Hello, " + name + "!"
	}
	fmt.Println(greet("Go"))

	// 立即执行
	result2 := func(x, y int) int {
		return x + y
	}(3, 4) // 立即调用
	fmt.Printf("立即执行: %d\n", result2)

	// 函数作为参数（回调）
	fmt.Println("\n函数作为参数:")
	numbers := []int{1, 2, 3, 4, 5}
	doubled := mapInts(numbers, func(n int) int {
		return n * 2
	})
	fmt.Printf("原始: %v\n", numbers)
	fmt.Printf("翻倍: %v\n", doubled)

	filtered := filterInts(numbers, func(n int) bool {
		return n%2 == 0
	})
	fmt.Printf("偶数: %v\n", filtered)

	// sort 包使用函数参数
	people := []string{"Charlie", "Alice", "Bob"}
	sort.Slice(people, func(i, j int) bool {
		return people[i] < people[j]
	})
	fmt.Printf("排序后: %v\n", people)

	// ========================================================================
	// 7. 闭包 (Closures)
	// ========================================================================
	fmt.Println("\n=== 7. 闭包 ===")

	// 📚 概念：闭包
	// 闭包 = 函数 + 它引用的外部变量
	// 闭包可以"记住"创建时的环境

	// 计数器
	counter := makeCounter()
	fmt.Printf("counter: %d\n", counter()) // 1
	fmt.Printf("counter: %d\n", counter()) // 2
	fmt.Printf("counter: %d\n", counter()) // 3

	// 每个闭包有自己的环境
	counter2 := makeCounter()
	fmt.Printf("counter2: %d (独立计数)\n", counter2()) // 1

	// 斐波那契数列生成器
	fib := fibonacci()
	fmt.Print("Fibonacci: ")
	for i := 0; i < 10; i++ {
		fmt.Printf("%d ", fib())
	}
	fmt.Println()

	// 闭包陷阱
	fmt.Println("\n⚠️  闭包陷阱:")
	closureTrap()

	// ========================================================================
	// 8. 函数作为返回值（高阶函数）
	// ========================================================================
	fmt.Println("\n=== 8. 高阶函数 ===")

	// 返回函数
	adder := makeAdder(10)
	fmt.Printf("adder(5) = %d\n", adder(5))   // 15
	fmt.Printf("adder(20) = %d\n", adder(20)) // 30

	// 中间件模式
	fmt.Println("\n中间件模式:")
	hello := func(name string) string {
		return "Hello, " + name
	}

	// 包装函数
	logged := withLogging(hello)
	fmt.Println(logged("World"))

	// ========================================================================
	// 9. 递归
	// ========================================================================
	fmt.Println("\n=== 9. 递归 ===")

	fmt.Printf("factorial(5) = %d\n", factorial(5))
	fmt.Printf("fibonacci_r(10) = %d\n", fibonacciRecursive(10))

	// 计算目录大小的递归模拟
	fmt.Println("\n目录遍历模拟:")
	root := File{"root", 0, []File{
		{Name: "a.txt", Size: 100},
		{Name: "dir1", Size: 0, Children: []File{
			{Name: "b.txt", Size: 200},
			{Name: "c.txt", Size: 300},
		}},
	}}
	fmt.Printf("总大小: %d\n", totalSize(root))

	// ========================================================================
	// 10. 实用函数模式
	// ========================================================================
	fmt.Println("\n=== 10. 实用模式 ===")

	// Options 模式（函数选项模式）
	srv := NewServer(
		WithPort(8080),
		WithHost("localhost"),
		WithTimeout(30),
	)
	fmt.Printf("Server: %+v\n", srv)
}

// ============================================================================
// 辅助函数
// ============================================================================

// 高阶函数：对切片中每个元素应用函数
func mapInts(nums []int, f func(int) int) []int {
	result := make([]int, len(nums))
	for i, n := range nums {
		result[i] = f(n)
	}
	return result
}

// 高阶函数：过滤切片
func filterInts(nums []int, predicate func(int) bool) []int {
	var result []int
	for _, n := range nums {
		if predicate(n) {
			result = append(result, n)
		}
	}
	return result
}

// 闭包：计数器工厂
func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// 闭包：加法器工厂
func makeAdder(x int) func(int) int {
	return func(y int) int {
		return x + y
	}
}

// 闭包：斐波那契生成器
func fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		result := a
		a, b = b, a+b
		return result
	}
}

// ⚠️ 闭包陷阱：循环变量捕获
func closureTrap() {
	// ❌ 错误写法：所有闭包捕获同一个变量 i
	funcs := make([]func(), 5)
	for i := 0; i < 5; i++ {
		funcs[i] = func() {
			fmt.Printf("%d ", i) // i 在循环结束后是 5
		}
	}
	fmt.Print("陷阱(全是5): ")
	for _, f := range funcs {
		f()
	}
	fmt.Println()

	// ✅ 正确写法1：通过参数传递
	for i := 0; i < 5; i++ {
		i := i // 创建新变量遮蔽循环变量
		funcs[i] = func() {
			fmt.Printf("%d ", i)
		}
	}
	fmt.Print("正确: ")
	for _, f := range funcs {
		f()
	}
	fmt.Println()
}

// 中间件：给函数添加日志
func withLogging(f func(string) string) func(string) string {
	return func(s string) string {
		fmt.Printf("  [LOG] 调用函数，参数: %s\n", s)
		result := f(s)
		fmt.Printf("  [LOG] 返回结果: %s\n", result)
		return result
	}
}

// 递归：阶乘
func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

// 递归：斐波那契（效率低，仅示例）
func fibonacciRecursive(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacciRecursive(n-1) + fibonacciRecursive(n-2)
}

// 递归：文件大小
type File struct {
	Name     string
	Size     int
	Children []File
}

func totalSize(f File) int {
	size := f.Size
	for _, child := range f.Children {
		size += totalSize(child)
	}
	return size
}

// ============================================================================
// 函数选项模式 (Functional Options Pattern)
// ============================================================================

type Server struct {
	Host    string
	Port    int
	Timeout int
}

type ServerOption func(*Server)

func WithHost(host string) ServerOption {
	return func(s *Server) { s.Host = host }
}

func WithPort(port int) ServerOption {
	return func(s *Server) { s.Port = port }
}

func WithTimeout(timeout int) ServerOption {
	return func(s *Server) { s.Timeout = timeout }
}

func NewServer(opts ...ServerOption) Server {
	// 默认值
	s := Server{
		Host:    "0.0.0.0",
		Port:    3000,
		Timeout: 60,
	}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

// 让编译器高兴（使用了 math 包）
var _ = math.Pi

// ============================================================================
// 💻 运行方式：
//   cd 06-functions && go run main.go
//
// 📝 练习：
// 1. 实现一个 reduce 函数（类似 JavaScript 的 reduce）
// 2. 用闭包实现一个简单的缓存（memoize）
// 3. 实现 compose 函数：将多个函数组合成一个
// ============================================================================
