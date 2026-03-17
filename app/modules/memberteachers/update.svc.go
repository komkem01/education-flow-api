package memberteachers

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest, addresses *[]ent.TeacherAddressInput) (*ent.MemberTeacher, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "memberteachers.service.update")
	defer span.End()

	if req.MemberID == nil && req.GenderID == nil && req.PrefixID == nil && req.Code == nil && req.CitizenID == nil && req.FirstNameTH == nil && req.LastNameTH == nil && req.FirstNameEN == nil && req.LastNameEN == nil && req.Phone == nil && req.Position == nil && req.AcademicStanding == nil && req.DepartmentID == nil && req.StartDate == nil && req.EndDate == nil && req.IsActive == nil && addresses == nil {
		return nil, fmt.Errorf("%w", ErrMemberTeacherConditionFail)
	}

	payload := &ent.MemberTeacherUpdate{
		MemberID:         req.MemberID,
		GenderID:         req.GenderID,
		PrefixID:         req.PrefixID,
		Code:             req.Code,
		CitizenID:        req.CitizenID,
		FirstNameTH:      req.FirstNameTH,
		LastNameTH:       req.LastNameTH,
		FirstNameEN:      req.FirstNameEN,
		LastNameEN:       req.LastNameEN,
		Phone:            req.Phone,
		Position:         req.Position,
		AcademicStanding: req.AcademicStanding,
		DepartmentID:     req.DepartmentID,
		IsActive:         req.IsActive,
	}

	if req.StartDate != nil {
		t, err := parseDate(*req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrMemberTeacherConditionFail, err)
		}
		payload.StartDate = &t
	}

	if req.EndDate != nil {
		t, err := parseOptionalDate(req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrMemberTeacherConditionFail, err)
		}
		payload.EndDate = t
	}

	teacher, err := s.db.UpdateMemberTeacherByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	if addresses != nil {
		normalizedAddresses, err := normalizeTeacherAddresses(*addresses)
		if err != nil {
			return nil, err
		}

		if err := s.db.ReplaceTeacherAddressesByMemberTeacherID(ctx, id, normalizedAddresses); err != nil {
			return nil, normalizeServiceError(err)
		}
	}

	return teacher, nil
}
