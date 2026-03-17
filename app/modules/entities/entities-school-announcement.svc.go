package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.SchoolAnnouncementEntity = (*Service)(nil)

func (s *Service) CreateSchoolAnnouncement(ctx context.Context, data *ent.SchoolAnnouncement) (*ent.SchoolAnnouncement, error) {
	item := data
	if _, err := s.db.NewInsert().Model(item).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) GetSchoolAnnouncementByID(ctx context.Context, id uuid.UUID) (*ent.SchoolAnnouncement, error) {
	item := new(ent.SchoolAnnouncement)
	if err := s.db.NewSelect().Model(item).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) ListSchoolAnnouncements(ctx context.Context, req *base.RequestPaginate, schoolID *uuid.UUID, status *ent.SchoolAnnouncementStatus, targetRole *string, isPinned *bool) ([]*ent.SchoolAnnouncement, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.SchoolAnnouncement, 0)
	query := s.db.NewSelect().Model(&items)

	if schoolID != nil {
		query.Where("school_id = ?", *schoolID)
	}
	if status != nil {
		query.Where("status = ?", *status)
	}
	if targetRole != nil {
		query.Where("target_role = ?", *targetRole)
	}
	if isPinned != nil {
		query.Where("is_pinned = ?", *isPinned)
	}

	if err := req.SetSearchBy(query, []string{"title", "content", "category", "created_by_name", "target_role", "status"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}

	if err := req.SetSortOrder(query, []string{"created_at", "announced_at", "published_at", "expires_at", "title", "status"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)

	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{
		Page:  req.GetPage(),
		Size:  req.GetSize(),
		Total: int64(total),
	}, nil
}

func (s *Service) UpdateSchoolAnnouncementByID(ctx context.Context, id uuid.UUID, data *ent.SchoolAnnouncementUpdate) (*ent.SchoolAnnouncement, error) {
	query := s.db.NewUpdate().Model(&ent.SchoolAnnouncement{}).Where("id = ?", id).Set("updated_at = now()")

	if data.SchoolID != nil {
		query.Set("school_id = ?", *data.SchoolID)
	}
	if data.AuthorMemberID != nil {
		query.Set("author_member_id = ?", *data.AuthorMemberID)
	}
	if data.Title != nil {
		query.Set("title = ?", *data.Title)
	}
	if data.Content != nil {
		query.Set("content = ?", *data.Content)
	}
	if data.Category != nil {
		query.Set("category = ?", *data.Category)
	}
	if data.Status != nil {
		query.Set("status = ?", *data.Status)
	}
	if data.ClearAnnouncedAt {
		query.Set("announced_at = NULL")
	} else if data.AnnouncedAt != nil {
		query.Set("announced_at = ?", *data.AnnouncedAt)
	}
	if data.ClearPublishedAt {
		query.Set("published_at = NULL")
	} else if data.PublishedAt != nil {
		query.Set("published_at = ?", *data.PublishedAt)
	}
	if data.ClearExpiresAt {
		query.Set("expires_at = NULL")
	} else if data.ExpiresAt != nil {
		query.Set("expires_at = ?", *data.ExpiresAt)
	}
	if data.CreatedByName != nil {
		query.Set("created_by_name = ?", *data.CreatedByName)
	}
	if data.ClearTargetRole {
		query.Set("target_role = NULL")
	} else if data.TargetRole != nil {
		query.Set("target_role = ?", *data.TargetRole)
	}
	if data.IsPinned != nil {
		query.Set("is_pinned = ?", *data.IsPinned)
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

	return s.GetSchoolAnnouncementByID(ctx, id)
}

func (s *Service) SoftDeleteSchoolAnnouncementByID(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.SchoolAnnouncement{}).Where("id = ?", id).Exec(ctx)
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
