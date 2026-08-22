package repository

import (
	"context"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/ctxkeys"
	"maicare_go/internal/domain"

	"github.com/google/uuid"
)

type AttachmentRepository struct {
	store *db.Store
}

func NewAttachmentRepository(store *db.Store) domain.AttachmentRepository {
	return &AttachmentRepository{store: store}
}

func (r *AttachmentRepository) CreateAttachment(ctx context.Context, params domain.CreateAttachmentParams) (*domain.AttachmentFile, error) {
	attachment, err := r.attachmentQuery(ctx, func(q *db.Queries) (db.AttachmentFile, error) {
		return q.CreateAttachment(ctx, db.CreateAttachmentParams{
			Uuid: params.ID,
			Name: params.Name,
			File: params.File,
			Size: params.Size,
			Tag:  params.Tag,
		})
	})
	if err != nil {
		return nil, err
	}
	return toDomainAttachment(attachment), nil
}

func (r *AttachmentRepository) GetAttachmentByID(ctx context.Context, id uuid.UUID) (*domain.AttachmentFile, error) {
	attachment, err := r.attachmentQuery(ctx, func(q *db.Queries) (db.AttachmentFile, error) {
		return q.GetActorAttachmentById(ctx, id)
	})
	if err != nil {
		return nil, err
	}
	return toDomainAttachment(attachment), nil
}

func (r *AttachmentRepository) SetAttachmentUsed(ctx context.Context, id uuid.UUID, used bool) (*domain.AttachmentFile, error) {
	attachment, err := r.attachmentQuery(ctx, func(q *db.Queries) (db.AttachmentFile, error) {
		return q.SetActorAttachmentAsUsedOrUnused(ctx, db.SetActorAttachmentAsUsedOrUnusedParams{Uuid: id, IsUsed: used})
	})
	if err != nil {
		return nil, err
	}
	return toDomainAttachment(attachment), nil
}

func (r *AttachmentRepository) DeleteAttachment(ctx context.Context, id uuid.UUID) (*domain.AttachmentFile, error) {
	attachment, err := r.attachmentQuery(ctx, func(q *db.Queries) (db.AttachmentFile, error) {
		return q.DeleteActorAttachment(ctx, id)
	})
	if err != nil {
		return nil, err
	}
	return toDomainAttachment(attachment), nil
}

func (r *AttachmentRepository) attachmentQuery(ctx context.Context, query func(*db.Queries) (db.AttachmentFile, error)) (db.AttachmentFile, error) {
	if _, ok := ctxkeys.ActorIdentityFromContext(ctx); !ok {
		return query(r.store.Queries)
	}
	return actorQuery(ctx, r.store, query)
}

func toDomainAttachment(attachment db.AttachmentFile) *domain.AttachmentFile {
	return &domain.AttachmentFile{
		UUID:      attachment.Uuid,
		Name:      attachment.Name,
		File:      attachment.File,
		Size:      attachment.Size,
		IsUsed:    attachment.IsUsed,
		CreatedAt: attachment.Created.Time,
	}
}
