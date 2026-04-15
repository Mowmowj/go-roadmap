// ============================================================================
// 第14章：测试 🟡 中级 | ⭐ 必须学习
// ============================================================================
// 测试文件：展示 Go 测试的各种模式和最佳实践
// ============================================================================

package main

import (
	"errors"
	"fmt"
	"testing"
)

// ============================================================================
// 1. 基础测试
// ============================================================================

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	if result != 5 {
		t.Errorf("Add(2, 3) = %f; want 5", result)
	}
}

func TestSubtract(t *testing.T) {
	result := Subtract(5, 3)
	if result != 2 {
		t.Errorf("Subtract(5, 3) = %f; want 2", result)
	}
}

// ============================================================================
// 2. 表驱动测试（Table-Driven Tests）— Go 最推荐的测试模式
// ============================================================================

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr bool
	}{
		{"正常除法", 10, 2, 5, false},
		{"除以零", 10, 0, 0, true},
		{"负数", -10, 2, -5, false},
		{"零除正", 0, 5, 0, false},
		{"小数", 7, 2, 3.5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Divide(%v, %v) error = %v, wantErr %v",
					tt.a, tt.b, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Divide(%v, %v) = %v, want %v",
					tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestIsPrime(t *testing.T) {
	tests := []struct {
		input int
		want  bool
	}{
		{-1, false},
		{0, false},
		{1, false},
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{10, false},
		{13, true},
		{97, true},
		{100, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("IsPrime(%d)", tt.input), func(t *testing.T) {
			if got := IsPrime(tt.input); got != tt.want {
				t.Errorf("IsPrime(%d) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestFibonacci(t *testing.T) {
	tests := map[string]struct {
		input int
		want  int
	}{
		"zero":     {0, 0},
		"one":      {1, 1},
		"two":      {2, 1},
		"five":     {5, 5},
		"ten":      {10, 55},
		"negative": {-1, 0},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := Fibonacci(tt.input); got != tt.want {
				t.Errorf("Fibonacci(%d) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// ============================================================================
// 3. 字符串测试
// ============================================================================

func TestReverse(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"hello", "olleh"},
		{"", ""},
		{"a", "a"},
		{"ab", "ba"},
		{"你好世界", "界世好你"}, // Unicode 支持
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Reverse(tt.input); got != tt.want {
				t.Errorf("Reverse(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"racecar", true},
		{"hello", false},
		{"Madam", true}, // 忽略大小写
		{"", true},
		{"a", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsPalindrome(tt.input); got != tt.want {
				t.Errorf("IsPalindrome(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestWordCount(t *testing.T) {
	input := "hello world hello Go world hello"
	got := WordCount(input)

	expected := map[string]int{
		"hello": 3,
		"world": 2,
		"go":    1,
	}

	for word, count := range expected {
		if got[word] != count {
			t.Errorf("WordCount: %q count = %d, want %d", word, got[word], count)
		}
	}
}

// ============================================================================
// 4. Mock 测试
// ============================================================================

// Mock UserRepository
type MockUserRepo struct {
	users map[int]*User
	err   error
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{
		users: map[int]*User{
			1: {ID: 1, Name: "Alice", Email: "alice@go.dev"},
			2: {ID: 2, Name: "Bob", Email: "bob@go.dev"},
		},
	}
}

func (m *MockUserRepo) FindByID(id int) (*User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, exists := m.users[id]
	if !exists {
		return nil, fmt.Errorf("user %d not found", id)
	}
	return user, nil
}

func (m *MockUserRepo) Save(user *User) error {
	if m.err != nil {
		return m.err
	}
	user.ID = len(m.users) + 1
	m.users[user.ID] = user
	return nil
}

func TestUserService_GetUser(t *testing.T) {
	repo := NewMockUserRepo()
	svc := NewUserService(repo)

	t.Run("找到用户", func(t *testing.T) {
		user, err := svc.GetUser(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.Name != "Alice" {
			t.Errorf("name = %q, want %q", user.Name, "Alice")
		}
	})

	t.Run("无效ID", func(t *testing.T) {
		_, err := svc.GetUser(0)
		if err == nil {
			t.Error("expected error for invalid id")
		}
	})

	t.Run("用户不存在", func(t *testing.T) {
		_, err := svc.GetUser(999)
		if err == nil {
			t.Error("expected error for non-existent user")
		}
	})
}

func TestUserService_CreateUser(t *testing.T) {
	t.Run("创建成功", func(t *testing.T) {
		repo := NewMockUserRepo()
		svc := NewUserService(repo)

		user, err := svc.CreateUser("Carol", "carol@go.dev")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.Name != "Carol" {
			t.Errorf("name = %q, want %q", user.Name, "Carol")
		}
	})

	t.Run("名字为空", func(t *testing.T) {
		repo := NewMockUserRepo()
		svc := NewUserService(repo)

		_, err := svc.CreateUser("", "email@go.dev")
		if err == nil {
			t.Error("expected error for empty name")
		}
	})

	t.Run("数据库错误", func(t *testing.T) {
		repo := NewMockUserRepo()
		repo.err = errors.New("db connection failed")
		svc := NewUserService(repo)

		_, err := svc.CreateUser("Test", "test@go.dev")
		if err == nil {
			t.Error("expected error for db failure")
		}
	})
}

// ============================================================================
// 5. 测试辅助函数 (Test Helper)
// ============================================================================

func assertNoError(t *testing.T, err error) {
	t.Helper() // 标记为辅助函数，错误行号会指向调用方
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestWithHelpers(t *testing.T) {
	result, err := Divide(10, 2)
	assertNoError(t, err)
	assertEqual(t, result, 5.0)
}

// ============================================================================
// 6. 基准测试 (Benchmark)
// ============================================================================

func BenchmarkFibonacci10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(10)
	}
}

func BenchmarkFibonacci20(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(20)
	}
}

func BenchmarkIsPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPrime(7919)
	}
}

func BenchmarkReverse(b *testing.B) {
	s := "Hello, World! 你好世界"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Reverse(s)
	}
}

// ============================================================================
// 7. Example 测试（同时作为文档和测试）
// ============================================================================

func ExampleAdd() {
	fmt.Println(Add(2, 3))
	// Output: 5
}

func ExampleReverse() {
	fmt.Println(Reverse("hello"))
	// Output: olleh
}

func ExampleIsPalindrome() {
	fmt.Println(IsPalindrome("racecar"))
	fmt.Println(IsPalindrome("hello"))
	// Output:
	// true
	// false
}

// ============================================================================
// 📚 测试最佳实践总结：
//
// 1. 使用表驱动测试处理多个用例
// 2. 使用子测试 t.Run() 组织测试
// 3. 使用 t.Helper() 标记辅助函数
// 4. 使用 t.Parallel() 并行运行独立测试
// 5. 测试失败用 t.Error/t.Errorf（继续执行）
//    致命失败用 t.Fatal/t.Fatalf（立即停止）
// 6. 通过接口 mock 外部依赖
// 7. 基准测试发现性能问题
// 8. Example 测试兼作文档
// 9. 目标覆盖率 80%+，但不要追求 100%
// 10. go test -race 检测数据竞争
// ============================================================================

// ============================================================================
// 💻 运行方式：
//   cd 14-testing
//   go test -v                      # 运行所有测试（详细）
//   go test -run TestDivide         # 运行特定测试
//   go test -run TestDivide/正常除法 # 运行特定子测试
//   go test -bench=.                # 基准测试
//   go test -bench=. -benchmem      # 基准测试+内存
//   go test -cover                  # 覆盖率
//   go test -v -count=1             # 跳过缓存
// ============================================================================
