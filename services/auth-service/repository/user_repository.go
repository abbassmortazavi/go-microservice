package repository

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"abbassmortazavi/go-microservice/services/auth-service/interfaces/repository_interface"
	"context"
	"database/sql"
	"errors"
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
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`
	_, err := u.db.ExecContext(ctx, query, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, name, email, password, created_at 
              FROM users WHERE email = $1`
	row := u.db.QueryRowContext(ctx, query, email)
	var user entity.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
func (u *UserRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
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

	permQuery := `
	        SELECT r.id, r.name
	        FROM roles r
	        JOIN user_roles ur ON r.id = ur.role_id
	        WHERE ur.user_id = $1
	    `

	rows, err := u.db.QueryContext(ctx, permQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []entity.Role
	for rows.Next() {
		var role entity.Role
		_ = rows.Scan(&role.ID, &role.Name)
		roles = append(roles, role)
	}

	user.Role = roles

	return &user, nil
}
