package repository

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:user,select:user"`

	Id           int       `bun:"id,pk,autoincrement" json:"id"`
	Name         string    `bun:"name" json:"name"`
	FamilyName   string    `bun:"family_name" json:"family_name"`
	Login        string    `bun:"login" json:"login"`
	Email        string    `bun:"email" json:"email"`
	Password     string    `bun:"password" json:"password"`
	Phone        string    `bun:"phone" json:"phone"`
	RegisteredAt time.Time `bun:"registered_at" json:"registered_at"`
	UpdatedAt    time.Time `bun:"updated_at" json:"updated_at"`
}
