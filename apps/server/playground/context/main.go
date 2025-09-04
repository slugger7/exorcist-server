package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	items := make(chan uuid.UUID, 3)
	var wg sync.WaitGroup
	wg.Add(1)
	go runJob(ctx, items, &wg)
	for i := 0; i <= 10; i = i + 1 {
		fmt.Println("For loop at ", i)
		if i == 5 {
			fmt.Println("Cancelling context")
			cancel()
			break
		}
		id, _ := uuid.NewRandom()
		fmt.Println("Pushing job", id)
		items <- id
	}
	fmt.Println("closing channel")
	close(items)
	wg.Wait()
}

func runJob(ctx context.Context, workItems chan uuid.UUID, wg *sync.WaitGroup) {
	fmt.Println("Waiting 1 second before running jobs")
	time.Sleep(time.Second * 1)
	for {
		fmt.Println("Starting job loop")
		select {
		case <-ctx.Done():
			fmt.Println("Done with what we were doing. Sleeping for 3 seconds")
			time.Sleep(time.Second * 3)
			fmt.Println("Done sleeping")
			wg.Done()
			return
		case item, ok := <-workItems:
			if !ok {
				fmt.Println("Channel closed")
				time.Sleep(time.Second * 2)
				fmt.Println("Done sleeping for 2 seconds")
				wg.Done()
				return
			}
			fmt.Println("Processing job", item)
			continue
		}
	}
}
