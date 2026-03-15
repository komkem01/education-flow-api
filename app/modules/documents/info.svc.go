package documents

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID, schoolID uuid.UUID) (*ent.Document, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "documents.service.get_by_id")
	defer span.End()

	item, err := s.db.GetDocumentByID(ctx, id, schoolID)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}
