package db

import (
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/entity"
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/repository_interface"
	"context"
	"database/sql"
	"errors"
	"log"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository_interface.UserRepositoryInterface {
	return &UserRepository{
		db: db,
	}
}
func (u *UserRepository) Create(ctx context.Context, user *entity.User) error {
	log.Println("creating user with repo")
	query := `INSERT INTO users (id, email, password, role, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := u.db.ExecContext(ctx, query, user.ID, user.Email, user.Password, user.Role, user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	row := u.db.QueryRowContext(ctx, query, email)
	var user entity.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
