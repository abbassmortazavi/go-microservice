package repository

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"abbassmortazavi/go-microservice/services/auth-service/interfaces/repository_interface"
	"context"
	"database/sql"
	"errors"
)

type TokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) repository_interface.TokenRepositoryInterface {
	return &TokenRepository{
		db: db,
	}
}

func (r *TokenRepository) FindByUserId(ctx context.Context, userId int) (*entity.Token, error) {
	token := &entity.Token{
		User: entity.User{},
	}

	query := `
	SELECT
		t.id,
		t.user_id,
		t.token_type,
		t.hash_token,
		t.expired_at,
		t.is_revoked,
		t.created_at,
		t.updated_at,

		u.id,
		u.name,
		u.email,
		u.created_at,
		u.updated_at
	FROM tokens t
	JOIN users u ON u.id = t.user_id
	WHERE t.user_id = $1
	LIMIT 1
	`
	err := r.db.QueryRowContext(ctx, query, userId).Scan(
		// token
		&token.ID,
		&token.UserID,
		&token.TokenType,
		&token.HashToken,
		&token.ExpiredAt,
		&token.IsRevoked,
		&token.CreatedAt,
		&token.UpdatedAt,

		// user
		&token.User.ID,
		&token.User.Name,
		&token.User.Email,
		&token.User.CreatedAt,
		&token.User.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return token, nil
}
func (r *TokenRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tokens WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
func (r *TokenRepository) RevokeAllUserTokens(ctx context.Context, userId int) error {
	query := `DELETE FROM tokens WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}
	return nil
}
func (r *TokenRepository) FindById(ctx context.Context, id int) (*entity.Token, error) {
	token := &entity.Token{}
	query := `SELECT * FROM tokens WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenType,
		&token.HashToken,
		&token.ExpiredAt,
		&token.IsRevoked,
		&token.CreatedAt,
		&token.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return token, nil
}
func (r *TokenRepository) Create(ctx context.Context, token *entity.Token) error {
	query := `INSERT INTO tokens (user_id, token_type, hash_token, expired_at, is_revoked)
	values ($1, $2, $3, $4, $5) RETURNING *`
	row := r.db.QueryRowContext(ctx, query, token.UserID, token.TokenType, token.HashToken, token.CreatedAt, token.IsRevoked)
	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.TokenType,
		&token.HashToken,
		&token.ExpiredAt,
		&token.IsRevoked,
		&token.CreatedAt,
		&token.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
func (r *TokenRepository) FindByToken(ctx context.Context, token string) (*entity.Token, error) {
	tokenModel := &entity.Token{}
	query := `SELECT * FROM tokens WHERE hash_token = $1`
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&tokenModel.ID,
		&tokenModel.UserID,
		&tokenModel.TokenType,
		&tokenModel.HashToken,
		&tokenModel.ExpiredAt,
		&tokenModel.IsRevoked,
		&tokenModel.CreatedAt,
		&tokenModel.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return tokenModel, nil
}
func (r *TokenRepository) Revoke(ctx context.Context, token string) error {
	query := `DELETE FROM tokens WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}
	return nil
}
