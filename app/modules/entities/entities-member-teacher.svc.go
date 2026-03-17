package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ entitiesinf.MemberTeacherEntity = (*Service)(nil)

func (s *Service) CreateMemberTeacher(ctx context.Context, teacher *ent.MemberTeacher) (*ent.MemberTeacher, error) {
	if _, err := s.db.NewInsert().Model(teacher).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return teacher, nil
}

func (s *Service) RegisterTeacher(ctx context.Context, data *ent.TeacherRegistrationInput) (*ent.TeacherRegistrationResult, error) {
	result := new(ent.TeacherRegistrationResult)

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

		teacher := &ent.MemberTeacher{
			MemberID:         member.ID,
			GenderID:         data.TeacherGenderID,
			PrefixID:         data.TeacherPrefixID,
			Code:             data.TeacherCode,
			CitizenID:        data.TeacherCitizenID,
			FirstNameTH:      data.TeacherFirstNameTH,
			LastNameTH:       data.TeacherLastNameTH,
			FirstNameEN:      data.TeacherFirstNameEN,
			LastNameEN:       data.TeacherLastNameEN,
			Phone:            data.TeacherPhone,
			Position:         data.TeacherPosition,
			AcademicStanding: data.TeacherAcademicStanding,
			DepartmentID:     data.TeacherDepartmentID,
			StartDate:        data.TeacherStartDate,
			EndDate:          data.TeacherEndDate,
			IsActive:         data.TeacherIsActive,
		}
		if _, err := tx.NewInsert().Model(teacher).Returning("*").Exec(ctx); err != nil {
			return err
		}

		if len(data.TeacherAddresses) > 0 {
			addresses := make([]*ent.TeacherAddress, 0, len(data.TeacherAddresses))
			for _, addr := range data.TeacherAddresses {
				addresses = append(addresses, &ent.TeacherAddress{
					MemberTeacherID: teacher.ID,
					HouseNo:         addr.HouseNo,
					Village:         addr.Village,
					Road:            addr.Road,
					Province:        addr.Province,
					District:        addr.District,
					Subdistrict:     addr.Subdistrict,
					PostalCode:      addr.PostalCode,
					IsPrimary:       addr.IsPrimary,
					SortOrder:       addr.SortOrder,
					IsActive:        true,
				})
			}

			if _, err := tx.NewInsert().Model(&addresses).Exec(ctx); err != nil {
				return err
			}
		}

		result.Member = member
		result.Teacher = teacher
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) GetMemberTeacherByID(ctx context.Context, id uuid.UUID) (*ent.MemberTeacher, error) {
	teacher := new(ent.MemberTeacher)
	if err := s.db.NewSelect().Model(teacher).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return teacher, nil
}

func (s *Service) ListTeacherAddressesByMemberTeacherID(ctx context.Context, memberTeacherID uuid.UUID) ([]*ent.TeacherAddress, error) {
	items := make([]*ent.TeacherAddress, 0)
	if err := s.db.NewSelect().
		Model(&items).
		Where("member_teacher_id = ?", memberTeacherID).
		Order("is_primary DESC").
		Order("sort_order ASC").
		Order("created_at ASC").
		Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Service) ReplaceTeacherAddressesByMemberTeacherID(ctx context.Context, memberTeacherID uuid.UUID, addresses []ent.TeacherAddressInput) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().
			Model(&ent.TeacherAddress{}).
			Where("member_teacher_id = ?", memberTeacherID).
			Exec(ctx); err != nil {
			return err
		}

		if len(addresses) == 0 {
			return nil
		}

		items := make([]*ent.TeacherAddress, 0, len(addresses))
		for _, addr := range addresses {
			items = append(items, &ent.TeacherAddress{
				MemberTeacherID: memberTeacherID,
				HouseNo:         addr.HouseNo,
				Village:         addr.Village,
				Road:            addr.Road,
				Province:        addr.Province,
				District:        addr.District,
				Subdistrict:     addr.Subdistrict,
				PostalCode:      addr.PostalCode,
				IsPrimary:       addr.IsPrimary,
				SortOrder:       addr.SortOrder,
				IsActive:        true,
			})
		}

		if _, err := tx.NewInsert().Model(&items).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) ListMemberTeachers(ctx context.Context, req *base.RequestPaginate, isActive *bool, memberID *uuid.UUID, departmentID *uuid.UUID) ([]*ent.MemberTeacher, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	teachers := make([]*ent.MemberTeacher, 0)
	query := s.db.NewSelect().Model(&teachers)

	// Hide not-yet-approved teacher registrations from the teacher directory.
	query.Where(`not exists (
		select 1
		from approval_requests ar
		where ar.subject_type = 'member_teacher'
		  and ar.request_type = 'teacher_registration'
		  and ar.current_status = 'pending'
		  and ar.subject_id = mt.id
	)`)

	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}
	if memberID != nil {
		query.Where("member_id = ?", *memberID)
	}
	if departmentID != nil {
		query.Where("department = ?", *departmentID)
	}

	if err := req.SetSearchBy(query, []string{"code", "citizen_id", "first_name_th", "last_name_th", "first_name_en", "last_name_en", "phone", "position", "academic_standing"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}

	if err := req.SetSortOrder(query, []string{"created_at", "code", "start_date", "first_name_th", "last_name_th", "is_active"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)

	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return teachers, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateMemberTeacherByID(ctx context.Context, id uuid.UUID, req *ent.MemberTeacherUpdate) (*ent.MemberTeacher, error) {
	query := s.db.NewUpdate().
		Model(&ent.MemberTeacher{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if req.MemberID != nil {
		query.Set("member_id = ?", *req.MemberID)
	}
	if req.GenderID != nil {
		query.Set("gender_id = ?", *req.GenderID)
	}
	if req.PrefixID != nil {
		query.Set("prefix_id = ?", *req.PrefixID)
	}
	if req.Code != nil {
		query.Set("code = ?", *req.Code)
	}
	if req.CitizenID != nil {
		query.Set("citizen_id = ?", *req.CitizenID)
	}
	if req.FirstNameTH != nil {
		query.Set("first_name_th = ?", *req.FirstNameTH)
	}
	if req.LastNameTH != nil {
		query.Set("last_name_th = ?", *req.LastNameTH)
	}
	if req.FirstNameEN != nil {
		query.Set("first_name_en = ?", *req.FirstNameEN)
	}
	if req.LastNameEN != nil {
		query.Set("last_name_en = ?", *req.LastNameEN)
	}
	if req.Phone != nil {
		query.Set("phone = ?", *req.Phone)
	}
	if req.Position != nil {
		query.Set("position = ?", *req.Position)
	}
	if req.AcademicStanding != nil {
		query.Set("academic_standing = ?", *req.AcademicStanding)
	}
	if req.DepartmentID != nil {
		query.Set("department = ?", *req.DepartmentID)
	}
	if req.StartDate != nil {
		query.Set("start_date = ?", *req.StartDate)
	}
	if req.EndDate != nil {
		query.Set("end_date = ?", *req.EndDate)
	}
	if req.IsActive != nil {
		query.Set("is_active = ?", *req.IsActive)
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

	return s.GetMemberTeacherByID(ctx, id)
}

func (s *Service) SoftDeleteMemberTeacherByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		teacher := new(ent.MemberTeacher)
		if err := tx.NewSelect().
			Model(teacher).
			Where("id = ?", id).
			Scan(ctx); err != nil {
			return err
		}

		res, err := tx.NewUpdate().
			Model(&ent.MemberTeacher{}).
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

		if _, err = tx.NewUpdate().
			Model(&ent.MemberTeacher{}).
			Set("deleted_at = now()").
			Set("updated_at = now()").
			Where("id = ?", id).
			Exec(ctx); err != nil {
			return err
		}

		if _, err = tx.NewUpdate().
			Model(&ent.Member{}).
			Set("is_active = ?", false).
			Set("updated_at = now()").
			Where("id = ?", teacher.MemberID).
			Exec(ctx); err != nil {
			return err
		}

		_, err = tx.NewUpdate().
			Model(&ent.Member{}).
			Set("deleted_at = now()").
			Set("updated_at = now()").
			Where("id = ?", teacher.MemberID).
			Exec(ctx)
		return err
	})
}
