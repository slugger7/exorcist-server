package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(1)

	go consumer(ch, &wg)

	fmt.Println("Pushing to channel first time")
	ch <- true
	fmt.Println("Pushing a second time")
	ch <- true

	close(ch)
	wg.Wait()
}

func consumer(ch chan bool, wg *sync.WaitGroup) {
	for {
		time.Sleep(time.Second)
		select {
		case val, ok := <-ch:
			if !ok {
				fmt.Println("Channel closed exiting")
				wg.Done()
				return
			}

			fmt.Println("Consumed value", val)
		}
	}
}
