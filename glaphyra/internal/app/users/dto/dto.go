package dto

import (
	"time"
)

const DefaultUsername = "Звездочет"

const (
	DefaultUser = 1
)

type User struct {
	TgID             int64
	Username         string
	Type             string
	Style            string
	ZodiacSign       string
	Gender           string
	RegistrationDate time.Time
	BirthDate        time.Time
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
	TypeID   int `json:"type_id"`
}

type UpdateUserRequest struct {
	TgID             int64
	Username         string
	Style            string
	Gender           string
	BirthDate        time.Time
	ZodiacSignID     int64 `json:"zodiac_sign_id"`
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
