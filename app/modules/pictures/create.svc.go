package pictures

import (
	"context"
	"fmt"
	"strings"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, schoolID uuid.UUID, uploadedByMemberID uuid.UUID, ownerMemberID *uuid.UUID, objectKey string, fileName string, contentType string, sizeBytes int64, metadata map[string]any) (*ent.Document, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "pictures.service.create")
	defer span.End()

	if sizeBytes < 0 || !isImageContentType(contentType) {
		return nil, fmt.Errorf("%w", ErrPictureConditionFail)
	}

	storage, err := s.resolvePictureStorage(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	if metadata == nil {
		metadata = map[string]any{}
	}
	metadata["kind"] = "picture"

	item, err := s.db.CreateDocument(ctx, &ent.Document{
		SchoolID:           schoolID,
		StorageID:          storage.ID,
		OwnerMemberID:      ownerMemberID,
		UploadedByMemberID: uploadedByMemberID,
		ObjectKey:          strings.TrimSpace(objectKey),
		FileName:           strings.TrimSpace(fileName),
		ContentType:        strings.TrimSpace(contentType),
		SizeBytes:          sizeBytes,
		Status:             ent.DocumentStatusPendingUpload,
		Metadata:           metadata,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

func isImageContentType(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "" {
		return false
	}
	parts := strings.SplitN(v, ";", 2)
	baseType := strings.TrimSpace(parts[0])
	return strings.HasPrefix(baseType, "image/")
}
