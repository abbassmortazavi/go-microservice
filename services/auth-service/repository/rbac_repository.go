package repository

import (
	"context"
	"database/sql"
	"log"
)

type RBACRepository struct {
	db *sql.DB
}

func NewRBACRepository(db *sql.DB) *RBACRepository {
	return &RBACRepository{
		db: db,
	}
}

func (r *RBACRepository) AssignRoleToUser(ctx context.Context, userId, roleId int64) error {
	query := `INSERT INTO user_roles (role_id, user_id) VALUES ($1, $2) RETURNING *;`
	row := r.db.QueryRow(query, roleId, userId)
	err := row.Scan(&roleId, &userId)
	if err != nil {
		return err
	}
	return nil
}
func (r *RBACRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionID int64) error {
	log.Println(permissionID, roleID)
	query := `insert into role_permissions (role_id, permission_id) VALUES ($1, $2) RETURNING *;`
	row := r.db.QueryRow(query, roleID, permissionID)
	err := row.Scan(&roleID, &permissionID)
	if err != nil {
		return err
	}
	return nil
}
func (r *RBACRepository) GetPermissionsByUserID(ctx context.Context, userID int64) ([]string, error) {
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
func (r *RBACRepository) CheckUserHasRole(ctx context.Context, roleID, userID int64) (bool, error) {
	q := "select * from user_roles where role_id = $1 and user_id = $2;"
	row := r.db.QueryRow(q, roleID, userID)
	err := row.Scan(&roleID, &userID)
	if err != nil {
		log.Println(555555)
		return false, err
	}

	return true, nil

}
