package main

import (
	"fmt"

	errs "github.com/slugger7/exorcist/internal/errors"
)

func main() {
	err := nestError(4)
	if err != nil {
		panic(err)
	}
}

func nestError(count int) error {
	if count != 0 {
		return errs.BuildError(nestError(count-1), "erorr at count %v", count)
	}
	return fmt.Errorf("base error")
}
