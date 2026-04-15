// ============================================================================
// 第09章：泛型 🟡 中级 | ⭐ 必须学习 (Go 1.18+)
// ============================================================================
// 本章内容：
// 1. 为什么需要泛型
// 2. 泛型函数
// 3. 泛型类型（结构体、接口）
// 4. 类型约束 (constraints)
// 5. 类型推断
// 6. 实用泛型模式
// ============================================================================
// ⚠️ 需要 Go 1.18+
// ============================================================================

package main

import (
	"fmt"
	"strings"
)

// ============================================================================
// 1. 为什么需要泛型？
// ============================================================================
// 📚 概念：
// 在泛型之前，通用操作只能用 interface{} + 类型断言，不安全且需类型转换
// 泛型允许编写适用于多种类型的函数和数据结构，同时保持类型安全

// 没有泛型的写法（Go 1.18 之前）
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// ============================================================================
// 2. 泛型函数
// ============================================================================

// 📚 语法：func 函数名[T 约束](参数) 返回值
// T 是类型参数，约束限制 T 可以是哪些类型

// 自定义 Ordered 约束（Go 1.21+ 可用 cmp.Ordered）
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

func Max[T Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// 多个类型参数
func Map[T any, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func Filter[T any](slice []T, fn func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func Reduce[T any, U any](slice []T, initial U, fn func(U, T) U) U {
	result := initial
	for _, v := range slice {
		result = fn(result, v)
	}
	return result
}

// Contains — 检查切片是否包含某元素
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

// Keys — 获取 map 的所有键
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values — 获取 map 的所有值
func Values[K comparable, V any](m map[K]V) []V {
	vals := make([]V, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}

// ============================================================================
// 3. 泛型类型
// ============================================================================

// 泛型栈
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// 泛型 Pair
type Pair[T, U any] struct {
	First  T
	Second U
}

func NewPair[T, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{First: first, Second: second}
}

// 泛型 Result（类似 Rust 的 Result）
type Result[T any] struct {
	Value T
	Err   error
}

func Ok[T any](value T) Result[T] {
	return Result[T]{Value: value}
}

func Err[T any](err error) Result[T] {
	return Result[T]{Err: err}
}

func (r Result[T]) IsOk() bool {
	return r.Err == nil
}

func (r Result[T]) Unwrap() T {
	if r.Err != nil {
		panic(fmt.Sprintf("called Unwrap on error: %v", r.Err))
	}
	return r.Value
}

// ============================================================================
// 4. 类型约束
// ============================================================================

// 📚 概念：类型约束 (Type Constraint)
// 约束是用接口定义的，限制类型参数可以是哪些类型

// 自定义约束
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// ~ 表示底层类型匹配（包括自定义类型）
type MyInt int // ~int 可以匹配 MyInt

func Sum[T Number](numbers []T) T {
	var total T
	for _, n := range numbers {
		total += n
	}
	return total
}

// 带方法的约束
type Stringer interface {
	String() string
}

func JoinStrings[T Stringer](items []T, sep string) string {
	strs := make([]string, len(items))
	for i, item := range items {
		strs[i] = item.String()
	}
	return strings.Join(strs, sep)
}

// Comparable 内置约束 — 支持 == 和 !=
// any 内置约束 — 等同于 interface{}

// ============================================================================
// 5. 泛型接口
// ============================================================================

type Collection[T any] interface {
	Add(item T)
	Get(index int) (T, bool)
	Size() int
}

// 实现泛型接口
type ArrayList[T any] struct {
	items []T
}

func NewArrayList[T any]() *ArrayList[T] {
	return &ArrayList[T]{}
}

func (a *ArrayList[T]) Add(item T) {
	a.items = append(a.items, item)
}

func (a *ArrayList[T]) Get(index int) (T, bool) {
	var zero T
	if index < 0 || index >= len(a.items) {
		return zero, false
	}
	return a.items[index], true
}

func (a *ArrayList[T]) Size() int {
	return len(a.items)
}

// ============================================================================
func main() {
	// 1. 为什么需要泛型
	fmt.Println("=== 1. 泛型之前 vs 之后 ===")
	fmt.Println("旧方式: maxInt(3, 5) =", maxInt(3, 5))
	fmt.Println("旧方式: maxFloat64(3.14, 2.71) =", maxFloat64(3.14, 2.71))
	fmt.Println("泛型:   Max(3, 5) =", Max(3, 5))
	fmt.Println("泛型:   Max(3.14, 2.71) =", Max(3.14, 2.71))
	fmt.Println("泛型:   Max(\"hello\", \"world\") =", Max("hello", "world"))

	// 2. 函数式操作
	fmt.Println("\n=== 2. Map/Filter/Reduce ===")
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Map: 平方
	squares := Map(nums, func(n int) int { return n * n })
	fmt.Println("Map (平方):", squares)

	// Filter: 偶数
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println("Filter (偶数):", evens)

	// Reduce: 求和
	sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
	fmt.Println("Reduce (求和):", sum)

	// Map: int -> string
	strs := Map(nums, func(n int) string { return fmt.Sprintf("#%d", n) })
	fmt.Println("Map (转字符串):", strs)

	// 3. Contains
	fmt.Println("\n=== 3. Contains ===")
	fmt.Println("Contains [1..10] 5:", Contains(nums, 5))
	fmt.Println("Contains [1..10] 11:", Contains(nums, 11))
	fmt.Println("Contains [\"a\",\"b\"] \"b\":", Contains([]string{"a", "b"}, "b"))

	// 4. Map 操作
	fmt.Println("\n=== 4. Keys/Values ===")
	scores := map[string]int{"Alice": 95, "Bob": 87, "Carol": 92}
	fmt.Println("Keys:", Keys(scores))
	fmt.Println("Values:", Values(scores))

	// 5. 泛型栈
	fmt.Println("\n=== 5. 泛型栈 ===")
	intStack := &Stack[int]{}
	intStack.Push(1)
	intStack.Push(2)
	intStack.Push(3)
	fmt.Printf("Stack size: %d\n", intStack.Size())
	if v, ok := intStack.Pop(); ok {
		fmt.Printf("Pop: %d\n", v)
	}
	if v, ok := intStack.Peek(); ok {
		fmt.Printf("Peek: %d\n", v)
	}

	strStack := &Stack[string]{}
	strStack.Push("hello")
	strStack.Push("world")
	if v, ok := strStack.Pop(); ok {
		fmt.Printf("String Pop: %s\n", v)
	}

	// 6. Pair
	fmt.Println("\n=== 6. Pair ===")
	p1 := NewPair("name", 42)
	fmt.Printf("Pair: (%v, %v)\n", p1.First, p1.Second)

	p2 := NewPair(3.14, true)
	fmt.Printf("Pair: (%v, %v)\n", p2.First, p2.Second)

	// 7. Result 类型
	fmt.Println("\n=== 7. Result ===")
	r1 := Ok(42)
	if r1.IsOk() {
		fmt.Println("Result OK:", r1.Unwrap())
	}

	r2 := Err[string](fmt.Errorf("something went wrong"))
	fmt.Printf("Result Error: %v\n", r2.Err)

	// 8. 类型约束
	fmt.Println("\n=== 8. 类型约束 ===")
	ints := []int{1, 2, 3, 4, 5}
	floats := []float64{1.1, 2.2, 3.3}
	fmt.Println("Sum ints:", Sum(ints))
	fmt.Println("Sum floats:", Sum(floats))

	// 自定义类型也适用（~int）
	type Score int
	scoreSlice := []Score{90, 85, 95, 88}
	fmt.Println("Sum scores:", Sum(scoreSlice))

	// 9. ArrayList
	fmt.Println("\n=== 9. 泛型集合 ===")
	list := NewArrayList[string]()
	list.Add("Go")
	list.Add("Rust")
	list.Add("Python")
	fmt.Printf("ArrayList size: %d\n", list.Size())
	if v, ok := list.Get(1); ok {
		fmt.Printf("Get(1): %s\n", v)
	}

	// 10. 类型推断
	fmt.Println("\n=== 10. 类型推断 ===")
	// 编译器可以推断类型参数
	fmt.Println(Max(1, 2))         // 推断 T = int
	fmt.Println(Max(1.5, 2.5))     // 推断 T = float64
	fmt.Println(Contains(ints, 3)) // 推断 T = int

	fmt.Println(`
📚 泛型使用指南：
1. 优先使用具体类型，只在真正需要时使用泛型
2. 好的使用场景：容器、通用算法、工具函数
3. 不好的使用场景：替代接口的多态行为
4. 保持约束尽量严格（不要默认用 any）
5. 标准库 slices 和 maps 包提供了很多泛型工具
6. cmp.Ordered 用于可比较大小的类型
7. comparable 用于可判等的类型`)
}

// ============================================================================
// 💻 运行方式：
//   cd 09-generics && go run main.go
//
// ⚠️ 需要 Go 1.18+ (推荐 Go 1.21+)
//
// 📝 练习：
// 1. 实现泛型 Queue（FIFO）
// 2. 实现泛型 Set（基于 map）
// 3. 实现泛型 BinarySearch 函数
// 4. 实现泛型 GroupBy 函数（按键分组）
// ============================================================================
