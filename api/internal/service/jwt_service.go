package service

import (
	"time"

	"github.com/kemenpppa/siarsip-api/internal/config"
	"github.com/kemenpppa/siarsip-api/internal/middleware"
	"gorm.io/gorm"
)

// JWTService mengelola operasi JWT.
type JWTService interface {
	GenerateToken(userID uint, username, nip, role string, rememberMe bool) (string, int64, error)
	BlacklistToken(jti string, expiresAt int64) error
}

type jwtServiceImpl struct {
	authMiddleware *middleware.AuthMiddleware
}

// NewJWTService membuat instance JWTService.
func NewJWTService(cfg *config.Config, db *gorm.DB) JWTService {
	return &jwtServiceImpl{
		authMiddleware: middleware.NewAuthMiddleware(cfg, db),
	}
}

func (s *jwtServiceImpl) GenerateToken(userID uint, username, nip, role string, rememberMe bool) (string, int64, error) {
	token, expiresAt, err := s.authMiddleware.GenerateToken(userID, username, nip, role, rememberMe)
	if err != nil {
		return "", 0, err
	}
	return token, expiresAt.Unix(), nil
}

func (s *jwtServiceImpl) BlacklistToken(jti string, expiresAt int64) error {
	t := time.Unix(expiresAt, 0)
	return s.authMiddleware.BlacklistToken(jti, t)
}
