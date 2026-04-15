// ============================================================================
// 第18章：PostgreSQL 数据库实战 🟡 中级 | 📌 推荐学习
// ============================================================================
// 本章内容：
// 1. database/sql 标准接口 + lib/pq 驱动
// 2. 连接管理（连接池配置）
// 3. 建表 / DDL 操作
// 4. CRUD — 增删改查
// 5. 事务（Transaction）
// 6. 预处理语句（Prepared Statement）
// 7. 批量操作
// 8. NULL 值处理（sql.NullString 等）
// 9. JSON/JSONB 字段
// 10. 错误处理最佳实践
// 11. 数据库迁移思路
// 12. Repository 模式封装
//
// 前置条件：
//   需要一个可用的 PostgreSQL 实例。
//   最简单的方式是用 Docker:
//     docker run --name go-pg -e POSTGRES_PASSWORD=postgres \
//       -e POSTGRES_DB=go_roadmap -p 5432:5432 -d postgres:16-alpine
//
// 运行方式：
//   cd 18-database && go run main.go
//   # 也可指定连接字符串：
//   cd 18-database && DATABASE_URL="postgres://user:pass@host:5432/db?sslmode=disable" go run main.go
// ============================================================================

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq" // PostgreSQL 驱动（通过 blank import 注册）
)

// ============================================================================
// 数据模型
// ============================================================================

// User 用户表模型
type User struct {
	ID        int64
	Name      string
	Email     string
	Age       int
	Bio       sql.NullString // 可为 NULL 的字段
	Tags      []string       // 存储为 PostgreSQL JSONB
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Post 文章表模型 — 演示外键关联
type Post struct {
	ID        int64
	UserID    int64
	Title     string
	Content   string
	Published bool
	Metadata  map[string]interface{} // JSONB 字段
	CreatedAt time.Time
}

// ============================================================================
// 1. 数据库连接
// ============================================================================

// connectDB 建立数据库连接并配置连接池
// 📌 database/sql 内置连接池，不需要第三方连接池库
func connectDB() *sql.DB {
	// 优先从环境变量读取，否则使用默认值
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/go_roadmap?sslmode=disable"
	}

	fmt.Println("📌 连接数据库...")
	fmt.Printf("   DSN: %s\n", maskPassword(dsn))

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ sql.Open 失败: %v", err)
	}

	// ========================================================================
	// 连接池配置 — 生产环境必须设置！
	// ========================================================================
	db.SetMaxOpenConns(25)                 // 最大打开连接数
	db.SetMaxIdleConns(5)                  // 最大空闲连接数
	db.SetConnMaxLifetime(5 * time.Minute) // 连接最大存活时间
	db.SetConnMaxIdleTime(1 * time.Minute) // 空闲连接最大存活时间

	fmt.Println(`
  📚 连接池参数说明：
  ┌─────────────────────┬───────┬────────────────────────────┐
  │ 参数                │ 建议值│ 说明                       │
  ├─────────────────────┼───────┼────────────────────────────┤
  │ MaxOpenConns        │ 25    │ 取决于 PG max_connections  │
  │ MaxIdleConns        │ 5     │ 通常为 MaxOpen 的 20%      │
  │ ConnMaxLifetime     │ 5min  │ 防止用到已失效的连接       │
  │ ConnMaxIdleTime     │ 1min  │ 及时释放空闲连接           │
  └─────────────────────┴───────┴────────────────────────────┘`)

	// Ping 验证连接是否真的可用
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ 无法连接数据库: %v\n请确保 PostgreSQL 已启动。\n最简单的方式:\n  docker run --name go-pg -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=go_roadmap -p 5432:5432 -d postgres:16-alpine", err)
	}
	fmt.Println("✅ 数据库连接成功!")

	return db
}

// maskPassword 隐藏 DSN 中的密码
func maskPassword(dsn string) string {
	// 简单处理: postgres://user:password@host -> postgres://user:****@host
	if idx := strings.Index(dsn, "://"); idx != -1 {
		rest := dsn[idx+3:]
		if atIdx := strings.Index(rest, "@"); atIdx != -1 {
			if colonIdx := strings.Index(rest[:atIdx], ":"); colonIdx != -1 {
				return dsn[:idx+3] + rest[:colonIdx+1] + "****" + rest[atIdx:]
			}
		}
	}
	return dsn
}

// ============================================================================
// 2. 建表 (DDL)
// ============================================================================

func createTables(db *sql.DB) {
	fmt.Println("\n=== 2. 建表 (DDL) ===")

	// 📌 使用 IF NOT EXISTS 保证幂等性（可重复执行）
	schema := `
	-- 启用 uuid 扩展（可选）
	-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	-- 用户表
	CREATE TABLE IF NOT EXISTS users (
		id         BIGSERIAL PRIMARY KEY,
		name       VARCHAR(100) NOT NULL,
		email      VARCHAR(255) NOT NULL UNIQUE,
		age        INTEGER DEFAULT 0,
		bio        TEXT,                          -- 可为 NULL
		tags       JSONB DEFAULT '[]'::jsonb,     -- JSONB 数组
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	-- 文章表（外键关联用户）
	CREATE TABLE IF NOT EXISTS posts (
		id         BIGSERIAL PRIMARY KEY,
		user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		title      VARCHAR(500) NOT NULL,
		content    TEXT NOT NULL DEFAULT '',
		published  BOOLEAN DEFAULT false,
		metadata   JSONB DEFAULT '{}'::jsonb,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	-- 索引 — 查询优化的关键
	CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
	CREATE INDEX IF NOT EXISTS idx_users_email   ON users(email);
	-- GIN 索引用于 JSONB 查询
	CREATE INDEX IF NOT EXISTS idx_users_tags    ON users USING GIN(tags);
	`

	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("❌ 建表失败: %v", err)
	}

	fmt.Println("✅ 表创建成功 (users, posts)")
	fmt.Println(`
  📚 PostgreSQL 数据类型映射：
  ┌───────────────┬────────────────┬──────────────────┐
  │ PostgreSQL    │ Go 类型        │ 说明             │
  ├───────────────┼────────────────┼──────────────────┤
  │ BIGSERIAL     │ int64          │ 自增主键         │
  │ VARCHAR(n)    │ string         │ 变长字符串       │
  │ TEXT          │ sql.NullString │ 可能为 NULL      │
  │ INTEGER       │ int            │ 整数             │
  │ BOOLEAN       │ bool           │ 布尔值           │
  │ TIMESTAMPTZ   │ time.Time      │ 带时区的时间     │
  │ JSONB         │ []byte / 自定义│ JSON 二进制存储  │
  └───────────────┴────────────────┴──────────────────┘`)
}

// ============================================================================
// 3. INSERT — 插入数据
// ============================================================================

func insertUsers(db *sql.DB) []int64 {
	fmt.Println("\n=== 3. INSERT — 插入数据 ===")

	// 先清理旧数据（方便重复运行演示）
	db.Exec("DELETE FROM posts")
	db.Exec("DELETE FROM users")

	var ids []int64

	// ----- 方式1: 插入单条 + RETURNING 获取自增 ID ------
	var id int64
	// 📌 使用 $1, $2 占位符（PostgreSQL 风格），防止 SQL 注入！
	// ⚠️ 绝对不要用 fmt.Sprintf 拼接 SQL！
	err := db.QueryRow(
		`INSERT INTO users (name, email, age, bio, tags)
		 VALUES ($1, $2, $3, $4, $5::jsonb)
		 RETURNING id`,
		"张三", "zhangsan@example.com", 28,
		"Go 语言开发者，热爱开源", // bio (非 NULL)
		`["go", "postgresql", "docker"]`,
	).Scan(&id)
	if err != nil {
		log.Fatalf("❌ 插入用户失败: %v", err)
	}
	ids = append(ids, id)
	fmt.Printf("  ✅ 插入用户: 张三 (ID=%d)\n", id)

	// ----- 方式2: 使用 sql.NullString 插入 NULL 值 ------
	err = db.QueryRow(
		`INSERT INTO users (name, email, age, bio, tags)
		 VALUES ($1, $2, $3, $4, $5::jsonb)
		 RETURNING id`,
		"李四", "lisi@example.com", 32,
		sql.NullString{Valid: false}, // bio = NULL
		`["python", "data"]`,
	).Scan(&id)
	if err != nil {
		log.Fatalf("❌ 插入用户失败: %v", err)
	}
	ids = append(ids, id)
	fmt.Printf("  ✅ 插入用户: 李四 (ID=%d, bio=NULL)\n", id)

	// ----- 方式3: 批量插入（一条 SQL，高性能）------
	fmt.Println("\n  📌 批量插入:")
	batchUsers := []struct {
		Name  string
		Email string
		Age   int
	}{
		{"王五", "wangwu@example.com", 25},
		{"赵六", "zhaoliu@example.com", 30},
		{"孙七", "sunqi@example.com", 22},
	}

	// 动态构建批量 INSERT
	valueStrings := make([]string, 0, len(batchUsers))
	valueArgs := make([]interface{}, 0, len(batchUsers)*3)
	for i, u := range batchUsers {
		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		valueArgs = append(valueArgs, u.Name, u.Email, u.Age)
	}

	query := fmt.Sprintf(
		`INSERT INTO users (name, email, age) VALUES %s RETURNING id`,
		strings.Join(valueStrings, ", "),
	)

	rows, err := db.Query(query, valueArgs...)
	if err != nil {
		log.Fatalf("❌ 批量插入失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var batchID int64
		rows.Scan(&batchID)
		ids = append(ids, batchID)
		fmt.Printf("    ✅ 批量插入 ID=%d\n", batchID)
	}

	fmt.Println(`
  📚 INSERT 要点：
  1. 始终使用 $1, $2 参数化查询 — 防止 SQL 注入
  2. 用 RETURNING id 获取自增 ID — 比 LastInsertId 更可靠
  3. 批量插入比逐条插入快 10-100 倍
  4. 可为 NULL 的字段用 sql.NullString / sql.NullInt64`)

	return ids
}

// ============================================================================
// 4. SELECT — 查询数据
// ============================================================================

func queryUsers(db *sql.DB) {
	fmt.Println("\n=== 4. SELECT — 查询数据 ===")

	// ----- 查询单行: QueryRow -----
	fmt.Println("  --- QueryRow: 查询单条 ---")
	var user User
	var tagsJSON []byte
	err := db.QueryRow(
		`SELECT id, name, email, age, bio, tags, created_at, updated_at
		 FROM users WHERE email = $1`, "zhangsan@example.com",
	).Scan(
		&user.ID, &user.Name, &user.Email, &user.Age,
		&user.Bio,      // sql.NullString 自动处理 NULL
		&tagsJSON,       // JSONB 读为 []byte
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		fmt.Println("  用户不存在")
	} else if err != nil {
		log.Fatalf("❌ 查询失败: %v", err)
	} else {
		json.Unmarshal(tagsJSON, &user.Tags)
		fmt.Printf("  ✅ 找到用户: ID=%d, Name=%s, Email=%s, Age=%d\n",
			user.ID, user.Name, user.Email, user.Age)
		if user.Bio.Valid {
			fmt.Printf("     Bio: %s\n", user.Bio.String)
		} else {
			fmt.Println("     Bio: <NULL>")
		}
		fmt.Printf("     Tags: %v\n", user.Tags)
	}

	// ----- 查询多行: Query -----
	fmt.Println("\n  --- Query: 查询多条 ---")
	rows, err := db.Query(
		`SELECT id, name, email, age, bio FROM users ORDER BY id`)
	if err != nil {
		log.Fatalf("❌ 查询失败: %v", err)
	}
	defer rows.Close() // 📌 必须关闭！否则连接泄漏

	fmt.Println("  所有用户:")
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.Bio); err != nil {
			log.Fatalf("❌ Scan 失败: %v", err)
		}
		bio := "<NULL>"
		if u.Bio.Valid {
			bio = u.Bio.String
		}
		fmt.Printf("    [%d] %s (%s) age=%d bio=%s\n",
			u.ID, u.Name, u.Email, u.Age, bio)
	}
	// 📌 检查迭代过程中是否发生错误
	if err := rows.Err(); err != nil {
		log.Fatalf("❌ rows 迭代错误: %v", err)
	}

	// ----- 条件查询 + LIKE + LIMIT -----
	fmt.Println("\n  --- 条件查询 ---")
	rows2, err := db.Query(
		`SELECT name, email FROM users WHERE age >= $1 ORDER BY age DESC LIMIT $2`,
		25, 3,
	)
	if err != nil {
		log.Fatalf("❌ 查询失败: %v", err)
	}
	defer rows2.Close()

	fmt.Println("  年龄 >= 25 的用户 (最多3条):")
	for rows2.Next() {
		var name, email string
		rows2.Scan(&name, &email)
		fmt.Printf("    %s (%s)\n", name, email)
	}

	fmt.Println(`
  📚 SELECT 要点：
  1. QueryRow — 查询单行，不存在返回 sql.ErrNoRows
  2. Query — 查询多行，必须 defer rows.Close()
  3. 每次 Scan 后检查 rows.Err()
  4. JSONB 用 []byte 接收，再 json.Unmarshal
  5. NULL 字段用 sql.NullString / sql.NullInt64`)
}

// ============================================================================
// 5. UPDATE — 更新数据
// ============================================================================

func updateUsers(db *sql.DB) {
	fmt.Println("\n=== 5. UPDATE — 更新数据 ===")

	// 更新单条记录
	result, err := db.Exec(
		`UPDATE users SET age = $1, bio = $2, updated_at = NOW()
		 WHERE email = $3`,
		29, "资深 Go 开发者", "zhangsan@example.com",
	)
	if err != nil {
		log.Fatalf("❌ 更新失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("  ✅ 更新了 %d 行 (张三 age -> 29)\n", rowsAffected)

	// 更新 JSONB 字段 — 追加 tag
	_, err = db.Exec(
		`UPDATE users SET tags = tags || $1::jsonb
		 WHERE email = $2`,
		`["kubernetes"]`, "zhangsan@example.com",
	)
	if err != nil {
		log.Fatalf("❌ 更新 JSONB 失败: %v", err)
	}
	fmt.Println("  ✅ 追加 tag: kubernetes")

	// 条件批量更新
	result, err = db.Exec(
		`UPDATE users SET bio = '新用户' WHERE bio IS NULL`)
	if err != nil {
		log.Fatalf("❌ 批量更新失败: %v", err)
	}
	affected, _ := result.RowsAffected()
	fmt.Printf("  ✅ 批量更新 bio IS NULL 的用户: %d 行\n", affected)

	fmt.Println(`
  📚 UPDATE 要点：
  1. 始终检查 RowsAffected() 确认是否更新成功
  2. JSONB 可用 || 追加, - 删除, jsonb_set() 修改嵌套
  3. updated_at = NOW() 记录更新时间
  4. 批量更新注意 WHERE 条件，避免全表更新`)
}

// ============================================================================
// 6. DELETE — 删除数据
// ============================================================================

func deleteUsers(db *sql.DB) {
	fmt.Println("\n=== 6. DELETE — 删除数据 ===")

	// 删除指定用户
	result, err := db.Exec(
		`DELETE FROM users WHERE email = $1`, "sunqi@example.com")
	if err != nil {
		log.Fatalf("❌ 删除失败: %v", err)
	}
	affected, _ := result.RowsAffected()
	fmt.Printf("  ✅ 删除了 %d 行 (孙七)\n", affected)

	// 验证
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	fmt.Printf("  剩余用户数: %d\n", count)

	fmt.Println(`
  📚 DELETE 要点：
  1. 始终带 WHERE 条件！不带 WHERE 会删除全表
  2. 先 SELECT 确认再 DELETE（生产环境）
  3. 软删除（推荐）：添加 deleted_at 字段，而不是真正删除
  4. CASCADE: 外键设置 ON DELETE CASCADE 自动删除关联行`)
}

// ============================================================================
// 7. 事务（Transaction）
// ============================================================================

func transactionDemo(db *sql.DB) {
	fmt.Println("\n=== 7. 事务 (Transaction) ===")

	// 📌 事务保证一组操作要么全部成功，要么全部回滚
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("❌ 开始事务失败: %v", err)
	}

	// 📌 defer + 错误处理模式 — 保证异常时回滚
	defer func() {
		if err != nil {
			tx.Rollback()
			fmt.Println("  ⚠️  事务已回滚")
		}
	}()

	// 事务中的操作1: 创建文章
	var postID int64
	err = tx.QueryRow(
		`INSERT INTO posts (user_id, title, content, published, metadata)
		 VALUES (
			(SELECT id FROM users WHERE email = $1),
			$2, $3, $4, $5::jsonb
		 ) RETURNING id`,
		"zhangsan@example.com",
		"Go 语言数据库编程指南",
		"本文介绍如何在 Go 中操作 PostgreSQL...",
		true,
		`{"category": "tutorial", "views": 0}`,
	).Scan(&postID)
	if err != nil {
		log.Printf("  ❌ 插入文章失败: %v", err)
		return
	}
	fmt.Printf("  ✅ 事务中插入文章 ID=%d\n", postID)

	// 事务中的操作2: 更新用户
	_, err = tx.Exec(
		`UPDATE users SET bio = $1, updated_at = NOW() WHERE email = $2`,
		"Go 开发者 & 技术博主", "zhangsan@example.com",
	)
	if err != nil {
		log.Printf("  ❌ 更新用户失败: %v", err)
		return
	}
	fmt.Println("  ✅ 事务中更新用户 bio")

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Fatalf("❌ 提交事务失败: %v", err)
	}
	fmt.Println("  ✅ 事务提交成功!")

	fmt.Println(`
  📚 事务要点：
  1. Begin() 开始 → 操作 → Commit() 提交
  2. 出错时 Rollback() 回滚
  3. defer + Rollback 模式保证异常安全
  4. 事务中用 tx.Query/tx.Exec，不要用 db.Query/db.Exec
  5. PostgreSQL 事务隔离级别:
     READ COMMITTED  (默认) — 防止脏读
     REPEATABLE READ         — 防止不可重复读
     SERIALIZABLE            — 最严格，防止幻读`)
}

// ============================================================================
// 8. 预处理语句 (Prepared Statement)
// ============================================================================

func preparedStatementDemo(db *sql.DB) {
	fmt.Println("\n=== 8. 预处理语句 (Prepared Statement) ===")

	// 📌 相同 SQL 多次执行时，Prepare 可以提升性能
	// 数据库只需解析一次 SQL，后续只传参数
	stmt, err := db.Prepare(
		`SELECT id, name, email FROM users WHERE age >= $1`)
	if err != nil {
		log.Fatalf("❌ Prepare 失败: %v", err)
	}
	defer stmt.Close() // 📌 必须关闭

	// 使用相同的 stmt 执行多次查询
	ages := []int{20, 25, 30}
	for _, age := range ages {
		rows, err := stmt.Query(age)
		if err != nil {
			log.Printf("  查询失败 (age>=%d): %v", age, err)
			continue
		}

		var names []string
		for rows.Next() {
			var id int64
			var name, email string
			rows.Scan(&id, &name, &email)
			names = append(names, name)
		}
		rows.Close()

		fmt.Printf("  age >= %d: %v (%d人)\n", age, names, len(names))
	}

	fmt.Println(`
  📚 Prepared Statement 要点：
  1. 相同 SQL 多次执行时使用，节省解析开销
  2. 自动防止 SQL 注入（参数和 SQL 分离）
  3. stmt 必须 Close()，否则资源泄漏
  4. 注意：db.Query/db.Exec 内部也会自动 prepare
  5. 连接池环境下，stmt 会在多个连接上重新 prepare`)
}

// ============================================================================
// 9. 关联查询 (JOIN)
// ============================================================================

func joinQueryDemo(db *sql.DB) {
	fmt.Println("\n=== 9. 关联查询 (JOIN) ===")

	// 再插入几篇文章用于演示
	db.Exec(`INSERT INTO posts (user_id, title, content, published) VALUES
		((SELECT id FROM users WHERE email='zhangsan@example.com'), '并发编程入门', '...', true),
		((SELECT id FROM users WHERE email='lisi@example.com'), 'Python 数据分析', '...', false)
	`)

	rows, err := db.Query(`
		SELECT u.name, u.email, p.title, p.published, p.created_at
		FROM users u
		INNER JOIN posts p ON u.id = p.user_id
		ORDER BY p.created_at DESC`)
	if err != nil {
		log.Fatalf("❌ JOIN 查询失败: %v", err)
	}
	defer rows.Close()

	fmt.Println("  用户及其文章:")
	for rows.Next() {
		var name, email, title string
		var published bool
		var createdAt time.Time
		rows.Scan(&name, &email, &title, &published, &createdAt)

		status := "草稿"
		if published {
			status = "已发布"
		}
		fmt.Printf("    %s | %-24s | %s | %s\n",
			name, title, status, createdAt.Format("2006-01-02 15:04"))
	}

	// 聚合查询
	fmt.Println("\n  --- 聚合查询 ---")
	rows2, err := db.Query(`
		SELECT u.name, COUNT(p.id) as post_count
		FROM users u
		LEFT JOIN posts p ON u.id = p.user_id
		GROUP BY u.id, u.name
		HAVING COUNT(p.id) > 0
		ORDER BY post_count DESC`)
	if err != nil {
		log.Fatalf("❌ 聚合查询失败: %v", err)
	}
	defer rows2.Close()

	fmt.Println("  各用户文章数 (有文章的):")
	for rows2.Next() {
		var name string
		var count int
		rows2.Scan(&name, &count)
		fmt.Printf("    %s: %d 篇\n", name, count)
	}
}

// ============================================================================
// 10. JSONB 操作
// ============================================================================

func jsonbDemo(db *sql.DB) {
	fmt.Println("\n=== 10. JSONB 高级操作 ===")

	// 更新 JSONB 嵌套字段
	_, err := db.Exec(`
		UPDATE posts SET metadata = jsonb_set(metadata, '{views}', '100')
		WHERE title = 'Go 语言数据库编程指南'`)
	if err != nil {
		log.Printf("  更新 JSONB 字段失败: %v", err)
	}

	// 查询 JSONB 字段
	var title string
	var metadataBytes []byte
	err = db.QueryRow(`
		SELECT title, metadata FROM posts
		WHERE metadata->>'category' = $1`, "tutorial",
	).Scan(&title, &metadataBytes)

	if err == nil {
		var metadata map[string]interface{}
		json.Unmarshal(metadataBytes, &metadata)
		fmt.Printf("  ✅ 标题: %s\n", title)
		fmt.Printf("     Metadata: %v\n", metadata)
	}

	// JSONB 数组查询（查 tags 包含 'go' 的用户）
	rows, err := db.Query(`
		SELECT name, tags FROM users
		WHERE tags @> '"go"'::jsonb`)
	if err != nil {
		log.Printf("  JSONB 查询失败: %v", err)
		return
	}
	defer rows.Close()

	fmt.Println("\n  包含 'go' tag 的用户:")
	for rows.Next() {
		var name string
		var tags []byte
		rows.Scan(&name, &tags)
		fmt.Printf("    %s: %s\n", name, string(tags))
	}

	fmt.Println(`
  📚 JSONB 操作速查：
  ┌──────────────────────┬────────────────────────────┐
  │ 操作                 │ SQL                        │
  ├──────────────────────┼────────────────────────────┤
  │ 获取字段(文本)       │ data->>'key'               │
  │ 获取字段(JSON)       │ data->'key'                │
  │ 修改嵌套字段         │ jsonb_set(data,'{k}','v')  │
  │ 删除字段             │ data - 'key'               │
  │ 追加到数组           │ data || '["new"]'::jsonb   │
  │ 包含查询             │ data @> '{"k":"v"}'        │
  │ 键存在查询           │ data ? 'key'               │
  └──────────────────────┴────────────────────────────┘`)
}

// ============================================================================
// 11. Repository 模式 — 封装数据库操作
// ============================================================================

// UserRepository 用户数据访问层
// 📌 将 SQL 操作封装到 Repository，业务层不直接写 SQL
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository 创建 UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByID 按 ID 查找用户
func (r *UserRepository) FindByID(id int64) (*User, error) {
	user := &User{}
	var tagsJSON []byte
	err := r.db.QueryRow(
		`SELECT id, name, email, age, bio, tags, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Age,
		&user.Bio, &tagsJSON, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("用户 ID=%d 不存在", id)
	}
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	json.Unmarshal(tagsJSON, &user.Tags)
	return user, nil
}

// FindAll 查询所有用户（分页）
func (r *UserRepository) FindAll(limit, offset int) ([]User, error) {
	rows, err := r.db.Query(
		`SELECT id, name, email, age, bio, created_at FROM users
		 ORDER BY id LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.Bio, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan 失败: %w", err)
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// Create 创建用户
func (r *UserRepository) Create(name, email string, age int) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id`,
		name, email, age,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("创建用户失败: %w", err)
	}
	return id, nil
}

func repositoryDemo(db *sql.DB) {
	fmt.Println("\n=== 11. Repository 模式 ===")

	repo := NewUserRepository(db)

	// 分页查询
	users, err := repo.FindAll(3, 0)
	if err != nil {
		log.Printf("  ❌ %v", err)
		return
	}
	fmt.Println("  前3个用户 (Repository 查询):")
	for _, u := range users {
		fmt.Printf("    [%d] %s (%s)\n", u.ID, u.Name, u.Email)
	}

	// 按 ID 查找
	if len(users) > 0 {
		user, err := repo.FindByID(users[0].ID)
		if err != nil {
			fmt.Printf("  ❌ %v\n", err)
		} else {
			fmt.Printf("  FindByID(%d): %s, tags=%v\n", user.ID, user.Name, user.Tags)
		}
	}

	fmt.Println(`
  📚 Repository 模式要点：
  1. 每个表/实体一个 Repository
  2. Repository 封装所有 SQL — 业务层只调用方法
  3. 返回领域对象，不暴露 sql.Rows
  4. 用 interface 定义 Repository 便于 Mock 测试
  5. 结构:
     main.go
     ├── models/     数据模型 (User, Post)
     ├── repository/ 数据访问层 (UserRepo, PostRepo)
     ├── service/    业务逻辑层
     └── handler/    HTTP 处理层`)
}

// ============================================================================
// 12. 清理 & 最佳实践总结
// ============================================================================

func cleanup(db *sql.DB) {
	fmt.Println("\n=== 12. 清理 & 最佳实践 ===")

	// 演示结束后清理表（可选）
	// db.Exec("DROP TABLE IF EXISTS posts")
	// db.Exec("DROP TABLE IF EXISTS users")
	fmt.Println("  📌 表保留用于后续实验，如需清理取消上方注释")

	// 查看连接池状态
	stats := db.Stats()
	fmt.Printf("\n  📊 连接池状态:\n")
	fmt.Printf("    打开连接数:  %d\n", stats.OpenConnections)
	fmt.Printf("    使用中:      %d\n", stats.InUse)
	fmt.Printf("    空闲:        %d\n", stats.Idle)
	fmt.Printf("    等待次数:    %d\n", stats.WaitCount)
	fmt.Printf("    等待时长:    %v\n", stats.WaitDuration)

	fmt.Println(`
  ✅ PostgreSQL + Go 最佳实践总结：

  连接管理：
  • 全局共享一个 *sql.DB（它是连接池，线程安全）
  • 合理配置 MaxOpenConns / MaxIdleConns
  • 用 db.Ping() 验证连接可用性

  查询安全：
  • 始终用 $1,$2 参数化查询 — 防止 SQL 注入
  • 绝不用 fmt.Sprintf 拼接 SQL 值

  资源管理：
  • rows 用完必须 Close()（推荐 defer）
  • stmt 用完必须 Close()
  • 检查 rows.Err() 捕获迭代中的错误

  事务：
  • 多步操作用事务包裹
  • defer Rollback + Commit 模式
  • 注意事务隔离级别

  性能优化：
  • 批量 INSERT 代替逐条插入
  • 合理创建索引（常查字段、外键）
  • JSONB 用 GIN 索引加速
  • 用 EXPLAIN ANALYZE 分析慢查询

  生产建议：
  • 数据库迁移用 golang-migrate 或 goose
  • 复杂项目考虑 sqlc（类型安全）或 GORM（ORM）
  • 监控连接池指标 (db.Stats())
  • 敏感信息用环境变量，不要硬编码`)
}

// ============================================================================
// main
// ============================================================================

func main() {
	fmt.Println("╔════════════════════════════════════════════╗")
	fmt.Println("║  第18章：PostgreSQL 数据库实战             ║")
	fmt.Println("╚════════════════════════════════════════════╝")

	// 1. 连接数据库
	db := connectDB()
	defer db.Close()

	// 2. 建表
	createTables(db)

	// 3. 插入数据
	insertUsers(db)

	// 4. 查询数据
	queryUsers(db)

	// 5. 更新数据
	updateUsers(db)

	// 6. 删除数据
	deleteUsers(db)

	// 7. 事务
	transactionDemo(db)

	// 8. 预处理语句
	preparedStatementDemo(db)

	// 9. 关联查询
	joinQueryDemo(db)

	// 10. JSONB 操作
	jsonbDemo(db)

	// 11. Repository 模式
	repositoryDemo(db)

	// 12. 清理 & 总结
	cleanup(db)

	fmt.Println("\n🎉 第18章完成！你已掌握 Go + PostgreSQL 数据库开发核心技能")
	fmt.Println(`
  💻 下一步练习：
  1. 给 Post 添加评论表 (comments)，练习多表关联
  2. 实现完整的 CRUD RESTful API (结合第15章)
  3. 添加数据库迁移 (golang-migrate)
  4. 尝试用 sqlc 从 SQL 自动生成 Go 代码
  5. 用 GORM 重写本章的 Repository 层`)
}
