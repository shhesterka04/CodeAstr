package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"glaphyra/internal/app/users/dto"
	"glaphyra/internal/app/users/model"
	"glaphyra/internal/pkg/db"
	errs "glaphyra/internal/pkg/errors"
	"glaphyra/internal/pkg/log"
)

type Repository interface {
	FindByID(ctx context.Context, userID int64) (*dto.User, error)
	Insert(ctx context.Context, insertReq *dto.InsertUserRequest) error
}

type repository struct {
	db db.DBops
}

func NewRepository(db db.DBops) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByID(ctx context.Context, userID int64) (*dto.User, error) {
	var user model.User

	err := r.db.Get(ctx, &user, "SELECT id, first_name, last_name, middle_name, created_at FROM users WHERE id=$1;", userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &errs.ObjectNotFoundError{}
		}
		return nil, log.Wrap(err)
	}

	return modelToDTO(&user), nil
}

func (r *repository) Insert(ctx context.Context, req *dto.InsertUserRequest) error {
	_, err := r.db.Exec(ctx, `INSERT INTO users(id, first_name, last_name, middle_name) VALUES(nextval('users_id'), $1, $2, $3);`, req.Name, req.LastName, req.MiddleName)

	if err != nil {
		return log.Wrap(err)
	}

	return nil
}

func modelToDTO(user *model.User) *dto.User {
	return &dto.User{
		ID:         user.ID,
		Name:       user.Name,
		LastName:   user.LastName,
		MiddleName: user.MiddleName.V,
		CreatedAt:  user.CreatedAt,
	}
}
