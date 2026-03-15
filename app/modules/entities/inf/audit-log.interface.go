package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
)

type AuditLogEntity interface {
	CreateAuditLog(ctx context.Context, data *ent.AuditLog) (*ent.AuditLog, error)
}
