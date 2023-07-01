package main

import (
	"context"
	reentrant_mutex "dynamic_param_lua/reentrant-lock"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var rdb *redis.Client

func init() {
	// 初始化连接池
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
		Username: "",
	})
}

/*
测试配置:
hash锁的field通过线程初始化时生成,执行过程中field不变,field是判断一个锁是否属于当前线程唯一标准
加锁失败后重试次数为20，重试间隔为50ms
通过随机生成的Tag来标识线程以及打印流程
*/

func main() {
	max := 10
	/*for i := 0; i < max; i++ {
		go reentrant_mutex.NewReentrantLock(rdb).MockBusiness()
	}*/

	lockKey := "Test-ReentrantLock"
	ctx := context.TODO()

	for i := 0; i < max; i++ {
		go businessOne(ctx, lockKey, i)
	}

	for i := 0; i < max; i++ {
		go businessTwo(ctx, lockKey, i)
	}

	time.Sleep(time.Second * time.Duration(max))
}

func businessOne(ctx context.Context, lockKey string, runNum int) {
	lock := reentrant_mutex.NewReentrantLock(
		rdb,
		reentrant_mutex.WithLockkey(lockKey),
		reentrant_mutex.WithLockfield("businessOne"),
		reentrant_mutex.WithLockexpiration(10*time.Second),
		reentrant_mutex.WithLockRetryInterval(50*time.Millisecond),
		reentrant_mutex.WithLockRetryTimes(120),
		reentrant_mutex.WithLockTag(fmt.Sprintf("businessOne-%d", runNum)))
	res, err := lock.TryLock(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer lock.Unlock(ctx)

	if res {
		fmt.Printf("businessOne-%d 加锁成功: %s\n", runNum, time.Now())
	} else {
		fmt.Printf("businessOne-%d 加锁失败: %s\n", runNum, time.Now())
	}

	time.Sleep(time.Millisecond * 5000)
}

func businessTwo(ctx context.Context, lockKey string, runNum int) {
	lock := reentrant_mutex.NewReentrantLock(
		rdb,
		reentrant_mutex.WithLockkey(lockKey),
		reentrant_mutex.WithLockfield("businessTwo"),
		reentrant_mutex.WithLockexpiration(10*time.Second),
		reentrant_mutex.WithLockRetryInterval(50*time.Millisecond),
		reentrant_mutex.WithLockRetryTimes(120),
		reentrant_mutex.WithLockTag(fmt.Sprintf("businessTwo-%d", runNum)))
	res, err := lock.TryLock(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer lock.Unlock(ctx)

	if res {
		fmt.Printf("businessTwo-%d 加锁成功: %s\n", runNum, time.Now())
	} else {
		fmt.Printf("businessTwo-%d 加锁失败: %s\n", runNum, time.Now())
	}
	time.Sleep(time.Millisecond * 5000)
}
