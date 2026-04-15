// ============================================================================
// 第07章：指针 🟡 中级 | ⭐ 必须学习
// ============================================================================
// 本章内容：
// 1. 指针基础
// 2. 指针与结构体
// 3. 指针与 Map/Slice
// 4. 内存管理与垃圾回收概述
// ============================================================================

package main

import "fmt"

func main() {
	// ========================================================================
	// 1. 指针基础
	// ========================================================================
	fmt.Println("=== 1. 指针基础 ===")

	// 📚 概念：指针
	// 指针存储变量的内存地址
	// & 取地址运算符：获取变量的地址
	// * 解引用运算符：获取指针指向的值
	// Go 指针没有指针运算（不能 p++ 移动指针），比 C 安全

	x := 42
	p := &x // p 是指向 x 的指针，类型 *int
	fmt.Printf("x 的值:    %d\n", x)
	fmt.Printf("x 的地址:  %p\n", &x)
	fmt.Printf("p 的值:    %p (和 &x 相同)\n", p)
	fmt.Printf("*p 的值:   %d (解引用，和 x 相同)\n", *p)

	// 通过指针修改值
	*p = 100
	fmt.Printf("修改 *p 后: x = %d\n", x)

	// new() 函数 - 分配内存，返回指针，值为零值
	pi := new(int)    // 分配 int，值为 0
	ps := new(string) // 分配 string，值为 ""
	fmt.Printf("new(int): %p, 值=%d\n", pi, *pi)
	fmt.Printf("new(string): %p, 值='%s'\n", ps, *ps)

	// nil 指针
	var np *int // nil 指针
	fmt.Printf("nil 指针: %v\n", np)
	// fmt.Println(*np) // ❌ panic: nil pointer dereference

	// 安全使用指针：先检查 nil
	if np != nil {
		fmt.Println(*np)
	} else {
		fmt.Println("指针为 nil，跳过")
	}

	// ========================================================================
	// 2. 指针与函数
	// ========================================================================
	fmt.Println("\n=== 2. 指针与函数 ===")

	// 值传递 vs 指针传递
	val := 10
	incrementByValue(val)
	fmt.Printf("值传递后: %d (没变)\n", val)

	incrementByPointer(&val)
	fmt.Printf("指针传递后: %d (被修改)\n", val)

	// 什么时候用指针参数？
	// 1. 需要修改调用者的变量
	// 2. 结构体很大，避免复制开销
	// 3. 和其他语言接口对接（如 C）

	// ========================================================================
	// 3. 指针与结构体
	// ========================================================================
	fmt.Println("\n=== 3. 指针与结构体 ===")

	// 结构体指针
	u := &User{Name: "Alice", Age: 30, Email: "alice@go.dev"}

	// Go 自动解引用：u.Name 等同于 (*u).Name
	fmt.Printf("Name: %s\n", u.Name)
	fmt.Printf("Age:  %d\n", u.Age)

	// 通过指针修改结构体
	u.Age = 31
	fmt.Printf("修改后 Age: %d\n", u.Age)

	// 函数中修改结构体
	user := User{Name: "Bob", Age: 25}
	fmt.Printf("修改前: %+v\n", user)
	updateUser(&user, "Bobby", 26)
	fmt.Printf("修改后: %+v\n", user)

	// 方法中的指针接收者（预览，详见第08章）
	fmt.Println("\n指针接收者预览:")
	c := Counter{value: 0}
	c.Increment()
	c.Increment()
	c.Increment()
	fmt.Printf("计数器: %d\n", c.Value())

	// ========================================================================
	// 4. 指针与 Map/Slice
	// ========================================================================
	fmt.Println("\n=== 4. 指针与 Map/Slice ===")

	// 📚 概念：
	// Slice、Map、Channel 本身就是引用类型
	// 它们内部包含指针，传参时不需要额外取地址
	// 但注意：append 可能会改变 slice 的底层数组

	// Map 是引用类型
	m := map[string]int{"a": 1}
	modifyMap(m)
	fmt.Printf("Map 被修改: %v\n", m) // {a: 1, b: 2}

	// Slice 引用特性
	s := []int{1, 2, 3}
	modifySliceElements(s)
	fmt.Printf("Slice 元素被修改: %v\n", s) // [100, 2, 3]

	// ⚠️ 但 append 在新函数中可能不生效
	appendToSlice(s)
	fmt.Printf("Slice 长度没变: %v (len=%d)\n", s, len(s))

	// 修复：传指针
	appendToSliceFixed(&s)
	fmt.Printf("Slice 被追加: %v (len=%d)\n", s, len(s))

	// ========================================================================
	// 5. 指针切片和指针Map
	// ========================================================================
	fmt.Println("\n=== 5. 结构体指针切片 ===")

	users := []*User{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Carol", Age: 35},
	}

	// 修改切片中的结构体
	for _, u := range users {
		u.Age += 1 // 直接修改，因为是指针
	}

	for _, u := range users {
		fmt.Printf("  %s: %d\n", u.Name, u.Age)
	}

	// ========================================================================
	// 6. 内存管理与垃圾回收
	// ========================================================================
	fmt.Println("\n=== 6. 内存管理概述 ===")

	// 📚 概念：Go 的内存管理
	//
	// 栈 (Stack)：
	// - 存储局部变量、函数参数
	// - 自动分配和回收，速度快
	// - 每个 goroutine 有自己的栈
	//
	// 堆 (Heap)：
	// - 存储需要在函数外存活的变量
	// - 由垃圾回收器 (GC) 管理
	// - 比栈慢
	//
	// 逃逸分析 (Escape Analysis)：
	// - 编译器自动决定变量在栈还是堆上分配
	// - 如果变量的引用"逃逸"出函数，就分配在堆上
	// - 用 go build -gcflags="-m" 查看逃逸分析结果

	// 变量逃逸到堆的例子
	p2 := createInt(42) // 返回局部变量的指针 → 逃逸到堆
	fmt.Printf("逃逸到堆: %d\n", *p2)

	// 不逃逸的例子
	localVar := 42
	useLocally(&localVar) // 指针没有逃逸
	fmt.Printf("栈上变量: %d\n", localVar)

	fmt.Println(`📚 垃圾回收 (GC) 要点：
1. Go 使用并发、三色标记清除 GC
2. GC 是自动的，不需要手动管理内存
3. 通过减少堆分配来优化性能：
   - 尽量使用值而不是指针
   - 预分配正确大小的 slice 和 map
   - 使用 sync.Pool 复用临时对象
   - 用 go tool pprof 分析内存使用
4. runtime.GC() 可以手动触发 GC（通常不需要）`)
}

// ============================================================================
// 辅助函数
// ============================================================================

func incrementByValue(n int) {
	n++ // 只修改副本
}

func incrementByPointer(n *int) {
	*n++ // 修改原值
}

func updateUser(u *User, name string, age int) {
	u.Name = name
	u.Age = age
}

type User struct {
	Name  string
	Age   int
	Email string
}

// 方法：指针接收者
type Counter struct {
	value int
}

func (c *Counter) Increment() {
	c.value++
}

func (c Counter) Value() int {
	return c.value
}

func modifyMap(m map[string]int) {
	m["b"] = 2
}

func modifySliceElements(s []int) {
	s[0] = 100
}

func appendToSlice(s []int) {
	s = append(s, 99) // 这个 s 是副本，不影响原切片
}

func appendToSliceFixed(s *[]int) {
	*s = append(*s, 99) // 通过指针修改原切片
}

// 返回局部变量指针（会逃逸到堆）
func createInt(n int) *int {
	result := n
	return &result // result 逃逸到堆
}

func useLocally(p *int) {
	*p += 1
}

// ============================================================================
// 💻 运行方式：
//   cd 07-pointers && go run main.go
//
// 📝 练习：
// 1. 实现一个链表（用指针连接节点）
// 2. 实现一个 swap 函数，通过指针交换两个变量的值
// 3. 运行 go build -gcflags="-m" 查看逃逸分析
// ============================================================================
