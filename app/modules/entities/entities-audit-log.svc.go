package entities

import (
	"context"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
)

var _ entitiesinf.AuditLogEntity = (*Service)(nil)

func (s *Service) CreateAuditLog(ctx context.Context, data *ent.AuditLog) (*ent.AuditLog, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}
