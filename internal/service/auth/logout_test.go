package auth_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/romreign/AuthService/internal/core/domain"
	core_errors "github.com/romreign/AuthService/internal/core/errors"
)

func TestAuthService_Logout(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()
	sessionID := uuid.New()

	session := domain.Session{
		ID:               sessionID,
		UserID:           userID,
		RefreshToken:     "refresh-token-hash",
		IsRevoked:        false,
		RefreshCreatedAt: time.Now().Add(-time.Hour),
		RefreshExpiresAt: time.Now().Add(time.Hour),
	}

	tests := []struct {
		name string

		getSessionErr error
		updateErr     error

		wantErr   bool
		wantErrIs error

		wantUpdate bool
	}{
		{
			name:       "success",
			wantUpdate: true,
		},
		{
			name:          "get session error",
			getSessionErr: errors.New("database unavailable"),
			wantErr:       true,
			wantErrIs:     core_errors.ErrNotFound,
			wantUpdate:    false,
		},
		{
			name:       "update session error",
			updateErr:  errors.New("database unavailable"),
			wantErr:    true,
			wantUpdate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				updateCalled   bool
				updatedSession domain.Session
			)

			repository := &mockAuthRepository{
				GetSessionByRefreshTokenFunc: func(
					ctx context.Context,
					session domain.Session,
				) (domain.Session, error) {
					if tt.getSessionErr != nil {
						return domain.Session{}, tt.getSessionErr
					}

					return session, nil
				},

				UpdateSessionFunc: func(
					ctx context.Context,
					session domain.Session,
				) (domain.Session, error) {
					updateCalled = true
					updatedSession = session

					if tt.updateErr != nil {
						return domain.Session{}, tt.updateErr
					}

					return session, nil
				},
			}

			service := NewAuthService(repository, nil)

			err := service.Logout(ctx, session)

			if tt.wantErr && err == nil {
				t.Fatal("Logout() error = nil, want error")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("Logout() unexpected error = %v", err)
			}

			if tt.wantErrIs != nil &&
				!errors.Is(err, tt.wantErrIs) {
				t.Errorf(
					"Logout() error = %v, want %v",
					err,
					tt.wantErrIs,
				)
			}

			if updateCalled != tt.wantUpdate {
				t.Errorf(
					"UpdateSession called = %v, want %v",
					updateCalled,
					tt.wantUpdate,
				)
			}

			if updateCalled && !updatedSession.IsRevoked {
				t.Error("session was not revoked")
			}
		})
	}
}
