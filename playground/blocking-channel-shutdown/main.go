package main

import (
	"context"
	"fmt"
	"sync"
)

func blockingChannel(ch chan bool, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("In Blocking channel")
	// Working example
	select {
	case <-ctx.Done():
		fmt.Println("Context done returning")
		return
	default:
		fmt.Println("About to write to channel")
		ch <- true
		fmt.Println("Written to channel")
	}

	// Broken example
	// fmt.Println("About to write to channel")
	// ch <- true
	// fmt.Println("Written to channel")

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan bool)
	var wg sync.WaitGroup

	fmt.Println("Creating first blocking channel goroutine")
	wg.Add(1)
	go blockingChannel(ch, ctx, &wg)

	fmt.Println("Creating second blocking channel goroutine")
	wg.Add(1)
	go blockingChannel(ch, ctx, &wg)

	fmt.Println("Reading from channel for the first time")
	<-ch

	fmt.Println("Cancelling context")
	cancel()

	wg.Wait()
}
