package store

import (
	"context"
	"database/sql"

	"github.com/ChristoferNVR2/gopher-social/internal/models"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *models.Post) error
	}
	Users interface {
		Create(context.Context, *models.User) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostsStore{db: db},
		Users: &UsersStore{db: db},
	}
}
