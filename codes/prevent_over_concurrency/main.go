package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"golang.org/x/time/rate"
	"math/rand"
	"sync"
	"time"
)

func main() {
	// 错误的并发控制
	badConcurrency()
	// 使用WaitGroup防止并发超限
	//useWaitGroup()
	// 使用Semaphore防止并发超限
	//useSemaphore()
	// 使用golang标准库的限流器
	//useRateLimit()
	// 使用通道实现的生产者和消费者队列
	useChannel()
}

func useWaitGroup() {

	batchSize := 50
	for {
		data, _ := queryDataWithSizeN(batchSize)
		if len(data) == 0 {
			fmt.Println("End of all data")
			break
		}
		var wg sync.WaitGroup
		for _, item := range data {
			wg.Add(1)
			go func(i int) {
				doSomething(i)
				wg.Done()
			}(item)
		}
		wg.Wait()

		fmt.Println("Next bunch of data")
	}
}

func useSemaphore() {
	var concurrentNum int64 = 10
	var weight int64 = 1
	var batchSize int = 50
	s := semaphore.NewWeighted(concurrentNum)
	for {
		data, _ := queryDataWithSizeN(batchSize)
		if len(data) == 0 {
			fmt.Println("End of all data")
			break
		}

		for _, item := range data {
			s.Acquire(context.Background(), weight)
			go func(i int) {
				doSomething(i)
				s.Release(weight)
			}(item)
		}

	}
}

func useRateLimit() {
	limiter := rate.NewLimiter(rate.Every(1*time.Second), 50)
	batchSize := 50
	for {
		data, _ :=queryDataWithSizeN(batchSize)
		if len(data) == 0 {
			fmt.Println("End of all data")
			break
		}

		for _, item := range data {
			// blocking until the bucket have sufficient token
			err := limiter.Wait(context.Background())
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			go func(i int) {
				doSomething(i)
			}(item)
		}
	}
}

func useChannel() {
	batchSize := 50
	dataChan := make(chan int)
	var wg sync.WaitGroup
	wg.Add(batchSize + 1)
	// 生产者
	go func() {
		for {
			data, _ := queryDataWithSizeN(batchSize)
			if len(data) == 0 {
				break
			}
			for _, item := range data {
				dataChan <- item
			}
		}
		close(dataChan)
		wg.Done()
	}()
        // 消费者
	go func() {
		for i := 0; i < 50; i++ {
			go func() {
				for {
					select {
					case v, ok := <- dataChan:
						if !ok {
							wg.Done()
							return
						}
						doSomething(v)
					}
				}
			}()
		}
	}()

	wg.Wait()
}

func badConcurrency() {
	batchSize := 50
	for {
		data, _ := queryDataWithSizeN(batchSize)
		if len(data) == 0 {
			break
		}

		for _, item := range data {
			go func(i int) {
				doSomething(i)
			}(item)
		}

		time.Sleep(time.Second * 1)
	}
}

func doSomething(i int) {
	time.Sleep(2 * time.Second)
	fmt.Println("End:", i)
}

func queryDataWithSizeN(size int) (dataList []int, err error) {
	rand.Seed(time.Now().Unix())
	dataList = rand.Perm(size)
	return
}
