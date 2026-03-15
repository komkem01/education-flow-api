package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.ApprovalRequestEntity = (*Service)(nil)

func (s *Service) CreateApprovalRequest(ctx context.Context, data *ent.ApprovalRequest) (*ent.ApprovalRequest, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetApprovalRequestByID(ctx context.Context, id uuid.UUID) (*ent.ApprovalRequest, error) {
	row := new(ent.ApprovalRequest)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListApprovalRequests(ctx context.Context, req *base.RequestPaginate, requestedBy *uuid.UUID, requestedByRole *ent.ApprovalActorRole, status *ent.ApprovalRequestStatus, requestType *string) ([]*ent.ApprovalRequest, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.ApprovalRequest, 0)
	query := s.db.NewSelect().Model(&items)

	if requestedBy != nil {
		query.Where("requested_by = ?", *requestedBy)
	}
	if requestedByRole != nil {
		query.Where("requested_by_role = ?", *requestedByRole)
	}
	if status != nil {
		query.Where("current_status = ?", *status)
	}
	if requestType != nil {
		query.Where("request_type = ?", *requestType)
	}

	if err := req.SetSearchBy(query, []string{"request_type", "subject_type"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "request_type", "current_status", "submitted_at"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateApprovalRequestByID(ctx context.Context, id uuid.UUID, data *ent.ApprovalRequestUpdate) (*ent.ApprovalRequest, error) {
	query := s.db.NewUpdate().
		Model(&ent.ApprovalRequest{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.RequestType != nil {
		query.Set("request_type = ?", *data.RequestType)
	}
	if data.SubjectType != nil {
		query.Set("subject_type = ?", *data.SubjectType)
	}
	if data.SubjectID != nil {
		query.Set("subject_id = ?", *data.SubjectID)
	}
	if data.RequestedBy != nil {
		query.Set("requested_by = ?", *data.RequestedBy)
	}
	if data.RequestedByRole != nil {
		query.Set("requested_by_role = ?", *data.RequestedByRole)
	}
	if data.Payload != nil {
		query.Set("payload = ?", *data.Payload)
	}
	if data.CurrentStatus != nil {
		query.Set("current_status = ?", *data.CurrentStatus)
	}
	if data.SubmittedAt != nil {
		query.Set("submitted_at = ?", *data.SubmittedAt)
	}
	if data.ResolvedAt != nil {
		query.Set("resolved_at = ?", *data.ResolvedAt)
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

	return s.GetApprovalRequestByID(ctx, id)
}

func (s *Service) DeleteApprovalRequestByID(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.ApprovalRequest{}).Where("id = ?", id).Exec(ctx)
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
