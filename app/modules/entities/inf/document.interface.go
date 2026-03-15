package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type DocumentEntity interface {
	CreateDocument(ctx context.Context, data *ent.Document) (*ent.Document, error)
	GetDocumentByID(ctx context.Context, id uuid.UUID, schoolID uuid.UUID) (*ent.Document, error)
	ListDocuments(ctx context.Context, req *base.RequestPaginate, schoolID uuid.UUID, ownerMemberID *uuid.UUID, status *ent.DocumentStatus, storageID *uuid.UUID) ([]*ent.Document, *base.ResponsePaginate, error)
	UpdateDocumentByID(ctx context.Context, id uuid.UUID, schoolID uuid.UUID, data *ent.DocumentUpdate) (*ent.Document, error)
	SoftDeleteDocumentByID(ctx context.Context, id uuid.UUID, schoolID uuid.UUID) error
}
