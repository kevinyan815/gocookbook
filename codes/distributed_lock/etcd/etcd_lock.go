package etcd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
)

func NewClient(endpoints []string) (*clientv3.Client, error){
	ctx, cancel := context.WithCancel(context.Background())
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
		DialTimeout: 200 * time.Millisecond,
		Context: ctx,
	})
	if err != nil {
		cancel()
		return nil, err
	}

	return cli, err
}

// 测试锁
func UseLock(cli *clientv3.Client, name string) {
	// 创建一个用于获取锁的Session
	s, err:= concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	// 在指定Key上创建一个互斥锁
	l := concurrency.NewMutex(s, "/distributed-lock/")
	ctx := context.TODO() // 也可以通过 ctx 指定获取锁等待的超时时间，实现LockWithTimeout
	// 获得锁 / 阻塞等待知道获得锁
	if err := l.Lock(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Acquired lock for ", name)
	fmt.Println("Do some work in", name)
	time.Sleep(5 * time.Second)
	if err := l.Unlock(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Released lock for ", name)
}

// 使用TryLock, 获取不到则退出
func UseTryLock(cli *clientv3.Client, name string) {
	// 创建一个用于获取锁的Session
	s, err:= concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	// 在指定Key上创建一个互斥锁
	l := concurrency.NewMutex(s, "/distributed-lock/")
	ctx := context.TODO() // 也可以通过 ctx 指定获取锁等待的超时时间，实现LockWithTimeout
	// 获得锁 / 阻塞等待知道获得锁
	if err := l.Lock(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Acquired lock for ", name)
	fmt.Println("Do some work in", name)
	time.Sleep(5 * time.Second)
	if err := l.Unlock(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Released lock for ", name)
}
