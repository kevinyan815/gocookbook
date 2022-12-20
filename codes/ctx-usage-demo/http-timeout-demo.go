package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	// 创建一个监听8000端口的HTTP Server
	http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// 开始处理请求
		fmt.Fprint(os.Stdout, "processing request\n")
		// select 接收定时和ctx.Done两个channel，哪个先来执行哪个。
		select {
		case <-time.After(2 * time.Second):
			// 用定时器模拟请求处理成功，如果2s后接收到了channel里的信息，返回请求成功
			w.Write([]byte("request processed"))
		case <-ctx.Done():
			// 如果请求被取消了，返回请求被取消
			fmt.Fprint(os.Stderr, "request cancelled\n")
		}
	}))
}
