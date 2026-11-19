package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ChristoferNVR2/gopher-social/internal/models"
	"github.com/redis/go-redis/v9"
)

type UsersStore struct {
	rdb *redis.Client
}

const userExpTime = 1 * time.Minute

func (s *UsersStore) Get(ctx context.Context, userID int64) (*models.UserWithRole, error) {
	cacheKey := fmt.Sprintf("user-%d", userID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, err
	}

	var user models.UserWithRole
	if data != "" {
		err = json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (s *UsersStore) Set(ctx context.Context, user *models.UserWithRole) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}
	cacheKey := fmt.Sprintf("user-%d", user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.Set(ctx, cacheKey, json, userExpTime).Err()
}
