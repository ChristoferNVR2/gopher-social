package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/ChristoferNVR2/gopher-social/internal/models"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (m *MockUserStore) Create(ctx context.Context, tx *sql.Tx, user *models.User) error {
	return nil
}

func (m *MockUserStore) GetByID(ctx context.Context, id int64) (*models.UserWithRole, error) {
	return &models.UserWithRole{User: &models.User{ID: id}}, nil
}

func (m *MockUserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return &models.User{}, nil
}

func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *models.User, token string, inviteExp time.Duration) error {
	return nil
}

func (m *MockUserStore) Activate(ctx context.Context, token string) error {
	return nil
}

func (m *MockUserStore) Delete(ctx context.Context, id int64) error {
	return nil
}
