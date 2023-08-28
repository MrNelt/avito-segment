package main

import (
	"context"
	"fmt"
	"log"
	"segment/pkg/config"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
)

func main() {
	ctx := context.Background()
	var cfg config.Config
	err := confita.NewLoader(
		env.NewBackend(),
	).Load(ctx, &cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to parse config: %w", err))
	}
	log.Println(cfg)
}
