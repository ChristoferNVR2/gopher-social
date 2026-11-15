package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ChristoferNVR2/gopher-social/internal/models"
)

type RolesStore struct {
	db *sql.DB
}

func (s *RolesStore) GetByName(ctx context.Context, name string) (*models.Role, error) {
	const query = `
		SELECT id, name, description, level
		FROM roles
		WHERE name = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var role models.Role

	err := s.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.Level,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &role, nil
}
