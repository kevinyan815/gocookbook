package main

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
	"time"
)

// 模拟一个Redis客户端
type client struct {
	// ... 其他的配置省略
	requestGroup singleflight.Group
}

// 普通查询
func (c *client) Get(key string) (interface{}, error) {
	fmt.Println("Querying Database")
	time.Sleep(time.Second)
	v := "Content of key" + key
	return  v, nil
}

// SingleFlight查询
func (c *client) SingleFlightGet(key string) (interface{}, error) {
	v, err, _ := c.requestGroup.Do(key, func() (interface{}, error) {
		return c.Get(key)

	})
	if err != nil {
		return nil, err
	}
	return v, err
}

func main() {
	redisClient := new(client)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {

		go func() {
			defer wg.Done()
			v, _ := redisClient.SingleFlightGet("Cyberpunk2077!!!")
			fmt.Println(v)
		}()
	}
	wg.Wait()
}
