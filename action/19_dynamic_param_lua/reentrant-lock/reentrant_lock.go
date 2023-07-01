package reentrant_mutex

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"sync"
	"time"
)

const KEY = "EXAMPLE_LOCK"

// ReentrantLock 可重入锁
type ReentrantLock struct {
	// redis连接池
	rdb *redis.Client
	// hash锁key
	key string
	// hash锁field, 具体的场景字段
	field string
	// 锁有效期
	expiration time.Duration
	// 用于模拟的初始递归层数
	recursionLevel int
	// 用于模拟的最大递归层数
	maxrecursionLevel int
	// 用于模拟的任务最小执行时间
	min int
	// 用于模拟的任务最大执行时间
	max int
	// 加锁失败的重试间隔
	retryInterval time.Duration
	// 加锁失败的重试次数
	retryTimes int
	// 继承*sync.Once的特性
	*sync.Once
	// 用于测试打印的线程标签
	tag string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 生成一个随机标签
func getRandtag(n int) string {
	var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	tag := make([]rune, n)
	for i := range tag {
		tag[i] = runes[rand.Intn(len(runes))]
	}
	return string(tag)
}

type ReentrantLockOption func(*ReentrantLock)

// NewReentrantLock 初始化
func NewReentrantLock(rdb *redis.Client, options ...ReentrantLockOption) *ReentrantLock {
	l := ReentrantLock{
		rdb:               rdb,
		key:               KEY, // 固定值
		field:             fmt.Sprintf("%d", rand.Int()),
		expiration:        time.Millisecond * 200,
		recursionLevel:    1,
		maxrecursionLevel: 1,
		min:               50,
		max:               100,
		retryInterval:     time.Millisecond * 50,
		retryTimes:        5,
		Once:              new(sync.Once),
		tag:               getRandtag(2),
	}
	for _, option := range options {
		option(&l)
	}
	return &l
}

// WithLockkey 设置锁key
func WithLockkey(lockkey string) ReentrantLockOption {
	return func(lock *ReentrantLock) {
		lock.key = lockkey
	}
}

// WithLockfield 设置锁field
func WithLockfield(lockFile string) ReentrantLockOption {
	return func(lock *ReentrantLock) {
		lock.field = lockFile
	}
}

// WithLockexpiration 设置锁过期时间
func WithLockexpiration(lockDuration time.Duration) ReentrantLockOption {
	return func(lock *ReentrantLock) {
		lock.expiration = lockDuration
	}
}

// WithLockRecursionLevel 模拟初始递归层数
func WithLockRecursionLevel(recursionLevel int) ReentrantLockOption {
	return func(lock *ReentrantLock) {
		lock.recursionLevel = recursionLevel
	}
}

// WithLockMaxrecursionLevel 模拟最大递归层数
func WithLockMaxrecursionLevel(maxrecursionLevel int) ReentrantLockOption {
	return func(lock *ReentrantLock) {
		lock.maxrecursionLevel = maxrecursionLevel
	}
}

// WithLockMax 模拟最大锁时间
func WithLockMax(max int) ReentrantLockOption {
	return func(lock *ReentrantLock) {
		lock.max = max
	}
}

// WithLockMin 模拟最小锁时间
func WithLockMin(min int) ReentrantLockOption {
	return func(lock *ReentrantLock) {
		lock.min = min
	}
}

// WithLockTag 标记
func WithLockTag(tag string) ReentrantLockOption {
	return func(lock *ReentrantLock) {
		lock.tag = tag
	}
}

// WithLockRetryInterval 加锁失败的重试间隔
func WithLockRetryInterval(retryInterval time.Duration) ReentrantLockOption {
	return func(lock *ReentrantLock) {
		lock.retryInterval = retryInterval
	}
}

// WithLockRetryTimes 加锁失败的重试次数
func WithLockRetryTimes(retryTimes int) ReentrantLockOption {
	return func(lock *ReentrantLock) {
		lock.retryTimes = retryTimes
	}
}

// MockBusiness 模拟分布式业务加锁场景
func (l *ReentrantLock) MockBusiness() {
	fmt.Printf("%s的第%d次调用,field:%d\n", l.tag, l.recursionLevel, l.field)

	// 初始化仅用于当前调用的ctx,避免在重入调用完成后执行cancel()导致的上层调用出现context canceled错误
	var ctx, cancel = context.WithCancel(context.Background())

	defer func() {
		// 延迟停止守护线程
		cancel()
	}()

	set, err := l.Lock(ctx)

	if err != nil {
		fmt.Println(l.tag + " 加锁失败:" + err.Error())
		return
	}

	// 加锁失败,重试
	if set == false {
		res, err := l.TryLock(ctx)
		if err != nil {
			fmt.Println(l.tag + " 重试加锁失败:" + err.Error())
			return
		}
		// 重试达到最大次数
		if res == false {
			fmt.Println(l.tag + " server unavailable, try again later")
			return
		}
	}

	fmt.Println(l.tag + "成功加锁")

	// 加锁成功,通过守护线程自动续期(此处可以异步执行,即使自动续期还没来得及执行业务就已经完成,也不会影响流程)
	go l.watchDog(ctx)

	fmt.Println(l.tag + "等待业务处理完成...")
	// 模拟处理业务(通过随机时间模拟业务延迟)
	time.Sleep(time.Duration(rand.Intn(l.max-l.min)+l.min) * time.Millisecond)

	// 模拟重入调用(测试锁的可重入)
	if l.recursionLevel <= l.maxrecursionLevel {
		l.recursionLevel += 1
		l.MockBusiness()
	}

	// 业务处理完成
	// 释放锁
	val, err := l.Unlock(ctx)
	if err != nil {
		fmt.Println(l.tag + "锁释放失败:" + err.Error())
		return
	}

	// 递归调用中的结果都是false,因为lua脚本中的if分支counter>0,没有释放
	fmt.Println(l.tag+"释放结果:", val)
}

// 守护线程(通过sync.Once.Do确保仅在线程第一次调用时执行自动续期)
func (l *ReentrantLock) watchDog(ctx context.Context) {
	l.Once.Do(func() {
		fmt.Printf("打开了%s的守护线程\n", l.tag)
		for {
			select {
			// 业务完成
			case <-ctx.Done():
				fmt.Printf("%s任务完成,关闭%s的自动续期\n", l.tag, l.key)
				return
				// 业务未完成
			default:
				// 自动续期
				l.rdb.PExpire(ctx, l.key, l.expiration)
				// 继续等待
				time.Sleep(l.expiration / 2)
			}
		}
	})
}

// Lock 加锁
func (l *ReentrantLock) Lock(ctx context.Context) (res bool, err error) {
	lua := `
-- KEYS[1]:锁对应的key
-- ARGV[1]:锁的expire
-- ARGV[2]:锁对应的计数器field(随机值,防止误解锁),记录当前线程已加锁的次数
-- 判断锁是否空闲
if (redis.call('EXISTS', KEYS[1]) == 0) then
    -- 线程首次加锁(锁的初始化,值和过期时间)
    redis.call('HINCRBY', KEYS[1], ARGV[2], 1);
    redis.call('PEXPIRE', KEYS[1], ARGV[1]);
    return 1;
end;
-- 判断当前线程是否持有锁(锁被某个线程持有,通常是程序第N次(N>1)在线程内调用时会执行到此处)
if (redis.call('HEXISTS', KEYS[1], ARGV[2]) == 1) then
    -- 调用次数递增
    redis.call('HINCRBY', KEYS[1], ARGV[2], 1);
    -- 不处理续期,通过守护线程续期
    return 1;
end;
-- 锁被其他线程占用,加锁失败
return 0;
`

	scriptkeys := []string{l.key}

	val, err := l.rdb.Eval(ctx, lua, scriptkeys, l.expiration.Milliseconds(), l.field).Result()
	if err != nil {
		return
	}

	res = val == int64(1)

	return
}

// Unlock 解锁
func (l *ReentrantLock) Unlock(ctx context.Context) (res bool, err error) {
	lua := `
-- KEYS[1]:锁对应的key
-- ARGV[1]:锁对应的计数器field(随机值,防止误解锁),记录当前线程已加锁的次数
-- 判断 hash set 是否存在
if (redis.call('HEXISTS', KEYS[1], ARGV[1]) == 0) then
    -- err = redis.Nil
    return nil;
end;
-- 计算当前可重入次数
local counter = redis.call('HINCRBY', KEYS[1], ARGV[1], -1);
if (counter > 0) then
-- 同一线程内部多次调用完成后尝试释放锁会进入此if分支
    return 0;
else
-- 同一线程最外层(第一次)调用完成后尝试释放锁会进入此if分支
-- 小于等于 0 代表内层嵌套调用已全部完成，可以解锁
    redis.call('DEL', KEYS[1]);
    return 1;
end;
-- err = redis.Nil
return nil;
`

	scriptkeys := []string{l.key}
	val, err := l.rdb.Eval(ctx, lua, scriptkeys, l.field).Result()
	if err != nil {
		return
	}

	res = val == int64(1)

	return
}

// TryLock 重试
func (l *ReentrantLock) TryLock(ctx context.Context) (res bool, err error) {
	i := 1
	for i <= l.retryTimes {
		res, err = l.Lock(ctx)

		if err != nil {
			return
		}

		if res == true {
			return
		}

		time.Sleep(l.retryInterval)
		i++
	}
	return
}
