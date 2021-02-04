package context

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

// this program demonstrates how to catch ctrl+C during tests
func Program() {
	// An initial context to do something with it
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		select {
		case <-c:
			cancel()
		case <-time.After(10 * time.Second):
			cancel()
		}
	}()

	Do(ctx)
}

func Do(ctx context.Context) {
	counter := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("I should clean something here")
			break
		default:
			fmt.Printf("hardworking in %d seconds already\n", counter)
			counter++
			time.Sleep(1 * time.Second)
		}
	}
}
