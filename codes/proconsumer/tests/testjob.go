package main

import (
	"bufio"
	"example.com/proconsumer"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	 ts := proconsumer.NewTaskScheduler(ProducerFunc, ConsumerFunc, 10)
	err := ts.Execute()
	fmt.Println(err, "11111")
}

func ProducerFunc(producer *proconsumer.Producer) error {
	dir, _ := os.Getwd()
	fmt.Println(dir)
	// 自己试一下故意写错文件路径
	file, err := os.Open(dir + "/tests/file.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {

		select {
		case <-producer.CancelCtx.Done():
			fmt.Println("producer end")
			return nil
		default:
		}
		// TODO 制造消息
		str := scanner.Text()
		i += 1
		if i % 100 == 0 {
			// 模拟耗时查询
			time.Sleep(time.Second * 3)
		}
		// TODO 制造错误情况, 测试下退出
		//if strings.HasPrefix(str, "rt") {
		//	err = fmt.Errorf("producer Error")
		//	return err
		//}
		// 发送数据给consumer
		producer.DataCh <- str
	}

	close(producer.DataCh) //查询完成关闭通道

	return nil
}

func ConsumerFunc(consumer *proconsumer.Consumer) error {
	for {

		select {
		case <-consumer.CancelCtx.Done():
			fmt.Println("consumer cancel end")
			return nil
		case v, ok := <-consumer.DataCh:
			if !ok {
				fmt.Println("consumer execution end")
				return nil
			}
			err := Do(v)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func Do(v interface{}) error {
	val, _ := v.(string)
	// TODO 制造错误情况, 测试下退出
	if strings.HasPrefix(val, "rt") {
		return fmt.Errorf("worker error")
	}
	fmt.Println(strings.ToUpper(val))
	return nil
}
