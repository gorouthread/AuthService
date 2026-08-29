package idempotency_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/romreign/AuthService/internal/core/domain"
)

func TestIdempotencyService_SaveIdempotencyKey(t *testing.T) {
	ctx := context.Background()

	key := uuid.New()
	ttl := 10 * time.Minute

	data := &domain.IdempotencyData{
		Method:     "POST",
		URL:        "/auth/login",
		StatusCode: 200,
		Body:       []byte(`{"access_token":"token"}`),
	}

	repositoryErr := errors.New("redis unavailable")

	tests := []struct {
		name string

		repositoryErr error

		wantErr   bool
		wantErrIs error

		wantKey      string
		wantTTL      time.Duration
		wantData     *domain.IdempotencyData
		wantSaveCall bool
	}{
		{
			name: "success",

			wantErr: false,

			wantKey:      "idempotency:" + key.String(),
			wantTTL:      ttl,
			wantData:     data,
			wantSaveCall: true,
		},
		{
			name: "repository error",

			repositoryErr: repositoryErr,

			wantErr:      true,
			wantSaveCall: true,

			wantKey:  "idempotency:" + key.String(),
			wantTTL:  ttl,
			wantData: data,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				saveCalled bool

				gotKey  string
				gotTTL  time.Duration
				gotData *domain.IdempotencyData
			)

			repository := &mockIdempotencyRepository{
				SaveIdempotencyKeyFunc: func(
					ctx context.Context,
					key string,
					data *domain.IdempotencyData,
					duration time.Duration,
				) error {
					saveCalled = true

					gotKey = key
					gotTTL = duration
					gotData = data

					return tt.repositoryErr
				},
			}

			service := NewIdempotencyService(
				repository,
				ttl,
			)

			err := service.SaveIdempotencyKey(
				ctx,
				key,
				data,
			)

			// Проверяем ошибку
			if tt.wantErr && err == nil {
				t.Fatal("SaveIdempotencyKey() error = nil, want error")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf(
					"SaveIdempotencyKey() unexpected error = %v",
					err,
				)
			}

			// Проверяем wrapping repository error
			if tt.wantErrIs != nil &&
				!errors.Is(err, tt.wantErrIs) {
				t.Errorf(
					"SaveIdempotencyKey() error = %v, want %v",
					err,
					tt.wantErrIs,
				)
			}

			// Проверяем вызов repository
			if saveCalled != tt.wantSaveCall {
				t.Errorf(
					"SaveIdempotencyKey() repository called = %v, want %v",
					saveCalled,
					tt.wantSaveCall,
				)
			}

			if !tt.wantSaveCall {
				return
			}

			// Проверяем Redis key
			if gotKey != tt.wantKey {
				t.Errorf(
					"repository key = %q, want %q",
					gotKey,
					tt.wantKey,
				)
			}

			// Проверяем TTL
			if gotTTL != tt.wantTTL {
				t.Errorf(
					"repository TTL = %v, want %v",
					gotTTL,
					tt.wantTTL,
				)
			}

			// Проверяем data
			if gotData == nil {
				t.Fatal("repository data = nil, want data")
			}

			if gotData.Method != tt.wantData.Method {
				t.Errorf(
					"Method = %q, want %q",
					gotData.Method,
					tt.wantData.Method,
				)
			}

			if gotData.URL != tt.wantData.URL {
				t.Errorf(
					"URL = %q, want %q",
					gotData.URL,
					tt.wantData.URL,
				)
			}

			if gotData.StatusCode != tt.wantData.StatusCode {
				t.Errorf(
					"StatusCode = %d, want %d",
					gotData.StatusCode,
					tt.wantData.StatusCode,
				)
			}

			if string(gotData.Body) != string(tt.wantData.Body) {
				t.Errorf(
					"Body = %s, want %s",
					gotData.Body,
					tt.wantData.Body,
				)
			}
		})
	}
}
