package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Bio  string
	Age  int
}

func main() {

	// Filepath is relative to where you're running
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	// Defining like this is an anonymous struct, declared inline.
	// Second set of curly braces is instantiating the struct
	user := User{
		Name: "Jon Calhoun",
		Bio:  `<script>alert("Haha, you have been h4x0r3d!");</script>`,
		Age:  123,
	}

	// Execute is how you process a template
	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}

}
