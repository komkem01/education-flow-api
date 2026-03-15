package auditlogs

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.AuditLog, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "auditlogs.service.get_by_id")
	defer span.End()

	item, err := s.db.GetAuditLogByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}
