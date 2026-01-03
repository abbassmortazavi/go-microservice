package repository

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"abbassmortazavi/go-microservice/services/auth-service/interfaces/repository_interface"
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
	query := `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)`
	_, err := u.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, name, email, role, password, created_at 
              FROM users WHERE email = $1`
	row := u.db.QueryRowContext(ctx, query, email)
	var user entity.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
func (u *UserRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	log.Println("finding user by id in repo", id)
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`
	var user entity.User
	err := u.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
