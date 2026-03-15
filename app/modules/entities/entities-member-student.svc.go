package entities

import (
	"context"
	"database/sql"
	"time"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ entitiesinf.MemberStudentEntity = (*Service)(nil)

func (s *Service) CreateMemberStudent(ctx context.Context, data *ent.MemberStudent) (*ent.MemberStudent, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) RegisterStudent(ctx context.Context, data *ent.StudentRegistrationInput) (*ent.StudentRegistrationResult, error) {
	result := new(ent.StudentRegistrationResult)

	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		member := &ent.Member{
			SchoolID:  data.MemberSchoolID,
			Email:     data.MemberEmail,
			Password:  data.MemberPasswordHash,
			Role:      data.MemberRole,
			IsActive:  data.MemberIsActive,
			LastLogin: data.MemberLastLogin,
		}
		if _, err := tx.NewInsert().Model(member).Returning("*").Exec(ctx); err != nil {
			return err
		}

		student := &ent.MemberStudent{
			MemberID:         member.ID,
			SchoolID:         data.StudentSchoolID,
			GenderID:         data.StudentGenderID,
			PrefixID:         data.StudentPrefixID,
			AdvisorTeacherID: data.StudentAdvisorTeacherID,
			StudentCode:      data.StudentCode,
			FirstNameTH:      data.StudentFirstNameTH,
			LastNameTH:       data.StudentLastNameTH,
			FirstNameEN:      data.StudentFirstNameEN,
			LastNameEN:       data.StudentLastNameEN,
			CitizenID:        data.StudentCitizenID,
			Phone:            data.StudentPhone,
			IsActive:         data.StudentIsActive,
		}
		if _, err := tx.NewInsert().Model(student).Returning("*").Exec(ctx); err != nil {
			return err
		}

		profile := &ent.StudentProfile{
			StudentID:             student.ID,
			BirthDate:             data.ProfileBirthDate,
			Nationality:           data.ProfileNationality,
			Religion:              data.ProfileReligion,
			AddressCurrent:        data.ProfileAddressCurrent,
			AddressRegistered:     data.ProfileAddressRegistered,
			EmergencyContactName:  data.ProfileEmergencyContactName,
			EmergencyContactPhone: data.ProfileEmergencyContactPhone,
		}
		if _, err := tx.NewInsert().Model(profile).Returning("*").Exec(ctx); err != nil {
			return err
		}

		result.Member = member
		result.Student = student
		result.Profile = profile

		if data.RequireApproval {
			now := time.Now()
			subjectID := student.ID
			approval := &ent.ApprovalRequest{
				RequestType:     "student_registration",
				SubjectType:     "member_student",
				SubjectID:       &subjectID,
				RequestedBy:     data.RequestedBy,
				RequestedByRole: data.RequestedByRole,
				Payload: map[string]any{
					"member_id":    member.ID.String(),
					"student_id":   student.ID.String(),
					"student_code": student.StudentCode,
					"school_id":    student.SchoolID.String(),
				},
				CurrentStatus: ent.ApprovalRequestStatusPending,
				SubmittedAt:   &now,
			}
			if data.RequestReason != nil && *data.RequestReason != "" {
				approval.Payload["reason"] = *data.RequestReason
			}
			if _, err := tx.NewInsert().Model(approval).Returning("*").Exec(ctx); err != nil {
				return err
			}

			action := &ent.ApprovalAction{
				RequestID:   approval.ID,
				Action:      ent.ApprovalActionTypeSubmit,
				ActedBy:     data.RequestedBy,
				ActedByRole: data.RequestedByRole,
				Comment:     data.RequestReason,
				Metadata: map[string]any{
					"source":     "memberstudents.register",
					"student_id": student.ID.String(),
				},
				CreatedAt: now,
			}
			if _, err := tx.NewInsert().Model(action).Returning("*").Exec(ctx); err != nil {
				return err
			}

			result.Approval = approval
			result.ApprovalAction = action
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) GetMemberStudentByID(ctx context.Context, id uuid.UUID) (*ent.MemberStudent, error) {
	row := new(ent.MemberStudent)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListMemberStudents(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, advisorTeacherID *uuid.UUID) ([]*ent.MemberStudent, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.MemberStudent, 0)
	query := s.db.NewSelect().Model(&items)

	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}
	if schoolID != nil {
		query.Where("school_id = ?", *schoolID)
	}
	if advisorTeacherID != nil {
		query.Where("advisor_teacher_id = ?", *advisorTeacherID)
	}

	if err := req.SetSearchBy(query, []string{"student_code", "first_name_th", "last_name_th", "first_name_en", "last_name_en"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "student_code", "first_name_th", "last_name_th", "is_active"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateMemberStudentByID(ctx context.Context, id uuid.UUID, data *ent.MemberStudentUpdate) (*ent.MemberStudent, error) {
	query := s.db.NewUpdate().
		Model(&ent.MemberStudent{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.MemberID != nil {
		query.Set("member_id = ?", *data.MemberID)
	}
	if data.SchoolID != nil {
		query.Set("school_id = ?", *data.SchoolID)
	}
	if data.GenderID != nil {
		query.Set("gender_id = ?", *data.GenderID)
	}
	if data.PrefixID != nil {
		query.Set("prefix_id = ?", *data.PrefixID)
	}
	if data.AdvisorTeacherID != nil {
		query.Set("advisor_teacher_id = ?", *data.AdvisorTeacherID)
	}
	if data.StudentCode != nil {
		query.Set("student_code = ?", *data.StudentCode)
	}
	if data.FirstNameTH != nil {
		query.Set("first_name_th = ?", *data.FirstNameTH)
	}
	if data.LastNameTH != nil {
		query.Set("last_name_th = ?", *data.LastNameTH)
	}
	if data.FirstNameEN != nil {
		query.Set("first_name_en = ?", *data.FirstNameEN)
	}
	if data.LastNameEN != nil {
		query.Set("last_name_en = ?", *data.LastNameEN)
	}
	if data.CitizenID != nil {
		query.Set("citizen_id = ?", *data.CitizenID)
	}
	if data.Phone != nil {
		query.Set("phone = ?", *data.Phone)
	}
	if data.IsActive != nil {
		query.Set("is_active = ?", *data.IsActive)
	}

	res, err := query.Exec(ctx)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, sql.ErrNoRows
	}

	return s.GetMemberStudentByID(ctx, id)
}

func (s *Service) SoftDeleteMemberStudentByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(&ent.MemberStudent{}).
			Set("is_active = ?", false).
			Set("updated_at = now()").
			Where("id = ?", id).
			Exec(ctx)
		if err != nil {
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return sql.ErrNoRows
		}

		_, err = tx.NewDelete().Model(&ent.MemberStudent{}).Where("id = ?", id).Exec(ctx)
		return err
	})
}
