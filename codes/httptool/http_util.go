package httptool

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func httpRequest(method string, url string, options ...Option) (code int, content string, err error) {
	start := time.Now()

	reqOpts := defaultRequestOptions() // 默认的请求选项
	for _, opt := range options {      // 在reqOpts上应用通过options设置的选项
		err = opt.apply(reqOpts); if err != nil {
			return
		}
	}
	// 创建请求对象
	req, err := http.NewRequest(method, url, strings.NewReader(reqOpts.data))
	// 记录请求日志
	defer func() {
		dur := int64(time.Since(start) / time.Second)
		if dur >= 500 {
			log.Println("Http"+method, url, "in", reqOpts.data, "out", content, "err", err, "dur/ms", dur)
		} else {
			log.Println("Http"+method, url, "in", reqOpts.data, "out", content, "err", err, "dur/ms", dur)
		}
	}()
	defer req.Body.Close()

	if len(reqOpts.headers) != 0 { // 设置请求头
		for key, value := range reqOpts.headers {
			req.Header.Add(key, value)
		}
	}
	if err != nil {
		return
	}
	// 发起请求
	client := &http.Client{Timeout: reqOpts.timeout}
	resp, error := client.Do(req)
	if error != nil {
		return 0, "", error
	}
	// 解析响应
	defer resp.Body.Close()
	code = resp.StatusCode
	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)

	return
}

// 发起GET请求
func HttpGet(url string, options ...Option) (code int, content string, err error) {
	return httpRequest("GET", url, options...)
}

// 发起POST请求，请求头指定Content-Type: application/json
func HttpJsonPost(url string, data string, timeout time.Duration) (code int, content string, err error) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	code, content, err = httpRequest(
		"POST", url, WithTimeout(timeout), WithHeaders(headers), WithData(data))

	return
}

// 针对可选的HTTP请求配置项，模仿gRPC使用的Options设计模式实现
type requestOption struct {
	timeout time.Duration
	data    string
	headers map[string]string
}

type Option interface {
	apply(option *requestOption) error
}

type optionFunc func(option *requestOption) error

func (f optionFunc) apply(opts *requestOption) error {
	return f(opts)
}


func defaultRequestOptions() *requestOption {
	return &requestOption{ // 默认请求选项
		timeout: 5 * time.Second,
		data:    "",
		headers: nil,
	}
}

func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		opts.timeout, err = timeout, nil
		return
	})
}

func WithHeaders(headers map[string]string) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		opts.headers, err = headers, nil
		return
	})
}

func WithData(data string) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		opts.data, err = data, nil
		return
	})
}
