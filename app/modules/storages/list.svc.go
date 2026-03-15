package storages

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, schoolID uuid.UUID, provider *string, isDefault *bool) ([]*ent.Storage, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "storages.service.list")
	defer span.End()

	var parsedProvider *ent.StorageProvider
	if provider != nil {
		v, ok := parseStorageProvider(*provider)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrStorageConditionFail)
		}
		parsedProvider = &v
	}

	items, page, err := s.db.ListStorages(ctx, req, schoolID, parsedProvider, isDefault)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrStorageConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}
