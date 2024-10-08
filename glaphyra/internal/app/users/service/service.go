package service

import (
	"context"

	"glaphyra/internal/app/users/dto"
	usersrepository "glaphyra/internal/app/users/repository"
	"glaphyra/internal/pkg/log"
)

type UserService interface {
	FindByID(ctx context.Context, tgID int32) (*dto.User, error)
	Create(ctx context.Context, req *dto.CreateUserRequest) error
	Update(ctx context.Context, req *dto.UpdateUserRequest) error
	Delete(ctx context.Context, tgID int32) error
}

type implementSrv struct {
	repo usersrepository.Repository
}

func NewUserService(repo usersrepository.Repository) UserService {
	return &implementSrv{
		repo: repo,
	}
}

func (i *implementSrv) FindByID(ctx context.Context, tgID int32) (*dto.User, error) {
	user, err := i.repo.FindByID(ctx, tgID)
	if err != nil {
		return nil, log.Wrap(err)
	}

	return user, nil
}

func (i *implementSrv) Create(ctx context.Context, req *dto.CreateUserRequest) error {
	err := i.repo.Create(ctx, req)
	if err != nil {
		return log.Wrap(err)
	}

	return nil
}

func (i *implementSrv) Update(ctx context.Context, req *dto.UpdateUserRequest) error {
	err := i.repo.UpdateByTgID(ctx, req)
	if err != nil {
		return log.Wrap(err)
	}

	return nil
}

func (i *implementSrv) Delete(ctx context.Context, tgID int32) error {
	err := i.repo.DeleteByID(ctx, tgID)
	if err != nil {
		return log.Wrap(err)
	}

	return nil
}
