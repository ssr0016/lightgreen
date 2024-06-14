package data

import (
	"database/sql"
	"fmt"
)

var (
	ErrRecordNotFound = fmt.Errorf("record not found")
	ErrEditConflict   = fmt.Errorf("edit conflict")
)

type Models struct {
	Movies MovieModel
	Users  UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
		Users:  UserModel{DB: db},
	}
}

func NewMockModels() Models {
	return Models{
		// Movies: MockMovieModel{},
	}
}
