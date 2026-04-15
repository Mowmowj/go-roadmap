// ============================================================================
// 第08章：方法与接口 🟡 中级 | ⭐ 必须学习
// ============================================================================
// 本章内容：
// 1. 方法基础（值接收者 vs 指针接收者）
// 2. 接口基础
// 3. 空接口 interface{}
// 4. 接口嵌入与组合
// 5. 类型断言与类型选择 (type switch)
// 6. 常见接口模式 (Stringer, error, io.Reader/Writer)
// ============================================================================

package main

import (
	"fmt"
	"math"
	"strings"
)

// ============================================================================
// 1. 方法基础
// ============================================================================

// 📚 概念：方法 (Method)
// 方法是带有特殊接收者(receiver)参数的函数
// 语法: func (接收者) 方法名(参数列表) 返回值
// Go 没有类(class)，通过方法给类型添加行为

// 定义一个类型
type Rectangle struct {
	Width, Height float64
}

// 值接收者方法 — 不修改原对象
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// 实现 Stringer 接口（类似其他语言的 toString）
func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle(%.1f x %.1f)", r.Width, r.Height)
}

// 指针接收者方法 — 可以修改原对象
func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

// 📚 值接收者 vs 指针接收者选择指南：
// 用指针接收者：
//   1. 方法需要修改接收者
//   2. 接收者是大型结构体（避免复制）
//   3. 保持一致性：如果某个方法用了指针接收者，其他方法也应该用
// 用值接收者：
//   1. 接收者是小型且不可变的
//   2. map、slice 等引用类型（本身就含指针）

// ============================================================================
// 2. 给非结构体类型定义方法
// ============================================================================

type Celsius float64
type Fahrenheit float64

func (c Celsius) ToFahrenheit() Fahrenheit {
	return Fahrenheit(c*9/5 + 32)
}

func (f Fahrenheit) ToCelsius() Celsius {
	return Celsius((f - 32) * 5 / 9)
}

func (c Celsius) String() string {
	return fmt.Sprintf("%.1f°C", float64(c))
}

// ============================================================================
// 3. 接口基础
// ============================================================================

// 📚 概念：接口 (Interface)
// 接口定义了一组方法签名，任何类型只要实现了所有方法就自动满足接口
// Go 的接口是隐式实现的（不需要 implements 关键字）
// 这是 Go 最强大的特性之一 — "组合优于继承"

type Shape interface {
	Area() float64
	Perimeter() float64
}

// Circle 也实现 Shape 接口
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle(r=%.1f)", c.Radius)
}

// Triangle 也实现 Shape 接口
type Triangle struct {
	A, B, C float64 // 三条边
}

func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

func (t Triangle) Area() float64 {
	// 海伦公式
	s := t.Perimeter() / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

// 使用接口的函数
func printShapeInfo(s Shape) {
	fmt.Printf("  形状: %v\n", s)
	fmt.Printf("  面积: %.2f\n", s.Area())
	fmt.Printf("  周长: %.2f\n", s.Perimeter())
}

// ============================================================================
// 4. 空接口 interface{} / any
// ============================================================================

// 📚 概念：空接口
// interface{} 不包含任何方法，所以所有类型都实现了空接口
// Go 1.18+ 提供了 any 作为 interface{} 的别名
// 类似其他语言中的 Object, void*, dynamic
// 尽量少用空接口，会失去类型安全

func describe(i interface{}) {
	fmt.Printf("  值: %v, 类型: %T\n", i, i)
}

// ============================================================================
// 5. 接口嵌入与组合
// ============================================================================

// 小接口组合成大接口 — Go 推崇小接口
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Closer interface {
	Close() error
}

// 组合接口
type ReadWriter interface {
	Reader
	Writer
}

type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}

// 实际例子：一个简单的缓冲区
type Buffer struct {
	data []byte
}

func (b *Buffer) Read(p []byte) (int, error) {
	n := copy(p, b.data)
	b.data = b.data[n:]
	return n, nil
}

func (b *Buffer) Write(p []byte) (int, error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

// Buffer 同时实现了 Reader 和 Writer，所以也满足 ReadWriter

// ============================================================================
// 6. 类型断言与类型选择
// ============================================================================

// 📚 概念：类型断言 (Type Assertion)
// 从接口值中提取具体类型
// 语法: value, ok := interfaceVar.(ConcreteType)
// 不带 ok 的写法: value := interfaceVar.(ConcreteType) — 失败会 panic

// 📚 概念：类型选择 (Type Switch)
// 比链式类型断言更优雅的方式

func classifyShape(s Shape) string {
	switch v := s.(type) {
	case Rectangle:
		return fmt.Sprintf("矩形 %s, 面积=%.1f", v, v.Area())
	case Circle:
		return fmt.Sprintf("圆形 %s, 面积=%.1f", v, v.Area())
	case Triangle:
		return fmt.Sprintf("三角形, 面积=%.1f", v.Area())
	default:
		return fmt.Sprintf("未知形状: %T", v)
	}
}

// ============================================================================
// 7. 常见内置接口
// ============================================================================

// fmt.Stringer — 类似 toString()
type Student struct {
	Name   string
	Grade  int
	Scores map[string]float64
}

func (s Student) String() string {
	var subjects []string
	for k, v := range s.Scores {
		subjects = append(subjects, fmt.Sprintf("%s:%.0f", k, v))
	}
	return fmt.Sprintf("%s (Grade %d) [%s]", s.Name, s.Grade, strings.Join(subjects, ", "))
}

// error 接口 — Go 的错误处理基础
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

func validateAge(age int) error {
	if age < 0 || age > 150 {
		return &ValidationError{Field: "age", Message: "must be between 0 and 150"}
	}
	return nil
}

// ============================================================================
// 8. 接口最佳实践
// ============================================================================

// 📚 Go 接口设计哲学：
// "Accept interfaces, return structs"
// 1. 函数参数用接口 → 灵活
// 2. 函数返回用具体类型 → 明确
// 3. 接口尽量小 → 一到两个方法最佳
// 4. 在使用方定义接口，不在实现方

// 例子：在消费方定义接口
type Notifier interface {
	Notify(message string) error
}

type EmailSender struct {
	Address string
}

func (e EmailSender) Notify(message string) error {
	fmt.Printf("  📧 发送邮件到 %s: %s\n", e.Address, message)
	return nil
}

type SMSSender struct {
	Phone string
}

func (s SMSSender) Notify(message string) error {
	fmt.Printf("  📱 发送短信到 %s: %s\n", s.Phone, message)
	return nil
}

// Alert 依赖接口而非具体类型
func Alert(n Notifier, msg string) {
	if err := n.Notify(msg); err != nil {
		fmt.Println("通知失败:", err)
	}
}

// ============================================================================
func main() {
	// 1. 方法演示
	fmt.Println("=== 1. 方法 ===")
	rect := Rectangle{Width: 10, Height: 5}
	fmt.Printf("%s 面积=%.1f 周长=%.1f\n", rect, rect.Area(), rect.Perimeter())

	rect.Scale(2)
	fmt.Printf("Scale(2)后: %s 面积=%.1f\n", rect, rect.Area())

	// 给非结构体类型定义方法
	temp := Celsius(100)
	fmt.Printf("%v = %v\n", temp, temp.ToFahrenheit())

	// 2. 接口多态
	fmt.Println("\n=== 2. 接口多态 ===")
	shapes := []Shape{
		Rectangle{Width: 5, Height: 3},
		Circle{Radius: 4},
		Triangle{A: 3, B: 4, C: 5},
	}

	for _, s := range shapes {
		printShapeInfo(s)
		fmt.Println()
	}

	// 3. 空接口
	fmt.Println("=== 3. 空接口 ===")
	describe(42)
	describe("hello")
	describe(true)
	describe(Rectangle{Width: 1, Height: 2})
	describe([]int{1, 2, 3})

	// 4. 接口组合
	fmt.Println("\n=== 4. 接口组合 ===")
	var rw ReadWriter = &Buffer{}
	rw.Write([]byte("Hello, Go interfaces!"))
	buf := make([]byte, 100)
	n, _ := rw.Read(buf)
	fmt.Printf("Buffer 读到: %s\n", string(buf[:n]))

	// 5. 类型断言
	fmt.Println("\n=== 5. 类型断言 ===")
	var s Shape = Circle{Radius: 5}

	// 安全断言（带 ok）
	if c, ok := s.(Circle); ok {
		fmt.Printf("是圆形，半径=%.1f\n", c.Radius)
	}

	// 断言失败
	if _, ok := s.(Rectangle); !ok {
		fmt.Println("不是矩形")
	}

	// 6. 类型选择
	fmt.Println("\n=== 6. 类型选择 ===")
	for _, shape := range shapes {
		fmt.Println(classifyShape(shape))
	}

	// 7. Stringer 接口
	fmt.Println("\n=== 7. Stringer 接口 ===")
	student := Student{
		Name:  "小明",
		Grade: 3,
		Scores: map[string]float64{
			"数学": 95, "英语": 88, "科学": 92,
		},
	}
	fmt.Println(student) // 自动调用 String()

	// 8. error 接口
	fmt.Println("\n=== 8. error 接口 ===")
	if err := validateAge(200); err != nil {
		fmt.Println("错误:", err)

		// 类型断言获取详细信息
		if ve, ok := err.(*ValidationError); ok {
			fmt.Printf("字段: %s, 消息: %s\n", ve.Field, ve.Message)
		}
	}

	// 9. 接口设计模式
	fmt.Println("\n=== 9. 接口设计模式 ===")
	Alert(EmailSender{Address: "user@example.com"}, "系统告警!")
	Alert(SMSSender{Phone: "138-0000-0000"}, "服务器宕机!")

	// 10. nil 接口 vs nil 值
	fmt.Println("\n=== 10. nil 接口陷阱 ===")
	// ⚠️ 接口值包含两部分：(type, value)
	// 只有两部分都是 nil 时，接口才等于 nil
	var myErr error // nil（type=nil, value=nil）
	fmt.Println("myErr == nil:", myErr == nil)

	var vp *ValidationError = nil
	myErr = vp                                     // 注意！type=*ValidationError, value=nil
	fmt.Println("赋值后 myErr == nil:", myErr == nil) // false! 陷阱！
	fmt.Println("这是因为接口内部的 type 不为 nil")

	fmt.Println(`
📚 关键总结：
1. 方法让类型拥有行为
2. 值接收者不修改，指针接收者可修改
3. 接口是隐式实现的（鸭子类型）
4. 小接口 + 组合 = Go 的多态
5. 空接口可存任意类型，但失去类型安全
6. 类型断言和类型选择用于从接口中恢复具体类型
7. 注意 nil 接口陷阱`)
}

// ============================================================================
// 💻 运行方式：
//   cd 08-methods-interfaces && go run main.go
//
// 📝 练习：
// 1. 实现一个 Storage 接口（Save/Load），分别用内存和文件实现
// 2. 实现 sort.Interface（Len/Less/Swap）给自定义类型排序
// 3. 实现一个 Logger 接口，支持不同输出目标
// ============================================================================
