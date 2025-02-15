package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	items := make(chan uuid.UUID, 3)
	go runJob(ctx, items)
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
}

func runJob(ctx context.Context, workItems chan uuid.UUID) {
	fmt.Println("Waiting 1 second before running jobs")
	time.Sleep(time.Second * 1)
	for {
		fmt.Println("Starting job loop")
		select {
		case <-ctx.Done():
			fmt.Println("Done with what we were doing")
			return
		case item, ok := <-workItems:
			if !ok {
				fmt.Println("It was all good")
			}
			fmt.Println("Processing job", item)
			continue
		}
	}
}
