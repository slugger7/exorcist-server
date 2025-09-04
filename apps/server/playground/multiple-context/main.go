package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(2)
	go doWork(ctx, &wg)
	go doWork(ctx, &wg)

	time.Sleep(time.Second * 3)

	cancel()

	wg.Wait()
}

func doWork(ctx context.Context, wg *sync.WaitGroup) {
	fmt.Println("Doing work")

	<-ctx.Done()

	fmt.Println("Done doing work because context is done")
	wg.Done()
}
