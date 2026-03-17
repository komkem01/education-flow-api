package approvals

import (
	"context"
	"fmt"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) applyApprovalEffects(ctx context.Context, req *ent.ApprovalRequest) error {
	if req == nil {
		return nil
	}

	switch req.RequestType {
	case "student_registration":
		return s.activateStudentRegistration(ctx, req)
	case "student_update":
		return s.applyStudentUpdate(ctx, req)
	case "management_registration":
		return s.activateManagementRegistration(ctx, req)
	case "management_update":
		return s.applyManagementUpdate(ctx, req)
	case "teacher_registration":
		return s.activateTeacherRegistration(ctx, req)
	case "teacher_update":
		return s.applyTeacherUpdate(ctx, req)
	default:
		return nil
	}
}

func payloadString(payload map[string]any, key string) (string, bool) {
	if payload == nil {
		return "", false
	}
	raw, ok := payload[key]
	if !ok {
		return "", false
	}
	v, ok := raw.(string)
	if !ok {
		return "", false
	}
	v = strings.TrimSpace(v)
	if v == "" {
		return "", false
	}
	return v, true
}

func payloadBool(payload map[string]any, key string) (bool, bool) {
	if payload == nil {
		return false, false
	}
	raw, ok := payload[key]
	if !ok {
		return false, false
	}
	v, ok := raw.(bool)
	if !ok {
		return false, false
	}
	return v, true
}

func (s *Service) applyStudentUpdate(ctx context.Context, req *ent.ApprovalRequest) error {
	studentID, err := payloadUUID(req.Payload, "student_id")
	if err != nil {
		return err
	}

	payload := &ent.MemberStudentUpdate{}
	if v, ok := payloadString(req.Payload, "member_id"); ok {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.MemberID = &parsed
	}
	if v, ok := payloadString(req.Payload, "school_id"); ok {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.SchoolID = &parsed
	}
	if v, ok := payloadString(req.Payload, "gender_id"); ok {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.GenderID = &parsed
	}
	if v, ok := payloadString(req.Payload, "prefix_id"); ok {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.PrefixID = &parsed
	}
	if v, ok := payloadString(req.Payload, "advisor_teacher_id"); ok {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.AdvisorTeacherID = &parsed
	}
	if v, ok := payloadString(req.Payload, "student_code"); ok {
		payload.StudentCode = &v
	}
	if v, ok := payloadString(req.Payload, "first_name_th"); ok {
		payload.FirstNameTH = &v
	}
	if v, ok := payloadString(req.Payload, "last_name_th"); ok {
		payload.LastNameTH = &v
	}
	if v, ok := payloadString(req.Payload, "first_name_en"); ok {
		payload.FirstNameEN = &v
	}
	if v, ok := payloadString(req.Payload, "last_name_en"); ok {
		payload.LastNameEN = &v
	}
	if v, ok := payloadString(req.Payload, "citizen_id"); ok {
		payload.CitizenID = &v
	}
	if v, ok := payloadString(req.Payload, "phone"); ok {
		payload.Phone = &v
	}
	if v, ok := payloadBool(req.Payload, "is_active"); ok {
		payload.IsActive = &v
		memberID, memberErr := payloadUUID(req.Payload, "member_id")
		if memberErr == nil {
			if _, err := s.memberDB.UpdateMemberByID(ctx, memberID, nil, nil, nil, nil, &v, nil); err != nil {
				return normalizeServiceError(err)
			}
		}
	}

	if _, err := s.studentDB.UpdateMemberStudentByID(ctx, studentID, payload); err != nil {
		return normalizeServiceError(err)
	}

	profileRows, _, err := s.profileDB.ListStudentProfiles(ctx, &base.RequestPaginate{Page: 1, Size: 1}, &studentID)
	if err != nil {
		return normalizeServiceError(err)
	}
	var profile *ent.StudentProfile
	if len(profileRows) > 0 {
		profile = profileRows[0]
	}

	profileUpdate := &ent.StudentProfileUpdate{}
	hasProfileUpdate := false
	if v, ok := payloadString(req.Payload, "birth_date"); ok {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		profileUpdate.BirthDate = &t
		hasProfileUpdate = true
	}
	if v, ok := payloadString(req.Payload, "nationality"); ok {
		profileUpdate.Nationality = &v
		hasProfileUpdate = true
	}
	if v, ok := payloadString(req.Payload, "religion"); ok {
		profileUpdate.Religion = &v
		hasProfileUpdate = true
	}
	if v, ok := payloadString(req.Payload, "address_current"); ok {
		profileUpdate.AddressCurrent = &v
		hasProfileUpdate = true
	}
	if v, ok := payloadString(req.Payload, "address_registered"); ok {
		profileUpdate.AddressRegistered = &v
		hasProfileUpdate = true
	}
	if v, ok := payloadString(req.Payload, "emergency_contact_name"); ok {
		profileUpdate.EmergencyContactName = &v
		hasProfileUpdate = true
	}
	if v, ok := payloadString(req.Payload, "emergency_contact_phone"); ok {
		profileUpdate.EmergencyContactPhone = &v
		hasProfileUpdate = true
	}

	if hasProfileUpdate {
		if profile != nil {
			if _, err := s.profileDB.UpdateStudentProfileByID(ctx, profile.ID, profileUpdate); err != nil {
				return normalizeServiceError(err)
			}
		} else {
			birthDate := profileUpdate.BirthDate
			if _, err := s.profileDB.CreateStudentProfile(ctx, &ent.StudentProfile{
				StudentID:             studentID,
				BirthDate:             birthDate,
				Nationality:           profileUpdate.Nationality,
				Religion:              profileUpdate.Religion,
				AddressCurrent:        profileUpdate.AddressCurrent,
				AddressRegistered:     profileUpdate.AddressRegistered,
				EmergencyContactName:  profileUpdate.EmergencyContactName,
				EmergencyContactPhone: profileUpdate.EmergencyContactPhone,
			}); err != nil {
				return normalizeServiceError(err)
			}
		}
	}

	healthRows, _, err := s.healthDB.ListStudentHealthProfiles(ctx, &base.RequestPaginate{Page: 1, Size: 1}, &studentID, nil)
	if err != nil {
		return normalizeServiceError(err)
	}
	var health *ent.StudentHealthProfile
	if len(healthRows) > 0 {
		health = healthRows[0]
	}

	healthUpdate := &ent.StudentHealthProfileUpdate{}
	hasHealthUpdate := false
	if v, ok := payloadString(req.Payload, "blood_type"); ok {
		healthUpdate.BloodType = &v
		hasHealthUpdate = true
	}
	if v, ok := payloadString(req.Payload, "allergy_info"); ok {
		healthUpdate.AllergyInfo = &v
		hasHealthUpdate = true
	}
	if v, ok := payloadString(req.Payload, "chronic_disease"); ok {
		healthUpdate.ChronicDisease = &v
		hasHealthUpdate = true
	}
	if v, ok := payloadString(req.Payload, "medical_note"); ok {
		healthUpdate.MedicalNote = &v
		hasHealthUpdate = true
	}

	if hasHealthUpdate {
		if health != nil {
			if _, err := s.healthDB.UpdateStudentHealthProfileByID(ctx, health.ID, healthUpdate); err != nil {
				return normalizeServiceError(err)
			}
		} else {
			if _, err := s.healthDB.CreateStudentHealthProfile(ctx, &ent.StudentHealthProfile{
				StudentID:      studentID,
				BloodType:      healthUpdate.BloodType,
				AllergyInfo:    healthUpdate.AllergyInfo,
				ChronicDisease: healthUpdate.ChronicDisease,
				MedicalNote:    healthUpdate.MedicalNote,
			}); err != nil {
				return normalizeServiceError(err)
			}
		}
	}

	return nil
}

func (s *Service) applyManagementUpdate(ctx context.Context, req *ent.ApprovalRequest) error {
	managementID, err := payloadUUID(req.Payload, "management_id")
	if err != nil {
		return err
	}

	payload := &ent.MemberManagementUpdate{}
	if v, ok := payloadString(req.Payload, "member_id"); ok {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.MemberID = &parsed
	}
	if v, ok := payloadString(req.Payload, "employee_code"); ok {
		payload.EmployeeCode = &v
	}
	if v, ok := payloadString(req.Payload, "gender_id"); ok {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.GenderID = &parsed
	}
	if v, ok := payloadString(req.Payload, "prefix_id"); ok {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.PrefixID = &parsed
	}
	if v, ok := payloadString(req.Payload, "first_name"); ok {
		payload.FirstName = &v
	}
	if v, ok := payloadString(req.Payload, "last_name"); ok {
		payload.LastName = &v
	}
	if v, ok := payloadString(req.Payload, "phone"); ok {
		payload.Phone = &v
	}
	if v, ok := payloadString(req.Payload, "position"); ok {
		payload.Position = &v
	}
	if v, ok := payloadString(req.Payload, "start_work_date"); ok {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.StartWorkDate = &t
	}
	if v, ok := payloadString(req.Payload, "school_department_id"); ok {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.SchoolDepartmentID = &parsed
	}
	if v, ok := payloadString(req.Payload, "department_id"); ok {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.DepartmentID = &parsed
	}
	if v, ok := payloadBool(req.Payload, "is_active"); ok {
		payload.IsActive = &v
		memberID, memberErr := payloadUUID(req.Payload, "member_id")
		if memberErr == nil {
			if _, err := s.memberDB.UpdateMemberByID(ctx, memberID, nil, nil, nil, nil, &v, nil); err != nil {
				return normalizeServiceError(err)
			}
		}
	}

	if _, err := s.managementDB.UpdateMemberManagementByID(ctx, managementID, payload); err != nil {
		return normalizeServiceError(err)
	}

	return nil
}

func (s *Service) activateStudentRegistration(ctx context.Context, req *ent.ApprovalRequest) error {
	memberID, err := payloadUUID(req.Payload, "member_id")
	if err != nil {
		return err
	}
	studentID, err := payloadUUID(req.Payload, "student_id")
	if err != nil {
		return err
	}

	active := true
	if _, err := s.memberDB.UpdateMemberByID(ctx, memberID, nil, nil, nil, nil, &active, nil); err != nil {
		return normalizeServiceError(err)
	}
	if _, err := s.studentDB.UpdateMemberStudentByID(ctx, studentID, &ent.MemberStudentUpdate{IsActive: &active}); err != nil {
		return normalizeServiceError(err)
	}

	return nil
}

func (s *Service) activateManagementRegistration(ctx context.Context, req *ent.ApprovalRequest) error {
	memberID, err := payloadUUID(req.Payload, "member_id")
	if err != nil {
		return err
	}
	managementID, err := payloadUUID(req.Payload, "management_id")
	if err != nil {
		return err
	}

	active := true
	if _, err := s.memberDB.UpdateMemberByID(ctx, memberID, nil, nil, nil, nil, &active, nil); err != nil {
		return normalizeServiceError(err)
	}
	if _, err := s.managementDB.UpdateMemberManagementByID(ctx, managementID, &ent.MemberManagementUpdate{IsActive: &active}); err != nil {
		return normalizeServiceError(err)
	}

	return nil
}

func (s *Service) activateTeacherRegistration(ctx context.Context, req *ent.ApprovalRequest) error {
	memberID, err := payloadUUID(req.Payload, "member_id")
	if err != nil {
		return err
	}
	teacherID, err := payloadUUID(req.Payload, "teacher_id")
	if err != nil {
		return err
	}

	active := true
	if _, err := s.memberDB.UpdateMemberByID(ctx, memberID, nil, nil, nil, nil, &active, nil); err != nil {
		return normalizeServiceError(err)
	}
	if _, err := s.teacherDB.UpdateMemberTeacherByID(ctx, teacherID, &ent.MemberTeacherUpdate{IsActive: &active}); err != nil {
		return normalizeServiceError(err)
	}

	return nil
}

func (s *Service) applyTeacherUpdate(ctx context.Context, req *ent.ApprovalRequest) error {
	teacherID, err := payloadUUID(req.Payload, "teacher_id")
	if err != nil {
		return err
	}

	payload := &ent.MemberTeacherUpdate{}
	if v, ok := payloadString(req.Payload, "member_id"); ok {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.MemberID = &parsed
	}
	if v, ok := payloadString(req.Payload, "gender_id"); ok {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.GenderID = &parsed
	}
	if v, ok := payloadString(req.Payload, "prefix_id"); ok {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.PrefixID = &parsed
	}
	if v, ok := payloadString(req.Payload, "code"); ok {
		payload.Code = &v
	}
	if v, ok := payloadString(req.Payload, "citizen_id"); ok {
		payload.CitizenID = &v
	}
	if v, ok := payloadString(req.Payload, "first_name_th"); ok {
		payload.FirstNameTH = &v
	}
	if v, ok := payloadString(req.Payload, "last_name_th"); ok {
		payload.LastNameTH = &v
	}
	if v, ok := payloadString(req.Payload, "first_name_en"); ok {
		payload.FirstNameEN = &v
	}
	if v, ok := payloadString(req.Payload, "last_name_en"); ok {
		payload.LastNameEN = &v
	}
	if v, ok := payloadString(req.Payload, "phone"); ok {
		payload.Phone = &v
	}
	if v, ok := payloadString(req.Payload, "position"); ok {
		payload.Position = &v
	}
	if v, ok := payloadString(req.Payload, "academic_standing"); ok {
		payload.AcademicStanding = &v
	}
	if v, ok := payloadString(req.Payload, "department_id"); ok {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.DepartmentID = &parsed
	}
	if v, ok := payloadString(req.Payload, "start_date"); ok {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.StartDate = &t
	}
	if v, ok := payloadString(req.Payload, "end_date"); ok {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		payload.EndDate = &t
	}
	if v, ok := payloadBool(req.Payload, "is_active"); ok {
		payload.IsActive = &v
		memberID, memberErr := payloadUUID(req.Payload, "member_id")
		if memberErr == nil {
			if _, err := s.memberDB.UpdateMemberByID(ctx, memberID, nil, nil, nil, nil, &v, nil); err != nil {
				return normalizeServiceError(err)
			}
		}
	}

	if _, err := s.teacherDB.UpdateMemberTeacherByID(ctx, teacherID, payload); err != nil {
		return normalizeServiceError(err)
	}

	rawAddresses, ok := req.Payload["addresses"]
	if !ok {
		return nil
	}

	rows, ok := rawAddresses.([]any)
	if !ok {
		return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}

	addresses := make([]ent.TeacherAddressInput, 0, len(rows))
	for i, raw := range rows {
		item, ok := raw.(map[string]any)
		if !ok {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}

		houseNo, ok := payloadString(item, "house_no")
		if !ok {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		province, ok := payloadString(item, "province")
		if !ok {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		district, ok := payloadString(item, "district")
		if !ok {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		subdistrict, ok := payloadString(item, "subdistrict")
		if !ok {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		postalCode, ok := payloadString(item, "postal_code")
		if !ok {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		isPrimary, ok := payloadBool(item, "is_primary")
		if !ok {
			return fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}

		row := ent.TeacherAddressInput{
			HouseNo:     houseNo,
			Province:    province,
			District:    district,
			Subdistrict: subdistrict,
			PostalCode:  postalCode,
			IsPrimary:   isPrimary,
			SortOrder:   i + 1,
		}
		if v, ok := payloadString(item, "village"); ok {
			row.Village = &v
		}
		if v, ok := payloadString(item, "road"); ok {
			row.Road = &v
		}
		addresses = append(addresses, row)
	}

	if err := s.teacherDB.ReplaceTeacherAddressesByMemberTeacherID(ctx, teacherID, addresses); err != nil {
		return normalizeServiceError(err)
	}

	return nil
}

func payloadUUID(payload map[string]any, key string) (uuid.UUID, error) {
	if payload == nil {
		return uuid.Nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}
	raw, ok := payload[key]
	if !ok {
		return uuid.Nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}
	v, ok := raw.(string)
	if !ok || v == "" {
		return uuid.Nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}
	id, err := uuid.Parse(v)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}
	return id, nil
}
