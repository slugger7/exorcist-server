package main

import (
	"fmt"
	"strconv"
)

func main() {
	floatString := "716477.5564"

	intStringConverted, err := strconv.ParseFloat(floatString, 5)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Converted float string to int %v\n", intStringConverted)
}
