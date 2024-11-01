package user

import (
	"fmt"

	"github.com/hookenz/app-template/api/db"
	"github.com/hookenz/app-template/api/utils/hash"
)

type UserView struct {
	Id        string
	Email     string
	SessionId string
}

func Authenticate(db db.Database, email, password string) (*UserView, error) {
	record, err := db.SelectUser(email)
	if err != nil {
		return nil, err
	}

	match, err := hash.Compare(password, record.Password)
	if err != nil {
		return nil, fmt.Errorf("authentication failure")
	}

	if !match {
		return nil, fmt.Errorf("authentication failure")
	}

	return &UserView{Id: record.Id, Email: record.Email}, nil
}
