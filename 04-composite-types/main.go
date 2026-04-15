// ============================================================================
// 第04章：复合类型 🟢 初级 | ⭐ 必须学习
// ============================================================================
// 本章内容：
// 1. 数组 (Array)
// 2. 切片 (Slice) - 核心数据结构！
// 3. Map
// 4. 结构体 (Struct)
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"sort"
)

func main() {
	// ========================================================================
	// 1. 数组 (Array) - 固定长度
	// ========================================================================
	fmt.Println("=== 1. 数组 Array ===")

	// 📚 概念：数组
	// - 固定长度，声明后不能改变
	// - 类型的一部分：[3]int 和 [5]int 是不同类型
	// - 值类型：赋值和传参都是复制
	// - 实际开发中很少直接使用，通常用切片代替

	// 声明方式
	var arr1 [5]int // 零值初始化
	arr2 := [3]string{"Go", "Java", "Python"}
	arr3 := [...]int{10, 20, 30, 40} // ... 让编译器计算长度

	fmt.Printf("arr1: %v (长度=%d)\n", arr1, len(arr1))
	fmt.Printf("arr2: %v (长度=%d)\n", arr2, len(arr2))
	fmt.Printf("arr3: %v (长度=%d)\n", arr3, len(arr3))

	// 访问和修改
	arr1[0] = 100
	arr1[4] = 500
	fmt.Printf("修改后 arr1: %v\n", arr1)

	// 数组是值类型！赋值会复制
	arr4 := arr3
	arr4[0] = 999
	fmt.Printf("原数组: %v (没变)\n", arr3)
	fmt.Printf("副本:   %v (改了)\n", arr4)

	// 多维数组
	matrix := [2][3]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	fmt.Printf("矩阵: %v\n", matrix)

	// ========================================================================
	// 2. 切片 (Slice) - Go 最重要的数据结构之一
	// ========================================================================
	fmt.Println("\n=== 2. 切片 Slice ===")

	// 📚 概念：切片
	// 切片是对底层数组的动态窗口，包含三要素：
	// - 指针 (pointer): 指向底层数组的某个元素
	// - 长度 (length):  切片中的元素个数 len()
	// - 容量 (capacity): 从指针开始到底层数组末尾的元素个数 cap()
	//
	// 切片是引用类型：赋值和传参共享底层数组

	// 创建方式
	// 方式1: 字面量
	s1 := []int{1, 2, 3, 4, 5}
	fmt.Printf("字面量: %v, len=%d, cap=%d\n", s1, len(s1), cap(s1))

	// 方式2: make([]T, length, capacity)
	s2 := make([]int, 3, 10)
	fmt.Printf("make:   %v, len=%d, cap=%d\n", s2, len(s2), cap(s2))

	// 方式3: 从数组/切片截取
	arr := [5]int{10, 20, 30, 40, 50}
	s3 := arr[1:4] // [20, 30, 40]，不包含索引4
	fmt.Printf("截取:   %v, len=%d, cap=%d\n", s3, len(s3), cap(s3))

	// 切片截取语法
	fmt.Println("\n切片截取:")
	data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("data:       %v\n", data)
	fmt.Printf("data[2:5]:  %v (索引2到4)\n", data[2:5])
	fmt.Printf("data[:3]:   %v (前3个)\n", data[:3])
	fmt.Printf("data[7:]:   %v (索引7到末尾)\n", data[7:])
	fmt.Printf("data[:]:    %v (全部)\n", data[:])

	// ⚠️ 切片是引用类型！
	fmt.Println("\n⚠️  切片引用特性:")
	original := []int{1, 2, 3}
	shared := original // 共享底层数组
	shared[0] = 999
	fmt.Printf("original: %v (被修改了!)\n", original)
	fmt.Printf("shared:   %v\n", shared)

	// append - 追加元素
	fmt.Println("\nappend 追加:")
	s := []int{1, 2, 3}
	s = append(s, 4)              // 追加一个
	s = append(s, 5, 6, 7)        // 追加多个
	s = append(s, []int{8, 9}...) // 追加另一个切片
	fmt.Printf("append 结果: %v, len=%d, cap=%d\n", s, len(s), cap(s))

	// 📚 概念：切片扩容 (Capacity Growth)
	// 当 append 导致 len > cap 时，Go 会：
	// 1. 分配新的更大的底层数组
	// 2. 复制旧数据到新数组
	// 3. 返回新切片（所以 append 的结果必须重新赋值！）
	fmt.Println("\n切片扩容演示:")
	growth := make([]int, 0)
	for i := 0; i < 10; i++ {
		growth = append(growth, i)
		fmt.Printf("  len=%d, cap=%d\n", len(growth), cap(growth))
	}

	// copy - 复制切片（创建独立副本）
	fmt.Println("\ncopy 复制:")
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, len(src))
	copied := copy(dst, src)
	dst[0] = 999
	fmt.Printf("src: %v (没变)\n", src)
	fmt.Printf("dst: %v (独立副本)\n", dst)
	fmt.Printf("复制了 %d 个元素\n", copied)

	// 删除元素（Go 没有内置 delete 切片元素的方法）
	fmt.Println("\n删除切片元素:")
	items := []string{"a", "b", "c", "d", "e"}
	// 删除索引2的元素 ("c")
	items = append(items[:2], items[3:]...)
	fmt.Printf("删除后: %v\n", items)

	// nil 切片 vs 空切片
	fmt.Println("\nnil vs 空切片:")
	var nilSlice []int          // nil 切片
	emptySlice := []int{}       // 空切片
	makeSlice := make([]int, 0) // 空切片
	fmt.Printf("nil:   %v, len=%d, nil=%t\n", nilSlice, len(nilSlice), nilSlice == nil)
	fmt.Printf("empty: %v, len=%d, nil=%t\n", emptySlice, len(emptySlice), emptySlice == nil)
	fmt.Printf("make:  %v, len=%d, nil=%t\n", makeSlice, len(makeSlice), makeSlice == nil)
	// 都可以正常 append
	nilSlice = append(nilSlice, 1)
	fmt.Printf("nil append: %v\n", nilSlice)

	// ========================================================================
	// 3. Map
	// ========================================================================
	fmt.Println("\n=== 3. Map ===")

	// 📚 概念：Map
	// - 键值对集合（哈希表）
	// - 引用类型
	// - 键必须是可比较类型（不能是 slice、map、func）
	// - 无序遍历

	// 创建方式
	// 方式1: 字面量
	scores := map[string]int{
		"Alice": 95,
		"Bob":   87,
		"Carol": 92,
	}
	fmt.Printf("scores: %v\n", scores)

	// 方式2: make
	ages := make(map[string]int)
	ages["Alice"] = 30
	ages["Bob"] = 25
	fmt.Printf("ages: %v\n", ages)

	// 访问
	fmt.Println("\n访问 Map:")
	fmt.Printf("Alice 的分数: %d\n", scores["Alice"])

	// 📚 概念：Comma-Ok 惯用法
	// 区分"键不存在返回零值"和"键存在但值为零值"
	score, ok := scores["Dave"] // Dave 不存在
	fmt.Printf("Dave: score=%d, exists=%t\n", score, ok)

	score, ok = scores["Alice"]
	fmt.Printf("Alice: score=%d, exists=%t\n", score, ok)

	// 惯用模式
	if val, ok := scores["Bob"]; ok {
		fmt.Printf("Bob 存在，分数: %d\n", val)
	} else {
		fmt.Println("Bob 不存在")
	}

	// 删除
	delete(scores, "Bob")
	fmt.Printf("删除 Bob 后: %v\n", scores)

	// 遍历（无序！）
	fmt.Println("\n遍历 Map（注意：每次顺序可能不同）:")
	for key, value := range scores {
		fmt.Printf("  %s: %d\n", key, value)
	}

	// 有序遍历：先排序键
	fmt.Println("\n有序遍历 Map:")
	keys := make([]string, 0, len(scores))
	for k := range scores {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("  %s: %d\n", k, scores[k])
	}

	// Map 嵌套
	students := map[string]map[string]int{
		"Alice": {"Math": 95, "English": 88},
		"Bob":   {"Math": 78, "English": 92},
	}
	fmt.Printf("\n嵌套 Map: %v\n", students)
	fmt.Printf("Alice 的数学: %d\n", students["Alice"]["Math"])

	// ========================================================================
	// 4. 结构体 (Struct)
	// ========================================================================
	fmt.Println("\n=== 4. 结构体 Struct ===")

	// 📚 概念：结构体
	// - Go 没有类 (class)，用结构体组织数据
	// - 值类型：赋值会复制
	// - 可以定义方法（后续章节详解）

	// 基本使用
	type Person struct {
		Name string
		Age  int
		City string
	}

	// 创建方式
	p1 := Person{Name: "Alice", Age: 30, City: "Beijing"}
	p2 := Person{"Bob", 25, "Shanghai"} // 按顺序（不推荐）
	var p3 Person                       // 零值
	fmt.Printf("p1: %+v\n", p1)
	fmt.Printf("p2: %+v\n", p2)
	fmt.Printf("p3: %+v (零值)\n", p3)

	// 访问和修改
	p1.Age = 31
	fmt.Printf("修改后 p1.Age = %d\n", p1.Age)

	// 指针创建
	p4 := &Person{Name: "Dave", Age: 35}
	fmt.Printf("p4: %+v\n", *p4)
	p4.Name = "David" // Go 自动解引用，不需要 (*p4).Name
	fmt.Printf("修改后 p4.Name = %s\n", p4.Name)

	// Struct Tags & JSON
	fmt.Println("\n--- Struct Tags & JSON ---")
	type User struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"-"`             // 不序列化
		Bio      string `json:"bio,omitempty"` // 空值时省略
	}

	user := User{
		ID:       1,
		Username: "gopher",
		Email:    "gopher@go.dev",
		Password: "secret123",
		Bio:      "",
	}

	// 序列化 (结构体 -> JSON)
	jsonData, _ := json.MarshalIndent(user, "", "  ")
	fmt.Printf("JSON 序列化:\n%s\n", string(jsonData))

	// 反序列化 (JSON -> 结构体)
	jsonStr := `{"id":2,"username":"rob","email":"rob@go.dev"}`
	var user2 User
	json.Unmarshal([]byte(jsonStr), &user2)
	fmt.Printf("JSON 反序列化: %+v\n", user2)

	// 结构体嵌入 (Embedding)
	fmt.Println("\n--- 结构体嵌入 Embedding ---")

	// 📚 概念：嵌入 (Embedding)
	// Go 使用嵌入实现组合（而非继承）
	// 嵌入类型的字段和方法被"提升"到外层结构体

	type Address struct {
		Street  string
		City    string
		Country string
	}

	type Employee struct {
		Person  // 嵌入 Person
		Address // 嵌入 Address
		Company string
		Salary  float64
	}

	emp := Employee{
		Person:  Person{Name: "Alice", Age: 30},
		Address: Address{Street: "123 Go St", City: "Beijing", Country: "China"},
		Company: "GoLang Inc",
		Salary:  15000,
	}

	// 可以直接访问嵌入字段（被提升了）
	fmt.Printf("姓名: %s\n", emp.Name)         // 等同于 emp.Person.Name
	fmt.Printf("城市: %s\n", emp.Address.City) // Person 和 Address 都有 City，需要明确指定
	fmt.Printf("公司: %s\n", emp.Company)
	fmt.Printf("完整: %+v\n", emp)

	// 匿名结构体（临时使用）
	fmt.Println("\n--- 匿名结构体 ---")
	point := struct {
		X, Y int
	}{10, 20}
	fmt.Printf("匿名结构体: %+v\n", point)

	// 结构体比较
	fmt.Println("\n--- 结构体比较 ---")
	pa := Person{Name: "Alice", Age: 30, City: "Beijing"}
	pb := Person{Name: "Alice", Age: 30, City: "Beijing"}
	pc := Person{Name: "Bob", Age: 25, City: "Shanghai"}
	fmt.Printf("pa == pb: %t (所有字段相等)\n", pa == pb)
	fmt.Printf("pa == pc: %t\n", pa == pc)
}

// ============================================================================
// 💻 运行方式：
//   cd 04-composite-types && go run main.go
//
// 📝 练习：
// 1. 实现一个 Stack（栈）数据结构，使用切片模拟
// 2. 用 Map 实现一个简单的词频统计
// 3. 定义一个 Book 结构体，包含 JSON 标签，实现序列化/反序列化
// 4. 使用结构体嵌入设计一个简单的动物层次（Animal -> Dog/Cat）
// ============================================================================
