package repository

import (
	"context"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type RegistrationUploadSessionRepository struct {
	store *db.Store
}

func NewRegistrationUploadSessionRepository(store *db.Store) domain.RegistrationUploadSessionRepository {
	return &RegistrationUploadSessionRepository{store: store}
}

func (r *RegistrationUploadSessionRepository) Create(ctx context.Context, tokenHash string, expiresAt time.Time) (*domain.RegistrationUploadSession, error) {
	session, err := r.store.CreateRegistrationUploadSession(ctx, db.CreateRegistrationUploadSessionParams{TokenHash: tokenHash, ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true}})
	if err != nil {
		return nil, err
	}
	return toDomainUploadSession(session), nil
}

func (r *RegistrationUploadSessionRepository) GetActive(ctx context.Context, tokenHash string) (*domain.RegistrationUploadSession, error) {
	session, err := r.store.GetActiveRegistrationUploadSession(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	return toDomainUploadSession(session), nil
}

func (r *RegistrationUploadSessionRepository) AddAttachment(ctx context.Context, sessionID, attachmentID uuid.UUID) error {
	return r.store.AddRegistrationUploadAttachment(ctx, db.AddRegistrationUploadAttachmentParams{ID: sessionID, Column2: attachmentID})
}

func (r *RegistrationUploadSessionRepository) HasAttachments(ctx context.Context, sessionID uuid.UUID, attachmentIDs []uuid.UUID) (bool, error) {
	return r.store.RegistrationUploadSessionHasAttachments(ctx, db.RegistrationUploadSessionHasAttachmentsParams{ID: sessionID, Column2: attachmentIDs})
}

func toDomainUploadSession(session db.RegistrationUploadSession) *domain.RegistrationUploadSession {
	return &domain.RegistrationUploadSession{
		ID:            session.ID,
		TokenHash:     session.TokenHash,
		AttachmentIDs: session.AttachmentIds,
		ExpiresAt:     session.ExpiresAt.Time,
	}
}
