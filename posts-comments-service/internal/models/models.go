package models

import "time"

type Like struct {
	PostId int32
	UserId int32
	Time   time.Time
}

type View struct {
	PostId int32
	UserId int32
	Time   time.Time
}

type Comment struct {
	PostId int32
	UserId int32
	Time   time.Time
}
