package studentregistrationcases

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"
	"eduflow/app/utils/hashing"

	"github.com/google/uuid"
)

type UpsertInput struct {
	CaseNo            *string
	SchoolID          *uuid.UUID
	StudentID         *uuid.UUID
	RegistrationType  *string
	EffectiveDate     *string
	Reason            *string
	Email             *string
	Password          *string
	Core              *CoreInput
	Addresses         []AddressInput
	Health            *HealthInput
	Guardians         []GuardianInput
	StudentGuardians  []StudentGuardianInput
	PreviousEducation *PreviousEducationInput
	FamilyEconomic    *FamilyEconomicInput
	Documents         []DocumentInput
}

type CoreInput struct {
	MemberID         *uuid.UUID
	StudentID        *uuid.UUID
	GenderID         uuid.UUID
	PrefixID         uuid.UUID
	AdvisorTeacherID *uuid.UUID
	FirstNameTH      string
	LastNameTH       string
	FirstNameEN      *string
	LastNameEN       *string
	CitizenID        *string
	Phone            *string
	IsActiveTarget   *bool
}

type AddressInput struct {
	AddressType string
	HouseNo     string
	Village     *string
	Road        *string
	Province    string
	District    string
	Subdistrict string
	PostalCode  string
	Country     *string
	IsPrimary   *bool
	SortOrder   *int
}

type HealthInput struct {
	BloodType            *string
	AllergyInfo          *string
	ChronicDisease       *string
	MedicalNote          *string
	DisabilityFlag       *bool
	DisabilityDetail     *string
	SpecialSupportFlag   *bool
	SpecialSupportDetail *string
}

type GuardianInput struct {
	RowID            *uuid.UUID
	GenderID         uuid.UUID
	PrefixID         uuid.UUID
	FirstNameTH      string
	LastNameTH       string
	FirstNameEN      *string
	LastNameEN       *string
	CitizenID        *string
	Phone            string
	Occupation       *string
	Employer         *string
	MonthlyIncome    *float64
	AnnualIncome     *float64
	EducationLevel   *string
	RelationshipText *string
	IsActiveTarget   *bool
	SortOrder        *int
}

type StudentGuardianInput struct {
	GuardianRowID      uuid.UUID
	Relationship       string
	IsMainGuardian     *bool
	CanPickup          *bool
	IsEmergencyContact *bool
	Note               *string
	SortOrder          *int
}

type PreviousEducationInput struct {
	PreviousSchoolName     *string
	PreviousSchoolProvince *string
	PreviousGradeLevel     *string
	GPA                    *float64
	TransferCertificateNo  *string
	TransferDate           *string
	TranscriptReceived     *bool
	Remarks                *string
}

type FamilyEconomicInput struct {
	HouseholdSize          *int
	HouseholdIncomeMonthly *float64
	IncomeBracket          *string
	ScholarshipFlag        *bool
	ScholarshipType        *string
	WelfareFlag            *bool
	WelfareType            *string
	DebtFlag               *bool
	DebtDetail             *string
}

type DocumentInput struct {
	DocType        string
	FileDocumentID *uuid.UUID
	FileName       *string
	MimeType       *string
	FileSizeBytes  *int64
	IsRequired     *bool
	IsVerified     *bool
	VerifiedBy     *uuid.UUID
	VerifiedAt     *string
	Note           *string
}

func (s *Service) Create(ctx context.Context, actorID uuid.UUID, actorRole ent.MemberRole, in *UpsertInput) (*ent.StudentRegistrationCaseDetail, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.create")
	defer span.End()

	approvalRole, ok := parseApprovalActorRoleFromMember(actorRole)
	if !ok {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseUnauthorized)
	}
	if in == nil || in.SchoolID == nil || in.RegistrationType == nil {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
	}
	regType, ok := parseRegistrationType(*in.RegistrationType)
	if !ok {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
	}

	now := time.Now()
	caseNo, err := utils.GenerateNumericCode("SRC", 6)
	if err != nil {
		return nil, err
	}
	caseNo = strings.TrimSpace(caseNo)
	if in.CaseNo != nil && strings.TrimSpace(*in.CaseNo) != "" {
		caseNo = strings.TrimSpace(*in.CaseNo)
	}
	effectiveDate, err := parseDatePtr(in.EffectiveDate)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
	}

	item, err := s.db.CreateStudentRegistrationCase(ctx, &ent.StudentRegistrationCase{
		CaseNo:           caseNo,
		SchoolID:         *in.SchoolID,
		StudentID:        in.StudentID,
		RegistrationType: regType,
		Status:           ent.StudentRegistrationCaseStatusDraft,
		RequestedBy:      actorID,
		RequestedByRole:  approvalRole,
		RequestedAt:      now,
		EffectiveDate:    effectiveDate,
		Reason:           normalizeStrPtr(in.Reason),
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	bundle, err := s.toBundle(item.ID, in)
	if err != nil {
		return nil, err
	}
	if err := s.db.ReplaceStudentRegistrationCaseBundle(ctx, item.ID, bundle); err != nil {
		return nil, normalizeServiceError(err)
	}

	detail, err := s.db.GetStudentRegistrationCaseDetailByID(ctx, item.ID)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return detail, nil
}

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, schoolID *uuid.UUID, studentID *uuid.UUID, status *string, registrationType *string) ([]*ent.StudentRegistrationCase, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.list")
	defer span.End()

	var st *ent.StudentRegistrationCaseStatus
	if status != nil && strings.TrimSpace(*status) != "" {
		parsed, ok := parseRegistrationStatus(*status)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
		}
		st = &parsed
	}

	var rt *ent.StudentRegistrationType
	if registrationType != nil && strings.TrimSpace(*registrationType) != "" {
		parsed, ok := parseRegistrationType(*registrationType)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
		}
		rt = &parsed
	}

	items, page, err := s.db.ListStudentRegistrationCases(ctx, req, schoolID, studentID, st, rt)
	if err != nil {
		return nil, nil, normalizeServiceError(err)
	}
	return items, page, nil
}

func (s *Service) Info(ctx context.Context, id uuid.UUID) (*ent.StudentRegistrationCaseDetail, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.info")
	defer span.End()

	item, err := s.db.GetStudentRegistrationCaseDetailByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, actorRole ent.MemberRole, in *UpsertInput) (*ent.StudentRegistrationCaseDetail, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.update")
	defer span.End()

	if _, ok := parseApprovalActorRoleFromMember(actorRole); !ok {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseUnauthorized)
	}

	current, err := s.db.GetStudentRegistrationCaseByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	if current.Status != ent.StudentRegistrationCaseStatusDraft && current.Status != ent.StudentRegistrationCaseStatusRejected {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
	}

	update := &ent.StudentRegistrationCaseUpdate{}
	if in != nil {
		if in.CaseNo != nil {
			v := strings.TrimSpace(*in.CaseNo)
			if v == "" {
				return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
			}
			update.CaseNo = &v
		}
		if in.SchoolID != nil {
			update.SchoolID = in.SchoolID
		}
		if in.StudentID != nil {
			update.StudentID = in.StudentID
		}
		if in.RegistrationType != nil {
			parsed, ok := parseRegistrationType(*in.RegistrationType)
			if !ok {
				return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
			}
			update.RegistrationType = &parsed
		}
		if in.EffectiveDate != nil {
			v, err := parseDatePtr(in.EffectiveDate)
			if err != nil {
				return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
			}
			if v != nil {
				update.EffectiveDate = v
			}
		}
		if in.Reason != nil {
			update.Reason = normalizeStrPtr(in.Reason)
		}
	}

	if _, err := s.db.UpdateStudentRegistrationCaseByID(ctx, id, update); err != nil {
		return nil, normalizeServiceError(err)
	}

	if in != nil {
		bundle, err := s.toBundle(id, in)
		if err != nil {
			return nil, err
		}
		if err := s.db.ReplaceStudentRegistrationCaseBundle(ctx, id, bundle); err != nil {
			return nil, normalizeServiceError(err)
		}
	}

	item, err := s.db.GetStudentRegistrationCaseDetailByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID, actorRole ent.MemberRole) error {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.delete")
	defer span.End()

	if _, ok := parseApprovalActorRoleFromMember(actorRole); !ok {
		return fmt.Errorf("%w", ErrStudentRegistrationCaseUnauthorized)
	}
	if err := s.db.SoftDeleteStudentRegistrationCaseByID(ctx, id); err != nil {
		return normalizeServiceError(err)
	}
	return nil
}

func (s *Service) Submit(ctx context.Context, id uuid.UUID, actorID uuid.UUID, actorRole ent.MemberRole, reason *string) (*ent.StudentRegistrationCase, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.submit")
	defer span.End()

	approvalRole, ok := parseApprovalActorRoleFromMember(actorRole)
	if !ok {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseUnauthorized)
	}
	detail, err := s.db.GetStudentRegistrationCaseDetailByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	if detail != nil && detail.Case != nil && detail.Case.RegistrationType == ent.StudentRegistrationTypeNewEnrollment {
		if detail.Case.StudentID == nil {
			if detail.Core == nil || detail.Core.PendingMemberEmail == nil || detail.Core.PendingMemberPasswordHash == nil {
				return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
			}
		}
	}
	item, err := s.db.SubmitStudentRegistrationCase(ctx, id, actorID, approvalRole, normalizeStrPtr(reason))
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

func (s *Service) Reject(ctx context.Context, id uuid.UUID, actorID uuid.UUID, actorRole ent.MemberRole, reason *string) (*ent.StudentRegistrationCase, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.reject")
	defer span.End()

	approvalRole, ok := parseApprovalActorRoleFromMember(actorRole)
	if !ok || approvalRole != ent.ApprovalActorRoleAdmin {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseUnauthorized)
	}
	item, err := s.db.RejectStudentRegistrationCase(ctx, id, actorID, approvalRole, normalizeStrPtr(reason))
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

func (s *Service) Approve(ctx context.Context, id uuid.UUID, actorID uuid.UUID, actorRole ent.MemberRole, comment *string, idempotencyKey *string, metadata map[string]any) (*ent.StudentRegistrationCase, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.approve")
	defer span.End()

	approvalRole, ok := parseApprovalActorRoleFromMember(actorRole)
	if !ok || approvalRole != ent.ApprovalActorRoleAdmin {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseUnauthorized)
	}
	item, err := s.db.ApproveAndApplyStudentRegistrationCase(ctx, id, actorID, approvalRole, normalizeStrPtr(comment), normalizeStrPtr(idempotencyKey), metadata)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

func (s *Service) Apply(ctx context.Context, id uuid.UUID, actorID uuid.UUID, actorRole ent.MemberRole, comment *string) (*ent.StudentRegistrationCase, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentregistrationcases.service.apply")
	defer span.End()

	approvalRole, ok := parseApprovalActorRoleFromMember(actorRole)
	if !ok || approvalRole != ent.ApprovalActorRoleAdmin {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseUnauthorized)
	}
	item, err := s.db.ApplyStudentRegistrationCase(ctx, id, actorID, approvalRole, normalizeStrPtr(comment))
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

func (s *Service) toBundle(caseID uuid.UUID, in *UpsertInput) (*ent.StudentRegistrationCaseBundle, error) {
	bundle := &ent.StudentRegistrationCaseBundle{}
	if in == nil {
		return bundle, nil
	}

	pendingMemberEmail, pendingMemberPasswordHash, err := preparePendingMemberCredentials(in.Email, in.Password)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
	}

	if in.Core != nil {
		if strings.TrimSpace(in.Core.FirstNameTH) == "" || strings.TrimSpace(in.Core.LastNameTH) == "" {
			return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
		}
		isActive := true
		if in.Core.IsActiveTarget != nil {
			isActive = *in.Core.IsActiveTarget
		}
		bundle.Core = &ent.StudentRegistrationStudentCore{
			CaseID:                    caseID,
			MemberID:                  in.Core.MemberID,
			StudentID:                 in.Core.StudentID,
			GenderID:                  in.Core.GenderID,
			PrefixID:                  in.Core.PrefixID,
			AdvisorTeacherID:          in.Core.AdvisorTeacherID,
			FirstNameTH:               strings.TrimSpace(in.Core.FirstNameTH),
			LastNameTH:                strings.TrimSpace(in.Core.LastNameTH),
			FirstNameEN:               normalizeStrPtr(in.Core.FirstNameEN),
			LastNameEN:                normalizeStrPtr(in.Core.LastNameEN),
			CitizenID:                 normalizeStrPtr(in.Core.CitizenID),
			Phone:                     normalizeStrPtr(in.Core.Phone),
			PendingMemberEmail:        pendingMemberEmail,
			PendingMemberPasswordHash: pendingMemberPasswordHash,
			IsActiveTarget:            isActive,
		}
	} else if pendingMemberEmail != nil || pendingMemberPasswordHash != nil {
		return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
	}

	for idx, row := range in.Addresses {
		addressType, ok := parseAddressType(row.AddressType)
		if !ok || strings.TrimSpace(row.HouseNo) == "" || strings.TrimSpace(row.Province) == "" || strings.TrimSpace(row.District) == "" || strings.TrimSpace(row.Subdistrict) == "" || strings.TrimSpace(row.PostalCode) == "" {
			return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
		}
		isPrimary := false
		if row.IsPrimary != nil {
			isPrimary = *row.IsPrimary
		}
		sortOrder := idx + 1
		if row.SortOrder != nil && *row.SortOrder > 0 {
			sortOrder = *row.SortOrder
		}
		country := "TH"
		if row.Country != nil && strings.TrimSpace(*row.Country) != "" {
			country = strings.TrimSpace(*row.Country)
		}
		bundle.Addresses = append(bundle.Addresses, &ent.StudentRegistrationAddress{
			CaseID:      caseID,
			AddressType: addressType,
			HouseNo:     strings.TrimSpace(row.HouseNo),
			Village:     normalizeStrPtr(row.Village),
			Road:        normalizeStrPtr(row.Road),
			Province:    strings.TrimSpace(row.Province),
			District:    strings.TrimSpace(row.District),
			Subdistrict: strings.TrimSpace(row.Subdistrict),
			PostalCode:  strings.TrimSpace(row.PostalCode),
			Country:     country,
			IsPrimary:   isPrimary,
			SortOrder:   sortOrder,
		})
	}

	if in.Health != nil {
		disabilityFlag := false
		if in.Health.DisabilityFlag != nil {
			disabilityFlag = *in.Health.DisabilityFlag
		}
		specialSupportFlag := false
		if in.Health.SpecialSupportFlag != nil {
			specialSupportFlag = *in.Health.SpecialSupportFlag
		}
		bundle.Health = &ent.StudentRegistrationHealth{
			CaseID:               caseID,
			BloodType:            normalizeStrPtr(in.Health.BloodType),
			AllergyInfo:          normalizeStrPtr(in.Health.AllergyInfo),
			ChronicDisease:       normalizeStrPtr(in.Health.ChronicDisease),
			MedicalNote:          normalizeStrPtr(in.Health.MedicalNote),
			DisabilityFlag:       disabilityFlag,
			DisabilityDetail:     normalizeStrPtr(in.Health.DisabilityDetail),
			SpecialSupportFlag:   specialSupportFlag,
			SpecialSupportDetail: normalizeStrPtr(in.Health.SpecialSupportDetail),
		}
	}

	for idx, row := range in.Guardians {
		if strings.TrimSpace(row.FirstNameTH) == "" || strings.TrimSpace(row.LastNameTH) == "" || strings.TrimSpace(row.Phone) == "" {
			return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
		}
		sortOrder := idx + 1
		if row.SortOrder != nil && *row.SortOrder > 0 {
			sortOrder = *row.SortOrder
		}
		isActive := true
		if row.IsActiveTarget != nil {
			isActive = *row.IsActiveTarget
		}
		rowID := uuid.Nil
		if row.RowID != nil {
			rowID = *row.RowID
		}
		bundle.Guardians = append(bundle.Guardians, &ent.StudentRegistrationGuardian{
			ID:               rowID,
			CaseID:           caseID,
			GenderID:         row.GenderID,
			PrefixID:         row.PrefixID,
			FirstNameTH:      strings.TrimSpace(row.FirstNameTH),
			LastNameTH:       strings.TrimSpace(row.LastNameTH),
			FirstNameEN:      normalizeStrPtr(row.FirstNameEN),
			LastNameEN:       normalizeStrPtr(row.LastNameEN),
			CitizenID:        normalizeStrPtr(row.CitizenID),
			Phone:            strings.TrimSpace(row.Phone),
			Occupation:       normalizeStrPtr(row.Occupation),
			Employer:         normalizeStrPtr(row.Employer),
			MonthlyIncome:    row.MonthlyIncome,
			AnnualIncome:     row.AnnualIncome,
			EducationLevel:   normalizeStrPtr(row.EducationLevel),
			RelationshipText: normalizeStrPtr(row.RelationshipText),
			IsActiveTarget:   isActive,
			SortOrder:        sortOrder,
		})
	}

	for idx, row := range in.StudentGuardians {
		if row.GuardianRowID == uuid.Nil || strings.TrimSpace(row.Relationship) == "" {
			return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
		}
		sortOrder := idx + 1
		if row.SortOrder != nil && *row.SortOrder > 0 {
			sortOrder = *row.SortOrder
		}
		isMain := false
		if row.IsMainGuardian != nil {
			isMain = *row.IsMainGuardian
		}
		canPickup := true
		if row.CanPickup != nil {
			canPickup = *row.CanPickup
		}
		eContact := false
		if row.IsEmergencyContact != nil {
			eContact = *row.IsEmergencyContact
		}
		bundle.StudentGuardians = append(bundle.StudentGuardians, &ent.StudentRegistrationStudentGuardian{
			CaseID:             caseID,
			GuardianRowID:      row.GuardianRowID,
			Relationship:       strings.TrimSpace(row.Relationship),
			IsMainGuardian:     isMain,
			CanPickup:          canPickup,
			IsEmergencyContact: eContact,
			Note:               normalizeStrPtr(row.Note),
			SortOrder:          sortOrder,
		})
	}

	if in.PreviousEducation != nil {
		transferDate, err := parseDatePtr(in.PreviousEducation.TransferDate)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
		}
		transcriptReceived := false
		if in.PreviousEducation.TranscriptReceived != nil {
			transcriptReceived = *in.PreviousEducation.TranscriptReceived
		}
		bundle.PreviousEducation = &ent.StudentRegistrationPreviousEducation{
			CaseID:                 caseID,
			PreviousSchoolName:     normalizeStrPtr(in.PreviousEducation.PreviousSchoolName),
			PreviousSchoolProvince: normalizeStrPtr(in.PreviousEducation.PreviousSchoolProvince),
			PreviousGradeLevel:     normalizeStrPtr(in.PreviousEducation.PreviousGradeLevel),
			GPA:                    in.PreviousEducation.GPA,
			TransferCertificateNo:  normalizeStrPtr(in.PreviousEducation.TransferCertificateNo),
			TransferDate:           transferDate,
			TranscriptReceived:     transcriptReceived,
			Remarks:                normalizeStrPtr(in.PreviousEducation.Remarks),
		}
	}

	if in.FamilyEconomic != nil {
		var incomeBracket *ent.RegistrationIncomeBracket
		if in.FamilyEconomic.IncomeBracket != nil {
			parsed, ok := parseIncomeBracket(*in.FamilyEconomic.IncomeBracket)
			if !ok {
				return nil, fmt.Errorf("%w", ErrStudentRegistrationCaseConditionFail)
			}
			incomeBracket = &parsed
		}
		scholarshipFlag := false
		if in.FamilyEconomic.ScholarshipFlag != nil {
			scholarshipFlag = *in.FamilyEconomic.ScholarshipFlag
		}
		welfareFlag := false
		if in.FamilyEconomic.WelfareFlag != nil {
			welfareFlag = *in.FamilyEconomic.WelfareFlag
		}
		debtFlag := false
		if in.FamilyEconomic.DebtFlag != nil {
			debtFlag = *in.FamilyEconomic.DebtFlag
		}
		bundle.FamilyEconomic = &ent.StudentRegistrationFamilyEconomic{
			CaseID:                 caseID,
			HouseholdSize:          in.FamilyEconomic.HouseholdSize,
			HouseholdIncomeMonthly: in.FamilyEconomic.HouseholdIncomeMonthly,
			IncomeBracket:          incomeBracket,
			ScholarshipFlag:        scholarshipFlag,
			ScholarshipType:        normalizeStrPtr(in.FamilyEconomic.ScholarshipType),
			WelfareFlag:            welfareFlag,
			WelfareType:            normalizeStrPtr(in.FamilyEconomic.WelfareType),
			DebtFlag:               debtFlag,
			DebtDetail:             normalizeStrPtr(in.FamilyEconomic.DebtDetail),
		}
	}

	for _, row := range in.Documents {
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
		bundle.Documents = append(bundle.Documents, &ent.StudentRegistrationDocument{
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

	return bundle, nil
}

func preparePendingMemberCredentials(email *string, password *string) (*string, *string, error) {
	normalizedEmail := normalizeStrPtr(email)
	normalizedPassword := normalizeStrPtr(password)

	if normalizedEmail == nil && normalizedPassword == nil {
		return nil, nil, nil
	}
	if normalizedEmail == nil || normalizedPassword == nil {
		return nil, nil, fmt.Errorf("credentials-required")
	}
	if _, err := mail.ParseAddress(*normalizedEmail); err != nil {
		return nil, nil, err
	}
	if len(*normalizedPassword) < 8 {
		return nil, nil, fmt.Errorf("password-too-short")
	}
	emailLower := strings.ToLower(strings.TrimSpace(*normalizedEmail))
	hashed, err := hashing.HashPassword(*normalizedPassword)
	if err != nil {
		return nil, nil, err
	}
	hash := string(hashed)
	return &emailLower, &hash, nil
}
