package service

import (
	"context"
	"glaphyra/internal/app/users/dto"
	usersrepository "glaphyra/internal/app/users/repository"
	"glaphyra/internal/pkg/log"
)

type UserService interface {
	FindByID(ctx context.Context, userID uint64) (*dto.User, error)
	Insert(ctx context.Context, req *dto.InsertUserRequest) error // TODO заменить на че-нить другое, пока заглушка потыкать
}

type implementSrv struct {
	repo usersrepository.Repository
}

func NewUserService(repo usersrepository.Repository) UserService {
	return &implementSrv{
		repo: repo,
	}
}

func (i *implementSrv) FindByID(ctx context.Context, _ uint64) (*dto.User, error) { // TODO сделать, пока заглушка
	user, err := i.repo.FindByID(ctx, 1)
	if err != nil {
		return nil, log.Wrap(err)
	}

	return user, nil
}

func (i *implementSrv) Insert(ctx context.Context, req *dto.InsertUserRequest) error {
	err := i.repo.Insert(ctx, req)
	if err != nil {
		return log.Wrap(err)
	}

	return nil
}
