package models

type Session struct {
	ID        int
	UserId    int
	TokenHash string
}

// type SessionService struct {
// 	DB *sql.DB
// }
