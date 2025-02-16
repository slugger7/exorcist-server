package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ch := make(chan bool)

	var writerWg sync.WaitGroup
	go reader(ch, &wg)
	fmt.Println("Sync start")
	ch <- true
	ch <- true
	fmt.Println("Sync done")

	fmt.Println("Stacked start")
	writerWg.Add(2)
	go fireAndForget(ch, &writerWg)
	go fireAndForget(ch, &writerWg)
	fmt.Println("Stacked end")

	writerWg.Wait()
	fmt.Println("Writers done")
	close(ch)
	wg.Wait()
	fmt.Println("Exiting")
}

func fireAndForget(ch chan bool, wg *sync.WaitGroup) {
	fmt.Println("firing")
	ch <- true
	fmt.Println("forgetting")
	wg.Done()
}

func reader(ch chan bool, wg *sync.WaitGroup) {
	for {
		select {
		case ch, ok := <-ch:
			if !ok {
				fmt.Println("Channel closed. Exititing")
				wg.Done()
				return
			}
			fmt.Println("Processing", ch)
			time.Sleep(1 * time.Second)
			fmt.Println("Done processing")
		}
	}
}
