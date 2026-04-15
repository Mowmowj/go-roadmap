// ============================================================================
// 第15章：生态系统与流行库 🟡 中级 | 📌 可选学习
// ============================================================================
// 本章内容：
// 1. Web 开发 (net/http, Gin, Echo, Fiber)
// 2. CLI 工具 (Cobra, urfave/cli)
// 3. 数据库 (database/sql, pgx, GORM)
// 4. 日志 (log/slog, zerolog, zap)
// 5. 配置管理 (Viper)
// 6. gRPC
// 7. 实时通信 (WebSocket)
// ============================================================================
// 📝 注意：本章以代码结构和概念介绍为主
// 第三方库需要单独安装：go get <package>
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// ============================================================================
// 1. net/http — Go 标准库的 HTTP 服务
// ============================================================================

// 📚 概念：
// Go 的标准库 net/http 已经足够构建生产级 HTTP 服务
// Go 1.22+ 增强了路由（支持方法和路径参数）
// 很多框架（如 Gin）底层也是基于 net/http

// 数据模型
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// 内存存储
type TodoStore struct {
	mu     sync.RWMutex
	todos  map[int]Todo
	nextID int
}

func NewTodoStore() *TodoStore {
	return &TodoStore{
		todos:  make(map[int]Todo),
		nextID: 1,
	}
}

func (s *TodoStore) GetAll() []Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	todos := make([]Todo, 0, len(s.todos))
	for _, t := range s.todos {
		todos = append(todos, t)
	}
	return todos
}

func (s *TodoStore) Create(title string) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	todo := Todo{
		ID:        s.nextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	s.todos[s.nextID] = todo
	s.nextID++
	return todo
}

func (s *TodoStore) Toggle(id int) (Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	todo, exists := s.todos[id]
	if !exists {
		return Todo{}, false
	}
	todo.Completed = !todo.Completed
	s.todos[id] = todo
	return todo, true
}

func (s *TodoStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.todos[id]; !exists {
		return false
	}
	delete(s.todos, id)
	return true
}

// HTTP handler 辅助函数
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// 中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("=== 第15章：Go 生态系统 ===")
	fmt.Println()

	// ========================================================================
	// 展示 net/http 服务代码（可选启动）
	// ========================================================================
	fmt.Println("=== 1. net/http REST API 示例 ===")
	fmt.Println()
	fmt.Println("以下是一个完整的 Todo REST API 服务。")
	fmt.Println("取消注释 startServer() 即可启动服务。")
	fmt.Println()

	printServerCode()

	// ========================================================================
	// 各个流行框架的代码对比
	// ========================================================================
	printFrameworkComparison()

	// ========================================================================
	// 数据库
	// ========================================================================
	printDatabaseGuide()

	// ========================================================================
	// CLI 工具
	// ========================================================================
	printCLIGuide()

	// ========================================================================
	// 日志
	// ========================================================================
	printLoggingGuide()

	// ========================================================================
	// gRPC
	// ========================================================================
	printGRPCGuide()

	// 取消注释以启动 HTTP 服务：
	// startServer()
}

func startServer() {
	store := NewTodoStore()

	// 预填充数据
	store.Create("学习 Go 基础")
	store.Create("学习 Go 并发")
	store.Create("构建 REST API")

	mux := http.NewServeMux()

	// 路由（Go 1.22+ 语法）
	mux.HandleFunc("GET /api/todos", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, store.GetAll())
	})

	mux.HandleFunc("POST /api/todos", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Title string `json:"title"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		if req.Title == "" {
			writeError(w, http.StatusBadRequest, "title is required")
			return
		}
		todo := store.Create(req.Title)
		writeJSON(w, http.StatusCreated, todo)
	})

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Go Todo API",
			"version": "1.0",
		})
	})

	// 应用中间件
	handler := loggingMiddleware(corsMiddleware(mux))

	addr := ":8080"
	fmt.Printf("\n🚀 服务启动: http://localhost%s\n", addr)
	fmt.Println("API 端点:")
	fmt.Println("  GET  /api/todos   - 获取所有待办")
	fmt.Println("  POST /api/todos   - 创建待办")
	fmt.Println()
	log.Fatal(http.ListenAndServe(addr, handler))
}

func printServerCode() {
	fmt.Println(`  // 创建路由 (Go 1.22+)
  mux := http.NewServeMux()
  
  mux.HandleFunc("GET /api/todos", handleGetTodos)
  mux.HandleFunc("POST /api/todos", handleCreateTodo)
  mux.HandleFunc("PUT /api/todos/{id}", handleUpdateTodo)
  mux.HandleFunc("DELETE /api/todos/{id}", handleDeleteTodo)
  
  // 中间件链
  handler := loggingMiddleware(corsMiddleware(mux))
  
  // 启动服务
  http.ListenAndServe(":8080", handler)`)
}

func printFrameworkComparison() {
	fmt.Println("\n=== 2. Web 框架对比 ===")
	fmt.Println(`
  ┌─────────────┬──────────┬───────────────────────────────────┐
  │ 框架         │ 性能      │ 特点                              │
  ├─────────────┼──────────┼───────────────────────────────────┤
  │ net/http    │ ⭐⭐⭐⭐  │ 标准库，Go 1.22+路由增强          │
  │ Gin         │ ⭐⭐⭐⭐⭐│ 最流行，性能好，中间件生态丰富      │
  │ Echo        │ ⭐⭐⭐⭐⭐│ 高性能，API 设计清晰               │
  │ Fiber       │ ⭐⭐⭐⭐⭐│ Express 风格，基于 fasthttp        │
  │ Chi         │ ⭐⭐⭐⭐  │ 轻量，兼容 net/http               │
  └─────────────┴──────────┴───────────────────────────────────┘

  // Gin 示例:
  r := gin.Default()
  r.GET("/api/todos", getTodos)
  r.POST("/api/todos", createTodo)
  r.Run(":8080")

  // Echo 示例:
  e := echo.New()
  e.GET("/api/todos", getTodos)
  e.POST("/api/todos", createTodo)
  e.Start(":8080")

  // Fiber 示例:
  app := fiber.New()
  app.Get("/api/todos", getTodos)
  app.Post("/api/todos", createTodo)
  app.Listen(":8080")

  // Chi 示例 (兼容 net/http):
  r := chi.NewRouter()
  r.Use(middleware.Logger)
  r.Get("/api/todos", getTodos)
  r.Post("/api/todos", createTodo)
  http.ListenAndServe(":8080", r)

  安装:
  go get github.com/gin-gonic/gin
  go get github.com/labstack/echo/v4
  go get github.com/gofiber/fiber/v2
  go get github.com/go-chi/chi/v5`)
}

func printDatabaseGuide() {
	fmt.Println("\n=== 3. 数据库 ===")
	fmt.Println(`
  📚 Go 数据库生态：
  
  1. database/sql (标准库) — 通用接口
     import (
         "database/sql"
         _ "github.com/lib/pq"  // PostgreSQL 驱动
     )
     db, _ := sql.Open("postgres", "...")
     rows, _ := db.Query("SELECT id, name FROM users")

  2. pgx — PostgreSQL 高性能驱动
     go get github.com/jackc/pgx/v5
     
     conn, _ := pgx.Connect(ctx, "postgres://...")
     var name string
     conn.QueryRow(ctx, "SELECT name FROM users WHERE id=$1", 1).Scan(&name)

  3. GORM — Go 最流行的 ORM
     go get gorm.io/gorm
     go get gorm.io/driver/postgres
     
     type User struct {
         gorm.Model
         Name  string
         Email string ` + "`gorm:\"uniqueIndex\"`" + `
     }
     
     db, _ := gorm.Open(postgres.Open(dsn))
     db.AutoMigrate(&User{})
     db.Create(&User{Name: "Alice", Email: "alice@go.dev"})
     
     var user User
     db.First(&user, 1)       // 按ID查找
     db.Where("name = ?", "Alice").First(&user)

  4. sqlx — database/sql 增强版
     go get github.com/jmoiron/sqlx
     
     type User struct {
         ID    int    ` + "`db:\"id\"`" + `
         Name  string ` + "`db:\"name\"`" + `
     }
     users := []User{}
     db.Select(&users, "SELECT * FROM users")

  5. sqlc — 从 SQL 生成类型安全 Go 代码
     go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
     // 写 SQL → 自动生成 Go 代码`)
}

func printCLIGuide() {
	fmt.Println("\n=== 4. CLI 工具 ===")
	os.Stdout.WriteString(`
  📚 构建命令行工具：
  
  1. Cobra — 最流行的 CLI 框架 (kubectl, hugo, gh 都用它)
     go get github.com/spf13/cobra
     
     var rootCmd = &cobra.Command{
         Use:   "myapp",
         Short: "My CLI app",
         Run: func(cmd *cobra.Command, args []string) {
             fmt.Println("Hello from CLI!")
         },
     }
     
     var serveCmd = &cobra.Command{
         Use:   "serve",
         Short: "Start the server",
         Run: func(cmd *cobra.Command, args []string) {
             port, _ := cmd.Flags().GetInt("port")
             fmt.Printf("Serving on :%d\n", port)
         },
     }
     
     func init() {
         rootCmd.AddCommand(serveCmd)
         serveCmd.Flags().IntP("port", "p", 8080, "port number")
     }

  2. urfave/cli — 另一个流行的 CLI 框架
     go get github.com/urfave/cli/v2

  3. Bubble Tea — 终端 UI 框架（TUI）
     go get github.com/charmbracelet/bubbletea
     // 构建交互式终端界面，如进度条、选择器等
`)
}

func printLoggingGuide() {
	fmt.Println("\n=== 5. 日志 ===")
	fmt.Println(`
  📚 Go 日志生态：
  
  1. log/slog (Go 1.21+ 标准库) — 推荐首选
     logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
     logger.Info("用户登录", "user", "alice", "ip", "192.168.1.1")
     // {"time":"...","level":"INFO","msg":"用户登录","user":"alice","ip":"192.168.1.1"}

  2. zerolog — 零分配 JSON 日志
     go get github.com/rs/zerolog
     
     log := zerolog.New(os.Stdout).With().Timestamp().Logger()
     log.Info().Str("user", "alice").Msg("登录成功")

  3. zap — Uber 开发的高性能日志
     go get go.uber.org/zap
     
     logger, _ := zap.NewProduction()
     logger.Info("登录", zap.String("user", "alice"))`)
}

func printGRPCGuide() {
	fmt.Println("\n=== 6. gRPC ===")
	fmt.Println(`
  📚 gRPC — 高性能 RPC 框架
  
  安装:
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  
  // user.proto
  syntax = "proto3";
  
  service UserService {
      rpc GetUser (GetUserRequest) returns (User);
      rpc ListUsers (ListUsersRequest) returns (stream User);
  }
  
  message User {
      int32 id = 1;
      string name = 2;
      string email = 3;
  }
  
  // 生成代码
  protoc --go_out=. --go-grpc_out=. user.proto`)

	fmt.Println(`
📚 生态系统选择建议：
1. Web: 简单项目用 net/http，复杂项目用 Gin/Echo
2. 数据库: 简单查询用 pgx/sqlx，复杂映射用 GORM
3. CLI: 用 Cobra（业界标准）
4. 日志: 新项目用 slog，需要零分配用 zerolog
5. 配置: 用 Viper（cobra 配套）
6. RPC: 微服务间通信用 gRPC
7. "标准库优先" — Go 的标准库覆盖面非常广`)
}

// ============================================================================
// 💻 运行方式：
//   cd 15-ecosystem && go run main.go
//
// 🚀 启动 HTTP 服务（取消注释 main 中的 startServer()）：
//   cd 15-ecosystem && go run main.go
//   curl http://localhost:8080/api/todos
//   curl -X POST http://localhost:8080/api/todos -d '{"title":"new todo"}'
//
// 📝 练习：
// 1. 启动 Todo API 并用 curl 测试所有端点
// 2. 用 Gin 重写 Todo API
// 3. 添加 SQLite 持久化（用 GORM）
// 4. 用 Cobra 写一个 CLI 工具
// ============================================================================
