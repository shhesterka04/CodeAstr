package main

import (
	"context"

	"glaphyra/config"
	"glaphyra/internal/app/users/repository"
	"glaphyra/internal/app/users/service"
	"glaphyra/internal/bot/bot"
	"glaphyra/internal/llm/handlers"
	"glaphyra/internal/pkg/db"
	"glaphyra/internal/pkg/log"
)

const cfgPath = "config.yaml"

func main() {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Error(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx, cfg.DbDSN)
	if err != nil {
		log.Error(err)
		return
	}
	defer database.Close()

	repo := repository.NewRepository(database)
	userService := service.NewUserService(repo)
	yaGptApi := handlers.NewYaGPTHandler(cfg)

	b, err := bot.NewBot(cfg.TgToken, userService, yaGptApi)
	if err != nil {
		log.Error(err)
		return
	}

	log.WriteLog("Bot started")

	b.Start()
}
