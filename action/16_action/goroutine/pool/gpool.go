package pool

/*
Go 的两个核心语法知识点：接口与并发原语。

它们分别是耦合设计与并发设计的主要参与者，Go 应用的骨架设计离不开它们。

围绕 Go 并发来做，实现一个轻量级线程池，也就是 Goroutine 池。

1、为什么要用到 Goroutine 池？
Goroutine 的时候，就说过：
相对于操作系统线程，Goroutine 的开销十分小，一个 Goroutine 的起始栈大小为 2KB，而且创建、切换与销毁的代价很低，可以创建成千上万甚至更多 Goroutine。

所以和其他语言不同的是，Go 应用通常可以为每个新建立的连接创建一个对应的新 Goroutine，甚至是为每个传入的请求生成一个 Goroutine 去处理。
这种设计还有一个好处，实现起来十分简单，Gopher 们在编写代码时也没有很高的心智负担。

不过，Goroutine 的开销虽然“廉价”，但也不是免费的。

最明显的，一旦规模化后，这种非零成本也会成为瓶颈。
以一个 Goroutine 分配 2KB 执行栈为例，100w Goroutine 就是 2GB 的内存消耗。

其次，Goroutine 从Go 1.4 版本开始采用了连续栈的方案，也就是每个 Goroutine 的执行栈都是一块连续内存，如果空间不足，运行时会分配一个更大的连续内存空间作为这个 Goroutine 的执行栈，将原栈内容拷贝到新分配的空间中来。

连续栈的方案，虽然能避免 Go 1.3 采用的分段栈会导致的“hot ”问题，但连续栈的原理也决定了，一旦 Goroutine 的执行栈发生了 grow，那么即便这个 Goroutine 不再需要那么大的栈空间，这个 Goroutine 的栈空间也不会被 Shrink（收缩）了，这些空间可能会处于长时间闲置的状态，直到 Goroutine 退出。

另外，随着 Goroutine 数量的增加，Go 运行时进行 Goroutine 调度的处理器消耗，也会随之增加，成为阻碍 Go 应用性能提升的重要因素。

那么面对这样的问题，常见的应对方式是什么呢？

Goroutine 池就是一种常见的解决方案。
这个方案的核心思想是对 Goroutine 的重用，也就是把 M 个计算任务调度到 N 个 Goroutine 上，而不是为每个计算任务分配一个独享的 Goroutine，从而提高计算资源的利用率。

接下来，就来真正实现一个简单的 Goroutine 池，它就是 workerpool。

2、workerpool 的实现原理
workerpool 的工作逻辑通常都很简单，所以即便是用于生产环境的 workerpool 实现，代码规模也都在千行左右。

当然，workerpool 有很多种实现方式，这里为了更好地演示 Go 并发模型的应用模式，以及并发原语间的协作，采用完全基于 channel+select 的实现方案，不使用其他数据结构，也不使用 sync 包提供的各种同步结构，比如 Mutex、RWMutex，以及 Cond 等。

workerpool 的实现主要分为三个部分：
1）、pool 的创建与销毁；

2）、pool 中 worker（Goroutine）的管理；

3）、task 的提交与调度。

pool 对 worker 的管理

capacity 是 pool 的一个属性，代表整个 pool 中 worker 的最大容量。
使用一个带缓冲的 channel：active，作为 worker 的“计数器”，这种 channel 使用模式就是计数信号量。

当 active channel 可写时，就创建一个 worker，用于处理用户通过 Schedule 函数提交的待处理的请求。
当 active channel 满了的时候，pool 就会停止 worker 的创建，直到某个 worker 因故退出，active channel 又空出一个位置时，pool 才会创建新的 worker 填补那个空位。

我们把用户要提交给 workerpool 执行的请求抽象为一个 Task。
Task 的提交与调度也很简单：Task 通过 Schedule 函数提交到一个 task channel 中，已经创建的 worker 将从这个 task channel 中读取 task 并执行。

好了！“Talk is cheap，show me the code”！接下来写一版 workerpool 的代码，来验证一下这里分析的原理是否可行。

3、workerpool 的一个最小可行实现
先建立 workerpool 目录作为实战项目的源码根目录，然后为这个项目创建 go module：
$mkdir workerpool

$cd workerpool

$go mod init github.com/bigwhite/workerpool

创建 pool.go 作为 workpool 包的主要源码文件。
在这个源码文件中，定义了 Pool 结构体类型，这个类型的实例代表一个 workerpool：
action/16_action/workerpool/pool.go

从运行的输出结果来看，workerpool 的最小可行实现的运行逻辑与原理图是一致的。

不过，目前的 workerpool 实现好比“铁板一块”，虽然可以通过 capacity 参数可以指定 workerpool 容量，但无法对 workerpool 的行为进行定制。

比如当 workerpool 中的 worker 数量已达上限，而且 worker 都在处理 task 时，用户调用 Schedule 方法将阻塞，如果用户不想阻塞在这里，以目前的实现是做不到的。

那可以怎么改进呢？可以尝试在上面实现的基础上，为 workerpool 添加功能选项（functional option）机制。

3、添加功能选项机制
功能选项机制，可以让某个包的用户可以根据自己的需求，通过设置不同功能选项来定制包的行为。
Go 语言中实现功能选项机制有多种方法，但 Go 社区目前使用最为广泛的一个方案，是 Go 语言之父 Rob Pike 在 2014 年在博文《自引用函数与选项设计》中论述的一种，这种方案也被后人称为“功能选项（functional option）”方案。

首先，我们将 workerpool1 目录拷贝一份形成 workerpool2 目录，将在这个目录下为 workerpool 包添加功能选项机制。
然后，在 workerpool2 目录下创建 option.go 文件，在这个文件中，定义用于代表功能选项的类型 Option：
type Option func(*Pool)

这个 Option 实质是一个接受 *Pool 类型参数的函数类型。
那么如何运用这个 Option 类型呢？
现在先要做的是，明确给 workerpool 添加什么功能选项。
这里为 workerpool 添加两个功能选项：Schedule 调用是否阻塞，以及是否预创建所有的 worker。

为了支持这两个功能选项，需要在 Pool v2 类型中增加两个 bool 类型的字段，字段的具体含义，也在代码中注释了：
type Pool struct {
    ... ...
    preAlloc bool // 是否在创建pool的时候就预创建workers，默认值为：false
    // 当pool满的情况下，新的Schedule调用是否阻塞当前goroutine。默认值：true
    // 如果block = false，则Schedule返回ErrNoWorkerAvailInPool
    block  bool
    ... ...
}

针对这两个字段，在 option.go 中添加两个功能选项，WithBlock 与 WithPreAllocWorkers：
func WithBlock(block bool) Option {
    return func(p *Pool) {
        p.block = block
    }
}

func WithPreAllocWorkers(preAlloc bool) Option {
    return func(p *Pool) {
        p.preAlloc = preAlloc
    }
}

这两个功能选项实质上是两个返回闭包函数的函数。

为了支持将这两个 Option 传给 workerpool，还需要改造一下 workerpool 包的 New 函数，改造后的 New 函数代码如下：
func New(capacity int, opts ...Option) *Pool {
}

新版 New 函数除了接受 capacity 参数之外，还在它的参数列表中增加了一个类型为 Option 的可变长参数 opts。
在 New 函数体中，我们通过一个 for 循环，将传入的 Option 运用到 Pool 类型的实例上。

新版 New 函数还会根据 preAlloc 的值来判断是否预创建所有的 worker，如果需要，就调用 newWorker 方法把所有 worker 都创建出来。
newWorker 的实现与上一版代码并没有什么差异。

但由于 preAlloc 选项的加入，Pool 的 run 方法的实现有了变化，来看一下：
action/16_action/workerpool2/pool.go



*/
