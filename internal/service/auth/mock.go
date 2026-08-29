package auth_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/romreign/AuthService/internal/core/domain"

	jwt "github.com/romreign/AuthService/pkg/jwt"
)

type mockJWTManager struct {
	CreateTokenFunc func(
		uuid.UUID,
		string,
		time.Duration,
	) (string, *jwt.UserClaims, error)

	VerifyTokenFunc func(string) (*jwt.UserClaims, error)
}

func (m *mockJWTManager) CreateToken(
	id uuid.UUID,
	role string,
	duration time.Duration,
) (string, *jwt.UserClaims, error) {
	return m.CreateTokenFunc(id, role, duration)
}

func (m *mockJWTManager) VerifyToken(
	token string,
) (*jwt.UserClaims, error) {
	return m.VerifyTokenFunc(token)
}

type mockAuthRepository struct {
	CreateUserFunc               func(context.Context, domain.User) (domain.User, error)
	GetUserByIDFunc              func(context.Context, uuid.UUID) (domain.User, error)
	GetUserByLoginFunc           func(context.Context, domain.User) (domain.User, error)
	CreateSessionFunc            func(context.Context, domain.Session) (domain.Session, error)
	UpdateSessionFunc            func(context.Context, domain.Session) (domain.Session, error)
	GetSessionByRefreshTokenFunc func(context.Context, domain.Session) (domain.Session, error)
}

func (m *mockAuthRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	return m.CreateUserFunc(ctx, user)
}

func (m *mockAuthRepository) GetUserByID(
	ctx context.Context,
	userID uuid.UUID,
) (domain.User, error) {
	return m.GetUserByIDFunc(ctx, userID)
}

func (m *mockAuthRepository) GetUserByLogin(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	return m.GetUserByLoginFunc(ctx, user)
}

func (m *mockAuthRepository) CreateSession(
	ctx context.Context,
	session domain.Session,
) (domain.Session, error) {
	return m.CreateSessionFunc(ctx, session)
}

func (m *mockAuthRepository) UpdateSession(
	ctx context.Context,
	session domain.Session,
) (domain.Session, error) {
	return m.UpdateSessionFunc(ctx, session)
}

func (m *mockAuthRepository) GetSessionByRefreshToken(
	ctx context.Context,
	session domain.Session,
) (domain.Session, error) {
	return m.GetSessionByRefreshTokenFunc(ctx, session)
}
