package main

import (
	"log"

	"github.com/ChristoferNVR2/gopher-social/internal/db"
	"github.com/ChristoferNVR2/gopher-social/internal/env"
	"github.com/ChristoferNVR2/gopher-social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/gopher-social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
