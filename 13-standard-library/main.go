// ============================================================================
// 第13章：标准库精选 🟡 中级 | ⭐ 必须学习
// ============================================================================
// 本章内容：
// 1. fmt — 格式化 I/O
// 2. os — 操作系统交互
// 3. io / bufio — I/O 操作
// 4. strings / strconv — 字符串处理
// 5. time — 时间处理
// 6. encoding/json — JSON 编解码
// 7. net/http — HTTP 客户端与服务端
// 8. flag — 命令行参数
// 9. sort — 排序
// 10. regexp — 正则表达式
// 11. filepath — 文件路径
// 12. log/slog — 结构化日志 (Go 1.21+)
// ============================================================================

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	// ========================================================================
	// 1. fmt — 格式化 I/O
	// ========================================================================
	fmt.Println("=== 1. fmt 包 ===")

	name := "Go"
	version := 1.21
	fmt.Println("Println:", name)
	fmt.Printf("Printf: %s v%.2f\n", name, version)

	// 常用格式化动词
	x := 42
	fmt.Printf("  %%v  = %v\n", x)       // 默认格式
	fmt.Printf("  %%d  = %d\n", x)       // 十进制
	fmt.Printf("  %%b  = %b\n", x)       // 二进制
	fmt.Printf("  %%o  = %o\n", x)       // 八进制
	fmt.Printf("  %%x  = %x\n", x)       // 十六进制
	fmt.Printf("  %%f  = %f\n", 3.14)    // 浮点数
	fmt.Printf("  %%e  = %e\n", 3.14)    // 科学计数法
	fmt.Printf("  %%s  = %s\n", "hello") // 字符串
	fmt.Printf("  %%q  = %q\n", "hello") // 带引号字符串
	fmt.Printf("  %%t  = %t\n", true)    // 布尔值
	fmt.Printf("  %%T  = %T\n", x)       // 类型
	fmt.Printf("  %%p  = %p\n", &x)      // 指针

	type Point struct{ X, Y int }
	p := Point{1, 2}
	fmt.Printf("  %%v  = %v\n", p)  // {1 2}
	fmt.Printf("  %%+v = %+v\n", p) // {X:1 Y:2}
	fmt.Printf("  %%#v = %#v\n", p) // main.Point{X:1, Y:2}

	// Sprintf 返回字符串
	s := fmt.Sprintf("Hello, %s! You are %d.", "World", 42)
	fmt.Println("  Sprintf:", s)

	// Fprintf 写入 io.Writer
	fmt.Fprintf(os.Stdout, "  Fprintf: %s\n", "写到标准输出")

	// ========================================================================
	// 2. strings — 字符串操作
	// ========================================================================
	fmt.Println("\n=== 2. strings 包 ===")

	str := "Hello, Go World!"
	fmt.Println("  Contains:", strings.Contains(str, "Go"))
	fmt.Println("  HasPrefix:", strings.HasPrefix(str, "Hello"))
	fmt.Println("  HasSuffix:", strings.HasSuffix(str, "World!"))
	fmt.Println("  Index:", strings.Index(str, "Go"))
	fmt.Println("  ToUpper:", strings.ToUpper(str))
	fmt.Println("  ToLower:", strings.ToLower(str))
	fmt.Println("  TrimSpace:", strings.TrimSpace("  hello  "))
	fmt.Println("  Replace:", strings.Replace(str, "Go", "Golang", 1))
	fmt.Println("  ReplaceAll:", strings.ReplaceAll("aaa", "a", "b"))
	fmt.Println("  Split:", strings.Split("a,b,c,d", ","))
	fmt.Println("  Join:", strings.Join([]string{"a", "b", "c"}, "-"))
	fmt.Println("  Repeat:", strings.Repeat("Go ", 3))
	fmt.Println("  Count:", strings.Count("hello", "l"))
	fmt.Println("  EqualFold:", strings.EqualFold("Go", "go")) // 忽略大小写

	// strings.Builder — 高效字符串拼接
	var builder strings.Builder
	for i := 0; i < 5; i++ {
		builder.WriteString(fmt.Sprintf("item%d ", i))
	}
	fmt.Println("  Builder:", builder.String())

	// strings.NewReader
	reader := strings.NewReader("Hello from reader")
	buf := make([]byte, 5)
	n, _ := reader.Read(buf)
	fmt.Printf("  Reader: %s\n", string(buf[:n]))

	// ========================================================================
	// 3. strconv — 类型转换
	// ========================================================================
	fmt.Println("\n=== 3. strconv 包 ===")

	// 字符串 → 数字
	num, _ := strconv.Atoi("42")
	fmt.Println("  Atoi:", num)

	f64, _ := strconv.ParseFloat("3.14", 64)
	fmt.Println("  ParseFloat:", f64)

	b, _ := strconv.ParseBool("true")
	fmt.Println("  ParseBool:", b)

	// 数字 → 字符串
	fmt.Println("  Itoa:", strconv.Itoa(42))
	fmt.Println("  FormatFloat:", strconv.FormatFloat(3.14, 'f', 2, 64))
	fmt.Println("  FormatBool:", strconv.FormatBool(true))

	// ========================================================================
	// 4. os — 操作系统交互
	// ========================================================================
	fmt.Println("\n=== 4. os 包 ===")

	// 环境变量
	fmt.Println("  HOME:", os.Getenv("HOME"))
	fmt.Println("  PATH存在:", os.Getenv("PATH") != "")

	// 命令行参数
	fmt.Println("  Args:", os.Args[0])

	// 获取工作目录
	wd, _ := os.Getwd()
	fmt.Println("  工作目录:", wd)

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "go-demo-*.txt")
	if err == nil {
		tmpFile.WriteString("Hello, temp file!")
		tmpFile.Close()
		fmt.Println("  临时文件:", tmpFile.Name())

		// 读取文件
		content, _ := os.ReadFile(tmpFile.Name())
		fmt.Println("  文件内容:", string(content))

		// 删除临时文件
		os.Remove(tmpFile.Name())
	}

	// 文件信息
	info, err := os.Stat(".")
	if err == nil {
		fmt.Printf("  当前目录: name=%s, isDir=%t\n", info.Name(), info.IsDir())
	}

	// ========================================================================
	// 5. io / bufio — I/O 操作
	// ========================================================================
	fmt.Println("\n=== 5. io/bufio ===")

	// io.Reader / io.Writer 是核心接口
	// 许多标准库类型都实现了这些接口

	// io.Copy
	src := strings.NewReader("io.Copy 示例数据")
	var dst bytes.Buffer
	io.Copy(&dst, src)
	fmt.Println("  io.Copy:", dst.String())

	// bufio.Scanner — 逐行读取
	input := "line1\nline2\nline3"
	scanner := bufio.NewScanner(strings.NewReader(input))
	fmt.Print("  Scanner 逐行: ")
	for scanner.Scan() {
		fmt.Printf("[%s] ", scanner.Text())
	}
	fmt.Println()

	// bufio.Scanner — 按单词读取
	wordScanner := bufio.NewScanner(strings.NewReader("hello world foo bar"))
	wordScanner.Split(bufio.ScanWords)
	fmt.Print("  Scanner 按词: ")
	for wordScanner.Scan() {
		fmt.Printf("[%s] ", wordScanner.Text())
	}
	fmt.Println()

	// ========================================================================
	// 6. time — 时间处理
	// ========================================================================
	fmt.Println("\n=== 6. time 包 ===")

	now := time.Now()
	fmt.Println("  当前时间:", now)
	fmt.Println("  Unix 时间戳:", now.Unix())

	// ⚠️ Go 使用特殊的时间格式化参考值: 2006-01-02 15:04:05
	// 记忆方法: 1月2日下午3点4分5秒2006年 (1-2-3-4-5-6)
	fmt.Println("  格式化:", now.Format("2006-01-02 15:04:05"))
	fmt.Println("  仅日期:", now.Format("2006/01/02"))
	fmt.Println("  仅时间:", now.Format("15:04:05"))
	fmt.Println("  RFC3339:", now.Format(time.RFC3339))

	// 解析时间
	t, _ := time.Parse("2006-01-02", "2024-01-15")
	fmt.Println("  解析时间:", t)

	// 时间运算
	future := now.Add(24 * time.Hour)
	fmt.Println("  明天:", future.Format("2006-01-02"))

	duration := future.Sub(now)
	fmt.Println("  时间差:", duration)

	// Duration
	fmt.Println("  5秒:", 5*time.Second)
	fmt.Println("  1.5小时:", time.Duration(1.5*float64(time.Hour)))

	// 计时
	start := time.Now()
	time.Sleep(10 * time.Millisecond)
	elapsed := time.Since(start)
	fmt.Printf("  耗时: %v\n", elapsed)

	// ========================================================================
	// 7. encoding/json — JSON 处理
	// ========================================================================
	fmt.Println("\n=== 7. encoding/json ===")

	type User struct {
		Name   string   `json:"name"`
		Age    int      `json:"age"`
		Email  string   `json:"email,omitempty"` // 空值时省略
		Tags   []string `json:"tags"`
		Admin  bool     `json:"admin"`
		Secret string   `json:"-"` // 忽略此字段
	}

	// 结构体 → JSON
	user := User{
		Name:   "Alice",
		Age:    30,
		Tags:   []string{"go", "developer"},
		Admin:  true,
		Secret: "hidden",
	}

	jsonBytes, _ := json.Marshal(user)
	fmt.Println("  Marshal:", string(jsonBytes))

	// 格式化输出
	prettyJSON, _ := json.MarshalIndent(user, "  ", "  ")
	fmt.Println("  MarshalIndent:")
	fmt.Println(" ", string(prettyJSON))

	// JSON → 结构体
	jsonStr := `{"name":"Bob","age":25,"tags":["rust","wasm"],"admin":false}`
	var user2 User
	json.Unmarshal([]byte(jsonStr), &user2)
	fmt.Printf("  Unmarshal: %+v\n", user2)

	// JSON → map（动态 JSON）
	var m map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &m)
	fmt.Printf("  Map: name=%v, age=%v\n", m["name"], m["age"])

	// ========================================================================
	// 8. net/http — HTTP 入门
	// ========================================================================
	fmt.Println("\n=== 8. net/http 预览 ===")

	// HTTP GET 请求（简化版）
	fmt.Println("  HTTP GET 示例 (伪代码):")
	fmt.Println(`
    resp, err := http.Get("https://api.example.com/data")
    if err != nil { ... }
    defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)`)

	// HTTP 服务端（仅展示代码结构）
	fmt.Println("\n  HTTP 服务端示例 (伪代码):")
	fmt.Println(`
    mux := http.NewServeMux()
    mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })
    http.ListenAndServe(":8080", mux)`)

	// 展示如何创建一个简单的 handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!")
	})
	_ = handler // 只展示，不启动服务器

	// ========================================================================
	// 9. sort — 排序
	// ========================================================================
	fmt.Println("\n=== 9. sort 包 ===")

	// 基本排序
	ints := []int{5, 3, 1, 4, 2}
	sort.Ints(ints)
	fmt.Println("  排序 int:", ints)

	strs := []string{"banana", "apple", "cherry"}
	sort.Strings(strs)
	fmt.Println("  排序 string:", strs)

	// 自定义排序
	type Person struct {
		Name string
		Age  int
	}
	people := []Person{
		{"Alice", 30}, {"Bob", 25}, {"Carol", 35},
	}
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age < people[j].Age
	})
	fmt.Println("  按年龄排序:", people)

	// 稳定排序
	sort.SliceStable(people, func(i, j int) bool {
		return people[i].Name < people[j].Name
	})
	fmt.Println("  按名字排序:", people)

	// 二分查找
	idx := sort.SearchInts(ints, 3)
	fmt.Printf("  二分查找 3: index=%d\n", idx)

	// ========================================================================
	// 10. regexp — 正则表达式
	// ========================================================================
	fmt.Println("\n=== 10. regexp 包 ===")

	// 编译正则
	re := regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`)

	text := "联系 alice@go.dev 或 bob@example.com 了解详情"
	fmt.Println("  匹配:", re.MatchString(text))
	fmt.Println("  查找:", re.FindString(text))
	fmt.Println("  查找全部:", re.FindAllString(text, -1))

	// 替换
	masked := re.ReplaceAllString(text, "[EMAIL]")
	fmt.Println("  替换:", masked)

	// 分组捕获
	re2 := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})`)
	match := re2.FindStringSubmatch("日期: 2024-01-15")
	if match != nil {
		fmt.Printf("  分组: 全=%s, 年=%s, 月=%s, 日=%s\n",
			match[0], match[1], match[2], match[3])
	}

	// ========================================================================
	// 11. filepath — 文件路径
	// ========================================================================
	fmt.Println("\n=== 11. filepath 包 ===")

	path := "/usr/local/bin/go"
	fmt.Println("  Dir:", filepath.Dir(path))
	fmt.Println("  Base:", filepath.Base(path))
	fmt.Println("  Ext:", filepath.Ext("/tmp/test.go"))
	fmt.Println("  Join:", filepath.Join("usr", "local", "bin"))

	// 遍历目录（只打印前几个文件）
	fmt.Print("  当前目录文件: ")
	count := 0
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || count >= 5 {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			fmt.Printf("%s ", path)
			count++
		}
		return nil
	})
	fmt.Println()

	// ========================================================================
	// 12. flag — 命令行参数
	// ========================================================================
	fmt.Println("\n=== 12. flag 包 ===")
	fmt.Println("  flag 示例（因为已解析，仅展示代码）:")
	os.Stdout.WriteString(`
    port := flag.Int("port", 8080, "服务端口")
    host := flag.String("host", "localhost", "服务地址")
    verbose := flag.Bool("v", false, "详细输出")
    flag.Parse()

    fmt.Printf("启动 %s:%d (verbose=%t)\n", *host, *port, *verbose)
    
    使用: go run main.go -port=3000 -host=0.0.0.0 -v
`)

	// 展示 flag 的基本用法但不真正解析（因为可能和其他参数冲突）
	fs := flag.NewFlagSet("demo", flag.ContinueOnError)
	port := fs.Int("port", 8080, "服务端口")
	fs.Parse([]string{"-port", "3000"})
	fmt.Printf("  Flag 解析: port=%d\n", *port)

	fmt.Println(`
📚 标准库使用建议：
1. 先查标准库，再找第三方（Go 标准库非常强大）
2. io.Reader/io.Writer 是核心抽象，理解它们
3. 时间格式化记住 "2006-01-02 15:04:05"
4. JSON tag 控制序列化行为
5. strings.Builder 比 += 拼接高效
6. regexp.MustCompile 在 init 时编译，运行时复用
7. 文档：https://pkg.go.dev/std`)
}

// ============================================================================
// 💻 运行方式：
//   cd 13-standard-library && go run main.go
//
// 📝 练习：
// 1. 写一个命令行工具，读取文件并统计单词频率
// 2. 写一个 JSON 配置文件加载器
// 3. 用 time.Ticker 实现定时任务
// 4. 用 regexp 实现一个简单的模板引擎
// ============================================================================
