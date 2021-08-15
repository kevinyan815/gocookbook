package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	messages := make(chan int, 10)

	// producer
	for i := 0; i < 10; i++ {
		messages <- i
	}
	valCtx := context.WithValue(context.Background(), "client-ip", "10.96.130.167")
	ctx, cancel := context.WithTimeout(valCtx, 5*time.Second)

	// consumer
	go func(ctx context.Context) {
		ticker := time.NewTicker(1 * time.Second)
		for _ = range ticker.C {
			select {
			case <-ctx.Done():
				fmt.Println("child process interrupt...")
				return
			default:
				i := <-messages
				timeoutCtx, _ := context.WithTimeout(ctx, 10 * time.Second)
				fmt.Printf("child receive message: %d\n", 1)
				go anotherFunc(timeoutCtx, i)
			}
		}
	}(ctx)

	defer close(messages)
	defer cancel()

	select {
	case <-ctx.Done():
		time.Sleep(1 * time.Second)
		fmt.Println("main process exit!")
	}
}

func anotherFunc(ctx context.Context, message int) {
	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		select {
		case <-ctx.Done():
			fmt.Println("descendants process interrupt...")
			return
		default:
			fmt.Printf("host-ip: %s, send message %d\n", ctx.Value("client-ip"), message)
		}
	}
}
