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

	name := "New User"
	email := "new@calhoun.io"
	row := db.QueryRow(`
	INSERT INTO users(name, email)
	VALUES($1, $2) RETURNING id;`, name, email)

	// Could call row.Err != nill here first, but if there is an error with the row, it will be returned with row.Scan.  So if using row.Scan, extra row.Err check isn't needed.
	// row.Scan gets the RETURNING value (could have multiple RETURNING, order matters for row.Scan)
	var id int
	err = row.Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("User created. id =", id)

}
