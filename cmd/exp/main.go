package main

import (
	"fmt"
	"html/template"
	"os"
)

type User struct {
	Name string
	Age  int
	Meta UserMeta
}

type UserMeta struct {
	Visits int
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
		Name: "Susan Smith",
		Age:  111,
		Meta: UserMeta{
			Visits: 4,
		},
	}

	fmt.Println(user.Meta.Visits)

	// Execute is how you process a template
	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}

}
