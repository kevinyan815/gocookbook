package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var queue []struct{}

func main() {
	var wg sync.WaitGroup
	c := sync.NewCond(&sync.Mutex{})

	for i := 0; i < 2; i ++ {
		wg.Add(1)

		go func(i int) {
			// this go routine wait for changes to the sharedRsc
			c.L.Lock()
			for len(queue) <= 0 {
				fmt.Println("goroutine" + strconv.Itoa(i) +" wait")
				c.Wait()
			}
			fmt.Println("goroutine" + strconv.Itoa(i), "pop data")
			queue = queue[1:]
			c.L.Unlock()
			wg.Done()
		}(i)

	}


	for i := 0; i < 2; i ++ {
		// 主goroutine延迟两秒准备好后把变量设置为true
		time.Sleep(2 * time.Second)
		c.L.Lock()
		fmt.Println("main goroutine push data")
		queue= append(queue, struct{}{})
		c.Broadcast()
		fmt.Println("main goroutine broadcast")
		c.L.Unlock()

	}

	wg.Wait()
}
