package dto

import (
	"time"

	"glaphyra/internal/pkg/zodiac_signs"
)

const DefaultUsername = "Звездочет ебучий"

type User struct {
	TgID             int32
	Username         string
	Type             string
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
	TgID     int32
	Username string
	Type     string
}

type UpdateUserRequest struct {
	TgID       int32
	Username   string
	Style      string
	Gender     string
	BirthDate  time.Time
	ZodiacSign zodiac_signs.ZodiacSign
	BirthTime  time.Time
	BirthPlace string
	Tokens     int32
}
