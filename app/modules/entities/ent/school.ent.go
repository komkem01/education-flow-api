package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type School struct {
	bun.BaseModel `bun:"table:schools,alias:s"`

	ID          uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	Name        string     `bun:"name,notnull,unique"`
	LogoURL     *string    `bun:"logo_url"`
	ThemeColor  *string    `bun:"theme_color,type:varchar(7)"`
	Address     *string    `bun:"address"`
	Description *string    `bun:"description"`
	CreatedAt   time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt   *time.Time `bun:"deleted_at,soft_delete"`
}
