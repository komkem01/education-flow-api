package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SchoolAnnouncementStatus string

const (
	SchoolAnnouncementStatusDraft     SchoolAnnouncementStatus = "draft"
	SchoolAnnouncementStatusPublished SchoolAnnouncementStatus = "published"
	SchoolAnnouncementStatusExpired   SchoolAnnouncementStatus = "expired"
)

type SchoolAnnouncement struct {
	bun.BaseModel `bun:"table:school_announcements,alias:sa"`

	ID             uuid.UUID                `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID       uuid.UUID                `bun:"school_id,type:uuid,notnull"`
	AuthorMemberID uuid.UUID                `bun:"author_member_id,type:uuid,notnull"`
	Title          *string                  `bun:"title"`
	Content        *string                  `bun:"content"`
	Category       *string                  `bun:"category"`
	Status         SchoolAnnouncementStatus `bun:"status,notnull,type:school_announcement_status_enum,default:'draft'"`
	AnnouncedAt    *time.Time               `bun:"announced_at"`
	PublishedAt    *time.Time               `bun:"published_at"`
	ExpiresAt      *time.Time               `bun:"expires_at"`
	CreatedByName  *string                  `bun:"created_by_name"`
	TargetRole     *string                  `bun:"target_role,type:school_announcement_target_role_enum"`
	IsPinned       bool                     `bun:"is_pinned,notnull,default:false"`
	CreatedAt      time.Time                `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt      time.Time                `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt      *time.Time               `bun:"deleted_at,soft_delete"`
}

type SchoolAnnouncementUpdate struct {
	SchoolID         *uuid.UUID
	AuthorMemberID   *uuid.UUID
	Title            *string
	Content          *string
	Category         *string
	Status           *SchoolAnnouncementStatus
	AnnouncedAt      *time.Time
	ClearAnnouncedAt bool
	PublishedAt      *time.Time
	ClearPublishedAt bool
	ExpiresAt        *time.Time
	ClearExpiresAt   bool
	CreatedByName    *string
	TargetRole       *string
	ClearTargetRole  bool
	IsPinned         *bool
}
