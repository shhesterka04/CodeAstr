package model

import (
	"database/sql"
	"time"
)

type User struct {
	TgID             int64               `db:"tg_id"`
	Username         string              `db:"username"`
	TypeID           sql.Null[int]       `db:"type_id"`
	StyleID          sql.Null[int]       `db:"style_id"`
	ZodiacSignID     sql.Null[int]       `db:"zodiac_sign_id"`
	Gender           sql.Null[string]    `db:"gender"`
	RegistrationDate time.Time           `db:"registration_date"`
	BirthDate        sql.Null[time.Time] `db:"birth_date"`
	BirthTime        sql.Null[string]    `db:"birth_time"`
	BirthPlace       sql.Null[string]    `db:"birth_place"`
	FriendCode       sql.Null[string]    `db:"friend_code"`
	Tokens           sql.Null[int32]     `db:"tokens"`
	FamilyStatus     sql.Null[string]    `db:"family_status"`
	TypeOfActivity   sql.Null[string]    `db:"type_of_activity"`
	NotificationTime sql.Null[int32]     `db:"notification_time"`
	LastActionTime   sql.Null[time.Time] `db:"last_action_time"`
	LanguageID       sql.Null[int64]     `db:"language_id"`
}
