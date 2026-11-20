package cache

import (
	"context"

	"github.com/ChristoferNVR2/gopher-social/internal/models"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (m *MockUserStore) Get(ctx context.Context, userID int64) (*models.UserWithRole, error) {
	return &models.UserWithRole{User: &models.User{}}, nil
}

func (m *MockUserStore) Set(ctx context.Context, user *models.UserWithRole) error {
	return nil
}
