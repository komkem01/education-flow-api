package pictures

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, schoolID uuid.UUID, ownerMemberID *uuid.UUID, status *string) ([]*ent.Document, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "pictures.service.list")
	defer span.End()

	storage, err := s.resolvePictureStorage(ctx, schoolID)
	if err != nil {
		return nil, nil, err
	}

	var parsedStatus *ent.DocumentStatus
	if status != nil {
		v, ok := parsePictureStatus(*status)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrPictureConditionFail)
		}
		parsedStatus = &v
	}

	storageID := storage.ID
	items, page, err := s.db.ListDocuments(ctx, req, schoolID, ownerMemberID, parsedStatus, &storageID)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrPictureConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}
