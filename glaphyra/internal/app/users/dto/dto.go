package dto

import (
	"time"

	"glaphyra/internal/pkg/zodiac_signs"
)

const DefaultUsername = "Звездочет"

type UserType string

const (
	DefaultUser = UserType("user")
)

type User struct {
	TgID             int64
	Username         string
	Type             UserType
	Style            string
	Gender           string
	RegistrationDate time.Time
	BirthDate        time.Time
	ZodiacSign       zodiac_signs.ZodiacSign
	BirthTime        string
	FamilyStatus     string
	TypeOfActivity   string
	BirthPlace       string
	FriendCode       string
	Tokens           int32
	NotificationTime int32
	LastActionTime   time.Time
	Language         string
}

type CreateUserRequest struct {
	TgID     int64
	Username string
	Language string
	Type     UserType
}

type UpdateUserRequest struct {
	TgID             int64
	Username         string
	Style            string
	Gender           string
	BirthDate        time.Time
	ZodiacSign       zodiac_signs.ZodiacSign
	BirthTime        string
	BirthPlace       string
	Tokens           int32
	FamilyStatus     string
	TypeOfActivity   string
	NotificationTime int32
	LastActionTime   time.Time
	Language         string
}

type ReflectionRecord struct {
	UserID     int64
	MoodRating int
	Emotions   string
	Activity   string
}
