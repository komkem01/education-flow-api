package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StorageProvider string

const (
	StorageProviderS3 StorageProvider = "s3"
)

type Storage struct {
	bun.BaseModel `bun:"table:storages,alias:st"`

	ID         uuid.UUID       `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID   uuid.UUID       `bun:"school_id,type:uuid,notnull"`
	Provider   StorageProvider `bun:"provider,notnull,default:'s3'"`
	Name       string          `bun:"name,notnull"`
	Endpoint   *string         `bun:"endpoint"`
	BucketName string          `bun:"bucket_name,notnull"`
	IsDefault  bool            `bun:"is_default,notnull,default:false"`
	Config     map[string]any  `bun:"config,type:jsonb"`
	CreatedAt  time.Time       `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt  time.Time       `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt  *time.Time      `bun:"deleted_at,soft_delete"`
}

type StorageUpdate struct {
	Provider   *StorageProvider
	Name       *string
	Endpoint   *string
	BucketName *string
	IsDefault  *bool
	Config     *map[string]any
}
