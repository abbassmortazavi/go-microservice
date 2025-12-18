package db

import (
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/entity"
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/repository"
	"context"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &UserRepository{
		db: db,
	}
}
func (u *UserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, email, password, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := u.db.ExecContext(ctx, query, user.ID, user.Email, user.Password, user.Role, user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}
