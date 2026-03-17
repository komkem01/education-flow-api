package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type SchoolAnnouncementEntity interface {
	CreateSchoolAnnouncement(ctx context.Context, data *ent.SchoolAnnouncement) (*ent.SchoolAnnouncement, error)
	GetSchoolAnnouncementByID(ctx context.Context, id uuid.UUID) (*ent.SchoolAnnouncement, error)
	ListSchoolAnnouncements(ctx context.Context, req *base.RequestPaginate, schoolID *uuid.UUID, status *ent.SchoolAnnouncementStatus, targetRole *string, isPinned *bool) ([]*ent.SchoolAnnouncement, *base.ResponsePaginate, error)
	UpdateSchoolAnnouncementByID(ctx context.Context, id uuid.UUID, data *ent.SchoolAnnouncementUpdate) (*ent.SchoolAnnouncement, error)
	SoftDeleteSchoolAnnouncementByID(ctx context.Context, id uuid.UUID) error
}
