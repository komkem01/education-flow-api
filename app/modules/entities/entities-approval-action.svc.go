package entities

import (
	"context"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.ApprovalActionEntity = (*Service)(nil)

func (s *Service) CreateApprovalAction(ctx context.Context, data *ent.ApprovalAction) (*ent.ApprovalAction, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetApprovalActionByIdempotencyKey(ctx context.Context, requestID uuid.UUID, action ent.ApprovalActionType, idempotencyKey string) (*ent.ApprovalAction, error) {
	row := new(ent.ApprovalAction)
	if err := s.db.NewSelect().
		Model(row).
		Where("request_id = ?", requestID).
		Where("action = ?", action).
		Where("idempotency_key = ?", idempotencyKey).
		Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListApprovalActions(ctx context.Context, req *base.RequestPaginate, requestID *uuid.UUID, actedBy *uuid.UUID, action *ent.ApprovalActionType) ([]*ent.ApprovalAction, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.ApprovalAction, 0)
	query := s.db.NewSelect().Model(&items)

	if requestID != nil {
		query.Where("request_id = ?", *requestID)
	}
	if actedBy != nil {
		query.Where("acted_by = ?", *actedBy)
	}
	if action != nil {
		query.Where("action = ?", *action)
	}

	if err := req.SetSearchBy(query, []string{"action", "comment"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "action"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}
