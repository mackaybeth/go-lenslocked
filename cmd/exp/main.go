package main

import (
	"errors"
	"fmt"
)

func main() {
	err := B()
	if errors.Is(err, ErrNotFound) {
		fmt.Println("not found")
	}
	// TODO: Determine if the `err` variable is an `ErrNotFound`
}

var ErrNotFound = errors.New("not found")

func A() error {
	return ErrNotFound
}

func B() error {
	err := A()
	if err != nil {
		return fmt.Errorf("b: %w", err)
	}
	return nil
}
