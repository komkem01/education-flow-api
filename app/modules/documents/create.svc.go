package documents

import (
	"context"
	"fmt"
	"strings"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, schoolID uuid.UUID, uploadedByMemberID uuid.UUID, ownerMemberID *uuid.UUID, bucketName *string, objectKey string, fileName string, contentType string, sizeBytes int64, metadata map[string]any) (*ent.Document, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "documents.service.create")
	defer span.End()

	if sizeBytes < 0 {
		return nil, fmt.Errorf("%w", ErrDocumentConditionFail)
	}

	targetBucket := s.resolveBucketName(bucketName)
	if targetBucket == "" {
		return nil, fmt.Errorf("%w", ErrDocumentConditionFail)
	}

	storage, err := s.storageDB.EnsureStorageByBucket(ctx, schoolID, targetBucket, nil)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

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
