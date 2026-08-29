package idempotency_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/romreign/AuthService/internal/core/domain"
)

func (s *IdempotencyService) SaveIdempotencyKey(
	ctx context.Context,
	key uuid.UUID,
	data *domain.IdempotencyData,
) error {
	keyRedis := "idempotency:" + key.String()
	if err := s.idempotencyRepository.SaveIdempotencyKey(
		ctx,
		keyRedis,
		data,
		s.ttl,
	); err != nil {
		return fmt.Errorf("save idempotency data: %w", err)
	}
	return nil
}
