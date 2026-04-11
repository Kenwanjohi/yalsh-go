package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateKey	 = errors.New("duplicate key")
)

type Models struct {
	Links LinkModel
	Users UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Links: LinkModel{DB: db},
		Users: UserModel{DB: db},
	}
}
