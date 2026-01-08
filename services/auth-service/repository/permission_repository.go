package repository

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"context"
	"database/sql"
)

type PermissionRepository struct {
	db *sql.DB
}

func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{
		db: db,
	}
}
func (p *PermissionRepository) Save(ctx context.Context, permission entity.Permission) error {
	query := `insert into permissions ( name) values ($1) returning id, name`
	_, err := p.db.ExecContext(ctx, query, permission.Name)
	if err != nil {
		return err
	}
	return nil
}
func (p *PermissionRepository) FindByID(ctx context.Context, permissionId int) (entity.Permission, error) {
	query := `select * from permissions where id = $1`
	row := p.db.QueryRowContext(ctx, query, permissionId)
	var permission entity.Permission
	err := row.Scan(&permission.ID, &permission.Name)
	if err != nil {
		return entity.Permission{}, err
	}
	return permission, nil
}
func (p *PermissionRepository) FindByName(ctx context.Context, name string) (entity.Permission, error) {
	query := `select * from permissions where name = $1`
	row := p.db.QueryRowContext(ctx, query, name)
	var permission entity.Permission
	err := row.Scan(&permission.ID, &permission.Name)
	if err != nil {
		return entity.Permission{}, err
	}
	return permission, nil
}
func (p *PermissionRepository) Lists(ctx context.Context) ([]entity.Permission, error) {
	query := `select * from permissions`
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var permissions []entity.Permission
	for rows.Next() {
		var permission entity.Permission
		err := rows.Scan(&permission.ID, &permission.Name)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}
func (p *PermissionRepository) Delete(ctx context.Context, permissionId int) error {
	query := `delete from permissions where id = $1`
	_, err := p.db.ExecContext(ctx, query, permissionId)
	if err != nil {
		return err
	}
	return nil
}
