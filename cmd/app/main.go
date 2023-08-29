package main

import (
	"context"
	"fmt"
	"log"
	"segment/pkg/config"
	"segment/pkg/storage"
	"segment/pkg/storage/postgres"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
)

func InitializeStorage(cfg *config.Config) (storage.IStorage, error) {

	db := postgres.NewStorage(cfg.Database)
	if err := db.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.Init().Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if db.Init().Error != nil {
		return nil, fmt.Errorf("failed to create extension: %w", db.Init().Error)
	}

	return db, nil
}

func main() {
	ctx := context.Background()
	var cfg config.Config
	err := confita.NewLoader(
		env.NewBackend(),
	).Load(ctx, &cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to parse config: %v", err))
	}
	log.Println(cfg)
	db, err := InitializeStorage(&cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to init database: %v", err))
	}

	err = db.MakeMigrations()
	if err != nil {
		log.Fatalf("failed to make migrations: %v", err)
	}

}
