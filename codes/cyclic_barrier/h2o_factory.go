package main

import (
	"context"
	"fmt"
	"github.com/marusama/cyclicbarrier"
	"golang.org/x/sync/semaphore"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// 定义水分子合成的辅助数据结构
type H2O struct {
	semaH *semaphore.Weighted // 氢原子的信号量
	semaO *semaphore.Weighted // 氧原子的信号量
	b     cyclicbarrier.CyclicBarrier // 循环栅栏，用来控制合成
}


func New() *H2O {
	return &H2O{
		semaH: semaphore.NewWeighted(2), //氢原子需要两个
		semaO: semaphore.NewWeighted(1), // 氧原子需要一个
		b:     cyclicbarrier.New(3),  // 需要三个原子才能合成
	}
}

// 用来存放水分子结果的channel
var ch chan string

// 往channel里存放一个H原子
func releaseHydrogen() {
	ch <- "H"
}

// 往channel里存放一个O原子
func releaseOxygen() {
	ch <- "O"
}

func (h2o *H2O) hydrogen(releaseHydrogen func()) {
	h2o.semaH.Acquire(context.Background(), 1)

	releaseHydrogen() // 输出H
	h2o.b.Await(context.Background()) //等待栅栏放行
	h2o.semaH.Release(1) // 释放氢原子空槽
}


func (h2o *H2O) oxygen(releaseOxygen func()) {
	h2o.semaO.Acquire(context.Background(), 1)

	releaseOxygen() // 输出O
	h2o.b.Await(context.Background()) //等待栅栏放行
	h2o.semaO.Release(1) // 释放氢原子空槽
}


func main() {
	// 300个原子，300个goroutine,每个goroutine并发的产生一个原子
	var N = 100
	ch = make(chan string, N*3)

	h2o := New()

	// 用来等待所有的goroutine完成
	var wg sync.WaitGroup
	wg.Add(N * 3)

	// 200个生产氢原子的goroutine
	for i := 0; i < 2*N; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			h2o.hydrogen(releaseHydrogen)
			wg.Done()
		}()
	}
	// 100个生产氧原子的goroutine
	for i := 0; i < N; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			h2o.oxygen(releaseOxygen)
			wg.Done()
		}()
	}

	//等待所有的goroutine执行完
	wg.Wait()

	// 结果中肯定是300个原子
	if len(ch) != N*3 {
		fmt.Println(fmt.Errorf("expect %d atom but got %d", N*3, len(ch)))
	}

	// 每三个原子一组，分别进行检查。要求这一组原子中必须包含两个氢原子和一个氧原子，这样才能正确组成一个水分子。
	var s = make([]string, 3)
	for i := 0; i < N; i++ {
		s[0] = <-ch
		s[1] = <-ch
		s[2] = <-ch
		sort.Strings(s)


		water := s[0] + s[1] + s[2]
		if water != "HHO" {
			fmt.Println(fmt.Errorf("expect a water molecule but got %s", water))
		}

		fmt.Println(water)
	}
}
