// ============================================================================
// 第12章：并发编程 🔴 高级 | ⭐ 必须学习
// ============================================================================
// 本章内容：
// 1. Goroutine
// 2. Channel（无缓冲/有缓冲）
// 3. select 语句
// 4. sync 包（Mutex, RWMutex, WaitGroup, Once）
// 5. context 包
// 6. 并发模式（Fan-in, Fan-out, Pipeline, Worker Pool）
// 7. 竞态检测
// ============================================================================

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// ========================================================================
	// 1. Goroutine 基础
	// ========================================================================
	fmt.Println("=== 1. Goroutine ===")

	// 📚 概念：Goroutine
	// - Go 的轻量级协程，由 Go 运行时调度
	// - 初始栈仅 ~2KB（线程通常 1-8MB）
	// - 可以轻松创建数十万个
	// - 用 go 关键字启动

	go func() {
		fmt.Println("  [goroutine] Hello from goroutine!")
	}()

	// 主 goroutine 需要等待，否则程序直接退出
	time.Sleep(10 * time.Millisecond)

	// 用 WaitGroup 等待
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("  [goroutine %d] 执行中\n", id)
		}(i) // ⚠️ 注意：将 i 作为参数传入，避免闭包捕获问题
	}
	wg.Wait()
	fmt.Println("  所有 goroutine 完成")

	// ========================================================================
	// 2. Channel 基础
	// ========================================================================
	fmt.Println("\n=== 2. Channel ===")

	// 📚 概念：Channel
	// - goroutine 之间通信的管道
	// - 类型安全的
	// - "Don't communicate by sharing memory; share memory by communicating."
	// - 无缓冲 channel：发送和接收必须同时就绪（同步）
	// - 有缓冲 channel：可以暂存数据

	// 无缓冲 channel
	ch := make(chan string)
	go func() {
		ch <- "Hello from channel!" // 发送
	}()
	msg := <-ch // 接收（阻塞直到有数据）
	fmt.Println("  收到:", msg)

	// 有缓冲 channel
	buffered := make(chan int, 3) // 缓冲大小为 3
	buffered <- 1
	buffered <- 2
	buffered <- 3
	// buffered <- 4 // 缓冲满了会阻塞
	fmt.Printf("  缓冲 channel: %d, %d, %d\n", <-buffered, <-buffered, <-buffered)

	// Channel 方向
	// chan<- int  只能发送
	// <-chan int  只能接收
	producer := func(ch chan<- int) {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch) // 关闭 channel，通知接收方没有更多数据
	}

	consumer := func(ch <-chan int) {
		for v := range ch { // range 自动在 channel 关闭时停止
			fmt.Printf("  消费: %d\n", v)
		}
	}

	dataCh := make(chan int, 5)
	go producer(dataCh)
	consumer(dataCh)

	// ========================================================================
	// 3. select 语句
	// ========================================================================
	fmt.Println("\n=== 3. select ===")

	// 📚 概念：select
	// - 同时等待多个 channel 操作
	// - 类似 switch，但 case 是 channel 操作
	// - 如果多个 case 就绪，随机选择一个
	// - default 使 select 非阻塞

	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	go func() {
		time.Sleep(5 * time.Millisecond)
		ch1 <- "来自 channel 1"
	}()
	go func() {
		time.Sleep(2 * time.Millisecond)
		ch2 <- "来自 channel 2"
	}()

	// 等待第一个结果
	select {
	case msg1 := <-ch1:
		fmt.Println("  收到:", msg1)
	case msg2 := <-ch2:
		fmt.Println("  收到:", msg2)
	}

	// 超时模式
	fmt.Println("\n  超时模式:")
	slowCh := make(chan string)
	go func() {
		time.Sleep(100 * time.Millisecond)
		slowCh <- "慢速结果"
	}()

	select {
	case result := <-slowCh:
		fmt.Println("  结果:", result)
	case <-time.After(50 * time.Millisecond):
		fmt.Println("  超时了！")
	}

	// 非阻塞操作
	fmt.Println("\n  非阻塞模式:")
	nonBlock := make(chan int, 1)
	select {
	case v := <-nonBlock:
		fmt.Println("  收到:", v)
	default:
		fmt.Println("  channel 为空，不阻塞")
	}

	// ========================================================================
	// 4. sync 包
	// ========================================================================
	fmt.Println("\n=== 4. sync 包 ===")

	// --- Mutex ---
	fmt.Println("\n--- Mutex ---")
	// 保护共享资源的互斥锁
	var (
		mu      sync.Mutex
		counter int
	)

	var wg2 sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg2.Wait()
	fmt.Printf("  Mutex 计数器: %d (期望 1000)\n", counter)

	// --- RWMutex ---
	fmt.Println("\n--- RWMutex ---")
	// 读写锁：多读单写
	var rwmu sync.RWMutex
	data := map[string]int{"a": 1, "b": 2}

	// 读操作可以并行
	var wg3 sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg3.Add(1)
		go func(id int) {
			defer wg3.Done()
			rwmu.RLock()
			defer rwmu.RUnlock()
			fmt.Printf("  [读者%d] data = %v\n", id, data)
		}(i)
	}
	wg3.Wait()

	// --- Once ---
	fmt.Println("\n--- Once ---")
	var once sync.Once
	initFunc := func() {
		fmt.Println("  初始化（只执行一次）")
	}
	for i := 0; i < 3; i++ {
		once.Do(initFunc) // 只有第一次调用会执行
	}

	// --- atomic ---
	fmt.Println("\n--- atomic ---")
	var atomicCounter int64
	var wg4 sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg4.Add(1)
		go func() {
			defer wg4.Done()
			atomic.AddInt64(&atomicCounter, 1)
		}()
	}
	wg4.Wait()
	fmt.Printf("  Atomic 计数器: %d (期望 1000)\n", atomicCounter)

	// ========================================================================
	// 5. context 包
	// ========================================================================
	fmt.Println("\n=== 5. context ===")

	// 📚 概念：context
	// - 用于跨 goroutine 传递截止时间、取消信号、请求级别的值
	// - 所有 IO 操作/长时间运行的操作都应该接受 context
	// - context 形成一棵树：父 context 取消时，所有子 context 也会取消

	// WithCancel
	fmt.Println("\n--- context.WithCancel ---")
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("  [worker] 收到取消信号，退出")
				return
			default:
				fmt.Println("  [worker] 工作中...")
				time.Sleep(10 * time.Millisecond)
			}
		}
	}(ctx)

	time.Sleep(35 * time.Millisecond)
	cancel() // 发送取消信号
	time.Sleep(15 * time.Millisecond)

	// WithTimeout
	fmt.Println("\n--- context.WithTimeout ---")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2()

	result2, err := longRunningTask(ctx2)
	if err != nil {
		fmt.Println("  超时:", err)
	} else {
		fmt.Println("  结果:", result2)
	}

	// WithValue (谨慎使用)
	fmt.Println("\n--- context.WithValue ---")
	const requestIDKey contextKey = "requestID"
	ctx3 := context.WithValue(context.Background(), requestIDKey, "req-12345")
	processRequest(ctx3, requestIDKey)

	// ========================================================================
	// 6. 并发模式
	// ========================================================================
	fmt.Println("\n=== 6. 并发模式 ===")

	// --- Pipeline ---
	fmt.Println("\n--- Pipeline 模式 ---")
	// 数据处理流水线：生成 → 加倍 → 过滤
	nums := generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	doubled := multiply(nums, 2)
	filtered := filterChan(doubled, func(n int) bool { return n > 10 })

	fmt.Print("  Pipeline 结果: ")
	for v := range filtered {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// --- Fan-out / Fan-in ---
	fmt.Println("\n--- Fan-out/Fan-in ---")
	jobs := generate(1, 2, 3, 4, 5, 6, 7, 8)

	// Fan-out: 多个 worker 处理同一来源
	worker1 := squareWorker(jobs)
	worker2 := squareWorker(jobs)
	worker3 := squareWorker(jobs)

	// Fan-in: 合并多个来源
	merged := fanIn(worker1, worker2, worker3)
	fmt.Print("  Fan-in 结果: ")
	for v := range merged {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// --- Worker Pool ---
	fmt.Println("\n--- Worker Pool ---")
	const numJobs = 8
	const numWorkers = 3

	jobsCh := make(chan int, numJobs)
	resultsCh := make(chan string, numJobs)

	// 启动 worker
	var wgPool sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wgPool.Add(1)
		go worker(w, jobsCh, resultsCh, &wgPool)
	}

	// 发送任务
	for j := 1; j <= numJobs; j++ {
		jobsCh <- j
	}
	close(jobsCh)

	// 收集结果
	go func() {
		wgPool.Wait()
		close(resultsCh)
	}()

	for result := range resultsCh {
		fmt.Println(" ", result)
	}

	// --- 速率限制器 ---
	fmt.Println("\n--- Rate Limiter ---")
	rateLimiter := time.NewTicker(20 * time.Millisecond)
	defer rateLimiter.Stop()

	requests := []int{1, 2, 3, 4, 5}
	for _, req := range requests {
		<-rateLimiter.C // 每 20ms 放行一个请求
		fmt.Printf("  处理请求 %d at %s\n", req, time.Now().Format("15:04:05.000"))
	}

	// ========================================================================
	// 7. 竞态检测
	// ========================================================================
	fmt.Println("\n=== 7. 竞态检测 ===")
	fmt.Println(`
  使用 -race 标志检测数据竞争：
    go run -race main.go
    go test -race ./...
    go build -race

  竞态检测会在运行时检查：
  - 多个 goroutine 同时访问同一变量
  - 至少一个是写操作
  - 没有同步机制保护

  修复方式：
  1. 使用 Mutex/RWMutex
  2. 使用 Channel
  3. 使用 atomic 包
  4. 重新设计避免共享状态`)

	fmt.Println(`
📚 并发编程关键总结：
1. "Do not communicate by sharing memory; share memory by communicating"
2. Goroutine 轻量，但也要控制数量
3. 始终确保 channel 被关闭（由发送方关闭）
4. 使用 context 传递取消信号和超时
5. 使用 -race 检测数据竞争
6. WaitGroup 等待一组 goroutine 完成
7. select 同时处理多个 channel
8. Pipeline 和 Fan-in/Fan-out 是常用模式`)
}

// ============================================================================
// 辅助函数
// ============================================================================

func longRunningTask(ctx context.Context) (string, error) {
	select {
	case <-time.After(100 * time.Millisecond): // 模拟耗时任务
		return "完成", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

type contextKey string

func processRequest(ctx context.Context, key contextKey) {
	if reqID, ok := ctx.Value(key).(string); ok {
		fmt.Printf("  处理请求 ID: %s\n", reqID)
	}
}

// Pipeline 辅助函数
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func multiply(in <-chan int, factor int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * factor
		}
		close(out)
	}()
	return out
}

func filterChan(in <-chan int, fn func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if fn(n) {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

func squareWorker(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func fanIn(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	merged := make(chan int)

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				merged <- v
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func worker(id int, jobs <-chan int, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		// 模拟工作
		duration := time.Duration(rand.Intn(50)) * time.Millisecond
		time.Sleep(duration)
		results <- fmt.Sprintf("Worker %d 完成任务 %d (耗时 %v)", id, j, duration)
	}
}

// ============================================================================
// 💻 运行方式：
//   cd 12-concurrency && go run main.go
//   cd 12-concurrency && go run -race main.go  # 竞态检测
//
// 📝 练习：
// 1. 实现一个并发网页爬虫（控制并发数）
// 2. 实现一个并发文件搜索器
// 3. 用 Pipeline 模式处理日志文件
// 4. 实现生产者-消费者模式
// ============================================================================
