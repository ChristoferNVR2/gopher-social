package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ChristoferNVR2/gopher-social/internal/models"
	"github.com/lib/pq"
)

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *models.Post) error {
	const query = `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostsStore) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	const query = `
		SELECT id, content, title, user_id, tags, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	var post models.Post

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (s *PostsStore) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM posts
		WHERE id = $1;
	`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostsStore) Update(ctx context.Context, post *models.Post) error {
	const query = `
		UPDATE posts
		SET title = $1, content $2
		WHERE id = $3
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		post.Title,
		post.Title,
		post.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
