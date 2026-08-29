package auth_service

import (
	"context"
	"errors"
	"testing"

	"github.com/romreign/AuthService/internal/core/domain"
	core_errors "github.com/romreign/AuthService/internal/core/errors"
)

func TestAuthService_Register(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string

		user domain.User

		repositoryErr error

		wantErr    bool
		wantErrIs  error
		wantCreate bool
		checkUser  func(t *testing.T, user domain.User)
	}{
		{
			name: "success",
			user: domain.User{
				ID:           domain.UninitializedUUID,
				Login:        "roman",
				PasswordHash: "password123",
				Role:         domain.UninitializedRole,
			},
			wantCreate: true,
			checkUser: func(t *testing.T, user domain.User) {
				if user.Login != "roman" {
					t.Errorf("Login = %q, want %q", user.Login, "roman")
				}

				if user.PasswordHash == "password123" {
					t.Error("password was not hashed")
				}

				if user.PasswordHash == "" {
					t.Error("password hash is empty")
				}

				if user.Role != domain.UninitializedRole {
					t.Errorf(
						"Role = %q, want %q",
						user.Role,
						domain.UninitializedRole,
					)
				}
			},
		},
		{
			name: "invalid user",
			user: domain.User{
				Login:        "ab",
				PasswordHash: "password123",
				Role:         domain.UninitializedRole,
			},
			wantErr:    true,
			wantErrIs:  core_errors.ErrInvalidArgument,
			wantCreate: false,
		},
		{
			name: "repository error",
			user: domain.User{
				Login:        "roman",
				PasswordHash: "password123",
				Role:         domain.UninitializedRole,
			},
			repositoryErr: errors.New("database unavailable"),
			wantErr:       true,
			wantCreate:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var createCalled bool

			repository := &mockAuthRepository{
				CreateUserFunc: func(
					ctx context.Context,
					user domain.User,
				) (domain.User, error) {
					createCalled = true

					if tt.checkUser != nil {
						tt.checkUser(t, user)
					}

					if tt.repositoryErr != nil {
						return domain.User{}, tt.repositoryErr
					}

					return user, nil
				},
			}

			service := NewAuthService(repository, nil)

			err := service.Register(ctx, tt.user)

			if tt.wantErr && err == nil {
				t.Fatal("Register() error = nil, want error")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("Register() unexpected error = %v", err)
			}

			if tt.wantErrIs != nil &&
				!errors.Is(err, tt.wantErrIs) {
				t.Errorf(
					"Register() error = %v, want %v",
					err,
					tt.wantErrIs,
				)
			}

			if createCalled != tt.wantCreate {
				t.Errorf(
					"CreateUser called = %v, want %v",
					createCalled,
					tt.wantCreate,
				)
			}
		})
	}
}
