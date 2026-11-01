package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ChristoferNVR2/gopher-social/internal/models"
)

var (
	ErrNotFound          = errors.New("record not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = 5 * time.Second
)

type Storage struct {
	Posts interface {
		Create(context.Context, *models.Post) error
		GetByID(context.Context, int64) (*models.Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *models.Post) error
		GetUserFeed(context.Context, int64) ([]models.PostWithMetadata, error)
	}
	Users interface {
		Create(context.Context, *models.User) error
		GetByID(context.Context, int64) (*models.User, error)
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]models.CommentWithUser, error)
		Create(context.Context, *models.Comment) error
	}
	Followers interface {
		Follow(context.Context, int64, int64) error
		Unfollow(context.Context, int64, int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostsStore{db},
		Users:     &UsersStore{db},
		Comments:  &CommentsStore{db},
		Followers: &FollowersStore{db},
	}
}
