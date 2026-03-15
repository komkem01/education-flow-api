package auditlogs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, actorID *uuid.UUID, actorRole *string, method *string, path *string, statusCode *int, from *time.Time, to *time.Time) ([]*ent.AuditLog, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "auditlogs.service.list")
	defer span.End()

	var normalizedMethod *string
	if method != nil {
		m, ok := parseHTTPMethod(*method)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrAuditLogConditionFail)
		}
		normalizedMethod = &m
	}

	if statusCode != nil && (*statusCode < 100 || *statusCode > 599) {
		return nil, nil, fmt.Errorf("%w", ErrAuditLogConditionFail)
	}

	if from != nil && to != nil && from.After(*to) {
		return nil, nil, fmt.Errorf("%w", ErrAuditLogConditionFail)
	}

	items, page, err := s.db.ListAuditLogs(ctx, req, actorID, actorRole, normalizedMethod, path, statusCode, from, to)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrAuditLogConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

func parseHTTPMethod(v string) (string, bool) {
	method := strings.ToUpper(strings.TrimSpace(v))
	switch method {
	case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS":
		return method, true
	default:
		return "", false
	}
}
