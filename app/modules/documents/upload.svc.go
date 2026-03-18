package documents

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func normalizeObjectKeyForUpload(v string) string {
	v = strings.TrimSpace(v)
	v = strings.ReplaceAll(v, "\\", "/")
	v = strings.TrimLeft(v, "/")
	v = strings.TrimSpace(path.Clean(v))
	if v == "." {
		return ""
	}
	if strings.HasPrefix(v, "../") || v == ".." {
		return ""
	}
	if strings.Contains(v, "/../") || strings.HasSuffix(v, "/..") {
		return ""
	}
	return v
}

func sanitizeObjectKeyFileName(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "file"
	}
	var b strings.Builder
	b.Grow(len(v))
	for _, r := range v {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '_' || r == '-' {
			b.WriteRune(r)
			continue
		}
		b.WriteRune('_')
	}
	name := strings.Trim(strings.TrimSpace(b.String()), "._-")
	if name == "" {
		return "file"
	}
	if len(name) > 120 {
		name = name[:120]
	}
	return name
}

func buildFallbackObjectKey(schoolID uuid.UUID, documentID uuid.UUID, fileName string) string {
	safeName := sanitizeObjectKeyFileName(fileName)
	return fmt.Sprintf("student-registration/%s/%d-%s-%s", schoolID.String(), time.Now().UnixMilli(), documentID.String(), safeName)
}

func uploadByPresignedURL(ctx context.Context, url string, contentType string, data []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, strings.TrimSpace(url), bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)
	req.ContentLength = int64(len(data))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 {
			return fmt.Errorf("presigned-upload-http-%d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
		}
		return fmt.Errorf("presigned-upload-http-%d", resp.StatusCode)
	}

	return nil
}

func (s *Service) UploadByID(ctx context.Context, id uuid.UUID, schoolID uuid.UUID, actorMemberID uuid.UUID, actorRole ent.MemberRole, fileName string, contentType string, sizeBytes int64, file io.Reader) (*ent.Document, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "documents.service.upload_by_id")
	defer span.End()

	if sizeBytes < 0 {
		return nil, fmt.Errorf("%w", ErrDocumentConditionFail)
	}

	uploadData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrDocumentConditionFail)
	}
	sizeBytes = int64(len(uploadData))

	item, err := s.db.GetDocumentByID(ctx, id, schoolID)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	isElevated := actorRole == ent.MemberRoleAdmin || actorRole == ent.MemberRoleSuperadmin
	if !isElevated {
		isOwner := item.OwnerMemberID != nil && *item.OwnerMemberID == actorMemberID
		if item.UploadedByMemberID != actorMemberID && !isOwner {
			return nil, fmt.Errorf("%w", ErrDocumentUnauthorized)
		}
	}

	targetBucket := s.resolveBucketName(nil)
	if storage, storageErr := s.storageDB.GetStorageByID(ctx, item.StorageID, schoolID); storageErr == nil && storage != nil {
		targetBucket = s.resolveBucketName(&storage.BucketName)
	}

	normalizedObjectKey := normalizeObjectKeyForUpload(item.ObjectKey)
	if normalizedObjectKey == "" {
		normalizedObjectKey = buildFallbackObjectKey(schoolID, id, fileName)
	}

	trimmedContentType := strings.TrimSpace(contentType)
	if trimmedContentType == "" {
		trimmedContentType = "application/octet-stream"
	}

	uploadURL, err := s.PresignUploadURL(ctx, normalizedObjectKey, &targetBucket, nil)
	if err != nil {
		return nil, err
	}

	if err = uploadByPresignedURL(ctx, uploadURL, trimmedContentType, uploadData); err != nil {
		normalizedObjectKey = buildFallbackObjectKey(schoolID, id, fileName)
		uploadURL, err = s.PresignUploadURL(ctx, normalizedObjectKey, &targetBucket, nil)
		if err != nil {
			return nil, err
		}
		if err = uploadByPresignedURL(ctx, uploadURL, trimmedContentType, uploadData); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDocumentConditionFail, err)
		}
	}

	status := ent.DocumentStatusActive
	fileName = strings.TrimSpace(fileName)
	contentType = trimmedContentType
	updated, err := s.db.UpdateDocumentByID(ctx, id, schoolID, &ent.DocumentUpdate{
		ObjectKey:   &normalizedObjectKey,
		FileName:    &fileName,
		ContentType: &contentType,
		SizeBytes:   &sizeBytes,
		Status:      &status,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return updated, nil
}
