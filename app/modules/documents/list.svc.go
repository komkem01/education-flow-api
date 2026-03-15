package documents

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, schoolID uuid.UUID, ownerMemberID *uuid.UUID, status *string, storageID *uuid.UUID) ([]*ent.Document, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "documents.service.list")
	defer span.End()

	var parsedStatus *ent.DocumentStatus
	if status != nil {
		v, ok := parseDocumentStatus(*status)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrDocumentConditionFail)
		}
		parsedStatus = &v
	}

	items, page, err := s.db.ListDocuments(ctx, req, schoolID, ownerMemberID, parsedStatus, storageID)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrDocumentConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}
