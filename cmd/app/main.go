package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"segment/pkg/config"
	"segment/pkg/delivery/router"
	"segment/pkg/repo"
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
	repo := repo.NewRepository(db)
	server := router.SetupRouter(repo)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), server)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Close(); err != nil {
		log.Fatalf("failed to close connection to database: %v", err)
	}
}
