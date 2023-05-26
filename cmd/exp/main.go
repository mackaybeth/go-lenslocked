package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName, cfg.SSLMode)
}

func main() {

	cfg := PostgresConfig{
		Host:     "localhost",
		Port:     "5433",
		User:     "baloo",
		Password: "junglebook",
		DbName:   "lenslocked",
		SSLMode:  "disable",
	}

	// pass in driver and conenction string
	// port matches what's in my docker compose file for the "port on my computer"
	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT NOT NULL
	  );
	  
	  CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		amount INT,
		description TEXT
	  );
	`)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tables created.")

	// name := "New User"
	// email := "new@calhoun.io"
	// // Using QueryRow instead of Exec because we're expecting a return value
	// row := db.QueryRow(`
	// INSERT INTO users(name, email)
	// VALUES($1, $2) RETURNING id;`, name, email)

	// // Could call row.Err != nill here first, but if there is an error with the row, it will be returned with row.Scan.  So if using row.Scan, extra row.Err check isn't needed.
	// // row.Scan gets the RETURNING value (could have multiple RETURNING, order matters for row.Scan)
	// var id int
	// err = row.Scan(&id)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("User created. id =", id)

	// id := 5
	// row := db.QueryRow(`
	// 	SELECT name, email
	// 	FROM users
	// 	WHERE id=$1;`, id)
	// var name, email string
	// err = row.Scan(&name, &email)

	// // QueryRow expects at least one row back and returns the first.  If there are no rows, it returns an error so we can check for that
	// if err == sql.ErrNoRows {
	// 	fmt.Println("Error, no rows!")
	// }
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("User information: name=%s, email=%s\n", name, email)

	userID := 1
	for i := 1; i <= 5; i++ {
		amount := i * 100
		desc := fmt.Sprintf("Fake order #%d", i)
		_, err := db.Exec(`
		  INSERT INTO orders(user_id, amount, description)
		  VALUES($1, $2, $3)`, userID, amount, desc)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Created fake orders.")

}
