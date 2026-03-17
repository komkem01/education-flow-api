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

	if req.MemberID == nil && req.EmployeeCode == nil && req.GenderID == nil && req.PrefixID == nil && req.FirstName == nil && req.LastName == nil && req.Phone == nil && req.Position == nil && req.StartWorkDate == nil && req.SchoolDepartmentID == nil && req.DepartmentID == nil && req.IsActive == nil {
		return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
	}

	phone, err := validateManagementPhone(req.Phone)
	if err != nil {
		return nil, err
	}

	payload := &ent.MemberManagementUpdate{
		EmployeeCode: req.EmployeeCode,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        phone,
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
	if req.GenderID != nil {
		parsed, err := uuid.Parse(*req.GenderID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
		}
		payload.GenderID = &parsed
	}
	if req.PrefixID != nil {
		parsed, err := uuid.Parse(*req.PrefixID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
		}
		payload.PrefixID = &parsed
	}
	if req.SchoolDepartmentID != nil {
		parsed, err := uuid.Parse(*req.SchoolDepartmentID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
		}
		schoolDepartment, err := s.schoolDepartmentDB.GetSchoolDepartmentByID(ctx, parsed)
		if err != nil {
			return nil, normalizeServiceError(err)
		}
		payload.SchoolDepartmentID = &schoolDepartment.ID
		payload.DepartmentID = &schoolDepartment.DepartmentID
	}
	if req.DepartmentID != nil {
		if req.SchoolDepartmentID == nil {
			return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
		}
		parsed, err := uuid.Parse(*req.DepartmentID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
		}
		if payload.DepartmentID != nil && *payload.DepartmentID != parsed {
			return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
		}
	}

	item, err := s.db.UpdateMemberManagementByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}
