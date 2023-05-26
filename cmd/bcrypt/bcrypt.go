package main

import (
	"fmt"
	"os"
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
}

func compare(password, hash string) {
	fmt.Printf("Compare the PW %q with the hash %q\n", password, hash)
}
