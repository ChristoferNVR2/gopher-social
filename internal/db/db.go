package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/ChristoferNVR2/gopher-social/internal/env"
)

func New(addr string, MaxOpenConns, MaxIdleConns int, MaxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)

	duration, _ := env.GetDuration(MaxIdleTime, 15*time.Minute)
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
