package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	noMoreData := false
	concurrentNum := 10

	for {
		if noMoreData {
			break
		}

		var wg sync.WaitGroup
		for i := 0; i < concurrentNum; i++ {
			wg.Add(1)
			go func(i int) {
				time.Sleep(2 * time.Second)
				fmt.Println("End:", i)
				wg.Done()
			}(i)
		}
		wg.Wait()

		time.Sleep( 200 * time.Millisecond)
		fmt.Println("Next bunch of things")
	}
}
