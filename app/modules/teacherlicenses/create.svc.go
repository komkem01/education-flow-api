package teacherlicenses

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, teacherID uuid.UUID, licenseNo string, issuedAt *string, expiresAt *string, licenseStatus string, issuedBy *string, note *string) (*ent.TeacherLicense, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacherlicenses.service.create")
	defer span.End()

	parsedStatus, ok := parseTeacherLicenseStatus(licenseStatus)
	if !ok {
		return nil, fmt.Errorf("%w", ErrTeacherLicenseConditionFail)
	}

	var issuedAtVal *time.Time
	if issuedAt != nil && *issuedAt != "" {
		parsed, err := time.Parse("2006-01-02", *issuedAt)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrTeacherLicenseConditionFail)
		}
		issuedAtVal = &parsed
	}

	var expiresAtVal *time.Time
	if expiresAt != nil && *expiresAt != "" {
		parsed, err := time.Parse("2006-01-02", *expiresAt)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrTeacherLicenseConditionFail)
		}
		expiresAtVal = &parsed
	}

	if issuedAtVal != nil && expiresAtVal != nil && expiresAtVal.Before(*issuedAtVal) {
		return nil, fmt.Errorf("%w", ErrTeacherLicenseConditionFail)
	}

	item, err := s.db.CreateTeacherLicense(ctx, &ent.TeacherLicense{
		TeacherID:     teacherID,
		LicenseNo:     licenseNo,
		IssuedAt:      issuedAtVal,
		ExpiresAt:     expiresAtVal,
		LicenseStatus: parsedStatus,
		IssuedBy:      issuedBy,
		Note:          note,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}
