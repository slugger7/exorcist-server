package errors

import "fmt"

func CheckError(err error) {
	if err != nil {
		fmt.Println("Found an error")
		panic(err)
	}
}
