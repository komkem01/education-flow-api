package pictures

import (
	"context"

	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Delete(ctx context.Context, id uuid.UUID, schoolID uuid.UUID) error {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "pictures.service.delete")
	defer span.End()

	item, err := s.GetByID(ctx, id, schoolID)
	if err != nil {
		return err
	}

	if err := s.db.SoftDeleteDocumentByID(ctx, item.ID, schoolID); err != nil {
		return normalizeServiceError(err)
	}

	return nil
}
