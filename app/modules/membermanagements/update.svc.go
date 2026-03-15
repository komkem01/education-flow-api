package membermanagements

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.MemberManagement, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "membermanagements.service.update")
	defer span.End()

	if req.MemberID == nil && req.EmployeeCode == nil && req.Position == nil && req.StartWorkDate == nil && req.DepartmentID == nil && req.IsActive == nil {
		return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
	}

	payload := &ent.MemberManagementUpdate{
		EmployeeCode: req.EmployeeCode,
		Position:     req.Position,
		IsActive:     req.IsActive,
	}

	if req.MemberID != nil {
		parsed, err := uuid.Parse(*req.MemberID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
		}
		payload.MemberID = &parsed
	}
	if req.StartWorkDate != nil {
		t, err := parseDate(*req.StartWorkDate)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrMemberManagementConditionFail, err)
		}
		payload.StartWorkDate = &t
	}
	if req.DepartmentID != nil {
		parsed, err := uuid.Parse(*req.DepartmentID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
		}
		payload.DepartmentID = &parsed
	}

	item, err := s.db.UpdateMemberManagementByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}
