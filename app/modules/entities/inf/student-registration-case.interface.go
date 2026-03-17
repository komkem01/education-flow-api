package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type StudentRegistrationCaseEntity interface {
	CreateStudentRegistrationCase(ctx context.Context, data *ent.StudentRegistrationCase) (*ent.StudentRegistrationCase, error)
	GetStudentRegistrationCaseByID(ctx context.Context, id uuid.UUID) (*ent.StudentRegistrationCase, error)
	ListStudentRegistrationCases(ctx context.Context, req *base.RequestPaginate, schoolID *uuid.UUID, studentID *uuid.UUID, status *ent.StudentRegistrationCaseStatus, registrationType *ent.StudentRegistrationType) ([]*ent.StudentRegistrationCase, *base.ResponsePaginate, error)
	UpdateStudentRegistrationCaseByID(ctx context.Context, id uuid.UUID, data *ent.StudentRegistrationCaseUpdate) (*ent.StudentRegistrationCase, error)
	SoftDeleteStudentRegistrationCaseByID(ctx context.Context, id uuid.UUID) error
	GetStudentRegistrationCaseDetailByID(ctx context.Context, id uuid.UUID) (*ent.StudentRegistrationCaseDetail, error)
	ReplaceStudentRegistrationCaseBundle(ctx context.Context, caseID uuid.UUID, bundle *ent.StudentRegistrationCaseBundle) error
	SubmitStudentRegistrationCase(ctx context.Context, caseID uuid.UUID, actorID uuid.UUID, actorRole ent.ApprovalActorRole, reason *string) (*ent.StudentRegistrationCase, error)
	RejectStudentRegistrationCase(ctx context.Context, caseID uuid.UUID, actorID uuid.UUID, actorRole ent.ApprovalActorRole, reason *string) (*ent.StudentRegistrationCase, error)
	ApproveAndApplyStudentRegistrationCase(ctx context.Context, caseID uuid.UUID, actorID uuid.UUID, actorRole ent.ApprovalActorRole, comment *string, idempotencyKey *string, metadata map[string]any) (*ent.StudentRegistrationCase, error)
	ApplyStudentRegistrationCase(ctx context.Context, caseID uuid.UUID, actorID uuid.UUID, actorRole ent.ApprovalActorRole, comment *string) (*ent.StudentRegistrationCase, error)
	GetStudentRegistrationPreviousEducationByCaseID(ctx context.Context, caseID uuid.UUID) (*ent.StudentRegistrationPreviousEducation, error)
	UpsertStudentRegistrationPreviousEducationByCaseID(ctx context.Context, caseID uuid.UUID, data *ent.StudentRegistrationPreviousEducation) (*ent.StudentRegistrationPreviousEducation, error)
	GetStudentRegistrationFamilyEconomicByCaseID(ctx context.Context, caseID uuid.UUID) (*ent.StudentRegistrationFamilyEconomic, error)
	UpsertStudentRegistrationFamilyEconomicByCaseID(ctx context.Context, caseID uuid.UUID, data *ent.StudentRegistrationFamilyEconomic) (*ent.StudentRegistrationFamilyEconomic, error)
	ListStudentRegistrationDocumentsByCaseID(ctx context.Context, caseID uuid.UUID) ([]*ent.StudentRegistrationDocument, error)
	ReplaceStudentRegistrationDocumentsByCaseID(ctx context.Context, caseID uuid.UUID, rows []*ent.StudentRegistrationDocument) error
	CreateStudentRegistrationRule(ctx context.Context, data *ent.StudentRegistrationRule) (*ent.StudentRegistrationRule, error)
	GetStudentRegistrationRuleByID(ctx context.Context, id uuid.UUID) (*ent.StudentRegistrationRule, error)
	ListStudentRegistrationRules(ctx context.Context, req *base.RequestPaginate, schoolID *uuid.UUID, registrationType *ent.StudentRegistrationType) ([]*ent.StudentRegistrationRule, *base.ResponsePaginate, error)
	UpdateStudentRegistrationRuleByID(ctx context.Context, id uuid.UUID, data *ent.StudentRegistrationRuleUpdate) (*ent.StudentRegistrationRule, error)
	SoftDeleteStudentRegistrationRuleByID(ctx context.Context, id uuid.UUID) error
}
