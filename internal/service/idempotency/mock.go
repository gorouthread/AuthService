package idempotency_service

import (
	"context"
	"time"

	"github.com/romreign/AuthService/internal/core/domain"
)

type mockIdempotencyRepository struct {
	GetIdempotencyKeyFunc func(
		ctx context.Context,
		key string,
	) (*domain.IdempotencyData, error)

	SaveIdempotencyKeyFunc func(
		ctx context.Context,
		key string,
		data *domain.IdempotencyData,
		duration time.Duration,
	) error
}

func (m *mockIdempotencyRepository) GetIdempotencyKey(
	ctx context.Context,
	key string,
) (*domain.IdempotencyData, error) {
	return m.GetIdempotencyKeyFunc(ctx, key)
}

func (m *mockIdempotencyRepository) SaveIdempotencyKey(
	ctx context.Context,
	key string,
	data *domain.IdempotencyData,
	duration time.Duration,
) error {
	return m.SaveIdempotencyKeyFunc(ctx, key, data, duration)
}
