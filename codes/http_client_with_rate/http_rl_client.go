package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

//RLHTTPClient Rate Limited HTTP Client
type RLHTTPClient struct {
	client      *http.Client
	Ratelimiter *rate.Limiter
}

//Do dispatches the HTTP request to the network
func (c *RLHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Comment out the below 5 lines to turn off ratelimiting
	ctx := context.Background()
	err := c.Ratelimiter.Wait(ctx) // This is a blocking call. Honors the rate limit
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//NewClient return http client with a ratelimiter
func NewClient(rl *rate.Limiter) *RLHTTPClient {
	c := &RLHTTPClient{
		client:      http.DefaultClient,
		Ratelimiter: rl,
	}
	return c
}

func main() {
	rl := rate.NewLimiter(rate.Every(10*time.Second), 50) // 50 request every 10 seconds
	c := NewClient(rl)
	reqURL := "https://api.btcmarkets.net/v3/markets/BTC-AUD/ticker"
	req, _ := http.NewRequest("GET", reqURL, nil)
	for i := 0; i < 300; i++ {
		resp, err := c.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(resp.StatusCode)
			return
		}
		if resp.StatusCode == 429 {
			fmt.Printf("Rate limit reached after %d requests", i)
			return
		}
	}
}
