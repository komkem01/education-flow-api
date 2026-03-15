package entities

import (
	"context"
	"time"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.AuditLogEntity = (*Service)(nil)

func (s *Service) CreateAuditLog(ctx context.Context, data *ent.AuditLog) (*ent.AuditLog, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetAuditLogByID(ctx context.Context, id uuid.UUID) (*ent.AuditLog, error) {
	row := new(ent.AuditLog)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListAuditLogs(ctx context.Context, req *base.RequestPaginate, actorID *uuid.UUID, actorRole *string, method *string, path *string, statusCode *int, from *time.Time, to *time.Time) ([]*ent.AuditLog, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.AuditLog, 0)
	query := s.db.NewSelect().Model(&items)

	if actorID != nil {
		query.Where("actor_id = ?", *actorID)
	}
	if actorRole != nil {
		query.Where("actor_role = ?", *actorRole)
	}
	if method != nil {
		query.Where("method = ?", *method)
	}
	if path != nil {
		query.Where("path ILIKE ?", "%"+*path+"%")
	}
	if statusCode != nil {
		query.Where("status_code = ?", *statusCode)
	}
	if from != nil {
		query.Where("created_at >= ?", *from)
	}
	if to != nil {
		query.Where("created_at <= ?", *to)
	}

	if err := req.SetSearchBy(query, []string{"method", "path", "route_path", "actor_role", "trace_id", "query_string", "error_message"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "status_code", "latency_ms", "method", "path"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) PurgeAuditLogsBefore(ctx context.Context, before time.Time) error {
	_, err := s.db.NewDelete().
		Model(&ent.AuditLog{}).
		Where("created_at < ?", before).
		Exec(ctx)
	return err
}
