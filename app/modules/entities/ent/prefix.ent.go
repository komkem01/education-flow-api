package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Prefix struct {
	bun.BaseModel `bun:"table:prefixes,alias:p"`

	ID        uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	GenderID  uuid.UUID  `bun:"gender_id,type:uuid,notnull"`
	Name      string     `bun:"name,notnull,unique"`
	IsActive  bool       `bun:"is_active,notnull,default:false"`
	CreatedAt time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
}
