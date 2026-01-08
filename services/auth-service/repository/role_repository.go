package repository

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"context"
	"database/sql"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}
func (r *RoleRepository) Save(ctx context.Context, role entity.Role) error {
	query := `INSERT INTO roles (name) VALUES ($1) RETURNING role_id`
	row := r.db.QueryRowContext(ctx, query, role.Name)
	if err := row.Scan(&role.ID); err != nil {
		return err
	}
	return nil
}
func (r *RoleRepository) FindById(ctx context.Context, roleId int) (*entity.Role, error) {
	query := `SELECT * FROM roles WHERE role_id=$1`
	row := r.db.QueryRowContext(ctx, query, roleId)
	var role entity.Role
	if err := row.Scan(&role.ID, &role.Name); err != nil {
		return nil, err
	}
	return &role, nil
}
func (r *RoleRepository) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	query := `SELECT * FROM roles WHERE name=$1`
	row := r.db.QueryRowContext(ctx, query, name)
	var role entity.Role
	if err := row.Scan(&role.ID, &role.Name); err != nil {
		return nil, err
	}
	return &role, nil
}
func (r *RoleRepository) Lists(ctx context.Context) ([]entity.Role, error) {
	query := `SELECT * FROM roles`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []entity.Role
	for rows.Next() {
		var role entity.Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}
func (r *RoleRepository) Delete(ctx context.Context, roleId int) error {
	query := `DELETE FROM roles WHERE role_id=$1`
	_, err := r.db.ExecContext(ctx, query, roleId)
	if err != nil {
		return err
	}
	return nil
}
