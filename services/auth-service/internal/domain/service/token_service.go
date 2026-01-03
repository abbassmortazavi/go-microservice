package service

import (
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/entity"
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/repository_interface"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/response"
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenServiceInterface interface {
	GenerateToken(userID int, name string) (response.TokenResponse, error)
	RefreshAccessToken(refreshToken string) (response.TokenResponse, error)
	ValidateToken(token string) (*Claims, error)
}

var JwtAuthenticator *JWT

func NewJwtAuthenticator(j string, repo repository_interface.TokenRepositoryInterface) *JWT {
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
	TokenRepository repository_interface.TokenRepositoryInterface
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

func (j *JWT) RefreshAccessToken(refreshToken string) (response.TokenResponse, error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return response.TokenResponse{}, err
	}
	if claims.TokenType != "refresh" {
		return response.TokenResponse{}, errors.New("invalid token")
	}
	return j.GenerateToken(claims.UserID, claims.Name)
}

func (j *JWT) ValidateToken(token string) (*Claims, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if !tkn.Valid {
		return nil, errors.New("invalid token")
	}

	//check in database
	res, err := j.TokenRepository.FindByToken(context.Background(), token)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}
