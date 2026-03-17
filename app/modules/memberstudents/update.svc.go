package memberstudents

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
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "memberstudents.service.update")
	defer span.End()

	if req.MemberID == nil && req.SchoolID == nil && req.GenderID == nil && req.PrefixID == nil && req.AdvisorTeacherID == nil && req.StudentCode == nil && req.FirstNameTH == nil && req.LastNameTH == nil && req.FirstNameEN == nil && req.LastNameEN == nil && req.CitizenID == nil && req.Phone == nil && req.BirthDate == nil && req.Nationality == nil && req.Religion == nil && req.AddressCurrent == nil && req.AddressRegistered == nil && req.EmergencyContactName == nil && req.EmergencyContactPhone == nil && req.BloodType == nil && req.AllergyInfo == nil && req.ChronicDisease == nil && req.MedicalNote == nil && req.IsActive == nil {
		return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
	}

	student, err := s.db.GetMemberStudentByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	requestedByRole := ent.ApprovalActorRoleAdmin
	switch actorRole {
	case ent.MemberRoleTeacher:
		requestedByRole = ent.ApprovalActorRoleTeacher
	case ent.MemberRoleAdmin, ent.MemberRoleSuperadmin:
		requestedByRole = ent.ApprovalActorRoleAdmin
	default:
		return nil, fmt.Errorf("%w", ErrMemberStudentUnauthorized)
	}

	payload := map[string]any{
		"student_id": student.ID.String(),
		"member_id":  student.MemberID.String(),
	}

	if req.MemberID != nil {
		parsed, err := uuid.Parse(*req.MemberID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
		}
		payload["member_id"] = parsed.String()
	}
	if req.SchoolID != nil {
		parsed, err := uuid.Parse(*req.SchoolID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
		}
		payload["school_id"] = parsed.String()
	}
	if req.GenderID != nil {
		parsed, err := uuid.Parse(*req.GenderID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
		}
		payload["gender_id"] = parsed.String()
	}
	if req.PrefixID != nil {
		parsed, err := uuid.Parse(*req.PrefixID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
		}
		payload["prefix_id"] = parsed.String()
	}
	if req.AdvisorTeacherID != nil {
		parsed, err := uuid.Parse(*req.AdvisorTeacherID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
		}
		payload["advisor_teacher_id"] = parsed.String()
	}
	if req.StudentCode != nil {
		v := strings.TrimSpace(*req.StudentCode)
		if v == "" {
			return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
		}
		payload["student_code"] = v
	}
	if req.FirstNameTH != nil {
		v := strings.TrimSpace(*req.FirstNameTH)
		if v == "" {
			return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
		}
		payload["first_name_th"] = v
	}
	if req.LastNameTH != nil {
		v := strings.TrimSpace(*req.LastNameTH)
		if v == "" {
			return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
		}
		payload["last_name_th"] = v
	}
	if req.FirstNameEN != nil {
		v := strings.TrimSpace(*req.FirstNameEN)
		if v != "" {
			payload["first_name_en"] = v
		}
	}
	if req.LastNameEN != nil {
		v := strings.TrimSpace(*req.LastNameEN)
		if v != "" {
			payload["last_name_en"] = v
		}
	}
	if req.CitizenID != nil {
		v, err := validateOptionalCitizenID(req.CitizenID)
		if err != nil {
			return nil, err
		}
		if v != nil {
			payload["citizen_id"] = *v
		}
	}
	if req.Phone != nil {
		v, err := validateOptionalPhone(req.Phone)
		if err != nil {
			return nil, err
		}
		if v != nil {
			payload["phone"] = *v
		}
	}
	if req.BirthDate != nil {
		if strings.TrimSpace(*req.BirthDate) != "" {
			parsed, err := time.Parse("2006-01-02", strings.TrimSpace(*req.BirthDate))
			if err != nil || parsed.After(time.Now()) {
				return nil, fmt.Errorf("%w", ErrInvalidBirthDate)
			}
			payload["birth_date"] = parsed.Format("2006-01-02")
		}
	}
	if req.Nationality != nil {
		v := strings.TrimSpace(*req.Nationality)
		if v != "" {
			payload["nationality"] = v
		}
	}
	if req.Religion != nil {
		v := strings.TrimSpace(*req.Religion)
		if v != "" {
			payload["religion"] = v
		}
	}
	if req.AddressCurrent != nil {
		v := strings.TrimSpace(*req.AddressCurrent)
		if v != "" {
			payload["address_current"] = v
		}
	}
	if req.AddressRegistered != nil {
		v := strings.TrimSpace(*req.AddressRegistered)
		if v != "" {
			payload["address_registered"] = v
		}
	}
	if req.EmergencyContactName != nil {
		v := strings.TrimSpace(*req.EmergencyContactName)
		if v != "" {
			payload["emergency_contact_name"] = v
		}
	}
	if req.EmergencyContactPhone != nil {
		v, err := validateOptionalPhone(req.EmergencyContactPhone)
		if err != nil {
			return nil, err
		}
		if v != nil {
			payload["emergency_contact_phone"] = *v
		}
	}
	if req.BloodType != nil {
		v := strings.TrimSpace(*req.BloodType)
		if v != "" {
			payload["blood_type"] = v
		}
	}
	if req.AllergyInfo != nil {
		v := strings.TrimSpace(*req.AllergyInfo)
		if v != "" {
			payload["allergy_info"] = v
		}
	}
	if req.ChronicDisease != nil {
		v := strings.TrimSpace(*req.ChronicDisease)
		if v != "" {
			payload["chronic_disease"] = v
		}
	}
	if req.MedicalNote != nil {
		v := strings.TrimSpace(*req.MedicalNote)
		if v != "" {
			payload["medical_note"] = v
		}
	}
	if req.IsActive != nil {
		payload["is_active"] = *req.IsActive
	}

	reason := normalizeOptional(req.ApprovalRequestReason)
	if reason != nil {
		payload["reason"] = *reason
	}

	now := time.Now()
	approval, err := s.approval.CreateApprovalRequest(ctx, &ent.ApprovalRequest{
		RequestType:     "student_update",
		SubjectType:     "member_student",
		SubjectID:       &student.ID,
		RequestedBy:     actorID,
		RequestedByRole: requestedByRole,
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
		ActedByRole: requestedByRole,
		Comment:     reason,
		Metadata: map[string]any{
			"source":     "memberstudents.update",
			"student_id": student.ID.String(),
		},
		CreatedAt: now,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return approval, nil
}
