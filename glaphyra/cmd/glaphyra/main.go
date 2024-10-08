package main

import (
	"context"
	"fmt"
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

	repoExample(ctx) // TODO remove

	b, err := bot.NewBot("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	b.Start()
}

func repoExample(ctx context.Context) {
	database, err := db.NewDB(ctx, "postgres://test:test@localhost:5432/test?sslmode=disable") // TODO add to config
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	repo := repository.NewRepository(database)
	userService := service.NewUserService(repo)

	err = userService.Create(ctx, &dto.CreateUserRequest{
		TgID: 1, Username: "Chmonya", Type: "user",
	},
	) // TODO убрать безобразие
	if err != nil {
		log.Fatal(err)
	}

	user, err := userService.FindByID(ctx, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)

	err = userService.Update(ctx, &dto.UpdateUserRequest{
		Gender:   "Мужик",
		Username: "superCmonya",
	})
	if err != nil {
		log.Fatal(err)
	}

	err = userService.Delete(ctx, 1)
	if err != nil {
		log.Fatal(err)
	}
}
