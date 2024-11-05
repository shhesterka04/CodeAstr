package main

import (
	"context"
	yaml_config "glaphyra/config"
	"glaphyra/internal/llm/handlers"
	"log"

	"glaphyra/internal/app/users/repository"
	"glaphyra/internal/app/users/service"
	"glaphyra/internal/bot/bot"
	"glaphyra/internal/pkg/db"
)

func main() {
	config, err := yaml_config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx, "postgres://test:test@localhost:5432/test?sslmode=disable") // TODO add to config
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	repo := repository.NewRepository(database)
	userService := service.NewUserService(repo)
	yaGptApi := handlers.NewYaGPTHandler(config)

	b, err := bot.NewBot("8096088977:AAEQJYPk2ihyfZ2badyeCfop_6oW8G78tRU", userService, yaGptApi)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	b.Start()
}
