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

var _ entitiesinf.MemberEntity = (*Service)(nil)

func (s *Service) CreateMember(ctx context.Context, schoolID uuid.UUID, email string, password string, role ent.MemberRole, isActive bool, lastLogin *time.Time) (*ent.Member, error) {
	item := &ent.Member{
		SchoolID:  schoolID,
		Email:     email,
		Password:  password,
		Role:      role,
		IsActive:  isActive,
		LastLogin: lastLogin,
	}

	if _, err := s.db.NewInsert().Model(item).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) GetMemberByID(ctx context.Context, id uuid.UUID) (*ent.Member, error) {
	item := new(ent.Member)
	if err := s.db.NewSelect().Model(item).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) GetMemberByEmail(ctx context.Context, email string) (*ent.Member, error) {
	item := new(ent.Member)
	if err := s.db.NewSelect().Model(item).Where("email = ?", email).Scan(ctx); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) ListMembers(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, role *ent.MemberRole) ([]*ent.Member, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.Member, 0)
	query := s.db.NewSelect().Model(&items)

	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}
	if schoolID != nil {
		query.Where("school_id = ?", *schoolID)
	}
	if role != nil {
		query.Where("role = ?", *role)
	}

	if err := req.SetSearchBy(query, []string{"email", "role"}); err != nil {
		return nil, nil, err
	}
	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "email", "role", "is_active", "last_login"}); err != nil {
		return nil, nil, err
	}
	req.SetOffsetLimit(query)

	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateMemberByID(ctx context.Context, id uuid.UUID, schoolID *uuid.UUID, email *string, password *string, role *ent.MemberRole, isActive *bool, lastLogin *time.Time) (*ent.Member, error) {
	query := s.db.NewUpdate().
		Model(&ent.Member{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if schoolID != nil {
		query.Set("school_id = ?", *schoolID)
	}
	if email != nil {
		query.Set("email = ?", *email)
	}
	if password != nil {
		query.Set("password = ?", *password)
	}
	if role != nil {
		query.Set("role = ?", *role)
	}
	if isActive != nil {
		query.Set("is_active = ?", *isActive)
	}
	if lastLogin != nil {
		query.Set("last_login = ?", *lastLogin)
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

	return s.GetMemberByID(ctx, id)
}

func (s *Service) SoftDeleteMemberByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(&ent.Member{}).
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

		_, err = tx.NewDelete().Model(&ent.Member{}).Where("id = ?", id).Exec(ctx)
		return err
	})
}
