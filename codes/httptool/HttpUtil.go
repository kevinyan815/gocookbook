package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)



func httpRequest(method string, url string, data string, headers map[string]string, timeout time.Duration) (code int, content string, err error) {
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if len(headers) != 0 {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}
	if err != nil {
		return
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: timeout}
	resp, error := client.Do(req)
	if error != nil {
		return 0, "", error
	}

	defer resp.Body.Close()
	code = resp.StatusCode
	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)

	return
}

// 发起GET请求
func HttpGet(url string, data string, timeout time.Duration) (code int, content string, err error) {
	return httpRequest("GET", url, data, nil, timeout )
}

func HttpGetWithHeaders(url string, data string, headers map[string]string, timeout time.Duration) (code int, content string, err error) {
	return httpRequest("GET", url, data, headers, timeout )
}

// 发起POST请求，请求头指定Content-Type: application/json
func HttpJsonPost(url string, data string, timeout time.Duration) (code int, content string, err error) {
	headers := map[string]string {
		"Content-Type": "application/json",
	}
	code, content, err = httpRequest("POST", url, data, headers, timeout)
	return
}

// 发起POST请求，请求头由调用方指定
func HttpPostWithHeaders(url string, data string, headers map[string]string, timeout time.Duration) (code int, content string, err error) {
	code, content, err = httpRequest("POST", url, data, headers, timeout)
	return
}

