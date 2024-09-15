package dto

import "time"

type User struct {
	ID         uint64
	Name       string
	LastName   string
	MiddleName string
	CreatedAt  time.Time
}

type InsertUserRequest struct {
	Name       string
	LastName   string
	MiddleName string
}
