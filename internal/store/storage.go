package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ChristoferNVR2/gopher-social/internal/models"
)

var (
	ErrNotFound = errors.New("record not found")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *models.Post) error
		GetByID(context.Context, int64) (*models.Post, error)
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
