package service

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"abbassmortazavi/go-microservice/services/auth-service/interfaces/repository_interface"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/response"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenServiceInterface interface {
	GenerateToken(user *entity.User) (response.TokenResponse, error)
	RefreshAccessToken(refreshToken string) (response.TokenResponse, error)
	ValidateToken(token string) (*Claims, error)
	FindByToken(token string) (*entity.Token, error)
}

var JwtAuthenticator *JWT

func NewJwtAuthenticator(j string, repo repository_interface.TokenRepositoryInterface, userRepo repository_interface.UserRepositoryInterface) *JWT {
	JwtAuthenticator = &JWT{
		SigningKey:      []byte(j),
		TokenRepository: repo,
		UserRepository:  userRepo,
	}
	return JwtAuthenticator
}

type Claims struct {
	User      entity.User `json:"user"`
	Name      string      `json:"name"`
	TokenType string      `json:"token_type"`
	jwt.RegisteredClaims
}

type JWT struct {
	SigningKey      []byte
	TokenRepository repository_interface.TokenRepositoryInterface
	UserRepository  repository_interface.UserRepositoryInterface
}

func (j *JWT) GenerateToken(user *entity.User) (response.TokenResponse, error) {
	accessExpiry := time.Now().Add(time.Minute * 20)
	refreshExpiry := time.Now().Add(time.Minute * 30)
	now := time.Now()
	ctx := context.Background()
	//user, err := j.FindByUserId(ctx, userID)
	userInfo := entity.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	claims := Claims{
		User:      userInfo,
		Name:      user.Name,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
			Subject:   strconv.FormatInt(user.ID, 10),
		},
	}

	refreshExpiryClaims := &Claims{
		User:      userInfo,
		Name:      user.Name,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
			Subject:   strconv.FormatInt(user.ID, 10),
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

	res, err := j.TokenRepository.FindByUserId(ctx, user.ID)
	if err != nil {
		return response.TokenResponse{}, err
	}
	if res != nil {
		//delete all current user tokens
		//err := j.TokenRepository.RevokeAllUserTokens(ctx, userID)
		//if err != nil {
		//	return response.TokenResponse{}, err
		//}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err := j.TokenRepository.RevokeAllUserTokens(ctx, user.ID)
			if err != nil {
				return
			}
		}()
	}

	reqAccessToken := entity.Token{
		UserID:    user.ID,
		TokenType: "access",
		HashToken: accessTokenString,
		ExpiredAt: accessExpiry,
		IsRevoked: false,
	}
	err = j.TokenRepository.Create(ctx, &reqAccessToken)
	if err != nil {
		return response.TokenResponse{}, err
	}
	reqRefreshToken := entity.Token{
		UserID:    user.ID,
		TokenType: "refresh",
		HashToken: refreshTokenString,
		ExpiredAt: refreshExpiry,
		IsRevoked: false,
	}
	err = j.TokenRepository.Create(ctx, &reqRefreshToken)
	if err != nil {
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
	user := entity.User{
		ID:    claims.User.ID,
		Name:  claims.User.Name,
		Email: claims.User.Email,
		Role:  claims.User.Role,
	}
	return j.GenerateToken(&user)
}

func (j *JWT) ValidateToken(token string) (*Claims, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
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

func (j *JWT) FindByUserId(ctx context.Context, userId int64) (*entity.User, error) {
	user, err := j.UserRepository.FindByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (j *JWT) FindByToken(token string) (*entity.Token, error) {
	res, err := j.TokenRepository.FindByToken(context.Background(), token)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("token is invalid")
	}
	return res, nil
}
