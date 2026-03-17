package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TeacherAddress struct {
	bun.BaseModel `bun:"table:teacher_addresses,alias:ta"`

	ID              uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	MemberTeacherID uuid.UUID  `bun:"member_teacher_id,notnull,type:uuid"`
	HouseNo         string     `bun:"house_no,notnull,type:varchar(50)"`
	Village         *string    `bun:"village,type:varchar(100)"`
	Road            *string    `bun:"road,type:varchar(255)"`
	Province        string     `bun:"province,notnull,type:varchar(150)"`
	District        string     `bun:"district,notnull,type:varchar(150)"`
	Subdistrict     string     `bun:"subdistrict,notnull,type:varchar(150)"`
	PostalCode      string     `bun:"postal_code,notnull,type:varchar(10)"`
	IsPrimary       bool       `bun:"is_primary,notnull,default:false"`
	SortOrder       int        `bun:"sort_order,notnull,default:0"`
	IsActive        bool       `bun:"is_active,notnull,default:true"`
	CreatedAt       time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt       time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt       *time.Time `bun:"deleted_at,soft_delete"`
}

type TeacherAddressInput struct {
	HouseNo     string
	Village     *string
	Road        *string
	Province    string
	District    string
	Subdistrict string
	PostalCode  string
	IsPrimary   bool
	SortOrder   int
}
