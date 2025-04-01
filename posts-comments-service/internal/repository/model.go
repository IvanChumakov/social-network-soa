package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type Post struct {
	bun.BaseModel `bun:"table:posts,select:posts"`

	Id          int32     `bun:"id,pk,autoincrement" json:"id"`
	Name        string    `bun:"name" json:"name"`
	Description string    `bun:"description" json:"description"`
	CreatorId   int32     `bun:"creator_id" json:"creator_id"`
	IsPrivate   bool      `bun:"is_private" json:"is_private"`
	CreatedAt   time.Time `bun:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at" json:"updated_at"`
	Tags        []string  `bun:"tags" json:"tags"`
}
