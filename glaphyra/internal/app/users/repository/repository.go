package repository

import (
	"context"
	"fmt"
	"time"

	"glaphyra/internal/app/users/dto"
	"glaphyra/internal/app/users/model"
	"glaphyra/internal/pkg/db"
	"glaphyra/internal/pkg/db/utils"
	errs "glaphyra/internal/pkg/errors"
	"glaphyra/internal/pkg/log"
	"glaphyra/internal/pkg/tables"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

const DefaultLanguage = "Русский"

type Repository interface {
	FindByID(ctx context.Context, TgID int64) (*dto.User, error)
	Create(ctx context.Context, req *dto.CreateUserRequest) error
	DeleteByID(ctx context.Context, TgID int64) error
	UpdateByTgID(ctx context.Context, req *dto.UpdateUserRequest) error
	SaveFeedback(ctx context.Context, tgID int64, feedback string) error
	SaveReflection(ctx context.Context, reflectionRecord dto.ReflectionRecord) error
	GetUsersByNotificationTime(ctx context.Context, notificationTime int) ([]int64, error)
	GetReflectionsByUserIDLastWeek(ctx context.Context, userID int64) ([]dto.ReflectionRecord, error)
	GetAllUsersIDs(ctx context.Context) ([]int64, error)
}

type repository struct {
	db db.DBops
}

func NewRepository(db db.DBops) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) getLanguageIDByName(languageName string) (int, error) {
	if languageName == "" {
		languageName = DefaultLanguage
	}

	var languageID int

	query, args, err := sq.Select("id").
		From(tables.Languages).
		Where(sq.Eq{"language": languageName}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, log.Wrap(err)
	}

	err = r.db.Get(context.Background(), &languageID, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			defaultQuery, defaultArgs, defaultErr := sq.Select("id").
				From(tables.Languages).
				Where(sq.Eq{"language": DefaultLanguage}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			if defaultErr != nil {
				return 0, log.Wrap(defaultErr)
			}

			defaultErr = r.db.Get(context.Background(), &languageID, defaultQuery, defaultArgs...)
			if defaultErr != nil {
				return 0, log.Wrap(defaultErr)
			}
		} else {
			return 0, log.Wrap(err)
		}
	}

	return languageID, nil
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

	return r.modelToDTO(&user), nil
}

func (r *repository) UpdateByTgID(ctx context.Context, req *dto.UpdateUserRequest) error {
	valuesMap := map[string]interface{}{}

	valuesMap = utils.SetIfNotZero(valuesMap, "username", req.Username)

	if req.Style != "" {
		styleID, err := r.getStyleIDByName(req.Style)
		if err != nil {
			return log.Wrap(err)
		}
		valuesMap["style_id"] = styleID
	}

	valuesMap = utils.SetIfNotZero(valuesMap, "zodiac_sign_id", req.ZodiacSignID)

	valuesMap = utils.SetIfNotZero(valuesMap, "gender", req.Gender)
	valuesMap = utils.SetIfNotZero(valuesMap, "birth_date", req.BirthDate)
	valuesMap = utils.SetIfNotZero(valuesMap, "birth_time", req.BirthTime)
	valuesMap = utils.SetIfNotZero(valuesMap, "birth_place", req.BirthPlace)
	valuesMap = utils.SetIfNotZero(valuesMap, "tokens", req.Tokens)
	valuesMap = utils.SetIfNotZero(valuesMap, "type_of_activity", req.TypeOfActivity)
	valuesMap = utils.SetIfNotZero(valuesMap, "notification_time", req.NotificationTime)
	valuesMap = utils.SetIfNotZero(valuesMap, "family_status", req.FamilyStatus)
	valuesMap = utils.SetIfNotZero(valuesMap, "last_action_time", req.LastActionTime)
	valuesMap = utils.SetIfNotZero(valuesMap, "language_id", req.Language)

	if len(valuesMap) == 0 {
		return nil
	}

	query, args, err := sq.Update(tables.Users).
		Where(sq.Eq{"tg_id": req.TgID}).
		SetMap(valuesMap).
		PlaceholderFormat(sq.Dollar).
		ToSql()
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

	languageID, err := r.getLanguageIDByName(req.Language)
	if err != nil {
		return log.Wrap(err)
	}

	valuesMap := map[string]interface{}{
		"tg_id":             req.TgID,
		"username":          username,
		"type_id":           req.TypeID,
		"registration_date": time.Now(),
		"friend_code":       friendCode,
		"tokens":            100,
		"notification_time": 18,
		"language_id":       languageID,
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
		"mood_rating": reflectionRecord.MoodRating,
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
func (r *repository) modelToDTO(user *model.User) *dto.User {
	userType, _ := r.getTypeNameByID(user.TypeID.V)
	userStyle, _ := r.getStyleNameByID(user.StyleID.V)
	userLanguage, _ := r.getLanguageNameByID(int(user.LanguageID.V))

	userDTO := dto.User{
		TgID:             user.TgID,
		Username:         user.Username,
		Type:             userType,
		Style:            userStyle,
		ZodiacSign:       userLanguage,
		RegistrationDate: user.RegistrationDate,
	}

	if user.Gender.Valid {
		userDTO.Gender = user.Gender.V
	}

	if user.BirthDate.Valid {
		userDTO.BirthDate = user.BirthDate.V
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

	if user.LastActionTime.Valid {
		userDTO.LastActionTime = user.LastActionTime.V
	}

	return &userDTO
}

func buildFindByIDQuery(tgID int64) sq.SelectBuilder {
	queryBuilder := sq.Select(
		"tg_id",
		"username",
		"type_id",
		"style_id",
		"zodiac_sign_id",
		"gender",
		"registration_date",
		"birth_date",
		"birth_time",
		"birth_place",
		"friend_code",
		"tokens",
		"family_status",
		"type_of_activity",
		"notification_time",
		"last_action_time",
		"language_id",
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

func (r *repository) GetAllUsersIDs(ctx context.Context) ([]int64, error) {
	var users []int64

	query, args, err := sq.Select("tg_id").
		From(tables.Users).
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

func (r *repository) getStyleIDByName(styleName string) (int, error) {
	if styleName == "" {
		return 0, errors.New("style name is empty")
	}

	var styleID int

	query, args, err := sq.Select("id").
		From(tables.Styles).
		Where(sq.Eq{"style": styleName}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, log.Wrap(err)
	}

	err = r.db.Get(context.Background(), &styleID, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("style not found: %s", styleName)
		}
		return 0, log.Wrap(err)
	}

	return styleID, nil
}

func (r *repository) getStyleNameByID(styleID int) (string, error) {
	var styleName string

	query, args, err := sq.Select("style").
		From(tables.Styles).
		Where(sq.Eq{"id": styleID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", log.Wrap(err)
	}

	err = r.db.Get(context.Background(), &styleName, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("style not found: %d", styleID)
		}
		return "", log.Wrap(err)
	}

	return styleName, nil
}

func (r *repository) getLanguageNameByID(languageID int) (string, error) {
	var languageName string

	query, args, err := sq.Select("language").
		From(tables.Languages).
		Where(sq.Eq{"id": languageID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", log.Wrap(err)
	}

	err = r.db.Get(context.Background(), &languageName, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("language not found: %d", languageID)
		}
		return "", log.Wrap(err)
	}

	return languageName, nil
}

func (r *repository) getTypeNameByID(typeID int) (string, error) {
	var typeName string

	query, args, err := sq.Select("type").
		From(tables.Users).
		Where(sq.Eq{"id": typeID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", log.Wrap(err)
	}

	err = r.db.Get(context.Background(), &typeName, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("type not found: %d", typeID)
		}
		return "", log.Wrap(err)
	}

	return typeName, nil
}

func (r *repository) getZodiacSignNameByID(zodiacSignID int) (string, error) {
	var zodiacSignName string

	query, args, err := sq.Select("zodiac_sign").
		From(tables.ZodiacSigns).
		Where(sq.Eq{"id": zodiacSignID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", log.Wrap(err)
	}

	err = r.db.Get(context.Background(), &zodiacSignName, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("zodiac sign not found: %d", zodiacSignID)
		}
		return "", log.Wrap(err)
	}

	return zodiacSignName, nil
}
