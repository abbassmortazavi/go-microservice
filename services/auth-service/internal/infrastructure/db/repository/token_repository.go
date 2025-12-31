package repository

import (
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/entity"
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/repository_interface"
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

func (r *TokenRepository) FindByUserId(ctx context.Context, userId int64) (*entity.Token, error) {
	token := &entity.Token{}
	query := `SELECT * FROM tokens WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, userId).Scan(
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
func (r *TokenRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tokens WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
func (r *TokenRepository) RevokeAllUserTokens(ctx context.Context, userId int64) error {
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
	/*query := `INSERT INTO tokens (user_id, token_type, hash_token, expired_at, is_revoked)
	values ($1, $2, $3, $4, $5) RETURNING *`*/
	return nil
}
func (r *TokenRepository) FindByToken(ctx context.Context, token string) (*entity.Token, error) {
	return nil, nil
}
func (r *TokenRepository) Revoke(ctx context.Context, token string) error {
	return nil
}
