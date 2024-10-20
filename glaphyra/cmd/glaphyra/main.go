package main

import (
	"context"
	"log"

	"glaphyra/internal/app/users/repository"
	"glaphyra/internal/app/users/service"
	"glaphyra/internal/bot/bot"
	"glaphyra/internal/pkg/db"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx, "postgres://test:test@localhost:5432/test?sslmode=disable") // TODO add to config
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	repo := repository.NewRepository(database)
	userService := service.NewUserService(repo)

	b, err := bot.NewBot("8096088977:AAEQJYPk2ihyfZ2badyeCfop_6oW8G78tRU", userService)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	b.Start()
}
