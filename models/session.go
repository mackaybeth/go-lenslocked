package models

import "database/sql"

type Session struct {
	ID     int
	UserId int
	// Token is only set when creating a new session.
	// When looking up a session, this will be left empty as we only
	// store the hash of a session token on our databases and we cannot
	// reverse it into a raw token.
	Token     string // This field is not in the DB
	TokenHash string
}
type SessionService struct {
	DB *sql.DB
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	// TODO: create the session token
	// TODO: implement SessionService.Create
	return nil, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// TODO: implement SessionService.User
	return nil, nil
}
