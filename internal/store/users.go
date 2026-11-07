package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ChristoferNVR2/gopher-social/internal/models"
)

var (
	ErrDuplicatedUsername = errors.New("username already exists")
	ErrDuplicatedEmail    = errors.New("email already exists")
)

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, tx *sql.Tx, user *models.User) error {
	const query = `
		INSERT INTO users (username, password, email)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password.Hash,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicatedEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicatedUsername
		default:
			return err
		}
	}

	return nil
}

func (s *UsersStore) GetByID(ctx context.Context, id int64) (*models.User, error) {
	const query = `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user models.User

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (s *UsersStore) CreateAndInvite(ctx context.Context, user *models.User, token string, invitationExp time.Duration) error {
	return withTx(ctx, s.db, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := s.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s *UsersStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, invitationExp time.Duration, userID int64) error {
	const query = `
		INSERT INTO user_invitations (user_id, token, expiry)
		VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID, token, time.Now().Add(invitationExp))
	if err != nil {
		return err
	}

	return nil
}
