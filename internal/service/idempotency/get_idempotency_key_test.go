package idempotency_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/romreign/AuthService/internal/core/domain"
	core_errors "github.com/romreign/AuthService/internal/core/errors"
)

func TestIdempotencyService_GetIdempotencyKey(t *testing.T) {
	ctx := context.Background()

	key := uuid.New()

	data := &domain.IdempotencyData{
		Method:     "POST",
		URL:        "/auth/login",
		StatusCode: 200,
		Body:       []byte(`{"access_token":"token"}`),
	}

	repositoryErr := errors.New("redis unavailable")

	tests := []struct {
		name string

		method string
		url    string

		repositoryData *domain.IdempotencyData
		repositoryErr  error

		wantData    *domain.IdempotencyData
		wantErr     bool
		wantErrIs   error
		wantGetCall bool

		wantKey string
	}{
		{
			name: "key not found",

			method: "POST",
			url:    "/auth/login",

			repositoryData: nil,

			wantData:    nil,
			wantErr:     false,
			wantGetCall: true,

			wantKey: "idempotency:" + key.String(),
		},
		{
			name: "success",

			method: "POST",
			url:    "/auth/login",

			repositoryData: data,

			wantData:    data,
			wantErr:     false,
			wantGetCall: true,

			wantKey: "idempotency:" + key.String(),
		},
		{
			name: "repository error",

			method: "POST",
			url:    "/auth/login",

			repositoryErr: repositoryErr,

			wantData:    nil,
			wantErr:     true,
			wantGetCall: true,

			wantKey: "idempotency:" + key.String(),
		},
		{
			name: "method conflict",

			method: "GET",
			url:    "/auth/login",

			repositoryData: data,

			wantData:    nil,
			wantErr:     true,
			wantErrIs:   core_errors.ErrConflict,
			wantGetCall: true,

			wantKey: "idempotency:" + key.String(),
		},
		{
			name: "url conflict",

			method: "POST",
			url:    "/auth/logout",

			repositoryData: data,

			wantData:    nil,
			wantErr:     true,
			wantErrIs:   core_errors.ErrConflict,
			wantGetCall: true,

			wantKey: "idempotency:" + key.String(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				getCalled bool
				gotKey    string
			)

			repository := &mockIdempotencyRepository{
				GetIdempotencyKeyFunc: func(
					ctx context.Context,
					key string,
				) (*domain.IdempotencyData, error) {
					getCalled = true
					gotKey = key

					if tt.repositoryErr != nil {
						return nil, tt.repositoryErr
					}

					return tt.repositoryData, nil
				},
			}

			service := NewIdempotencyService(
				repository,
				10*time.Minute,
			)

			gotData, err := service.GetIdempotencyKey(
				ctx,
				key,
				tt.method,
				tt.url,
			)

			// Проверяем ошибку
			if tt.wantErr && err == nil {
				t.Fatal("GetIdempotencyKey() error = nil, want error")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf(
					"GetIdempotencyKey() unexpected error = %v",
					err,
				)
			}

			// Проверяем errors.Is
			if tt.wantErrIs != nil &&
				!errors.Is(err, tt.wantErrIs) {
				t.Errorf(
					"GetIdempotencyKey() error = %v, want %v",
					err,
					tt.wantErrIs,
				)
			}

			// Проверяем вызов repository
			if getCalled != tt.wantGetCall {
				t.Errorf(
					"GetIdempotencyKey() repository called = %v, want %v",
					getCalled,
					tt.wantGetCall,
				)
			}

			if !tt.wantGetCall {
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

			// Если ожидаем nil
			if tt.wantData == nil {
				if gotData != nil {
					t.Errorf(
						"GetIdempotencyKey() data = %+v, want nil",
						gotData,
					)
				}

				return
			}

			// Проверяем возвращённые данные
			if gotData == nil {
				t.Fatal(
					"GetIdempotencyKey() data = nil, want data",
				)
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
