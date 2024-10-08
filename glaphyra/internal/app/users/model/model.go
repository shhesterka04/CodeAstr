package model

import (
	"database/sql"
	"time"
)

type User struct {
	TgID             int32               `db:"tg_id"`
	Username         string              `db:"username"`
	Type             string              `db:"type"`
	Style            sql.Null[string]    `db:"style"`
	Gender           sql.Null[string]    `db:"gender"`
	RegistrationDate time.Time           `db:"registration_date"`
	BirthDate        sql.Null[time.Time] `db:"birth_date"`
	ZodiacSign       sql.Null[string]    `db:"zodiac_sign"`
	BirthTime        sql.Null[time.Time] `db:"birth_time"`
	BirthPlace       sql.Null[string]    `db:"birth_place"`
	FriendCode       sql.Null[string]    `db:"friend_code"`
	Tokens           sql.Null[int32]     `db:"tokens"`
}
