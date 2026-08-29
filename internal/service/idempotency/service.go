package idempotency_service

import (
	"context"
	"time"

	"github.com/romreign/AuthService/internal/core/domain"
)

type IdempotencyRepository interface {
	GetIdempotencyKey(ctx context.Context, key string) (*domain.IdempotencyData, error)
	SaveIdempotencyKey(ctx context.Context, key string, data *domain.IdempotencyData, duration time.Duration) error
}

type IdempotencyService struct {
	idempotencyRepository IdempotencyRepository
	ttl                   time.Duration
}

func NewIdempotencyService(
	idempotencyRepository IdempotencyRepository,
	ttl time.Duration,
) *IdempotencyService {
	return &IdempotencyService{
		idempotencyRepository: idempotencyRepository,
		ttl:                   ttl,
	}
}
