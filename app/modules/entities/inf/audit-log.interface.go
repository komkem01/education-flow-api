package entitiesinf

import (
	"context"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type AuditLogEntity interface {
	CreateAuditLog(ctx context.Context, data *ent.AuditLog) (*ent.AuditLog, error)
	GetAuditLogByID(ctx context.Context, id uuid.UUID) (*ent.AuditLog, error)
	ListAuditLogs(ctx context.Context, req *base.RequestPaginate, actorID *uuid.UUID, actorRole *string, method *string, path *string, statusCode *int, from *time.Time, to *time.Time) ([]*ent.AuditLog, *base.ResponsePaginate, error)
	PurgeAuditLogsBefore(ctx context.Context, before time.Time) error
}
