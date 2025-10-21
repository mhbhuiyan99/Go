package models

import (
	"database/sql"
)

type Models struct {
	Users UsersModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Users: UsersModel{DB: db},
	}
}