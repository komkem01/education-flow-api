package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.SchoolEntity = (*Service)(nil)

func (s *Service) CreateSchool(ctx context.Context, name string, logoURL *string, themeColor *string, address *string, description *string) (*ent.School, error) {
	school := &ent.School{
		Name:        name,
		LogoURL:     logoURL,
		ThemeColor:  themeColor,
		Address:     address,
		Description: description,
	}

	if _, err := s.db.NewInsert().Model(school).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return school, nil
}

func (s *Service) GetSchoolByID(ctx context.Context, id uuid.UUID) (*ent.School, error) {
	school := new(ent.School)
	if err := s.db.NewSelect().Model(school).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return school, nil
}

func (s *Service) ListSchools(ctx context.Context, req *base.RequestPaginate) ([]*ent.School, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	schools := make([]*ent.School, 0)
	query := s.db.NewSelect().Model(&schools)

	if err := req.SetSearchBy(query, []string{"name", "theme_color", "address", "description"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}

	if err := req.SetSortOrder(query, []string{"created_at", "name", "theme_color"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)

	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return schools, &base.ResponsePaginate{
		Page:  req.GetPage(),
		Size:  req.GetSize(),
		Total: int64(total),
	}, nil
}

func (s *Service) UpdateSchoolByID(ctx context.Context, id uuid.UUID, name *string, logoURL *string, themeColor *string, address *string, description *string) (*ent.School, error) {
	query := s.db.NewUpdate().
		Model(&ent.School{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if name != nil {
		query.Set("name = ?", *name)
	}
	if logoURL != nil {
		query.Set("logo_url = ?", *logoURL)
	}
	if themeColor != nil {
		query.Set("theme_color = ?", *themeColor)
	}
	if address != nil {
		query.Set("address = ?", *address)
	}
	if description != nil {
		query.Set("description = ?", *description)
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

	return s.GetSchoolByID(ctx, id)
}

func (s *Service) SoftDeleteSchoolByID(ctx context.Context, id uuid.UUID) error {
	school := &ent.School{}
	res, err := s.db.NewDelete().Model(school).Where("id = ?", id).Exec(ctx)
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
