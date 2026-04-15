// ============================================================================
// 第14章：测试 🟡 中级 | ⭐ 必须学习
// ============================================================================
// 本文件供 14-testing/ 下的测试文件导入
// ============================================================================

package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// ============================================================================
// 被测试的函数和类型
// ============================================================================

// 计算器
func Add(a, b float64) float64      { return a + b }
func Subtract(a, b float64) float64 { return a - b }
func Multiply(a, b float64) float64 { return a * b }

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// 是否为质数
func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// 斐波那契
func Fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// 字符串工具
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	return s == Reverse(s)
}

func WordCount(s string) map[string]int {
	counts := make(map[string]int)
	words := strings.Fields(s)
	for _, w := range words {
		counts[strings.ToLower(w)]++
	}
	return counts
}

// 用户服务（用于演示 mock）
type UserRepository interface {
	FindByID(id int) (*User, error)
	Save(user *User) error
}

type User struct {
	ID    int
	Name  string
	Email string
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id int) (*User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user id")
	}
	return s.repo.FindByID(id)
}

func (s *UserService) CreateUser(name, email string) (*User, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	user := &User{Name: name, Email: email}
	if err := s.repo.Save(user); err != nil {
		return nil, fmt.Errorf("saving user: %w", err)
	}
	return user, nil
}

func main() {
	fmt.Println("=== 第14章：测试 ===")
	fmt.Println(`
本章的代码在 main.go 和 main_test.go 中

运行测试：
  cd 14-testing
  go test -v                  # 详细输出
  go test -run TestAdd        # 运行特定测试
  go test -cover              # 测试覆盖率
  go test -coverprofile=c.out # 生成覆盖率文件
  go tool cover -html=c.out   # 浏览器查看覆盖率
  go test -bench=.            # 运行基准测试
  go test -bench=. -benchmem  # 基准测试+内存分配
  go test -race               # 竞态检测
  go test -count=1            # 禁用缓存
  go test ./...               # 测试所有包

测试文件规则：
  1. 文件名以 _test.go 结尾
  2. 函数名以 Test 开头，接受 *testing.T
  3. 基准测试以 Benchmark 开头，接受 *testing.B
  4. 示例以 Example 开头
  5. 模糊测试以 Fuzz 开头，接受 *testing.F (Go 1.18+)
  6. TestMain 可以控制测试的启动和退出

请查看 main_test.go 获取完整测试示例！`)
}
