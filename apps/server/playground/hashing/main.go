package main

import (
	"log"

	errs "github.com/slugger7/exorcist/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	errs.PanicError(err)

	log.Println(string(hashedPassword))
}
