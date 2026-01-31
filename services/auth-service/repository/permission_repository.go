package repository

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"context"
	"database/sql"
	"math"
	"strings"
)

type PermissionRepository struct {
	db *sql.DB
}

func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{
		db: db,
	}
}
func (p *PermissionRepository) Create(ctx context.Context, permission entity.Permission) (*entity.Permission, error) {
	query := `insert into permissions ( name) values ($1) returning id, name`
	var savedPermission entity.Permission
	err := p.db.QueryRowContext(ctx, query, permission.Name).Scan(&savedPermission.ID, &savedPermission.Name)
	if err != nil {
		return nil, err
	}
	return &savedPermission, nil
}
func (p *PermissionRepository) FindByID(ctx context.Context, permissionId int64) (entity.Permission, error) {
	query := `select id,name from permissions where id = $1`
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
func (p *PermissionRepository) Lists(ctx context.Context, page, perPage int64, orderBy, sortBy, search string) ([]entity.Permission, entity.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage
	if orderBy == "" {
		orderBy = "id"
	}

	sortBy = strings.ToLower(sortBy)
	if sortBy != "asc" && sortBy != "desc" {
		sortBy = "desc"
	}
	searchTerm := ""
	if search != "" {
		searchTerm = "%" + search + "%"
	}

	query := `select * from permissions where ($1 = '' or name ilike $1) order by $2 $3 limit $4 offset $5`
	rows, err := p.db.QueryContext(ctx, query, searchTerm, orderBy, sortBy, perPage, perPage, offset)
	if err != nil {
		return nil, entity.PaginationMeta{}, err
	}
	defer rows.Close()
	var permissions []entity.Permission
	for rows.Next() {
		var permission entity.Permission
		err := rows.Scan(&permission.ID, &permission.Name)
		if err != nil {
			return nil, entity.PaginationMeta{}, err
		}
		permissions = append(permissions, permission)
	}

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM permissions`
	err = p.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, entity.PaginationMeta{}, err
	}

	totalPages := int64(math.Ceil(float64(total) / float64(perPage)))
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	paginationMeta := entity.PaginationMeta{
		Page:        page,
		PerPage:     perPage,
		Total:       total,
		HasNextPage: page < totalPages,
		HasPrevPage: page > 1,
	}

	return permissions, paginationMeta, nil
}
func (p *PermissionRepository) Delete(ctx context.Context, permissionId int64) error {
	query := `delete from permissions where id = $1`
	_, err := p.db.ExecContext(ctx, query, permissionId)
	if err != nil {
		return err
	}
	return nil
}

func (p *PermissionRepository) Update(ctx context.Context, id int64, name string) (entity.Permission, error) {
	//TODO implement me
	panic("implement me")
}
