// ============================================================================
// 第10章：错误处理 🟡 中级 | ⭐ 必须学习
// ============================================================================
// 本章内容：
// 1. error 接口
// 2. 创建错误 (errors.New, fmt.Errorf)
// 3. 自定义错误类型
// 4. 错误包装与展开 (wrapping/unwrapping)
// 5. errors.Is 与 errors.As
// 6. panic 与 recover
// 7. 错误处理最佳实践
// ============================================================================

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// ============================================================================
// 1. error 接口
// ============================================================================

// 📚 概念：Go 的错误处理
// Go 没有 try/catch/throw 异常机制
// error 是一个内置接口：
//   type error interface {
//       Error() string
//   }
// 函数通过返回 error 来表示失败，调用者必须显式处理
// 这种方式让错误处理变得显式和清晰

// ============================================================================
// 2. 哨兵错误 (Sentinel Errors)
// ============================================================================

// 📚 哨兵错误 — 包级别的预定义错误变量
// 命名规范：ErrXxx
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
)

// ============================================================================
// 3. 自定义错误类型
// ============================================================================

// 自定义错误类型可以携带更多上下文信息
type HTTPError struct {
	StatusCode int
	Message    string
	URL        string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s (url: %s)", e.StatusCode, e.Message, e.URL)
}

type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: field '%s' with value '%v': %s",
		e.Field, e.Value, e.Message)
}

// 多个错误
type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	msgs := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		msgs[i] = err.Error()
	}
	return fmt.Sprintf("%d errors: [%s]", len(e.Errors), joinStrings(msgs, "; "))
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

// ============================================================================
// 4. 错误包装 (Error Wrapping)
// ============================================================================

// 📚 概念：错误包装
// Go 1.13+ 引入 %w 动词，用于包装错误
// 包装后的错误保留了原始错误链，可以用 errors.Is/As 检查

func readConfig(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		// 用 %w 包装错误，保留原始错误链
		return "", fmt.Errorf("reading config %s: %w", filename, err)
	}
	return string(data), nil
}

func loadSettings() (string, error) {
	config, err := readConfig("/nonexistent/config.yaml")
	if err != nil {
		// 可以多层包装
		return "", fmt.Errorf("loading settings: %w", err)
	}
	return config, nil
}

// ============================================================================
// 5. 业务逻辑中的错误处理
// ============================================================================

type User struct {
	Name  string
	Email string
	Age   int
}

func validateUser(u User) error {
	var errs []error

	if u.Name == "" {
		errs = append(errs, &ValidationError{Field: "name", Value: u.Name, Message: "cannot be empty"})
	}
	if u.Age < 0 || u.Age > 150 {
		errs = append(errs, &ValidationError{Field: "age", Value: u.Age, Message: "must be between 0 and 150"})
	}
	if u.Email == "" {
		errs = append(errs, &ValidationError{Field: "email", Value: u.Email, Message: "cannot be empty"})
	}

	if len(errs) > 0 {
		return &MultiError{Errors: errs}
	}
	return nil
}

func findUser(id int) (*User, error) {
	users := map[int]*User{
		1: {Name: "Alice", Email: "alice@go.dev", Age: 30},
		2: {Name: "Bob", Email: "bob@go.dev", Age: 25},
	}

	user, exists := users[id]
	if !exists {
		return nil, fmt.Errorf("user id=%d: %w", id, ErrNotFound)
	}
	return user, nil
}

func parseID(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("parsing id %q: %w", s, ErrInvalidInput)
	}
	if id <= 0 {
		return 0, fmt.Errorf("id must be positive, got %d: %w", id, ErrInvalidInput)
	}
	return id, nil
}

// ============================================================================
// 6. panic 和 recover
// ============================================================================

// 📚 概念：panic & recover
// panic: 不可恢复的错误（类似其他语言的异常），会停止当前 goroutine
// recover: 在 defer 中捕获 panic，防止程序崩溃
// 规则：
//   - 普通错误用 error 返回
//   - 只在真正不可恢复的情况用 panic
//   - 库代码永远不应该 panic（应该返回 error）
//   - recover 只在 defer 中有效

func safeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()
	return a / b, nil // b=0 会 panic
}

func mustParseInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("mustParseInt(%q): %v", s, err))
	}
	return n
}

// ============================================================================
func main() {
	// 1. 基础错误处理
	fmt.Println("=== 1. 基础错误处理 ===")
	n, err := strconv.Atoi("42")
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("解析成功:", n)
	}

	_, err = strconv.Atoi("not-a-number")
	if err != nil {
		fmt.Println("解析失败:", err)
	}

	// 2. 哨兵错误
	fmt.Println("\n=== 2. 哨兵错误 ===")
	user, err := findUser(1)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Printf("找到用户: %+v\n", user)
	}

	_, err = findUser(999)
	if err != nil {
		fmt.Println("错误:", err)
		// errors.Is 检查错误链中是否有目标错误
		if errors.Is(err, ErrNotFound) {
			fmt.Println("→ 这是 NotFound 错误")
		}
	}

	// 3. 错误包装与展开
	fmt.Println("\n=== 3. 错误包装 ===")
	_, err = loadSettings()
	if err != nil {
		fmt.Println("错误:", err)
		// 展开错误链
		fmt.Println("Unwrap:", errors.Unwrap(err))
		// errors.Is 可以穿透多层包装
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("→ 根本原因: 文件不存在")
		}
	}

	// 4. errors.As — 提取特定错误类型
	fmt.Println("\n=== 4. errors.As ===")
	_, err = strconv.Atoi("xyz")
	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		fmt.Printf("NumError: Func=%s, Num=%s\n", numErr.Func, numErr.Num)
	}

	// 5. 自定义验证错误
	fmt.Println("\n=== 5. 自定义错误类型 ===")
	badUser := User{Name: "", Email: "", Age: -5}
	if err := validateUser(badUser); err != nil {
		fmt.Println("验证失败:", err)
	}

	goodUser := User{Name: "Alice", Email: "alice@go.dev", Age: 30}
	if err := validateUser(goodUser); err != nil {
		fmt.Println("验证失败:", err)
	} else {
		fmt.Println("验证通过!")
	}

	// 6. HTTP 错误
	fmt.Println("\n=== 6. HTTP 错误 ===")
	httpErr := &HTTPError{StatusCode: 404, Message: "Page Not Found", URL: "/api/users/999"}
	fmt.Println(httpErr)

	// 7. parseID 示例
	fmt.Println("\n=== 7. 链式错误处理 ===")
	if id, err := parseID("42"); err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("解析ID:", id)
	}

	if _, err := parseID("abc"); err != nil {
		fmt.Println("错误:", err)
		if errors.Is(err, ErrInvalidInput) {
			fmt.Println("→ 输入无效")
		}
	}

	// 8. panic 和 recover
	fmt.Println("\n=== 8. panic & recover ===")

	// safeDivide 内部用 recover 捕获除零 panic
	result, err := safeDivide(10, 3)
	fmt.Printf("10/3 = %d, err = %v\n", result, err)

	result, err = safeDivide(10, 0)
	fmt.Printf("10/0: result=%d, err = %v\n", result, err)

	// must 模式 — 只在初始化时使用
	fmt.Println("\n=== must 模式 ===")
	port := mustParseInt("8080")
	fmt.Println("端口:", port)
	// mustParseInt("abc") // 这会 panic

	// 9. defer + recover 模式
	fmt.Println("\n=== 9. defer + recover ===")
	fmt.Println("调用可能 panic 的函数...")
	safeCall(func() {
		fmt.Println("  正常执行")
	})
	safeCall(func() {
		panic("something bad happened!")
	})
	fmt.Println("程序继续运行（panic 被恢复）")

	fmt.Println(`
📚 错误处理最佳实践：
1. 错误是值：可以组合、包装、编程处理
2. 用 fmt.Errorf + %w 包装错误，保留错误链
3. 用 errors.Is 检查特定错误，用 errors.As 提取错误类型
4. 哨兵错误用于可预期的错误情况
5. 自定义错误类型携带上下文信息
6. 不要忽略错误（_ = someFunc() 是坏习惯）
7. 在错误消息中添加上下文，便于调试
8. panic 只用于不可恢复的程序错误
9. 库代码不应 panic，应返回 error
10. "errors are values" — Rob Pike`)
}

func safeCall(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("  捕获 panic: %v\n", r)
		}
	}()
	fn()
}

// ============================================================================
// 💻 运行方式：
//   cd 10-error-handling && go run main.go
//
// 📝 练习：
// 1. 实现一个文件解析器，返回带行号的错误信息
// 2. 实现自定义 AppError，包含错误码、消息、堆栈信息
// 3. 实现 retry 函数，在错误时重试 N 次
// ============================================================================
