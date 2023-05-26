package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	switch os.Args[1] {
	case "hash":
		// hash the password
		hash(os.Args[2])

	case "compare":
		compare(os.Args[2], os.Args[3])

	default:
		fmt.Printf("Invalid command: %v\n", os.Args[1])
	}

}

func hash(password string) {
	fmt.Printf("todo: hash the password %q\n", password)

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("error hashing: %v\n", err)
		return
	}

	// bcrypt returns a base-64 enconded byte slice so we can just cast to a string
	hash := string(hashedBytes)
	fmt.Println(hash)

}

func compare(password, hash string) {
	fmt.Printf("Compare the PW %q with the hash %q\n", password, hash)

	// bcrypt.CompareHAshAndPassword

}
