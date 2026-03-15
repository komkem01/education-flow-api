package storages

import (
	"context"
	"fmt"
	"strings"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, schoolID uuid.UUID, provider string, name string, endpoint *string, bucketName string, isDefault *bool, config map[string]any) (*ent.Storage, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "storages.service.create")
	defer span.End()

	parsedProvider, ok := parseStorageProvider(provider)
	if !ok {
		return nil, fmt.Errorf("%w", ErrStorageConditionFail)
	}

	name = strings.TrimSpace(name)
	bucketName = strings.TrimSpace(bucketName)
	if name == "" || bucketName == "" {
		return nil, fmt.Errorf("%w", ErrStorageConditionFail)
	}

	var normalizedEndpoint *string
	if endpoint != nil {
		trimmed := strings.TrimSpace(*endpoint)
		if trimmed != "" {
			normalizedEndpoint = &trimmed
		}
	}

	defaultValue := false
	if isDefault != nil {
		defaultValue = *isDefault
	}

	item, err := s.db.CreateStorage(ctx, &ent.Storage{
		SchoolID:   schoolID,
		Provider:   parsedProvider,
		Name:       name,
		Endpoint:   normalizedEndpoint,
		BucketName: bucketName,
		IsDefault:  defaultValue,
		Config:     config,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}
