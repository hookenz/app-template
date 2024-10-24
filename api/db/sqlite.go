package db

import (
	"fmt"

	"github.com/hookenz/app-template/api/utils/hash"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const createUserSQL = `CREATE TABLE IF NOT EXISTS user (
    id TEXT PRIMARY KEY NOT NULL,
    email TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL
);`

const createSessionSQL = `CREATE TABLE IF NOT EXISTS session (
	id TEXT UNIQUE PRIMARY KEY NOT NULL,
	user_id TEXT NOT NULL,
	ip_address TEXT,
	active INTEGER DEFAULT 1,
	last_activity DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES user(id) 
)`

type SqliteStore struct {
	filename string
	db       *sqlx.DB
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

func (s *SqliteStore) SelectUser(email string) (UserRecord, error) {
	row := s.db.QueryRow(`SELECT id, email, password from USER WHERE email = ?`, email)
	user := UserRecord{}
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	return user, err
}

func (s *SqliteStore) InsertUser(email, password string) error {
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("error generating user id: %w", err)
	}

	hash, err := hash.Create(password)
	if err != nil {
		return fmt.Errorf("error creating password hash: %w", err)
	}

	_, err = s.db.Exec(`INSERT INTO user (id, email, password) 
						VALUES (?, ?, ?)`, id, email, hash)
	return err
}

func (s *SqliteStore) ChangeUserPassword(email, password string) error {
	_, err := s.db.Exec(`UPDATE user SET password = ? WHERE email = ?`, password, email)
	return err
}

func (s *SqliteStore) CreateSession(userId, ipAddress string) (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("error generating session id: %w", err)
	}

	_, err = s.db.Exec(`INSERT INTO session (id, user_id, ip_address) 
						VALUES (?, ?, ?)`, id, userId, ipAddress)
	return id.String(), err
}

func (s *SqliteStore) GetSession(id string) (SessionRecord, error) {
	var session SessionRecord
	err := s.db.Get(&session, `SELECT id, user_id, ip_address, active, last_activity 
								FROM session WHERE id = ?`, id)
	return session, err
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
