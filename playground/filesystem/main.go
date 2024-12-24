package main

import (
	"fmt"
	"log"

	"github.com/slugger7/exorcist/internal/media"
)

func main() {
	fmt.Println("Finding values")
	values, err := media.GetFilesByExtensions(".", []string{".mp4", ".m4v", ".mkv", ".avi", ".wmv", ".flv", ".webm", ".f4v", ".mpg", ".m2ts", ".mov"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Printing out results")
	for _, v := range values {
		fmt.Println(v)
	}
}
