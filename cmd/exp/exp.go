package main

import (
	"fmt"

	"github.com/mackaybeth/lenslocked/models"
)

// All of these values will vary depending on your mail service.
const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 587
	username = "da5feb4b8b5f65"
	password = "9f1778ccdb80cc"
)

func main() {

	// email := models.Email{
	// 	From:      "test@lenslocked.com",
	// 	To:        "jon@calhoun.io",
	// 	Subject:   "This is a test email",
	// 	Plaintext: "This is the body of the email",
	// 	HTML:      `<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>`,
	// }

	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})

	// err := es.Send(email)
	err := es.ForgotPassword("jon@calhoun.io", "https://lenslocked.com/reset-pw?token=123rty")
	if err != nil {
		panic(err)
	}
	fmt.Println("Reset PW Email sent")
	// msg := mail.NewMessage()
	// msg.SetHeader("To", to)
	// msg.SetHeader("From", from)
	// msg.SetHeader("Subject", subject)
	// msg.SetBody("text/plain", plaintext)
	// msg.AddAlternative("text/html", html)

	// // Takes the contents of the writer and writes them out to stdout
	// msg.WriteTo(os.Stdout)

	// dialer := mail.NewDialer(host, port, username, password)

	// err := dialer.DialAndSend(msg)
	// if err != nil {
	// 	// TODO: Handle the error correctly
	// 	panic(err)
	// }
	// fmt.Println("message sent")
}

// func main() {
// ctx := stdctx.Background()

// user := models.User{
// 	Email: "jon@calhoun.io",
// }
// ctx = context.WithUser(ctx, &user)

// retrievedUser := context.User(ctx)
// fmt.Println(retrievedUser.Email)

// cfg := models.DefaultPostgresConfig()
// db, err := models.Open(cfg)
// if err != nil {
// 	panic(err)
// }
// defer db.Close()
// err = db.Ping()
// if err != nil {
// 	panic(err)
// }
// fmt.Println("Connected!")

// _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
// 	id SERIAL PRIMARY KEY,
// 	name TEXT,
// 	email TEXT NOT NULL
//   );

//   CREATE TABLE IF NOT EXISTS orders (
// 	id SERIAL PRIMARY KEY,
// 	user_id INT NOT NULL,
// 	amount INT,
// 	description TEXT
//   );
// `)
// if err != nil {
// 	panic(err)
// }
// fmt.Println("Tables created.")

// us := models.UserService{
// 	DB: db,
// }
// user, err := us.Create("jon@calhoun.io", "joncalhoun")
// if err != nil {
// 	panic(err)
// }
// fmt.Println(user)

// user, err = us.Create("bob@bob.com", "bob123")
// if err != nil {
// 	panic(err)
// }
// fmt.Println(user)

// user, err = us.Create("fake@calhoun.com", "fake123")
// if err != nil {
// 	panic(err)
// }
// fmt.Println(user)

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

// userID := 1
// for i := 1; i <= 5; i++ {
// 	amount := i * 100
// 	desc := fmt.Sprintf("Fake order #%d", i)
// 	_, err := db.Exec(`
// 	  INSERT INTO orders(user_id, amount, description)
// 	  VALUES($1, $2, $3)`, userID, amount, desc)
// 	if err != nil {
// 		panic(err)
// 	}
// }
// fmt.Println("Created fake orders.")

// type Order struct {
// 	ID          int
// 	UserID      int
// 	Amount      int
// 	Description string
// }

// var orders []Order

// userID := 1
// // Returns 0 to N rows into the 'rows' object
// rows, err := db.Query(`
// 	SELECT id, amount, description
// 	FROM orders
// 	WHERE user_id=$1`, userID)
// if err != nil {
// 	panic(err)
// }
// defer rows.Close()

// // rows doesn't point to anything at first, need to call Next to load at least the first one
// for rows.Next() {

// 	var order Order
// 	order.UserID = userID
// 	err := rows.Scan(&order.ID, &order.Amount, &order.Description)
// 	if err != nil {
// 		panic(err)
// 	}
// 	orders = append(orders, order)
// }
// err = rows.Err()
// if err != nil {
// 	panic(err)
// }

// fmt.Println("Orders:", orders)

// }
