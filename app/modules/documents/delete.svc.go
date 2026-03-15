package documents

import (
	"context"

	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Delete(ctx context.Context, id uuid.UUID, schoolID uuid.UUID) error {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "documents.service.delete")
	defer span.End()

	if err := s.db.SoftDeleteDocumentByID(ctx, id, schoolID); err != nil {
		return normalizeServiceError(err)
	}

	return nil
}
