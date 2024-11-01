package session

import (
	"time"

	"github.com/hookenz/app-template/api/db"
)

type Session struct {
	Id           string
	UserId       string
	IpAddress    string
	LastActivity time.Time
}

func New(db db.Database, userID, ipAddress string) (*Session, error) {
	sr, err := db.CreateSession(userID, ipAddress)
	if err != nil {
		return nil, err
	}

	return &Session{
		Id:           sr.Id,
		UserId:       sr.UserId,
		IpAddress:    sr.IpAddress,
		LastActivity: sr.LastActivity,
	}, nil
}

func Get(userID string) (Session, error) {
	return Session{}, nil
}

func (s *Session) Expired() bool {
	return false
}

func (s *Session) Delete() error {
	return nil
}
