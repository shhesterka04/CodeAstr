package dto

import (
	"time"

	"glaphyra/internal/pkg/zodiac_signs"
)

const DefaultUsername = "Звездочет ебучий"

type UserType string

const (
	DefaultUser = UserType("user")
	Admin       = UserType("admin")
)

const (
	FriendlyStyle = "Дружелюбный стиль"
	SeriousStyle  = "Серьезный стиль"
	FunnyStyle    = "Шутливый стиль"
)

type User struct {
	TgID             int32
	Username         string
	Type             UserType
	Style            string
	Gender           string
	RegistrationDate time.Time
	BirthDate        time.Time
	ZodiacSign       zodiac_signs.ZodiacSign
	BirthTime        time.Time
	BirthPlace       string
	FriendCode       string
	Tokens           int32
}

type CreateUserRequest struct {
	TgID     int64
	Username string
	Type     UserType
}

type UpdateUserRequest struct {
	TgID       int64
	Username   string
	Style      string
	Gender     string
	BirthDate  time.Time
	ZodiacSign zodiac_signs.ZodiacSign
	BirthTime  time.Time
	BirthPlace string
	Tokens     int32
}
