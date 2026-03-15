package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.AuthSessionEntity = (*Service)(nil)

func (s *Service) CreateAuthSession(ctx context.Context, data *ent.AuthSession) (*ent.AuthSession, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetAuthSessionByToken(ctx context.Context, token string) (*ent.AuthSession, error) {
	row := new(ent.AuthSession)
	if err := s.db.NewSelect().Model(row).Where("token = ?", token).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) GetAuthSessionByRefreshToken(ctx context.Context, refreshToken string) (*ent.AuthSession, error) {
	row := new(ent.AuthSession)
	if err := s.db.NewSelect().Model(row).Where("refresh_token = ?", refreshToken).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) UpdateAuthSessionByID(ctx context.Context, id uuid.UUID, data *ent.AuthSessionUpdate) (*ent.AuthSession, error) {
	query := s.db.NewUpdate().
		Model(&ent.AuthSession{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.Token != nil {
		query.Set("token = ?", *data.Token)
	}
	if data.RefreshToken != nil {
		query.Set("refresh_token = ?", *data.RefreshToken)
	}
	if data.ExpireAt != nil {
		query.Set("expire_at = ?", *data.ExpireAt)
	}
	if data.RefreshExpireAt != nil {
		query.Set("refresh_expire_at = ?", *data.RefreshExpireAt)
	}

	res, err := query.Exec(ctx)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, sql.ErrNoRows
	}

	row := new(ent.AuthSession)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) DeleteAuthSessionByToken(ctx context.Context, token string) error {
	res, err := s.db.NewDelete().Model(&ent.AuthSession{}).Where("token = ?", token).Exec(ctx)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Service) DeleteAuthSessionsByMemberID(ctx context.Context, memberID uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.AuthSession{}).Where("member_id = ?", memberID).Exec(ctx)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}
