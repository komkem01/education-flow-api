package membermanagements

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func parseDate(v string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid-date-format: %w", err)
	}
	return t, nil
}

func (s *Service) Create(ctx context.Context, memberID uuid.UUID, employeeCode string, position string, startWorkDate string, departmentID uuid.UUID, isActive bool) (*ent.MemberManagement, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "membermanagements.service.create")
	defer span.End()

	startAt, err := parseDate(startWorkDate)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMemberManagementConditionFail, err)
	}

	item, err := s.db.CreateMemberManagement(ctx, &ent.MemberManagement{
		MemberID:      memberID,
		EmployeeCode:  employeeCode,
		Position:      position,
		StartWorkDate: startAt,
		DepartmentID:  departmentID,
		IsActive:      isActive,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}
