package core_service_jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager interface {
	CreateToken(id uuid.UUID, role string, duration time.Duration) (string, *UserClaims, error)
	VerifyToken(tokenStr string) (*UserClaims, error)
}

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(cfg Config) *JWTMaker {
	return &JWTMaker{
		secretKey: cfg.SecretKey,
	}
}

func (jm *JWTMaker) CreateToken(
	userID uuid.UUID,
	role string,
	duration time.Duration,
) (string, *UserClaims, error) {
	claims, err := NewUserClaims(userID, role, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(jm.secretKey))
	if err != nil {
		return "", nil, fmt.Errorf("signing token: %w", err)
	}

	return tokenStr, claims, nil
}

func (jw *JWTMaker) VerifyToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}

		return []byte(jw.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
