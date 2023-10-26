package main

// request body不进行 body.Close() 操作导致移除的先决条件
// 既不执行 ioutil.ReadAll(resp.Body) 也不执行resp.Body.Close()，并且不设置http.Client内timeout的时候，就会导致协程泄露
func main() {
    tr := &http.Transport{
        MaxIdleConns:    100,
        IdleConnTimeout: 3 * time.Second,
    }

    n := 5
    for i := 0; i < n; i++ {
        req, _ := http.NewRequest("POST", "https://www.baidu.com", nil)
        req.Header.Add("content-type", "application/json")
        client := &http.Client{
            Transport: tr,
        }
        resp, _ := client.Do(req)
        _ = resp
    }
    time.Sleep(time.Second * 5)
    fmt.Printf("goroutine num is %d\n", runtime.NumGoroutine())
}
