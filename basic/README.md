# Go 基础知识：最佳实践、常见问题与核心原理

## 目录
- [Slice (切片)](#slice-切片)
- [Map (映射)](#map-映射)
- [Channel (通道)](#channel-通道)

---

## Slice (切片)

### 核心原理

Slice 的底层结构：
```go
type slice struct {
    array unsafe.Pointer  // 指向底层数组的指针
    len   int             // 当前长度
    cap   int             // 容量
}
```

**关键特性：**
- Slice 是对底层数组的引用类型
- 包含三个字段：指针、长度、容量
- 多个 slice 可以共享同一个底层数组

### 最佳实践

#### 1. 预分配容量
```go
// ❌ 不好：频繁扩容
s := []int{}
for i := 0; i < 10000; i++ {
    s = append(s, i)
}

// ✅ 好：预分配容量
s := make([]int, 0, 10000)
for i := 0; i < 10000; i++ {
    s = append(s, i)
}
```

#### 2. 使用 append 的返回值
```go
// ✅ 正确
s = append(s, elem)

// ❌ 错误：忽略返回值
append(s, elem)  // 扩容后 s 不会更新
```

#### 3. 避免切片泄漏
```go
// ❌ 问题：保留了整个底层数组
func getFirstTwo(data []byte) []byte {
    return data[:2]  // 底层数组不会被 GC
}

// ✅ 解决：复制需要的数据
func getFirstTwo(data []byte) []byte {
    result := make([]byte, 2)
    copy(result, data[:2])
    return result
}
```

### 常见问题

#### 问题 1：切片共享底层数组
```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3]  // [2, 3]
b[0] = 999

fmt.Println(a)  // [1, 999, 3, 4, 5] - a 也被修改了！
```

**解决方案：** 使用 copy 创建独立副本
```go
b := make([]int, 2)
copy(b, a[1:3])
b[0] = 999
fmt.Println(a)  // [1, 2, 3, 4, 5] - a 不受影响
```

#### 问题 2：append 后的意外行为
```go
a := make([]int, 0, 5)
a = append(a, 1, 2, 3)
b := append(a, 4)
c := append(a, 5)

fmt.Println(b)  // [1, 2, 3, 5] - 不是预期的 [1, 2, 3, 4]！
```

**原因：** b 和 c 共享底层数组，c 覆盖了 b 的值

#### 问题 3：循环中的切片陷阱
```go
// ❌ 错误：所有指针指向同一个变量
var result []*int
s := []int{1, 2, 3}
for _, v := range s {
    result = append(result, &v)  // v 的地址不变
}
// 所有元素都是 3

// ✅ 正确
for i := range s {
    result = append(result, &s[i])
}
```

### 扩容机制

- 当 `len == cap` 时，append 会触发扩容
- Go 1.18+ 扩容策略：
  - cap < 256: 新容量 = 旧容量 * 2
  - cap >= 256: 新容量 ≈ 旧容量 * 1.25 + 192

---

## Map (映射)

### 核心原理

Map 的底层实现：
- 使用哈希表（hash table）
- 采用拉链法解决哈希冲突
- 由多个 bucket 组成，每个 bucket 可存储 8 个 key-value 对

**关键特性：**
- Map 是引用类型
- 非线程安全
- 遍历顺序是随机的
- key 必须是可比较类型

### 最佳实践

#### 1. 预分配容量
```go
// ❌ 不好
m := make(map[string]int)

// ✅ 好：预分配容量减少扩容
m := make(map[string]int, 1000)
```

#### 2. 检查 key 是否存在
```go
// ❌ 不好：无法区分零值和不存在
value := m[key]

// ✅ 好：使用 comma ok 模式
value, ok := m[key]
if ok {
    // key 存在
}
```

#### 3. 并发安全
```go
// ❌ 不安全：并发读写会 panic
m := make(map[string]int)
go func() { m["key"] = 1 }()
go func() { _ = m["key"] }()

// ✅ 方案 1：使用 sync.Map
var m sync.Map
m.Store("key", 1)
value, ok := m.Load("key")

// ✅ 方案 2：使用互斥锁
type SafeMap struct {
    mu sync.RWMutex
    m  map[string]int
}

func (sm *SafeMap) Set(key string, value int) {
    sm.mu.Lock()
    sm.m[key] = value
    sm.mu.Unlock()
}

func (sm *SafeMap) Get(key string) (int, bool) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    value, ok := sm.m[key]
    return value, ok
}
```

### 常见问题

#### 问题 1：Map 的零值是 nil
```go
var m map[string]int
m["key"] = 1  // panic: assignment to entry in nil map

// 必须初始化
m = make(map[string]int)
m["key"] = 1  // OK
```

#### 问题 2：Map 不能取地址
```go
type User struct {
    Name string
    Age  int
}

m := make(map[string]User)
m["alice"] = User{Name: "Alice", Age: 20}

// ❌ 错误：cannot take the address of m["alice"]
m["alice"].Age = 21

// ✅ 解决方案 1：重新赋值
u := m["alice"]
u.Age = 21
m["alice"] = u

// ✅ 解决方案 2：使用指针
m2 := make(map[string]*User)
m2["alice"] = &User{Name: "Alice", Age: 20}
m2["alice"].Age = 21  // OK
```

#### 问题 3：遍历时删除元素
```go
// ✅ 安全：Go 允许在遍历时删除
m := map[string]int{"a": 1, "b": 2, "c": 3}
for k := range m {
    if k == "b" {
        delete(m, k)  // 安全
    }
}

// ❌ 但不要在遍历时添加元素并期望遍历到
for k := range m {
    m[k+"_new"] = 1  // 新元素可能不会被遍历到
}
```

#### 问题 4：Map 作为函数参数
```go
// Map 是引用类型，函数内修改会影响外部
func modify(m map[string]int) {
    m["key"] = 100
}

m := map[string]int{"key": 1}
modify(m)
fmt.Println(m["key"])  // 100
```

### 扩容机制

- 负载因子（load factor）> 6.5 时触发扩容
- 扩容为原来的 2 倍
- 采用渐进式扩容，避免一次性复制大量数据

---

## Channel (通道)

### 核心原理

Channel 的底层结构：
```go
type hchan struct {
    qcount   uint           // 当前队列中的元素个数
    dataqsiz uint           // 环形队列的大小
    buf      unsafe.Pointer // 环形队列指针
    elemsize uint16         // 元素大小
    closed   uint32         // 是否关闭
    sendx    uint           // 发送索引
    recvx    uint           // 接收索引
    recvq    waitq          // 等待接收的 goroutine 队列
    sendq    waitq          // 等待发送的 goroutine 队列
    lock     mutex          // 互斥锁
}
```

**关键特性：**
- Channel 是引用类型
- 用于 goroutine 之间的通信
- 可以是有缓冲或无缓冲
- 支持关闭操作

### 最佳实践

#### 1. 选择合适的 Channel 类型
```go
// 无缓冲 channel：同步通信
ch := make(chan int)

// 有缓冲 channel：异步通信
ch := make(chan int, 10)

// 单向 channel：限制操作
func send(ch chan<- int) {  // 只能发送
    ch <- 1
}

func receive(ch <-chan int) {  // 只能接收
    <-ch
}
```

#### 2. 关闭 Channel 的原则
```go
// ✅ 原则：只由发送方关闭 channel
func producer(ch chan<- int) {
    defer close(ch)  // 发送完毕后关闭
    for i := 0; i < 10; i++ {
        ch <- i
    }
}

func consumer(ch <-chan int) {
    for v := range ch {  // range 会在 channel 关闭后退出
        fmt.Println(v)
    }
}

// ❌ 不要在接收方关闭 channel
// ❌ 不要关闭已关闭的 channel（会 panic）
// ❌ 不要向已关闭的 channel 发送数据（会 panic）
```

#### 3. 使用 select 处理多个 Channel
```go
select {
case v := <-ch1:
    fmt.Println("received from ch1:", v)
case v := <-ch2:
    fmt.Println("received from ch2:", v)
case ch3 <- 100:
    fmt.Println("sent to ch3")
case <-time.After(time.Second):
    fmt.Println("timeout")
default:
    fmt.Println("no channel ready")
}
```

#### 4. 优雅退出模式
```go
// 使用 context 控制 goroutine 退出
func worker(ctx context.Context, ch <-chan int) {
    for {
        select {
        case v := <-ch:
            process(v)
        case <-ctx.Done():
            return  // 优雅退出
        }
    }
}

// 使用 done channel
func worker(ch <-chan int, done <-chan struct{}) {
    for {
        select {
        case v := <-ch:
            process(v)
        case <-done:
            return
        }
    }
}
```

### 常见问题

#### 问题 1：Channel 的零值是 nil
```go
var ch chan int
ch <- 1    // 永久阻塞
v := <-ch  // 永久阻塞

// nil channel 的特性：
// - 发送和接收都会永久阻塞
// - 在 select 中，nil channel 的 case 永远不会被选中
```

#### 问题 2：死锁
```go
// ❌ 死锁：无缓冲 channel，没有接收者
ch := make(chan int)
ch <- 1  // fatal error: all goroutines are asleep - deadlock!

// ✅ 解决方案 1：使用 goroutine
go func() {
    ch <- 1
}()
<-ch

// ✅ 解决方案 2：使用有缓冲 channel
ch := make(chan int, 1)
ch <- 1
<-ch
```

#### 问题 3：Goroutine 泄漏
```go
// ❌ 问题：goroutine 永远阻塞
func leak() {
    ch := make(chan int)
    go func() {
        val := <-ch  // 永远等待，goroutine 泄漏
        fmt.Println(val)
    }()
    // ch 没有发送数据，goroutine 永远不会退出
}

// ✅ 解决：使用 context 或 done channel
func noLeak(ctx context.Context) {
    ch := make(chan int)
    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        case <-ctx.Done():
            return  // 可以退出
        }
    }()
}
```

#### 问题 4：向已关闭的 Channel 操作
```go
ch := make(chan int, 1)
close(ch)

// ❌ panic: send on closed channel
ch <- 1

// ✅ 接收已关闭的 channel 返回零值
v := <-ch  // v = 0, 不会阻塞

// ✅ 使用 comma ok 检查
v, ok := <-ch
if !ok {
    fmt.Println("channel closed")
}

// ❌ panic: close of closed channel
close(ch)
```

#### 问题 5：Channel 的容量选择
```go
// 无缓冲：强同步，发送和接收必须同时准备好
ch := make(chan int)

// 小缓冲：适合生产消费速率接近的场景
ch := make(chan int, 10)

// 大缓冲：适合突发流量，但要注意内存占用
ch := make(chan int, 10000)

// ❌ 不要用过大的缓冲掩盖设计问题
```

### Channel 使用模式

#### 1. Fan-out（扇出）
```go
func fanOut(input <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        outputs[i] = worker(input)
    }
    return outputs
}

func worker(input <-chan int) <-chan int {
    output := make(chan int)
    go func() {
        defer close(output)
        for v := range input {
            output <- process(v)
        }
    }()
    return output
}
```

#### 2. Fan-in（扇入）
```go
func fanIn(inputs ...<-chan int) <-chan int {
    output := make(chan int)
    var wg sync.WaitGroup
    
    for _, input := range inputs {
        wg.Add(1)
        go func(ch <-chan int) {
            defer wg.Done()
            for v := range ch {
                output <- v
            }
        }(input)
    }
    
    go func() {
        wg.Wait()
        close(output)
    }()
    
    return output
}
```

#### 3. Pipeline（管道）
```go
func pipeline() {
    // 生成数据
    gen := func(nums ...int) <-chan int {
        out := make(chan int)
        go func() {
            defer close(out)
            for _, n := range nums {
                out <- n
            }
        }()
        return out
    }
    
    // 处理数据
    square := func(in <-chan int) <-chan int {
        out := make(chan int)
        go func() {
            defer close(out)
            for n := range in {
                out <- n * n
            }
        }()
        return out
    }
    
    // 组合管道
    nums := gen(1, 2, 3, 4)
    squared := square(nums)
    
    for v := range squared {
        fmt.Println(v)
    }
}
```

#### 4. 超时控制
```go
select {
case result := <-ch:
    return result
case <-time.After(time.Second):
    return errors.New("timeout")
}
```

#### 5. 信号量模式
```go
// 限制并发数
type Semaphore chan struct{}

func NewSemaphore(max int) Semaphore {
    return make(Semaphore, max)
}

func (s Semaphore) Acquire() {
    s <- struct{}{}
}

func (s Semaphore) Release() {
    <-s
}

// 使用
sem := NewSemaphore(10)
for i := 0; i < 100; i++ {
    sem.Acquire()
    go func() {
        defer sem.Release()
        // 最多 10 个并发
        doWork()
    }()
}
```

---

## 性能优化建议

### Slice
- 预分配容量避免频繁扩容
- 使用 `s[:0]` 重用底层数组
- 大切片截取小片段时注意内存泄漏

### Map
- 预分配容量
- 避免频繁创建和销毁
- 考虑使用 `sync.Map` 处理并发场景

### Channel
- 选择合适的缓冲大小
- 避免 goroutine 泄漏
- 使用 `sync.Pool` 复用 channel

## 调试技巧

### 检测 Goroutine 泄漏
```go
import "runtime"

before := runtime.NumGoroutine()
// 执行代码
after := runtime.NumGoroutine()
if after > before {
    fmt.Printf("goroutine leak: %d\n", after-before)
}
```

### 检测竞态条件
```bash
go run -race main.go
go test -race ./...
```

### 性能分析
```go
import _ "net/http/pprof"

go func() {
    http.ListenAndServe("localhost:6060", nil)
}()

// 访问 http://localhost:6060/debug/pprof/
```

---

## 总结

### Slice 核心要点
- 理解底层数组共享机制
- 正确使用 append 返回值
- 预分配容量提升性能
- 注意切片泄漏问题

### Map 核心要点
- 必须初始化才能使用
- 非线程安全，并发需加锁
- 遍历顺序随机
- 不能直接修改 struct 字段

### Channel 核心要点
- 只由发送方关闭
- 注意 goroutine 泄漏
- 合理选择缓冲大小
- 使用 select 处理多路复用
- nil channel 会永久阻塞
