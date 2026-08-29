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
	"github.com/romreign/AuthService/pkg/utils"
)

func TestAuthService_Login(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()

	dbPasswordHash := func() string {
		hash, err := utils.HashPassword("password123")
		if err != nil {
			panic(err)
		}

		return hash
	}

	tests := []struct {
		name string

		user domain.User

		userFromDB domain.User

		getUserErr       error
		accessTokenErr   error
		refreshTokenErr  error
		createSessionErr error

		wantErr   bool
		wantErrIs error

		wantCreateSession bool
	}{
		{
			name: "success",
			user: domain.User{
				Login:        "roman",
				PasswordHash: "password123",
				Role:         domain.UninitializedRole,
			},
			userFromDB: domain.User{
				ID:           userID,
				Login:        "roman",
				PasswordHash: dbPasswordHash(),
				Role:         "user",
			},
			wantCreateSession: true,
		},
		{
			name: "invalid user",
			user: domain.User{
				Login:        "ab",
				PasswordHash: "password123",
			},
			wantErr:   true,
			wantErrIs: core_errors.ErrInvalidArgument,
		},
		{
			name: "get user repository error",
			user: domain.User{
				Login:        "roman",
				PasswordHash: "password123",
			},
			getUserErr: errors.New("database unavailable"),
			wantErr:    true,
			wantErrIs:  core_errors.ErrNotFound,
		},
		{
			name: "wrong password",
			user: domain.User{
				Login:        "roman",
				PasswordHash: "wrong-password",
			},
			userFromDB: domain.User{
				ID:           userID,
				Login:        "roman",
				PasswordHash: dbPasswordHash(),
				Role:         "user",
			},
			wantErr:   true,
			wantErrIs: core_errors.ErrNotFound,
		},
		{
			name: "access token error",
			user: domain.User{
				Login:        "roman",
				PasswordHash: "password123",
			},
			userFromDB: domain.User{
				ID:           userID,
				Login:        "roman",
				PasswordHash: dbPasswordHash(),
				Role:         "user",
			},
			accessTokenErr: errors.New("jwt error"),
			wantErr:        true,
		},
		{
			name: "refresh token error",
			user: domain.User{
				Login:        "roman",
				PasswordHash: "password123",
			},
			userFromDB: domain.User{
				ID:           userID,
				Login:        "roman",
				PasswordHash: dbPasswordHash(),
				Role:         "user",
			},
			refreshTokenErr: errors.New("jwt error"),
			wantErr:         true,
		},
		{
			name: "create session error",
			user: domain.User{
				Login:        "roman",
				PasswordHash: "password123",
			},
			userFromDB: domain.User{
				ID:           userID,
				Login:        "roman",
				PasswordHash: dbPasswordHash(),
				Role:         "user",
			},
			createSessionErr:  errors.New("database error"),
			wantErr:           true,
			wantCreateSession: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var createSessionCalled bool

			repository := &mockAuthRepository{
				GetUserByLoginFunc: func(
					ctx context.Context,
					user domain.User,
				) (domain.User, error) {
					if tt.getUserErr != nil {
						return domain.User{}, tt.getUserErr
					}

					return tt.userFromDB, nil
				},

				CreateSessionFunc: func(
					ctx context.Context,
					session domain.Session,
				) (domain.Session, error) {
					createSessionCalled = true

					if tt.createSessionErr != nil {
						return domain.Session{}, tt.createSessionErr
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
							IssuedAt:  external_jwt.NewNumericDate(time.Now()),
							ExpiresAt: external_jwt.NewNumericDate(time.Now().Add(duration)),
						},
					}

					if duration == 15*time.Minute &&
						tt.accessTokenErr != nil {
						return "", nil, tt.accessTokenErr
					}

					if duration == 24*time.Hour &&
						tt.refreshTokenErr != nil {
						return "", nil, tt.refreshTokenErr
					}

					if duration == 15*time.Minute {
						return "access-token", claims, nil
					}

					return "refresh-token", claims, nil
				},
			}

			service := NewAuthService(repository, jwtManager)

			session, err := service.Login(ctx, tt.user)

			if tt.wantErr && err == nil {
				t.Fatal("Login() error = nil, want error")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("Login() unexpected error = %v", err)
			}

			if tt.wantErrIs != nil &&
				!errors.Is(err, tt.wantErrIs) {
				t.Errorf(
					"Login() error = %v, want %v",
					err,
					tt.wantErrIs,
				)
			}

			if createSessionCalled != tt.wantCreateSession {
				t.Errorf(
					"CreateSession called = %v, want %v",
					createSessionCalled,
					tt.wantCreateSession,
				)
			}

			if !tt.wantErr {
				if session.AccessToken != "access-token" {
					t.Errorf(
						"AccessToken = %q, want %q",
						session.AccessToken,
						"access-token",
					)
				}

				if session.RefreshToken != "refresh-token" {
					t.Errorf(
						"RefreshToken = %q, want %q",
						session.RefreshToken,
						"refresh-token",
					)
				}

				if session.UserID != userID {
					t.Errorf(
						"UserID = %v, want %v",
						session.UserID,
						userID,
					)
				}
			}
		})
	}
}
