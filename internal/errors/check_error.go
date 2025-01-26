package errors

import (
	"log"
)

func CheckError(err error) {
	if err != nil {
		log.Println("Found an error")
		panic(err)
	}
}
