package models

import "time"

type RegistrationInfo struct {
	UserID     int
	Registered time.Time
	Login      string
}
