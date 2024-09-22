package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID         uint64           `db:"id"`
	Name       string           `db:"first_name"`
	MiddleName sql.Null[string] `db:"middle_name"`
	LastName   string           `db:"last_name"`
	CreatedAt  time.Time        `db:"created_at"`
}
