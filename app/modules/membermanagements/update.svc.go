package membermanagements

import (
	"context"
	"fmt"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, actorID uuid.UUID, actorRole ent.MemberRole, req *UpdateRequest) (*ent.ApprovalRequest, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "membermanagements.service.update")
	defer span.End()

	if req.MemberID == nil && req.EmployeeCode == nil && req.GenderID == nil && req.PrefixID == nil && req.FirstName == nil && req.LastName == nil && req.Phone == nil && req.Position == nil && req.StartWorkDate == nil && req.SchoolDepartmentID == nil && req.DepartmentID == nil && req.IsActive == nil {
		return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
	}

	if actorRole != ent.MemberRoleAdmin && actorRole != ent.MemberRoleSuperadmin {
		return nil, fmt.Errorf("%w", ErrMemberManagementUnauthorized)
	}

	item, err := s.db.GetMemberManagementByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	phone, err := validateManagementPhone(req.Phone)
	if err != nil {
		return nil, err
	}

	payload := map[string]any{
		"management_id": item.ID.String(),
		"member_id":     item.MemberID.String(),
	}

	if req.MemberID != nil {
		parsed, err := uuid.Parse(*req.MemberID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
		}
		payload["member_id"] = parsed.String()
	}
	if req.EmployeeCode != nil {
		v := strings.TrimSpace(*req.EmployeeCode)
		if v == "" {
			return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
		}
		payload["employee_code"] = v
	}
	if req.GenderID != nil {
		parsed, err := uuid.Parse(*req.GenderID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
		}
		payload["gender_id"] = parsed.String()
	}
	if req.PrefixID != nil {
		parsed, err := uuid.Parse(*req.PrefixID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
		}
		payload["prefix_id"] = parsed.String()
	}
	if req.FirstName != nil {
		v := strings.TrimSpace(*req.FirstName)
		if v != "" {
			payload["first_name"] = v
		}
	}
	if req.LastName != nil {
		v := strings.TrimSpace(*req.LastName)
		if v != "" {
			payload["last_name"] = v
		}
	}
	if phone != nil {
		payload["phone"] = *phone
	}
	if req.Position != nil {
		v := strings.TrimSpace(*req.Position)
		if v == "" {
			return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
		}
		payload["position"] = v
	}
	if req.StartWorkDate != nil {
		t, err := parseDate(strings.TrimSpace(*req.StartWorkDate))
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrMemberManagementConditionFail, err)
		}
		payload["start_work_date"] = t.Format("2006-01-02")
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
		payload["school_department_id"] = schoolDepartment.ID.String()
		payload["department_id"] = schoolDepartment.DepartmentID.String()
	}
	if req.DepartmentID != nil {
		parsed, err := uuid.Parse(*req.DepartmentID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
		}
		payload["department_id"] = parsed.String()
	}
	if req.IsActive != nil {
		payload["is_active"] = *req.IsActive
	}

	reason := normalizeManagementOptional(req.RequestReason)
	if reason != nil {
		payload["reason"] = *reason
	}

	now := time.Now()
	approval, err := s.approvalDB.CreateApprovalRequest(ctx, &ent.ApprovalRequest{
		RequestType:     "management_update",
		SubjectType:     "member_management",
		SubjectID:       &item.ID,
		RequestedBy:     actorID,
		RequestedByRole: ent.ApprovalActorRoleAdmin,
		Payload:         payload,
		CurrentStatus:   ent.ApprovalRequestStatusPending,
		SubmittedAt:     &now,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	_, err = s.actionDB.CreateApprovalAction(ctx, &ent.ApprovalAction{
		RequestID:   approval.ID,
		Action:      ent.ApprovalActionTypeSubmit,
		ActedBy:     actorID,
		ActedByRole: ent.ApprovalActorRoleAdmin,
		Comment:     reason,
		Metadata: map[string]any{
			"source":        "membermanagements.update",
			"management_id": item.ID.String(),
		},
		CreatedAt: now,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return approval, nil
}
