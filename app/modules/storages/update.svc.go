package storages

import (
	"context"
	"fmt"
	"strings"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, schoolID uuid.UUID, provider *string, name *string, endpoint *string, bucketName *string, isDefault *bool, config *map[string]any) (*ent.Storage, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "storages.service.update")
	defer span.End()

	if provider == nil && name == nil && endpoint == nil && bucketName == nil && isDefault == nil && config == nil {
		return nil, fmt.Errorf("%w", ErrStorageConditionFail)
	}

	updateData := &ent.StorageUpdate{}
	if provider != nil {
		parsedProvider, ok := parseStorageProvider(*provider)
		if !ok {
			return nil, fmt.Errorf("%w", ErrStorageConditionFail)
		}
		updateData.Provider = &parsedProvider
	}
	if name != nil {
		v := strings.TrimSpace(*name)
		if v == "" {
			return nil, fmt.Errorf("%w", ErrStorageConditionFail)
		}
		updateData.Name = &v
	}
	if endpoint != nil {
		v := strings.TrimSpace(*endpoint)
		updateData.Endpoint = &v
	}
	if bucketName != nil {
		v := strings.TrimSpace(*bucketName)
		if v == "" {
			return nil, fmt.Errorf("%w", ErrStorageConditionFail)
		}
		updateData.BucketName = &v
	}
	if isDefault != nil {
		updateData.IsDefault = isDefault
	}
	if config != nil {
		updateData.Config = config
	}

	item, err := s.db.UpdateStorageByID(ctx, id, schoolID, updateData)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}
