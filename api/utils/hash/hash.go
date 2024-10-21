package hash

import (
	"github.com/alexedwards/argon2id"
)

func Create(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func Compare(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
