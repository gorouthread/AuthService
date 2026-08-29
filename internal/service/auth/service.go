package auth_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/romreign/AuthService/internal/core/domain"
	jwt "github.com/romreign/AuthService/pkg/jwt"
)

type AuthRepositoryPostgres interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error)
	GetUserByLogin(ctx context.Context, user domain.User) (domain.User, error)

	CreateSession(ctx context.Context, session domain.Session) (domain.Session, error)
	UpdateSession(ctx context.Context, session domain.Session) (domain.Session, error)
	GetSessionByRefreshToken(ctx context.Context, session domain.Session) (domain.Session, error)
}

type AuthService struct {
	authRepositoryPostgres AuthRepositoryPostgres
	jwtManager             jwt.JWTManager
}

func NewAuthService(
	authRepositoryPostgres AuthRepositoryPostgres,
	jwtManager jwt.JWTManager,
) *AuthService {
	return &AuthService{
		authRepositoryPostgres: authRepositoryPostgres,
		jwtManager:             jwtManager,
	}
}
