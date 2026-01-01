package auth

import (
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/entity"
	"abbassmortazavi/go-microservice/services/auth-service/internal/infrastructure/db/repository"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/response"
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewJwtAuthenticator(j string, repo repository.TokenRepository) *JWT {
	JwtAuthenticator = &JWT{
		SigningKey:      []byte(j),
		TokenRepository: repo,
	}
	return JwtAuthenticator
}

type Claims struct {
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type JWT struct {
	SigningKey      []byte
	TokenRepository repository.TokenRepository
}

func (j *JWT) GenerateToken(userID int, name string) (response.TokenResponse, error) {
	accessExpiry := time.Now().Add(time.Minute * 5)
	refreshExpiry := time.Now().Add(time.Minute * 10)
	claims := Claims{
		UserID:    userID,
		Name:      name,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
		},
	}

	refreshExpiryClaims := &Claims{
		UserID:    userID,
		Name:      name,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString(j.SigningKey)
	if err != nil {
		return response.TokenResponse{}, err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshExpiryClaims)
	refreshTokenString, err := refreshToken.SignedString(j.SigningKey)
	if err != nil {
		return response.TokenResponse{}, err
	}

	ctx := context.Background()
	res, err := j.TokenRepository.FindByUserId(ctx, userID)
	if err != nil {
		log.Println(err)
		return response.TokenResponse{}, err
	}
	if res != nil {
		//delete all current user tokens
		err := j.TokenRepository.RevokeAllUserTokens(ctx, userID)
		if err != nil {
			log.Println(err)
			return response.TokenResponse{}, err
		}
	}

	reqAccessToken := entity.Token{
		UserID:    userID,
		TokenType: "access",
		HashToken: accessTokenString,
		ExpiredAt: accessExpiry,
		IsRevoked: false,
	}
	err = j.TokenRepository.Create(ctx, &reqAccessToken)
	if err != nil {
		log.Println(err)
		return response.TokenResponse{}, err
	}
	reqRefreshToken := entity.Token{
		UserID:    userID,
		TokenType: "refresh",
		HashToken: refreshTokenString,
		ExpiredAt: refreshExpiry,
		IsRevoked: false,
	}
	err = j.TokenRepository.Create(ctx, &reqRefreshToken)
	if err != nil {
		log.Println(err)
		return response.TokenResponse{}, err
	}

	return response.TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessExpiry.Unix(),
	}, nil

}
