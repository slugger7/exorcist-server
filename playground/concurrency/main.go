package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	numberChan := make(chan int)
	wg.Add(3)
	go routine(&wg, 1, numberChan)
	go routine(&wg, 2, numberChan)

	numberChan <- 666
	numberChan <- 777
	close(numberChan)
	go routine(&wg, 3, numberChan)
	wg.Wait()
}

func routine(wg *sync.WaitGroup, routine int, ch chan int) {
	defer wg.Done()
	time.Sleep(time.Duration(routine) * time.Second)
	num, ok := <-ch
	if !ok {
		fmt.Println("Channel closed")
		return
	}
	fmt.Printf("Routine %v with num %v\n", routine, num)
}
