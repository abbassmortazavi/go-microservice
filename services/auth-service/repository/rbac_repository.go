package repository

import (
	"context"
	"database/sql"
)

type RBACRepository struct {
	db *sql.DB
}

func NewRBACRepository(db *sql.DB) *RBACRepository {
	return &RBACRepository{
		db: db,
	}
}

func (r *RBACRepository) AssignRoleToUser(ctx context.Context, userId, roleId int) error {
	query := `INSERT INTO user_roles (role_id, user_id) VALUES ($1, $2) RETURNING *;`
	row := r.db.QueryRow(query, roleId, userId)
	err := row.Scan(&roleId, &userId)
	if err != nil {
		return err
	}
	return nil
}
func (r *RBACRepository) AssignPermissionToRole(ctx context.Context, permissionID, roleId int) error {
	query := `insert into role_permissions (role_id, permission_id) VALUES ($1, $2) RETURNING *;`
	row := r.db.QueryRow(query, roleId, permissionID)
	err := row.Scan(&roleId, &permissionID)
	if err != nil {
		return err
	}
	return nil
}
func (r *RBACRepository) GetPermissionsByUserID(ctx context.Context, userID int) ([]string, error) {
	query := `SELECT DISTINCT p.name From permissions p 
        join role_permissions rp on rp.permission_id = p.id
        join user_roles ur on ur.role_id = rp.role_id
         WHERE user_id = $1;`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var permissions []string
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}
