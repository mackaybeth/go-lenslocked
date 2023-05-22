package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	// pass in driver and conenction string
	// port matches what's in my docker compose file for the "port on my computer"
	db, err := sql.Open("pgx", "host=localhost port=5433 user=baloo password=junglebook dbname=lenslocked sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")
}
