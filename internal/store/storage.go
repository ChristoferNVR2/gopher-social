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
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]models.PostWithMetadata, error)
	}
	Users interface {
		GetByID(context.Context, int64) (*models.User, error)
		GetByEmail(context.Context, string) (*models.User, error)
		Create(context.Context, *sql.Tx, *models.User) error
		CreateAndInvite(context.Context, *models.User, string, time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]models.CommentWithUser, error)
		Create(context.Context, *models.Comment) error
	}
	Followers interface {
		Follow(context.Context, int64, int64) error
		Unfollow(context.Context, int64, int64) error
	}
	Roles interface {
		GetByName(context.Context, string) (*models.Role, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostsStore{db},
		Users:     &UsersStore{db},
		Comments:  &CommentsStore{db},
		Followers: &FollowersStore{db},
		Roles:     &RolesStore{db},
	}
}

func withTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.New("tx err: " + err.Error() + ", rb err: " + rbErr.Error())
		}
		return err
	}

	return tx.Commit()
}
