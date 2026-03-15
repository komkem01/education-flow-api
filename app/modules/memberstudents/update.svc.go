package memberstudents

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.MemberStudent, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "memberstudents.service.update")
	defer span.End()

	if req.MemberID == nil && req.SchoolID == nil && req.GenderID == nil && req.PrefixID == nil && req.AdvisorTeacherID == nil && req.StudentCode == nil && req.FirstNameTH == nil && req.LastNameTH == nil && req.FirstNameEN == nil && req.LastNameEN == nil && req.CitizenID == nil && req.Phone == nil && req.IsActive == nil {
		return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
	}

	payload := &ent.MemberStudentUpdate{
		StudentCode: req.StudentCode,
		FirstNameTH: req.FirstNameTH,
		LastNameTH:  req.LastNameTH,
		FirstNameEN: req.FirstNameEN,
		LastNameEN:  req.LastNameEN,
		CitizenID:   req.CitizenID,
		Phone:       req.Phone,
		IsActive:    req.IsActive,
	}

	if req.MemberID != nil {
		parsed, err := uuid.Parse(*req.MemberID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
		}
		payload.MemberID = &parsed
	}
	if req.SchoolID != nil {
		parsed, err := uuid.Parse(*req.SchoolID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
		}
		payload.SchoolID = &parsed
	}
	if req.GenderID != nil {
		parsed, err := uuid.Parse(*req.GenderID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
		}
		payload.GenderID = &parsed
	}
	if req.PrefixID != nil {
		parsed, err := uuid.Parse(*req.PrefixID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
		}
		payload.PrefixID = &parsed
	}
	if req.AdvisorTeacherID != nil {
		parsed, err := uuid.Parse(*req.AdvisorTeacherID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
		}
		payload.AdvisorTeacherID = &parsed
	}

	item, err := s.db.UpdateMemberStudentByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}
