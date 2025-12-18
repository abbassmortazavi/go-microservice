package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	GenerateAccessToken(userID, role string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(token string) (*jwt.RegisteredClaims, error)
}

type JWTSecret struct {
	secret []byte
}

type Claims struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func NewJWTSecret(secret []byte) *JWTSecret {
	return &JWTSecret{secret: secret}
}
func (s *JWTSecret) GenerateAccessToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(time.Minute * 5).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *JWTSecret) GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}
func (s *JWTSecret) ValidateToken(token string) (*jwt.RegisteredClaims, error) {
	tkn, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tkn.Claims.(*jwt.RegisteredClaims)
	if !ok || !tkn.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
