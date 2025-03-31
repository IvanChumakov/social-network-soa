package models

import (
	"time"
)

type RegisterModel struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserModel struct {
	Id           int       `json:"id" default:"0""`
	Name         string    `json:"name" example:"" default:""`
	FamilyName   string    `json:"family_name" example:"" default:""`
	Login        string    `json:"login" example:"" default:""`
	Email        string    `json:"email" example:"" default:""`
	Password     string    `json:"password" example:"" default:""`
	Phone        string    `json:"phone" example:"" default:""`
	RegisteredAt time.Time `json:"registered_at" example:"2023-10-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2023-10-01T00:00:00Z"`
}
