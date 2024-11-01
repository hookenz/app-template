package db

import "time"

type Database interface {
	Open() error

	InsertUser(email, password string) error
	SelectUser(email string) (UserRecord, error)

	CreateSession(userId, ipAddress string) (*SessionRecord, error)
	GetSession(id string) (*SessionRecord, error)
}

type UserRecord struct {
	Id       string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type SessionRecord struct {
	Id           string    `db:"id"`
	UserId       string    `db:"user_id"`
	IpAddress    string    `db:"ip_address"`
	Active       bool      `db:"active"`
	LastActivity time.Time `db:"last_activity"`
}
