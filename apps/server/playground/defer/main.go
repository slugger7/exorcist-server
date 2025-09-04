package main

import "fmt"

func main() {
	fmt.Println("Before defer main")
	defer fmt.Println("Defer in main")
	fmt.Println("After defer in main")

	deferredInFunction()
}

func deferredInFunction() {
	fmt.Println("Before defer en function")
	defer fmt.Println("Defferred in function") // will run once funcion exits
	fmt.Println("After defer in function")
}
