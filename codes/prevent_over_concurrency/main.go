package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
)

func main() {
	// 使用WaitGroup防止并发超限
	// useWaitGroup()
	// 使用Semaphore防止并发超限
	useSemaphore()
}

func useWaitGroup() {

	noMoreData := false
	var concurrentNum int64 = 10

	for {
		if noMoreData {
			break
		}

		var wg sync.WaitGroup
		var i int64 = 0
		for ; i < concurrentNum; i++ {
			wg.Add(1)
			go func(i int64) {
				doSomething(i)
				wg.Done()
			}(i)
		}
		wg.Wait()

		time.Sleep(200 * time.Millisecond)
		fmt.Println("Next bunch of things")
	}
}

func useSemaphore() {
	noMoreData := false
	var concurrentNum int64 = 10
	var weight int64 = 1
	s := semaphore.NewWeighted(concurrentNum)
	var i int64 = 1
	for {

		if noMoreData {
			break
		}

		go func(i int64) {
			s.Acquire(context.Background(), weight)
			doSomething(i)
			s.Release(weight)
		}(i)

		i++
	}
}

func doSomething(i int64) {
	time.Sleep(2 * time.Second)
	fmt.Println("End:", i)
}
