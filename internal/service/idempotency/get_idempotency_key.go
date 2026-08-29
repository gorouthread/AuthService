package idempotency_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/romreign/AuthService/internal/core/domain"
	core_errors "github.com/romreign/AuthService/internal/core/errors"
)

func (s *IdempotencyService) GetIdempotencyKey(ctx context.Context, key uuid.UUID, method string, url string) (*domain.IdempotencyData, error) {
	keyRedis := "idempotency:" + key.String()
	response, err := s.idempotencyRepository.GetIdempotencyKey(ctx, keyRedis)
	if err != nil {
		return nil, fmt.Errorf("get idempotency data: %w", err)
	}

	if response == nil {
		return nil, nil
	}

	if response.Method != method || response.URL != url {
		return nil, fmt.Errorf("idempotency key already used for another request: %w", core_errors.ErrConflict)
	}

	return response, nil
}
