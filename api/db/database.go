package db

type Database interface {
	Open() error

	InsertUser(name, password string) error
	SelectUser(name string) (UserRecord, error)

	CreateSession() (string, error)
	GetSession(id string) (bool, error)
}

type UserRecord struct {
	ID       string
	Name     string
	Password string
}
