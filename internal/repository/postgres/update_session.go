package auth_repository_postgres

import (
	"context"
	"fmt"

	"github.com/romreign/AuthService/internal/core/domain"
)

func (r *AuthRepositoryPostgres) UpdateSession(
	ctx context.Context,
	session domain.Session,
) (domain.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE auth_msginx.sessions
		SET
			refresh_token_hash = $1,
			user_id = $2,
			is_revoked = $3,
			created_at = $4,
			expires_at = $5,
		WHERE id = $6
		RETURNING
			id,
			refresh_token_hash,
			user_id,
			is_revoked,
			created_at,
			expires_at;
	`

	row := r.pool.QueryRow(ctx, query,
		[]byte(session.RefreshToken),
		session.UserID,
		session.IsRevoked,
		session.RefreshCreatedAt,
		session.RefreshExpiresAt,
		session.ID,
	)

	var sessionModel Session
	if err := row.Scan(
		&sessionModel.ID,
		&sessionModel.RefreshTokenHash,
		&sessionModel.UserID,
		&sessionModel.IsRevoked,
		&sessionModel.CreatedAt,
		&sessionModel.ExpiresAt,
	); err != nil {
		return domain.Session{}, fmt.Errorf("failed to scan session row: %w", err)
	}

	session = domain.NewSession(
		sessionModel.ID,
		sessionModel.UserID,
		string(sessionModel.RefreshTokenHash),
		sessionModel.IsRevoked,
		sessionModel.CreatedAt,
		sessionModel.ExpiresAt,
	)

	return session, nil
}
