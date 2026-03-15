package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.TeacherEmergencyContactEntity = (*Service)(nil)

func (s *Service) CreateTeacherEmergencyContact(ctx context.Context, data *ent.TeacherEmergencyContact) (*ent.TeacherEmergencyContact, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetTeacherEmergencyContactByID(ctx context.Context, id uuid.UUID) (*ent.TeacherEmergencyContact, error) {
	row := new(ent.TeacherEmergencyContact)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListTeacherEmergencyContacts(ctx context.Context, req *base.RequestPaginate, memberTeacherID *uuid.UUID, isPrimary *bool) ([]*ent.TeacherEmergencyContact, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.TeacherEmergencyContact, 0)
	query := s.db.NewSelect().Model(&items)

	if memberTeacherID != nil {
		query.Where("member_teacher_id = ?", *memberTeacherID)
	}
	if isPrimary != nil {
		query.Where("is_primary = ?", *isPrimary)
	}

	if err := req.SetSearchBy(query, []string{"emergency_contact_name", "relationship", "phone_primary", "phone_secondary"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "emergency_contact_name", "is_primary"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateTeacherEmergencyContactByID(ctx context.Context, id uuid.UUID, data *ent.TeacherEmergencyContactUpdate) (*ent.TeacherEmergencyContact, error) {
	query := s.db.NewUpdate().
		Model(&ent.TeacherEmergencyContact{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.MemberTeacherID != nil {
		query.Set("member_teacher_id = ?", *data.MemberTeacherID)
	}
	if data.EmergencyContactName != nil {
		query.Set("emergency_contact_name = ?", *data.EmergencyContactName)
	}
	if data.Relationship != nil {
		query.Set("relationship = ?", *data.Relationship)
	}
	if data.PhonePrimary != nil {
		query.Set("phone_primary = ?", *data.PhonePrimary)
	}
	if data.PhoneSecondary != nil {
		query.Set("phone_secondary = ?", *data.PhoneSecondary)
	}
	if data.CanDecideMedical != nil {
		query.Set("can_decide_medical = ?", *data.CanDecideMedical)
	}
	if data.IsPrimary != nil {
		query.Set("is_primary = ?", *data.IsPrimary)
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

	return s.GetTeacherEmergencyContactByID(ctx, id)
}

func (s *Service) SoftDeleteTeacherEmergencyContactByID(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.TeacherEmergencyContact{}).Where("id = ?", id).Exec(ctx)
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
	return nil
}
