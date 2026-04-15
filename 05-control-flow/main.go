// ============================================================================
// 第05章：控制流 🟢 初级 | ⭐ 必须学习
// ============================================================================
// 本章内容：
// 1. if / if-else / if-else if
// 2. switch 语句
// 3. for 循环（Go 唯一的循环语句）
// 4. for range 遍历
// 5. break, continue, goto
// ============================================================================

package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func main() {
	// ========================================================================
	// 1. if / if-else
	// ========================================================================
	fmt.Println("=== 1. if / if-else ===")

	// 基本 if
	age := 20
	if age >= 18 {
		fmt.Println("成年人")
	}

	// if-else
	score := 75
	if score >= 90 {
		fmt.Println("优秀")
	} else if score >= 80 {
		fmt.Println("良好")
	} else if score >= 60 {
		fmt.Println("及格")
	} else {
		fmt.Println("不及格")
	}

	// Go 特有：if 初始化语句
	// 📚 变量 err 只在 if-else 块内可见
	if num := 42; num%2 == 0 {
		fmt.Printf("%d 是偶数\n", num)
	} else {
		fmt.Printf("%d 是奇数\n", num)
	}
	// fmt.Println(num) // ❌ 编译错误！num 不在此作用域

	// 常见模式：错误检查
	if err := doWork(); err != nil {
		fmt.Printf("出错了: %v\n", err)
	} else {
		fmt.Println("工作完成")
	}

	// ========================================================================
	// 2. switch 语句
	// ========================================================================
	fmt.Println("\n=== 2. switch ===")

	// 基本 switch
	day := "Wednesday"
	switch day {
	case "Monday":
		fmt.Println("星期一")
	case "Tuesday":
		fmt.Println("星期二")
	case "Wednesday":
		fmt.Println("星期三 ✓")
	case "Thursday", "Friday": // 多个值匹配
		fmt.Println("星期四或五")
	default:
		fmt.Println("周末")
	}

	// 📚 概念：Go 的 switch 不需要 break
	// 每个 case 自动 break，不会 fallthrough
	// 如果需要穿透，使用 fallthrough 关键字
	//⚠️ fallthrough 是直接执行下一个 case 的代码，
	//不会再判断下一个 case 的条件是否成立。
	myscore := 95

	switch {
	case myscore >= 90:
		fmt.Println("优秀")
		fallthrough
	case myscore >= 60:
		fmt.Println("及格")
	default:
		fmt.Println("不及格")
	}
	// 输出：
	// 优秀
	// 及格
	// 但不会输出不及格，因为没有 fallthrough

	// switch 初始化语句
	switch n := 15; {
	case n < 0:
		fmt.Println("负数")
	case n == 0:
		fmt.Println("零")
	case n < 10:
		fmt.Println("个位数")
	case n < 100:
		fmt.Printf("%d 是两位数\n", n)
	default:
		fmt.Println("大数")
	}

	// 无条件 switch（替代长 if-else 链）
	hour := time.Now().Hour()
	switch {
	case hour < 6:
		fmt.Println("凌晨")
	case hour < 12:
		fmt.Println("上午")
	case hour < 18:
		fmt.Println("下午")
	default:
		fmt.Println("晚上")
	}

	// fallthrough（很少使用）
	fmt.Println("\nfallthrough 示例:")
	switch num := 1; num {
	case 1:
		fmt.Println("一")
		fallthrough // 继续执行下一个 case
	case 2:
		fmt.Println("二（被 fallthrough 执行）")
		// 不再 fallthrough，停在这里
	case 3:
		fmt.Println("三（不会执行）")
	}

	// Type switch（后续接口章节会详细讲）
	fmt.Println("\nType Switch:")
	checkType(42)
	checkType("hello")
	checkType(3.14)
	checkType(true)
	checkType([]int{1, 2})

	// ========================================================================
	// 3. for 循环
	// ========================================================================
	fmt.Println("\n=== 3. for 循环 ===")

	// 📚 概念：Go 只有 for 循环
	// 没有 while, do-while，全部用 for 实现

	// 标准 for（类似 C 的 for）
	fmt.Print("标准 for: ")
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// while 风格
	fmt.Print("while 风格: ")
	count := 0
	for count < 5 {
		fmt.Printf("%d ", count)
		count++
	}
	fmt.Println()

	// 无限循环
	fmt.Print("无限循环(带break): ")
	n := 0
	for {
		if n >= 5 {
			break
		}
		fmt.Printf("%d ", n)
		n++
	}
	fmt.Println()

	// ========================================================================
	// 4. for range 遍历
	// ========================================================================
	fmt.Println("\n=== 4. for range ===")

	// 遍历切片
	fruits := []string{"apple", "banana", "cherry"}
	fmt.Println("遍历切片:")
	for index, value := range fruits {
		fmt.Printf("  [%d] %s\n", index, value)
	}

	// 只要索引
	fmt.Print("只要索引: ")
	for i := range fruits {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// 只要值（用 _ 忽略索引）
	fmt.Print("只要值: ")
	for _, fruit := range fruits {
		fmt.Printf("%s ", fruit)
	}
	fmt.Println()

	// 遍历 Map
	fmt.Println("\n遍历 Map:")
	colors := map[string]string{
		"red":   "#FF0000",
		"green": "#00FF00",
		"blue":  "#0000FF",
	}
	for name, hex := range colors {
		fmt.Printf("  %s = %s\n", name, hex)
	}

	// 遍历字符串（按 rune 遍历）
	fmt.Println("\n遍历字符串:")
	greeting := "Hello, 世界!"
	for i, ch := range greeting {
		fmt.Printf("  byte=%d rune='%c'\n", i, ch)
	}

	// 遍历 channel（后续并发章节详解）

	// ========================================================================
	// 5. break, continue, goto
	// ========================================================================
	fmt.Println("\n=== 5. break, continue, goto ===")

	// break - 跳出循环
	fmt.Print("break: ")
	for i := 0; i < 10; i++ {
		if i == 5 {
			break
		}
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// continue - 跳过当次迭代
	fmt.Print("continue (跳过偶数): ")
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// 带标签的 break（跳出外层循环）
	fmt.Println("\n带标签 break (跳出嵌套循环):")
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				fmt.Println("  在 (1,1) 处跳出外层循环")
				break outer
			}
			fmt.Printf("  (%d,%d)\n", i, j)
		}
	}

	// goto（不推荐使用，但了解一下）
	// 📚 概念：goto 在 Go 中是合法的，但通常被认为是不好的实践
	// 主要用于跳出深层嵌套
	fmt.Println("\ngoto 示例 (不推荐):")
	i := 0
	if i == 0 {
		goto skipPrint
	}
	fmt.Println("这行不会执行")
skipPrint:
	fmt.Println("goto 跳到这里了")

	// ========================================================================
	// 综合实例：猜数字游戏
	// ========================================================================
	fmt.Println("\n=== 综合实例：猜数字逻辑 ===")
	target := rand.Intn(100) + 1
	guesses := []int{50, 25, 75, target}

	for attempt, guess := range guesses {
		fmt.Printf("第 %d 次猜测: %d -> ", attempt+1, guess)
		switch {
		case guess < target:
			fmt.Println("太小了")
		case guess > target:
			fmt.Println("太大了")
		default:
			fmt.Printf("🎉 猜对了！答案是 %d\n", target)
		}
	}

	// ========================================================================
	// 实用模式：字符串处理
	// ========================================================================
	fmt.Println("\n=== 实用模式 ===")

	// FizzBuzz
	fmt.Print("FizzBuzz: ")
	for i := 1; i <= 20; i++ {
		switch {
		case i%15 == 0:
			fmt.Print("FizzBuzz ")
		case i%3 == 0:
			fmt.Print("Fizz ")
		case i%5 == 0:
			fmt.Print("Buzz ")
		default:
			fmt.Printf("%d ", i)
		}
	}
	fmt.Println()

	// 构建金字塔
	fmt.Println("\n金字塔:")
	rows := 5
	for i := 1; i <= rows; i++ {
		spaces := strings.Repeat(" ", rows-i)
		stars := strings.Repeat("* ", i)
		fmt.Printf("%s%s\n", spaces, stars)
	}
}

func doWork() error {
	return nil // 模拟成功
}

func checkType(v interface{}) {
	switch v.(type) {
	case int:
		fmt.Printf("  %v 是 int\n", v)
	case string:
		fmt.Printf("  %v 是 string\n", v)
	case float64:
		fmt.Printf("  %v 是 float64\n", v)
	case bool:
		fmt.Printf("  %v 是 bool\n", v)
	default:
		fmt.Printf("  %v 是未知类型 %T\n", v, v)
	}
}

// ============================================================================
// 💻 运行方式：
//   cd 05-control-flow && go run main.go
//
// 📝 练习：
// 1. 用 switch 实现一个简单的计算器（+-*/）
// 2. 用 for 实现九九乘法表
// 3. 实现一个 FizzBuzz 变种：3=Fizz, 5=Buzz, 7=Woof
// ============================================================================
