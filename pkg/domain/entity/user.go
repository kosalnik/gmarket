package entity

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `db:"id"`
	Login    string    `login:"login"`
	Password string    `password:"password"`
}
