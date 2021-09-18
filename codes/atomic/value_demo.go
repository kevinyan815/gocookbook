package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Rectangle struct {
	length int
	width  int
}

var rect atomic.Value

func update(width, length int) {
	rectLocal := new(Rectangle)
	rectLocal.width = width
	rectLocal.length = length
	rect.Store(rectLocal)
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(10)
	// 10 个协程并发更新
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			update(i, i+5)
		}()
	}
	wg.Wait()
	_r := rect.Load().(*Rectangle)
	fmt.Printf("rect.width=%d\nrect.length=%d\n", _r.width, _r.length)
}
