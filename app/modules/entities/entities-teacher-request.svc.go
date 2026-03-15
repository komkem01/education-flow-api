package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.TeacherRequestEntity = (*Service)(nil)

func (s *Service) CreateTeacherRequest(ctx context.Context, data *ent.TeacherRequest) (*ent.TeacherRequest, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetTeacherRequestByID(ctx context.Context, id uuid.UUID) (*ent.TeacherRequest, error) {
	row := new(ent.TeacherRequest)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListTeacherRequests(ctx context.Context, req *base.RequestPaginate, teacherID *uuid.UUID, requestType *ent.TeacherRequestType, status *ent.TeacherRequestStatus) ([]*ent.TeacherRequest, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.TeacherRequest, 0)
	query := s.db.NewSelect().Model(&items)

	if teacherID != nil {
		query.Where("teacher_id = ?", *teacherID)
	}
	if requestType != nil {
		query.Where("request_type = ?", *requestType)
	}
	if status != nil {
		query.Where("status = ?", *status)
	}

	if err := req.SetSearchBy(query, []string{"request_reason"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "request_type", "status", "approved_at"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateTeacherRequestByID(ctx context.Context, id uuid.UUID, data *ent.TeacherRequestUpdate) (*ent.TeacherRequest, error) {
	query := s.db.NewUpdate().
		Model(&ent.TeacherRequest{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.TeacherID != nil {
		query.Set("teacher_id = ?", *data.TeacherID)
	}
	if data.RequestType != nil {
		query.Set("request_type = ?", *data.RequestType)
	}
	if data.RequestData != nil {
		query.Set("request_data = ?", *data.RequestData)
	}
	if data.RequestReason != nil {
		query.Set("request_reason = ?", *data.RequestReason)
	}
	if data.Status != nil {
		query.Set("status = ?", *data.Status)
	}
	if data.ApprovedBy != nil {
		query.Set("approved_by = ?", *data.ApprovedBy)
	}
	if data.ApprovedAt != nil {
		query.Set("approved_at = ?", *data.ApprovedAt)
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

	return s.GetTeacherRequestByID(ctx, id)
}

func (s *Service) DeleteTeacherRequestByID(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.TeacherRequest{}).Where("id = ?", id).Exec(ctx)
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
