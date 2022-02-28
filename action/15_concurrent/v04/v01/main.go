package main

/*
1、sync 包低级同步原语可以用在哪？

一般情况下，优先使用 CSP 并发模型进行并发程序设计，但是在下面一些场景中，依然需要 sync 包提供的低级同步原语。

1）首先是需要高性能的临界区（critical section）同步机制场景。

在 Go 中，channel 并发原语也可以用于对数据对象访问的同步，可以把 channel 看成是一种高级的同步原语，它自身的实现也是建构在低级同步原语之上的。
也正因为如此，channel 自身的性能与低级同步原语相比要略微逊色，开销要更大。

这里，关于 sync.Mutex 和 channel 各自实现的临界区同步机制，我做了一个简单的性能基准测试对比，通过对比结果，可以很容易看出两者的性能差异：
sync_test.go

运行这个对比测试（Go 1.17），得到：
$go test -bench .

结果：
goos: darwin
goarch: amd64
pkg: github.com/lcy2013/sync_mutex_channel_test
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkCriticalSectionSyncByMutex-8                   79013780                15.36 ns/op
BenchmarkCriticalSectionSyncByMutexInParallel-8         19623278                60.60 ns/op
BenchmarkCriticalSectionSyncByChan-8                    23533351                47.93 ns/op
BenchmarkCriticalSectionSyncByChanInParallel-8           4129382               293.8 ns/op
PASS
ok      github.com/lcy2013/sync_mutex_channel_test      5.697s

通过这个对比实验可以看到，无论是在单 Goroutine 情况下，还是在并发测试情况下，sync.Mutex实现的同步机制的性能，都要比 channel 实现的高出三倍多。

因此，通常在需要高性能的临界区（critical section）同步机制的情况下，sync 包提供的低级同步原语更为适合。

2）第二种就是在不想转移结构体对象所有权，但又要保证结构体内部状态数据的同步访问的场景。

基于 channel 的并发设计，有一个特点：在 Goroutine 间通过 channel 转移数据对象的所有权。
所以，只有拥有数据对象所有权（从 channel 接收到该数据）的 Goroutine 才可以对该数据对象进行状态变更。

如果你的设计中没有转移结构体对象所有权，但又要保证结构体内部状态数据在多个 Goroutine 之间同步访问，那么你可以使用 sync 包提供的低级同步原语来实现，比如最常用的sync.Mutex。

2、sync 包中同步原语使用的注意事项
在 sync 包的注释中（在$GOROOT/src/sync/mutex.go文件的头部注释），看到这样一行说明：

// Values containing the types defined in this package should not be copied.
翻译过来就是：“不应复制那些包含了此包中类型的值”。

在 sync 包的其他源文件中，同样看到类似的一些注释：
// $GOROOT/src/sync/mutex.go
// A Mutex must not be copied after first use. （禁止复制首次使用后的Mutex）

// $GOROOT/src/sync/rwmutex.go
// A RWMutex must not be copied after first use.（禁止复制首次使用后的RWMutex）

// $GOROOT/src/sync/cond.go
// A Cond must not be copied after first use.（禁止复制首次使用后的Cond）
... ...

那么，为什么首次使用 Mutex 等 sync 包中定义的结构类型后，不应该再对它们进行复制操作呢？

以 Mutex 这个同步原语为例，看看它的实现是怎样的。

Go 标准库中 sync.Mutex 的定义是这样的：
// $GOROOT/src/sync/mutex.go
type Mutex struct {
    state int32
    sema  uint32
}

可以看到，Mutex 的定义非常简单，由两个整型字段 state 和 sema 组成：
1、state：表示当前互斥锁的状态；
2、sema：用于控制锁状态的信号量。

初始情况下，Mutex 的实例处于 Unlocked 状态（state 和 sema 均为 0）。
对 Mutex 实例的复制也就是两个整型字段的复制。
一旦发生复制，原变量与副本就是两个单独的内存块，各自发挥同步作用，互相就没有了关联。

如果发生复制后，你仍然认为原变量与副本保护的是同一个数据对象，那可就大错特错了。
我们来看一个例子：
action/15_concurrent/v04/v01/copymutx/main.go

在这个例子中，使用一个 sync.Mutex 类型变量 mu 来同步对整型变量 i 的访问。
创建一个新 Goroutine：g1，g1 通过函数参数得到 mu 的一份拷贝 mu1，然后 g1 会通过 mu1 来同步对整型变量 i 的访问。

那么，g0 通过 mu 和 g1 通过 mu 的拷贝 mu1，是否能实现对同一个变量 i 的同步访问呢？来看看运行这个示例的运行结果：
g0: i = 1
g1: i = 1

从结果来看，这个程序并没有实现对 i 的同步访问，第 17 行 g1 对 mu1 的加锁操作，并没能阻塞第 27 行 g0 对 mu 的加锁。
于是，g1 刚刚将 i 赋值为 10 后，g0 就又将 i 赋值为 1 了。

出现这种结果的原因就是我们前面分析的情况，一旦 Mutex 类型变量被拷贝，原变量与副本就各自发挥作用，互相没有关联了。
甚至，如果拷贝的时机不对，比如在一个 mutex 处于 locked 的状态时对它进行了拷贝，就会对副本进行加锁操作，将导致加锁的 Goroutine 永远阻塞下去。

通过前面这个例子，可以很直观地看到：如果对使用过的、sync 包中的类型的示例进行复制，并使用了复制后得到的副本，将导致不可预期的结果。
所以，在使用 sync 包中的类型的时候，推荐通过闭包方式，或者是传递类型实例（或包裹该类型的类型实例）的地址（指针）的方式进行。
这就是使用 sync 包时最值得注意的事项。

3、互斥锁（Mutex）还是读写锁（RWMutex）？

sync 包提供了两种用于临界区同步的原语：互斥锁（Mutex）和读写锁（RWMutex）。
它们都是零值可用的数据类型，也就是不需要显式初始化就可以使用，并且使用方法都比较简单。

在上面的示例中，已经看到了 Mutex 的应用方法，这里再总结一下：
var mu sync.Mutex
mu.Lock()   // 加锁
doSomething()
mu.Unlock() // 解锁

一旦某个 Goroutine 调用的 Mutex 执行 Lock 操作成功，它将成功持有这把互斥锁。
这个时候，如果有其他 Goroutine 执行 Lock 操作，就会阻塞在这把互斥锁上，
直到持有这把锁的 Goroutine 调用 Unlock 释放掉这把锁后，才会抢到这把锁的持有权并进入临界区。

由此，可以得到使用互斥锁的两个原则：
1）尽量减少在锁中的操作。这可以减少其他因 Goroutine 阻塞而带来的损耗与延迟。
2）一定要记得调用 Unlock 解锁。忘记解锁会导致程序局部死锁，甚至是整个程序死锁，会导致严重的后果。同时，也可以结合 defer，优雅地执行解锁操作。

读写锁与互斥锁用法大致相同，只不过多了一组加读锁和解读锁的方法：
var rwmu sync.RWMutex
rwmu.RLock()   //加读锁
readSomething()
rwmu.RUnlock() //解读锁

rwmu.Lock()    //加写锁
changeSomething()
rwmu.Unlock()  //解写锁

写锁与 Mutex 的行为十分类似，一旦某 Goroutine 持有写锁，其他 Goroutine 无论是尝试加读锁，还是加写锁，都会被阻塞在写锁上。

但读锁就宽松多了，一旦某个 Goroutine 持有读锁，它不会阻塞其他尝试加读锁的 Goroutine，但加写锁的 Goroutine 依然会被阻塞住。

通常，互斥锁（Mutex）是临时区同步原语的首选，它常被用来对结构体对象的内部状态、缓存等进行保护，是使用最为广泛的临界区同步原语。
相比之下，读写锁的应用就没那么广泛了，只活跃于它擅长的场景下。

那读写锁（RWMutex）究竟擅长在哪种场景下呢？先来看一组基准测试：
action/15_concurrent/v04/v01/rwmutx/main_test.go

$ go test -bench .
goos: darwin
goarch: amd64
pkg: github.com/lcy2013/sync_mutex_channel_test/rwmutx
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkWriteSyncByMutex-8             20367465                59.72 ns/op
BenchmarkReadSyncByMutex-8              19002603                65.84 ns/op
BenchmarkReadSyncByRWMutex-8            30787990                39.09 ns/op
BenchmarkWriteSyncByRWMutex-8           15851838                73.22 ns/op
PASS
ok      github.com/lcy2013/sync_mutex_channel_test/rwmutx       5.531s

通过测试结果对比，得到了一些结论：
并发量较小的情况下，Mutex 性能最好；随着并发量增大，Mutex 的竞争激烈，导致加锁和解锁性能下降；

RWMutex 的读锁性能并没有随着并发量的增大，而发生较大变化，性能始终恒定在 40ns 左右；

在并发量较大的情况下，RWMutex 的写锁性能和 Mutex、RWMutex 读锁相比，是最差的，并且随着并发量增大，RWMutex 写锁性能有继续下降趋势。

由此，可以看出读写锁适合应用在具有一定并发量且读多写少的场合。在大量并发读的情况下，多个 Goroutine 可以同时持有读锁，从而减少在锁竞争中等待的时间。

而互斥锁，即便是读请求的场合，同一时刻也只能有一个 Goroutine 持有锁，其他 Goroutine 只能阻塞在加锁操作上等待被调度。

接下来，继续看条件变量 sync.Cond。

4、条件变量
sync.Cond是传统的条件变量原语概念在 Go 语言中的实现。
可以把一个条件变量理解为一个容器，这个容器中存放着一个或一组等待着某个条件成立的 Goroutine。
当条件成立后，这些处于等待状态的 Goroutine 将得到通知，并被唤醒继续进行后续的工作。
这与百米飞人大战赛场上，各位运动员等待裁判员的发令枪声的情形十分类似。

条件变量是同步原语的一种，如果没有条件变量，开发人员可能需要在 Goroutine 中通过连续轮询的方式，检查某条件是否为真，这种连续轮询非常消耗资源，因为 Goroutine 在这个过程中是处于活动状态的，但它的工作又没有进展。

这里先看一个用sync.Mutex 实现对条件轮询等待的例子：
action/15_concurrent/v04/v01/condmutex/v1/main.go

轮询的方式开销大，轮询间隔设置的不同，条件检查的及时性也会受到影响。

sync.Cond为 Goroutine 在这个场景下提供了另一种可选的、资源消耗更小、使用体验更佳的同步方式。使用条件变量原语，我们可以在实现相同目标的同时，避免对条件的轮询。

用sync.Cond对上面的例子进行改造，改造后的代码如下：
action/15_concurrent/v04/v01/condmutex/v2/main.go

sync.Cond实例的初始化，需要一个满足实现了sync.Locker接口的类型实例，通常使用sync.Mutex。

条件变量需要这个互斥锁来同步临界区，保护用作条件的数据。
加锁后，各个等待条件成立的 Goroutine 判断条件是否成立，如果不成立，则调用sync.Cond的 Wait 方法进入等待状态。
Wait 方法在 Goroutine 挂起前会进行 Unlock 操作。

当 main goroutine 将ready置为 true，并调用sync.Cond的 Broadcast 方法后，各个阻塞的 Goroutine 将被唤醒，并从 Wait 方法中返回。
Wait 方法返回前，Wait 方法会再次加锁让 Goroutine 进入临界区。
接下来 Goroutine 会再次对条件数据进行判定，如果条件成立，就会解锁并进入下一个工作阶段；如果条件依旧不成立，那么会再次进入循环体，并调用 Wait 方法挂起等待。

和sync.Mutex 、sync.RWMutex等相比，sync.Cond 应用的场景更为有限，只有在需要“等待某个条件成立”的场景下，Cond 才有用武之地。

其实，面向 CSP 并发模型的 channel 原语和面向传统共享内存并发模型的 sync 包提供的原语，已经能够满足 Go 语言应用并发设计中 **99.9%的并发同步需求了。
而剩余那0.1%** 的需求，可以使用 Go 标准库提供的 atomic 包来实现。

5、原子操作（atomic operations）
atomic 包是 Go 语言给用户提供的原子操作原语的相关接口。
原子操作（atomic operations）是相对于普通指令操作而言的。

以一个整型变量自增的语句为例说明一下：
var a int
a++

a++ 这行语句需要 3 条普通机器指令来完成变量 a 的自增：
1）LOAD：将变量从内存加载到 CPU 寄存器；
2）ADD：执行加法指令；
3）STORE：将结果存储回原内存地址中。

这 3 条普通指令在执行过程中是可以被中断的。
而原子操作的指令是不可中断的，它就好比一个事务，要么不执行，一旦执行就一次性全部执行完毕，中间不可分割。
也正因为如此，原子操作也可以被用于共享数据的并发同步。

原子操作由底层硬件直接提供支持，是一种硬件实现的指令级的“事务”，因此相对于操作系统层面和 Go 运行时层面提供的同步技术而言，它更为原始。

atomic 包封装了 CPU 实现的部分原子操作指令，为用户层提供体验良好的原子操作函数，因此 atomic 包中提供的原语更接近硬件底层，也更为低级，它也常被用于实现更为高级的并发同步技术，比如 channel 和 sync 包中的同步原语。

以 atomic.SwapInt64 函数在 x86_64 平台上的实现为例，看看这个函数的实现方法：
// $GOROOT/src/sync/atomic/doc.go
func SwapInt64(addr *int64, new int64) (old int64)

// $GOROOT/src/sync/atomic/asm.s
TEXT ·SwapInt64(SB),NOSPLIT,$0
        JMP     runtime∕internal∕atomic·Xchg64(SB)

// $GOROOT/src/runtime/internal/asm_amd64.s
TEXT runtime∕internal∕atomic·Xchg64(SB), NOSPLIT, $0-24
        MOVQ    ptr+0(FP), BX
        MOVQ    new+8(FP), AX
        XCHGQ   AX, 0(BX)
        MOVQ    AX, ret+16(FP)
        RET

从函数 SwapInt64 的实现中，可以看到：它基本就是对 x86_64 CPU 实现的原子操作指令XCHGQ的直接封装。

原子操作的特性，让 atomic 包也可以被用作对共享数据的并发同步，那么和更为高级的 channel 以及 sync 包中原语相比，究竟该怎么选择呢？

先来看看 atomic 包提供了哪些能力。

atomic 包提供了两大类原子操作接口，
一类是针对整型变量的，包括有符号整型、无符号整型以及对应的指针类型；
另外一类是针对自定义类型的。
因此，第一类原子操作接口的存在让 atomic 包天然适合去实现某一个共享整型变量的并发同步。

再看一个例子：
action/15_concurrent/v04/v01/atomic/atomic_test.go

$ go test -bench .

goos: darwin
goarch: amd64
pkg: github.com/lcy2013/sync_mutex_channel_test/atomic
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkAddSyncByAtomic-8      53621286                21.93 ns/op
BenchmarkReadSyncByAtomic-8     1000000000               0.2719 ns/op
BenchmarkAddSyncByRWMutex-8     14705517                76.90 ns/op
BenchmarkReadSyncByRWMutex-8    34780426                35.30 ns/op
PASS
ok      github.com/lcy2013/sync_mutex_channel_test/atomic       4.459s

通过这个运行结果，可以得出一些结论：

1）读写锁的性能随着并发量增大的情况，与 sync.RWMutex 一致；

2）利用原子操作的无锁并发写的性能，随着并发量增大几乎保持恒定；

3）利用原子操作的无锁并发读的性能，随着并发量增大有持续提升的趋势，并且性能是读锁的约 200 倍。

通过这些结论，大致可以看到 atomic 原子操作的特性：随着并发量提升，使用 atomic 实现的共享变量的并发读写性能表现更为稳定，尤其是原子读操作，和 sync 包中的读写锁原语比起来，atomic 表现出了更好的伸缩性和高性能。

由此，也可以看出 atomic 包更适合一些对性能十分敏感、并发量较大且读多写少的场合。

不过，atomic 原子操作可用来同步的范围有比较大限制，只能同步一个整型变量或自定义类型变量。
如果要对一个复杂的临界区数据进行同步，那么首选的依旧是 sync 包中的原语。

总结：
虽然 Go 推荐基于通信来共享内存的并发设计风格，但 Go 并没有彻底抛弃对基于共享内存并发模型的支持，Go 通过标准库的 sync 包以及 atomic 包提供了低级同步原语。这些原语有着它们自己的应用场景。

如果我们考虑使用低级同步原语，一般都是因为低级同步原语可以提供更佳的性能表现，性能基准测试结果告诉我们，使用低级同步原语的性能可以高出 channel 许多倍。在性能敏感的场景下，我们依然离不开这些低级同步原语。

在使用 sync 包提供的同步原语之前，一定要牢记这些原语使用的注意事项：不要复制首次使用后的 Mutex/RWMutex/Cond 等。一旦复制，你将很大可能得到意料之外的运行结果。

sync 包中的低级同步原语各有各的擅长领域，你可以记住：
1）在具有一定并发量且读多写少的场合使用 RWMutex；
2）在需要“等待某个条件成立”的场景下使用 Cond；
3）当你不确定使用什么原语时，那就使用 Mutex 吧。

如果你对同步的性能有极致要求，且并发量较大，读多写少，那么可以考虑一下 atomic 包提供的原子操作函数。





*/
