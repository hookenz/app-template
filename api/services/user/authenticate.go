package user

import (
	"fmt"

	"github.com/hookenz/app-template/api/db"
	"github.com/hookenz/app-template/api/utils/hash"
)

type UserView struct {
	Name      string
	SessionID string
	UserID    string
}

func Authenticate(db db.Database, email, password string) (UserView, error) {
	user := UserView{}
	record, err := db.SelectUser(email)
	if err != nil {
		return user, err
	}

	match, err := hash.Compare(password, record.Password)
	if err != nil {
		return user, fmt.Errorf("authentication failure")
	}

	if !match {
		return user, fmt.Errorf("authentication failure")
	}

	user.Name = record.Email
	user.UserID = record.ID
	return user, nil
}
