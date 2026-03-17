package studentregistrationcases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) GetPreviousEducation(ctx context.Context, caseID uuid.UUID) (*ent.StudentRegistrationPreviousEducation, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.get_previous_education")
	defer span.End()

	item, err := s.db.GetStudentRegistrationPreviousEducationByCaseID(ctx, caseID)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

func (s *Service) UpsertPreviousEducation(ctx context.Context, caseID uuid.UUID, actorRole ent.MemberRole, in *PreviousEducationInput) (*ent.StudentRegistrationPreviousEducation, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.upsert_previous_education")
	defer span.End()

	if _, ok := parseApprovalActorRoleFromMember(actorRole); !ok {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseUnauthorized)
	}
	if err := s.ensureCaseEditable(ctx, caseID); err != nil {
		return nil, err
	}

	if in == nil {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
	}
	transferDate, err := parseDatePtr(in.TransferDate)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
	}
	transcriptReceived := false
	if in.TranscriptReceived != nil {
		transcriptReceived = *in.TranscriptReceived
	}

	item, err := s.db.UpsertStudentRegistrationPreviousEducationByCaseID(ctx, caseID, &ent.StudentRegistrationPreviousEducation{
		CaseID:                 caseID,
		PreviousSchoolName:     normalizeStrPtr(in.PreviousSchoolName),
		PreviousSchoolProvince: normalizeStrPtr(in.PreviousSchoolProvince),
		PreviousGradeLevel:     normalizeStrPtr(in.PreviousGradeLevel),
		GPA:                    in.GPA,
		TransferCertificateNo:  normalizeStrPtr(in.TransferCertificateNo),
		TransferDate:           transferDate,
		TranscriptReceived:     transcriptReceived,
		Remarks:                normalizeStrPtr(in.Remarks),
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

func (s *Service) GetFamilyEconomic(ctx context.Context, caseID uuid.UUID) (*ent.StudentRegistrationFamilyEconomic, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.get_family_economic")
	defer span.End()

	item, err := s.db.GetStudentRegistrationFamilyEconomicByCaseID(ctx, caseID)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

func (s *Service) UpsertFamilyEconomic(ctx context.Context, caseID uuid.UUID, actorRole ent.MemberRole, in *FamilyEconomicInput) (*ent.StudentRegistrationFamilyEconomic, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.upsert_family_economic")
	defer span.End()

	if _, ok := parseApprovalActorRoleFromMember(actorRole); !ok {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseUnauthorized)
	}
	if err := s.ensureCaseEditable(ctx, caseID); err != nil {
		return nil, err
	}
	if in == nil {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
	}

	var incomeBracket *ent.RegistrationIncomeBracket
	if in.IncomeBracket != nil {
		parsed, ok := parseIncomeBracket(*in.IncomeBracket)
		if !ok {
			return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
		}
		incomeBracket = &parsed
	}
	scholarshipFlag := false
	if in.ScholarshipFlag != nil {
		scholarshipFlag = *in.ScholarshipFlag
	}
	welfareFlag := false
	if in.WelfareFlag != nil {
		welfareFlag = *in.WelfareFlag
	}
	debtFlag := false
	if in.DebtFlag != nil {
		debtFlag = *in.DebtFlag
	}

	item, err := s.db.UpsertStudentRegistrationFamilyEconomicByCaseID(ctx, caseID, &ent.StudentRegistrationFamilyEconomic{
		CaseID:                 caseID,
		HouseholdSize:          in.HouseholdSize,
		HouseholdIncomeMonthly: in.HouseholdIncomeMonthly,
		IncomeBracket:          incomeBracket,
		ScholarshipFlag:        scholarshipFlag,
		ScholarshipType:        normalizeStrPtr(in.ScholarshipType),
		WelfareFlag:            welfareFlag,
		WelfareType:            normalizeStrPtr(in.WelfareType),
		DebtFlag:               debtFlag,
		DebtDetail:             normalizeStrPtr(in.DebtDetail),
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

func (s *Service) ListDocuments(ctx context.Context, caseID uuid.UUID) ([]*ent.StudentRegistrationDocument, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.list_documents")
	defer span.End()

	rows, err := s.db.ListStudentRegistrationDocumentsByCaseID(ctx, caseID)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return rows, nil
}

func (s *Service) ReplaceDocuments(ctx context.Context, caseID uuid.UUID, actorRole ent.MemberRole, rows []DocumentInput) ([]*ent.StudentRegistrationDocument, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.replace_documents")
	defer span.End()

	if _, ok := parseApprovalActorRoleFromMember(actorRole); !ok {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseUnauthorized)
	}
	if err := s.ensureCaseEditable(ctx, caseID); err != nil {
		return nil, err
	}

	docs := make([]*ent.StudentRegistrationDocument, 0, len(rows))
	for _, row := range rows {
		docType, ok := parseRegistrationDocumentType(row.DocType)
		if !ok {
			return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
		}
		verifiedAt, err := parseDatePtr(row.VerifiedAt)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
		}
		isRequired := false
		if row.IsRequired != nil {
			isRequired = *row.IsRequired
		}
		isVerified := false
		if row.IsVerified != nil {
			isVerified = *row.IsVerified
		}
		docs = append(docs, &ent.StudentRegistrationDocument{
			CaseID:         caseID,
			DocType:        docType,
			FileDocumentID: row.FileDocumentID,
			FileName:       normalizeStrPtr(row.FileName),
			MimeType:       normalizeStrPtr(row.MimeType),
			FileSizeBytes:  row.FileSizeBytes,
			IsRequired:     isRequired,
			IsVerified:     isVerified,
			VerifiedBy:     row.VerifiedBy,
			VerifiedAt:     verifiedAt,
			Note:           normalizeStrPtr(row.Note),
		})
	}

	if err := s.db.ReplaceStudentRegistrationDocumentsByCaseID(ctx, caseID, docs); err != nil {
		return nil, normalizeServiceError(err)
	}
	return s.ListDocuments(ctx, caseID)
}

func (s *Service) ListRules(ctx context.Context, req *base.RequestPaginate, schoolID *uuid.UUID, registrationType *string) ([]*ent.StudentRegistrationRule, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.list_rules")
	defer span.End()

	var rt *ent.StudentRegistrationType
	if registrationType != nil && strings.TrimSpace(*registrationType) != "" {
		parsed, ok := parseRegistrationType(*registrationType)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
		}
		rt = &parsed
	}

	items, page, err := s.db.ListStudentRegistrationRules(ctx, req, schoolID, rt)
	if err != nil {
		return nil, nil, normalizeServiceError(err)
	}
	return items, page, nil
}

func (s *Service) GetRule(ctx context.Context, id uuid.UUID) (*ent.StudentRegistrationRule, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.get_rule")
	defer span.End()

	item, err := s.db.GetStudentRegistrationRuleByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

type RuleUpsertInput struct {
	SchoolID          uuid.UUID
	RegistrationType  string
	FieldCode         string
	IsRequired        *bool
	ValidationRegex   *string
	ValidationMessage *string
	ActiveFrom        *string
	ActiveTo          *string
}

func (s *Service) CreateRule(ctx context.Context, actorRole ent.MemberRole, in *RuleUpsertInput) (*ent.StudentRegistrationRule, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.create_rule")
	defer span.End()

	if actorRole != ent.MemberRoleAdmin && actorRole != ent.MemberRoleSuperadmin {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseUnauthorized)
	}
	if in == nil || strings.TrimSpace(in.FieldCode) == "" {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
	}
	regType, ok := parseRegistrationType(in.RegistrationType)
	if !ok {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
	}
	activeFrom := time.Now()
	if in.ActiveFrom != nil && strings.TrimSpace(*in.ActiveFrom) != "" {
		parsed, err := time.Parse("2006-01-02", strings.TrimSpace(*in.ActiveFrom))
		if err != nil {
			return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
		}
		activeFrom = parsed
	}
	var activeTo *time.Time
	if in.ActiveTo != nil && strings.TrimSpace(*in.ActiveTo) != "" {
		parsed, err := time.Parse("2006-01-02", strings.TrimSpace(*in.ActiveTo))
		if err != nil {
			return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
		}
		activeTo = &parsed
	}
	isRequired := false
	if in.IsRequired != nil {
		isRequired = *in.IsRequired
	}

	item, err := s.db.CreateStudentRegistrationRule(ctx, &ent.StudentRegistrationRule{
		SchoolID:          in.SchoolID,
		RegistrationType:  regType,
		FieldCode:         strings.TrimSpace(in.FieldCode),
		IsRequired:        isRequired,
		ValidationRegex:   normalizeStrPtr(in.ValidationRegex),
		ValidationMessage: normalizeStrPtr(in.ValidationMessage),
		ActiveFrom:        activeFrom,
		ActiveTo:          activeTo,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

func (s *Service) UpdateRule(ctx context.Context, id uuid.UUID, actorRole ent.MemberRole, in *RuleUpsertInput) (*ent.StudentRegistrationRule, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.update_rule")
	defer span.End()

	if actorRole != ent.MemberRoleAdmin && actorRole != ent.MemberRoleSuperadmin {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseUnauthorized)
	}
	if in == nil {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
	}

	update := &ent.StudentRegistrationRuleUpdate{}
	if in.SchoolID != uuid.Nil {
		update.SchoolID = &in.SchoolID
	}
	if strings.TrimSpace(in.RegistrationType) != "" {
		regType, ok := parseRegistrationType(in.RegistrationType)
		if !ok {
			return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
		}
		update.RegistrationType = &regType
	}
	if strings.TrimSpace(in.FieldCode) != "" {
		v := strings.TrimSpace(in.FieldCode)
		update.FieldCode = &v
	}
	if in.IsRequired != nil {
		update.IsRequired = in.IsRequired
	}
	if in.ValidationRegex != nil {
		update.ValidationRegex = normalizeStrPtr(in.ValidationRegex)
	}
	if in.ValidationMessage != nil {
		update.ValidationMessage = normalizeStrPtr(in.ValidationMessage)
	}
	if in.ActiveFrom != nil && strings.TrimSpace(*in.ActiveFrom) != "" {
		parsed, err := time.Parse("2006-01-02", strings.TrimSpace(*in.ActiveFrom))
		if err != nil {
			return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
		}
		update.ActiveFrom = &parsed
	}
	if in.ActiveTo != nil && strings.TrimSpace(*in.ActiveTo) != "" {
		parsed, err := time.Parse("2006-01-02", strings.TrimSpace(*in.ActiveTo))
		if err != nil {
			return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
		}
		update.ActiveTo = &parsed
	}

	item, err := s.db.UpdateStudentRegistrationRuleByID(ctx, id, update)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

func (s *Service) DeleteRule(ctx context.Context, id uuid.UUID, actorRole ent.MemberRole) error {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.delete_rule")
	defer span.End()

	if actorRole != ent.MemberRoleAdmin && actorRole != ent.MemberRoleSuperadmin {
		return fmt.Errorf("%w", ErrStudentRegistrationCaseUnauthorized)
	}
	if err := s.db.SoftDeleteStudentRegistrationRuleByID(ctx, id); err != nil {
		return normalizeServiceError(err)
	}
	return nil
}

func (s *Service) ensureCaseEditable(ctx context.Context, caseID uuid.UUID) error {
	item, err := s.db.GetStudentRegistrationCaseByID(ctx, caseID)
	if err != nil {
		return normalizeServiceError(err)
	}
	if item.Status != ent.StudentRegistrationCaseStatusDraft && item.Status != ent.StudentRegistrationCaseStatusRejected {
		return fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
	}
	return nil
}
