package repository

import (
	"context"
	"time"

	"glaphyra/internal/app/users/dto"
	"glaphyra/internal/app/users/model"
	"glaphyra/internal/pkg/db"
	"glaphyra/internal/pkg/db/utils"
	errs "glaphyra/internal/pkg/errors"
	"glaphyra/internal/pkg/log"
	"glaphyra/internal/pkg/tables"
	"glaphyra/internal/pkg/zodiac_signs"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type Repository interface {
	FindByID(ctx context.Context, TgID int64) (*dto.User, error)
	Create(ctx context.Context, req *dto.CreateUserRequest) error
	DeleteByID(ctx context.Context, TgID int64) error
	UpdateByTgID(ctx context.Context, req *dto.UpdateUserRequest) error
	SaveFeedback(ctx context.Context, tgID int64, feedback string) error
	SaveReflection(ctx context.Context, reflectionRecord dto.ReflectionRecord) error
	GetUsersByNotificationTime(ctx context.Context, notificationTime int) ([]int64, error)
	GetReflectionsByUserIDLastWeek(ctx context.Context, userID int64) ([]dto.ReflectionRecord, error)
}

type repository struct {
	db db.DBops
}

func NewRepository(db db.DBops) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByID(ctx context.Context, TgID int64) (*dto.User, error) {
	var user model.User

	query, args, err := buildFindByIDQuery(TgID).ToSql()
	if err != nil {
		return nil, log.Wrap(err)
	}

	err = r.db.Get(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &errs.ObjectNotFoundError{}
		}
		return nil, log.Wrap(err)
	}

	return modelToDTO(&user), nil
}

func (r *repository) UpdateByTgID(ctx context.Context, req *dto.UpdateUserRequest) error {
	valuesMap := map[string]interface{}{}

	valuesMap = utils.SetIfNotZero(valuesMap, "username", req.Username)
	valuesMap = utils.SetIfNotZero(valuesMap, "style", req.Style)
	valuesMap = utils.SetIfNotZero(valuesMap, "gender", req.Gender)
	valuesMap = utils.SetIfNotZero(valuesMap, "birth_date", req.BirthDate)
	valuesMap = utils.SetIfNotZero(valuesMap, "zodiac_sign", req.ZodiacSign)
	valuesMap = utils.SetIfNotZero(valuesMap, "birth_time", req.BirthTime)
	valuesMap = utils.SetIfNotZero(valuesMap, "birth_place", req.BirthPlace)
	valuesMap = utils.SetIfNotZero(valuesMap, "tokens", req.Tokens)
	valuesMap = utils.SetIfNotZero(valuesMap, "type_of_activity", req.TypeOfActivity)
	valuesMap = utils.SetIfNotZero(valuesMap, "notification_time", req.NotificationTime)
	valuesMap = utils.SetIfNotZero(valuesMap, "family_status", req.FamilyStatus)
	valuesMap = utils.SetIfNotZero(valuesMap, "last_action_time", req.LastActionTime)

	if len(valuesMap) == 0 {
		return nil
	}

	query, args, err := sq.Update(tables.Users).Where(sq.Eq{"tg_id": req.TgID}).SetMap(valuesMap).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return log.Wrap(err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return log.Wrap(err)
	}

	return nil
}

func (r *repository) Create(ctx context.Context, req *dto.CreateUserRequest) error {
	friendCode, err := uuid.NewUUID()
	if err != nil {
		return log.Wrap(err)
	}

	username := dto.DefaultUsername
	if len(req.Username) != 0 {
		username = req.Username
	}

	valuesMap := map[string]interface{}{
		"tg_id":             req.TgID,
		"username":          username,
		"type":              req.Type,
		"registration_date": time.Now(),
		"friend_code":       friendCode,
		"tokens":            100,
		"notification_time": 18,
	}

	query, args, err := sq.Insert(tables.Users).PlaceholderFormat(sq.Dollar).SetMap(valuesMap).ToSql()
	if err != nil {
		return log.Wrap(err)
	}

	_, err = r.db.Exec(ctx, query, args...)

	if err != nil {
		return log.Wrap(err)
	}

	return nil
}

func (r *repository) SaveFeedback(ctx context.Context, tgID int64, feedback string) error {
	valuesMap := map[string]interface{}{
		"tg_id":      tgID,
		"feedback":   feedback,
		"created_at": time.Now(),
	}

	query, args, err := sq.Insert(tables.Feedback).PlaceholderFormat(sq.Dollar).SetMap(valuesMap).ToSql()
	if err != nil {
		return log.Wrap(err)
	}

	_, err = r.db.Exec(ctx, query, args...)

	if err != nil {
		return log.Wrap(err)
	}

	return nil
}

func (r *repository) SaveReflection(ctx context.Context, reflectionRecord dto.ReflectionRecord) error {
	valuesMap := map[string]interface{}{
		"tg_id":       reflectionRecord.UserID,
		"mood_rating": reflectionRecord.Mark,
		"emotions":    reflectionRecord.Emotions,
		"activity":    reflectionRecord.Activity,
		"created_at":  time.Now(),
	}

	query, args, err := sq.Insert(tables.Reflections).PlaceholderFormat(sq.Dollar).SetMap(valuesMap).ToSql()
	if err != nil {
		return log.Wrap(err)
	}

	_, err = r.db.Exec(ctx, query, args...)

	if err != nil {
		return log.Wrap(err)
	}

	return nil
}

func (r *repository) DeleteByID(ctx context.Context, TgID int64) error {
	query, args, err := sq.Delete(tables.Users).PlaceholderFormat(sq.Dollar).Where(sq.Eq{"tg_id": TgID}).ToSql()
	if err != nil {
		return log.Wrap(err)
	}

	_, err = r.db.Exec(
		ctx,
		query,
		args...)
	if err != nil {
		return log.Wrap(err)
	}

	return nil
}

func modelToDTO(user *model.User) *dto.User {
	userDTO := dto.User{
		TgID:             user.TgID,
		Username:         user.Username,
		Type:             dto.UserType(user.Type),
		RegistrationDate: user.RegistrationDate,
	}

	if user.Style.Valid {
		userDTO.Style = user.Style.V
	}

	if user.Gender.Valid {
		userDTO.Gender = user.Gender.V
	}

	if user.BirthDate.Valid {
		userDTO.BirthDate = user.BirthDate.V
	}

	if user.ZodiacSign.Valid {
		userDTO.ZodiacSign = zodiac_signs.ZodiacSign(user.ZodiacSign.V)
	}

	if user.BirthTime.Valid {
		userDTO.BirthTime = user.BirthTime.V
	}

	if user.BirthPlace.Valid {
		userDTO.BirthPlace = user.BirthPlace.V
	}

	if user.FriendCode.Valid {
		userDTO.FriendCode = user.FriendCode.V
	}

	if user.Tokens.Valid {
		userDTO.Tokens = user.Tokens.V
	}

	if user.TypeOfActivity.Valid {
		userDTO.TypeOfActivity = user.TypeOfActivity.V
	}

	if user.FamilyStatus.Valid {
		userDTO.FamilyStatus = user.FamilyStatus.V
	}

	if user.NotificationTime.Valid {
		userDTO.NotificationTime = user.NotificationTime.V
	}

	if user.NotificationTime.Valid {
		userDTO.LastActionTime = user.LastActionTime.V
	}

	return &userDTO
}

func buildFindByIDQuery(tgID int64) sq.SelectBuilder {
	queryBuilder := sq.Select(
		"tg_id",
		"username",
		"type",
		"style",
		"gender",
		"registration_date",
		"birth_date",
		"birth_place",
		"zodiac_sign",
		"birth_time",
		"friend_code",
		"tokens",
		"family_status",
		"type_of_activity",
		"notification_time",
	).
		From(tables.Users).
		Where(sq.Eq{"tg_id": tgID}).
		PlaceholderFormat(sq.Dollar)

	return queryBuilder
}

func (r *repository) GetUsersByNotificationTime(ctx context.Context, notificationTime int) ([]int64, error) {
	var users []int64

	query, args, err := sq.Select("tg_id").
		From(tables.Users).
		Where(sq.Eq{"notification_time": notificationTime}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, log.Wrap(err)
	}

	err = r.db.Select(ctx, &users, query, args...)
	if err != nil {
		return nil, log.Wrap(err)
	}

	return users, nil
}

func (r *repository) GetReflectionsByUserIDLastWeek(ctx context.Context, userID int64) ([]dto.ReflectionRecord, error) {
	var reflections []dto.ReflectionRecord

	query, args, err := sq.Select(
		"mood_rating",
		"emotions",
		"activity",
	).
		From(tables.Reflections).
		Where(sq.Eq{"tg_id": userID}).
		Where("created_at > now() - interval '1 week'").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, log.Wrap(err)
	}

	err = r.db.Select(ctx, &reflections, query, args...)
	if err != nil {
		return nil, log.Wrap(err)
	}

	return reflections, nil
}
