package teacherlicenses

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.TeacherLicense, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacherlicenses.service.update")
	defer span.End()

	if req.TeacherID == nil && req.LicenseNo == nil && req.IssuedAt == nil && req.ExpiresAt == nil && req.LicenseStatus == nil && req.IssuedBy == nil && req.Note == nil {
		return nil, fmt.Errorf("%w", ErrTeacherLicenseConditionFail)
	}

	payload := &ent.TeacherLicenseUpdate{
		LicenseNo: req.LicenseNo,
		IssuedBy:  req.IssuedBy,
		Note:      req.Note,
	}

	if req.TeacherID != nil {
		parsed, err := uuid.Parse(*req.TeacherID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrTeacherLicenseConditionFail)
		}
		payload.TeacherID = &parsed
	}
	if req.LicenseStatus != nil {
		parsed, ok := parseTeacherLicenseStatus(*req.LicenseStatus)
		if !ok {
			return nil, fmt.Errorf("%w", ErrTeacherLicenseConditionFail)
		}
		payload.LicenseStatus = &parsed
	}
	if req.IssuedAt != nil {
		if *req.IssuedAt == "" {
			return nil, fmt.Errorf("%w", ErrTeacherLicenseConditionFail)
		}
		parsed, err := time.Parse("2006-01-02", *req.IssuedAt)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrTeacherLicenseConditionFail)
		}
		payload.IssuedAt = &parsed
	}
	if req.ExpiresAt != nil {
		if *req.ExpiresAt == "" {
			return nil, fmt.Errorf("%w", ErrTeacherLicenseConditionFail)
		}
		parsed, err := time.Parse("2006-01-02", *req.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrTeacherLicenseConditionFail)
		}
		payload.ExpiresAt = &parsed
	}

	if payload.IssuedAt != nil && payload.ExpiresAt != nil && payload.ExpiresAt.Before(*payload.IssuedAt) {
		return nil, fmt.Errorf("%w", ErrTeacherLicenseConditionFail)
	}

	item, err := s.db.UpdateTeacherLicenseByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}
