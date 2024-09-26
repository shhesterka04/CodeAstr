package main

import (
	"context"
	"log"

	"glaphyra/internal/app/users/dto"
	"glaphyra/internal/app/users/repository"
	"glaphyra/internal/app/users/service"
	"glaphyra/internal/pkg/bot"
	"glaphyra/internal/pkg/db"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx, "postgres://test:test@localhost:5432/test?sslmode=disable") // TODO add
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	repo := repository.NewRepository(database)
	userService := service.NewUserService(repo)

	err = userService.Insert(ctx, &dto.InsertUserRequest{Name: "Дмитрий", LastName: "Полевой", MiddleName: "ОпенСиВишевич"}) // TODO убрать безобразие
	if err != nil {
		log.Fatal(err)
	}

	user, err := userService.FindByID(ctx, 1)
	if err != nil {
		log.Fatal(err)
	}

	b, err := bot.NewBot("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	b.Start()

	log.Println(user)
}
