package handler

import (
	"github.com/hookenz/app-template/api/db"
)

type Handler struct {
	db db.Database
}

func NewHandler(db db.Database) *Handler {
	return &Handler{
		db: db,
	}
}
