package user

import (
	"fmt"

	"github.com/hookenz/app-template/api/db"
	"github.com/hookenz/app-template/api/utils/hash"
)

type UserView struct {
	Name      string
	SessionID string
}

func Authenticate(db db.Database, username, password string) (UserView, error) {
	user := UserView{}
	record, err := db.SelectUser(username)
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

	return user, nil
}
