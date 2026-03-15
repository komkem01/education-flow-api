package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type TeacherEmergencyContactEntity interface {
	CreateTeacherEmergencyContact(ctx context.Context, data *ent.TeacherEmergencyContact) (*ent.TeacherEmergencyContact, error)
	GetTeacherEmergencyContactByID(ctx context.Context, id uuid.UUID) (*ent.TeacherEmergencyContact, error)
	ListTeacherEmergencyContacts(ctx context.Context, req *base.RequestPaginate, memberTeacherID *uuid.UUID, isPrimary *bool) ([]*ent.TeacherEmergencyContact, *base.ResponsePaginate, error)
	UpdateTeacherEmergencyContactByID(ctx context.Context, id uuid.UUID, data *ent.TeacherEmergencyContactUpdate) (*ent.TeacherEmergencyContact, error)
	SoftDeleteTeacherEmergencyContactByID(ctx context.Context, id uuid.UUID) error
}
