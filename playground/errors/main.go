package main

import (
	"errors"
	"fmt"
)

func main() {
	err := thatThrowsError()
	if err != nil {
		panic(err)
	}
}

func thatThrowsError() error {
	err := fmt.Errorf("some error that has a value of %v", 666)
	return errors.Join(err, fmt.Errorf("Error number two %v", 777))
}
