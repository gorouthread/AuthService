package auth_service

import (
	"context"
	"errors"
	"testing"
	"time"

	external_jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/romreign/AuthService/internal/core/domain"
	core_errors "github.com/romreign/AuthService/internal/core/errors"
	jwt "github.com/romreign/AuthService/pkg/jwt"
)

func TestAuthService_Refresh(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()

	validSession := domain.Session{
		ID:               uuid.New(),
		UserID:           userID,
		RefreshToken:     "refresh-token-hash",
		IsRevoked:        false,
		RefreshCreatedAt: time.Now().Add(-time.Hour),
		RefreshExpiresAt: time.Now().Add(time.Hour),
	}

	user := domain.User{
		ID:           userID,
		Login:        "roman",
		PasswordHash: "hash",
		Role:         "user",
	}

	tests := []struct {
		name string

		session domain.Session

		getSessionErr error
		updateErr     error
		getUserErr    error
		accessErr     error
		refreshErr    error
		createErr     error

		expired bool
		revoked bool

		wantErr   bool
		wantErrIs error

		wantUpdate bool
		wantCreate bool
	}{
		{
			name:       "success",
			session:    validSession,
			wantUpdate: true,
			wantCreate: true,
		},
		{
			name:          "get session error",
			session:       validSession,
			getSessionErr: errors.New("database error"),
			wantErr:       true,
			wantErrIs:     core_errors.ErrNotFound,
		},
		{
			name:      "expired refresh token",
			session:   validSession,
			expired:   true,
			wantErr:   true,
			wantErrIs: core_errors.ErrUnauthorized,
		},
		{
			name:      "revoked session",
			session:   validSession,
			revoked:   true,
			wantErr:   true,
			wantErrIs: core_errors.ErrUnauthorized,
		},
		{
			name:       "update session error",
			session:    validSession,
			updateErr:  errors.New("database error"),
			wantErr:    true,
			wantUpdate: true,
		},
		{
			name:       "get user error",
			session:    validSession,
			getUserErr: errors.New("database error"),
			wantErr:    true,
			wantErrIs:  core_errors.ErrNotFound,
			wantUpdate: true,
		},
		{
			name:       "access token error",
			session:    validSession,
			accessErr:  errors.New("jwt error"),
			wantErr:    true,
			wantUpdate: true,
		},
		{
			name:       "refresh token error",
			session:    validSession,
			refreshErr: errors.New("jwt error"),
			wantErr:    true,
			wantUpdate: true,
		},
		{
			name:       "create session error",
			session:    validSession,
			createErr:  errors.New("database error"),
			wantErr:    true,
			wantUpdate: true,
			wantCreate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				updateCalled bool
				createCalled bool

				updatedSession domain.Session
				createdSession domain.Session
			)

			sessionFromDB := tt.session

			if tt.expired {
				sessionFromDB.RefreshExpiresAt = time.Now().Add(-time.Hour)
			}

			if tt.revoked {
				sessionFromDB.IsRevoked = true
			}

			repository := &mockAuthRepository{
				GetSessionByRefreshTokenFunc: func(
					ctx context.Context,
					session domain.Session,
				) (domain.Session, error) {
					if tt.getSessionErr != nil {
						return domain.Session{}, tt.getSessionErr
					}

					return sessionFromDB, nil
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

				GetUserByIDFunc: func(
					ctx context.Context,
					id uuid.UUID,
				) (domain.User, error) {
					if tt.getUserErr != nil {
						return domain.User{}, tt.getUserErr
					}

					return user, nil
				},

				CreateSessionFunc: func(
					ctx context.Context,
					session domain.Session,
				) (domain.Session, error) {
					createCalled = true
					createdSession = session

					if tt.createErr != nil {
						return domain.Session{}, tt.createErr
					}

					session.ID = uuid.New()

					return session, nil
				},
			}

			jwtManager := &mockJWTManager{
				CreateTokenFunc: func(
					id uuid.UUID,
					role string,
					duration time.Duration,
				) (string, *jwt.UserClaims, error) {

					claims := &jwt.UserClaims{
						UserID: id,
						Role:   role,
						RegisteredClaims: external_jwt.RegisteredClaims{
							IssuedAt: external_jwt.NewNumericDate(time.Now()),
							ExpiresAt: external_jwt.NewNumericDate(
								time.Now().Add(duration),
							),
						},
					}

					if duration == 15*time.Minute {
						if tt.accessErr != nil {
							return "", nil, tt.accessErr
						}

						return "new-access-token", claims, nil
					}

					if tt.refreshErr != nil {
						return "", nil, tt.refreshErr
					}

					return "new-refresh-token", claims, nil
				},
			}

			service := NewAuthService(repository, jwtManager)

			result, err := service.Refresh(ctx, tt.session)

			if tt.wantErr && err == nil {
				t.Fatal("Refresh() error = nil, want error")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf(
					"Refresh() unexpected error = %v",
					err,
				)
			}

			if tt.wantErrIs != nil &&
				!errors.Is(err, tt.wantErrIs) {
				t.Errorf(
					"Refresh() error = %v, want %v",
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

			if createCalled != tt.wantCreate {
				t.Errorf(
					"CreateSession called = %v, want %v",
					createCalled,
					tt.wantCreate,
				)
			}

			if tt.wantUpdate && !updatedSession.IsRevoked {
				t.Error(
					"old session was not revoked",
				)
			}

			if !tt.wantErr {
				if result.AccessToken != "new-access-token" {
					t.Errorf(
						"AccessToken = %q, want %q",
						result.AccessToken,
						"new-access-token",
					)
				}

				if result.RefreshToken != "new-refresh-token" {
					t.Errorf(
						"RefreshToken = %q, want %q",
						result.RefreshToken,
						"new-refresh-token",
					)
				}

				if createCalled &&
					createdSession.RefreshToken == "new-refresh-token" {
					t.Error(
						"new session must store hashed refresh token",
					)
				}

				if result.UserID != userID {
					t.Errorf(
						"UserID = %v, want %v",
						result.UserID,
						userID,
					)
				}
			}
		})
	}
}
