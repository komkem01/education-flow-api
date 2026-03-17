package memberteachers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, actorID uuid.UUID, actorRole ent.MemberRole, req *UpdateRequest, addresses *[]ent.TeacherAddressInput) (*ent.ApprovalRequest, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "memberteachers.service.update")
	defer span.End()

	if req.MemberID == nil && req.GenderID == nil && req.PrefixID == nil && req.Code == nil && req.CitizenID == nil && req.FirstNameTH == nil && req.LastNameTH == nil && req.FirstNameEN == nil && req.LastNameEN == nil && req.Phone == nil && req.Position == nil && req.AcademicStanding == nil && req.DepartmentID == nil && req.StartDate == nil && req.EndDate == nil && req.IsActive == nil && addresses == nil {
		return nil, fmt.Errorf("%w", ErrMemberTeacherConditionFail)
	}

	if actorRole != ent.MemberRoleAdmin && actorRole != ent.MemberRoleSuperadmin {
		return nil, fmt.Errorf("%w", ErrMemberTeacherUnauthorized)
	}

	teacher, err := s.db.GetMemberTeacherByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	payload := map[string]any{
		"teacher_id": teacher.ID.String(),
		"member_id":  teacher.MemberID.String(),
	}

	if req.MemberID != nil {
		payload["member_id"] = req.MemberID.String()
	}
	if req.GenderID != nil {
		payload["gender_id"] = req.GenderID.String()
	}
	if req.PrefixID != nil {
		payload["prefix_id"] = req.PrefixID.String()
	}
	if req.Code != nil {
		v := strings.TrimSpace(*req.Code)
		if v == "" {
			return nil, fmt.Errorf("%w", ErrMemberTeacherConditionFail)
		}
		payload["code"] = v
	}
	if req.CitizenID != nil {
		v, err := validateTeacherCitizenID(*req.CitizenID)
		if err != nil {
			return nil, err
		}
		payload["citizen_id"] = v
	}
	if req.FirstNameTH != nil {
		v := strings.TrimSpace(*req.FirstNameTH)
		if v == "" {
			return nil, fmt.Errorf("%w", ErrMemberTeacherConditionFail)
		}
		payload["first_name_th"] = v
	}
	if req.LastNameTH != nil {
		v := strings.TrimSpace(*req.LastNameTH)
		if v == "" {
			return nil, fmt.Errorf("%w", ErrMemberTeacherConditionFail)
		}
		payload["last_name_th"] = v
	}
	if req.FirstNameEN != nil {
		v := strings.TrimSpace(*req.FirstNameEN)
		if v == "" {
			return nil, fmt.Errorf("%w", ErrMemberTeacherConditionFail)
		}
		payload["first_name_en"] = v
	}
	if req.LastNameEN != nil {
		v := strings.TrimSpace(*req.LastNameEN)
		if v == "" {
			return nil, fmt.Errorf("%w", ErrMemberTeacherConditionFail)
		}
		payload["last_name_en"] = v
	}
	if req.Phone != nil {
		v, err := validateTeacherPhone(*req.Phone)
		if err != nil {
			return nil, err
		}
		payload["phone"] = v
	}
	if req.Position != nil {
		v := strings.TrimSpace(*req.Position)
		if v == "" {
			return nil, fmt.Errorf("%w", ErrMemberTeacherConditionFail)
		}
		payload["position"] = v
	}
	if req.AcademicStanding != nil {
		v := strings.TrimSpace(*req.AcademicStanding)
		if v == "" {
			v = "-"
		}
		payload["academic_standing"] = v
	}
	if req.DepartmentID != nil {
		payload["department_id"] = req.DepartmentID.String()
	}
	if req.StartDate != nil {
		t, err := parseDate(*req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrMemberTeacherConditionFail, err)
		}
		payload["start_date"] = t.Format("2006-01-02")
	}
	if req.EndDate != nil {
		t, err := parseOptionalDate(req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrMemberTeacherConditionFail, err)
		}
		if t != nil {
			payload["end_date"] = t.Format("2006-01-02")
		}
	}
	if req.IsActive != nil {
		payload["is_active"] = *req.IsActive
	}

	if addresses != nil {
		normalizedAddresses, err := normalizeTeacherAddresses(*addresses)
		if err != nil {
			return nil, err
		}
		addressPayload := make([]map[string]any, 0, len(normalizedAddresses))
		for _, addr := range normalizedAddresses {
			row := map[string]any{
				"house_no":    addr.HouseNo,
				"province":    addr.Province,
				"district":    addr.District,
				"subdistrict": addr.Subdistrict,
				"postal_code": addr.PostalCode,
				"is_primary":  addr.IsPrimary,
				"sort_order":  addr.SortOrder,
			}
			if addr.Village != nil {
				row["village"] = *addr.Village
			}
			if addr.Road != nil {
				row["road"] = *addr.Road
			}
			addressPayload = append(addressPayload, row)
		}
		payload["addresses"] = addressPayload
	}

	reason := normalizeTeacherOptional(req.ApprovalRequestReason)
	if reason != nil {
		payload["reason"] = *reason
	}

	now := time.Now()
	approval, err := s.approval.CreateApprovalRequest(ctx, &ent.ApprovalRequest{
		RequestType:     "teacher_update",
		SubjectType:     "member_teacher",
		SubjectID:       &teacher.ID,
		RequestedBy:     actorID,
		RequestedByRole: ent.ApprovalActorRoleAdmin,
		Payload:         payload,
		CurrentStatus:   ent.ApprovalRequestStatusPending,
		SubmittedAt:     &now,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	_, err = s.action.CreateApprovalAction(ctx, &ent.ApprovalAction{
		RequestID:   approval.ID,
		Action:      ent.ApprovalActionTypeSubmit,
		ActedBy:     actorID,
		ActedByRole: ent.ApprovalActorRoleAdmin,
		Comment:     reason,
		Metadata: map[string]any{
			"source":     "memberteachers.update",
			"teacher_id": teacher.ID.String(),
		},
		CreatedAt: now,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return approval, nil
}
