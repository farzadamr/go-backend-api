package main

import (
	"log"

	"github.com/farzadamr/go-backend-api/api"
	"github.com/farzadamr/go-backend-api/config"
	"github.com/farzadamr/go-backend-api/internal/infra/database"
	"github.com/farzadamr/go-backend-api/internal/infra/migration"
	"github.com/farzadamr/go-backend-api/pkg/logging"
)

func main() {
	cfg := LoadAndParseConfig()
	if err := logging.Init(cfg); err != nil {
		log.Fatalf("logger can not initialized: %w", err)
	}

	err := database.InitDb(cfg)
	defer database.CloseDb()
	if err != nil {
		log.Fatal("connecting to database failed")
	}
	migration.Up_1()

	if err := api.Run(cfg); err != nil {
		log.Fatalf("api: %v", err)
	}
	logging.Info("server is running", "port", cfg.HTTP.Port)
}

func LoadAndParseConfig() *config.Config {
	loader := config.NewLoader()
	if err := loader.LoadEnv(); err != nil {
		log.Fatalf("failed to load env: %v", err)
		return nil
	}

	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("invalid config: %v", err)
		return nil
	}

	return cfg
}
