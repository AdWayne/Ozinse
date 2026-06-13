package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
	jwt.RegisteredClaims
}

type Service struct {
	accessSecret  []byte
	refreshSecret []byte
}

func NewService(accessSecret, refreshSecret string) *Service {
	return &Service{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

func (s *Service) GenerateAccessToken(userID, roleID int) (string, error) {
	claims := Claims{
		UserID: userID,
		RoleID: roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.accessSecret)
}

func (s *Service) GenerateRefreshToken(userID int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.refreshSecret)
}

func (s *Service) ValidateAccessToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return s.accessSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func (s *Service) ValidateRefreshToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return s.refreshSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}