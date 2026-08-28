package auth_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/romreign/AuthService/internal/core/domain"
	core_service_jwt "github.com/romreign/AuthService/internal/core/service/jwt"
)

type AuthRepositoryPostgres interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error)
	GetUserByLogin(ctx context.Context, user domain.User) (domain.User, error)

	CreateSession(ctx context.Context, session domain.Session) (domain.Session, error)
	UpdateSession(ctx context.Context, session domain.Session) (domain.Session, error)
	GetSessionByRefreshToken(ctx context.Context, session domain.Session) (domain.Session, error)
}

type AuthRepositoryRedis interface {
	SaveAccessToken(ctx context.Context, token string, data AuthData, ttl time.Duration) error
	SaveIdempotencyKey(ctx context.Context, idemKey string, resp IdempotencyResponse, ttl time.Duration) error
	GetIdempotencyKey(ctx context.Context, idemKey string) (*IdempotencyResponse, error)
}

type AuthService struct {
	authRepositoryPostgres AuthRepositoryPostgres
	authRepositoryRedis    AuthRepositoryRedis
	jwtManager             core_service_jwt.JWTManager
}

func NewAuthService(
	authRepositoryPostgres AuthRepositoryPostgres,
	authRepositoryRedis AuthRepositoryRedis,
	jwtManager core_service_jwt.JWTManager,
) *AuthService {
	return &AuthService{
		authRepositoryPostgres: authRepositoryPostgres,
		authRepositoryRedis:    authRepositoryRedis,
		jwtManager:             jwtManager,
	}
}
