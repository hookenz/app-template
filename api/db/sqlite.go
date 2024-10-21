package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hookenz/moneygo/api/utils/hash"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/michaeljs1990/sqlitestore"
)

const maxAge = 3600

const createUserSQL = `CREATE TABLE IF NOT EXISTS user (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL
);`

const createSessionSQL = `CREATE TABLE IF NOT EXISTS session (
	id TEXT UNIQUE PRIMARY KEY NOT NULL
)`

const selectUserSQL = `SELECT id, name, password from USER WHERE name = ?`
const insertUserSQL = `INSERT INTO USER (id, name, password) VALUES (?, ?, ?)`
const insertSessionSQL = `INSERT INTO SESSION (id) VALUES (?)`
const selectSessionSQL = `SELECT id from SESSION where id = ?`

type SqliteStore struct {
	filename     string
	db           *sqlx.DB
	sessionStore *sqlitestore.SqliteStore
}

func NewSqliteStore(filename string) Database {
	return &SqliteStore{
		filename: filename,
	}
}

func (s *SqliteStore) Open() error {
	if s.db != nil {
		return nil
	}

	var err error
	s.db, err = sqlx.Open("sqlite3", s.filename)
	if err != nil {
		return fmt.Errorf("error opening sqlite database file: %w", err)
	}

	err = s.createTables()
	if err != nil {
		return fmt.Errorf("error creating sqlite database: %w", err)
	}

	return nil
}

func (s *SqliteStore) SelectUser(name string) (UserRecord, error) {
	row := s.db.QueryRow(selectUserSQL, name)
	user := UserRecord{}
	err := row.Scan(&user.ID, &user.Name, &user.Password)
	return user, err
}

func (s *SqliteStore) InsertUser(name, password string) error {
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("error generating user id: %w", err)
	}

	hash, err := hash.Create(password)
	if err != nil {
		return fmt.Errorf("error creating password hash: %w", err)
	}

	_, err = s.db.Exec(insertUserSQL, id, name, hash)
	return err
}

func (s *SqliteStore) CreateSession() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("error generating session id: %w", err)
	}

	_, err = s.db.Exec(insertSessionSQL, id)
	return id.String(), nil
}

func (s *SqliteStore) GetSession(id string) (bool, error) {
	row := s.db.QueryRow(selectSessionSQL, id)
	var sessionid string
	err := row.Scan(&sessionid)
	found := err == nil
	return found, err
}

func (s *SqliteStore) createTables() error {
	err := s.createTableUser()
	if err != nil {
		return err
	}

	err = s.createTableSession()
	if err != nil {
		return err
	}
	return err
}

func (s *SqliteStore) createTableUser() error {
	_, err := s.db.Exec(createUserSQL)
	if err != nil {
		return fmt.Errorf("error creating table user: %w", err)
	}

	return nil
}

func (s *SqliteStore) createTableSession() error {
	_, err := s.db.Exec(createSessionSQL)
	if err != nil {
		return fmt.Errorf("error creating table session: %w", err)
	}

	return nil
}
