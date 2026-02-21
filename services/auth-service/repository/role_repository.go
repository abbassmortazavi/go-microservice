package repository

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}
func (r *RoleRepository) Save(ctx context.Context, role entity.Role) (*entity.Role, error) {
	query := `INSERT INTO roles (name) VALUES ($1) returning id, name`
	var savedRole entity.Role
	err := r.db.QueryRowContext(ctx, query, role.Name).Scan(&savedRole.ID, &savedRole.Name)
	if err != nil {
		return nil, err
	}
	return &savedRole, nil
}
func (r *RoleRepository) FindById(ctx context.Context, roleId int64) (*entity.Role, error) {
	query := `
	        SELECT
	            r.id,
	            r.name,
	            COALESCE(
	                JSON_AGG(
	                    JSON_BUILD_OBJECT(
	                        'id', p.id,
	                        'name', p.name
	                    )
	                ) FILTER (WHERE p.id IS NOT NULL),
	                '[]'::JSON
	            ) as permissions
	        FROM roles r
	        LEFT JOIN role_permissions rp ON r.id = rp.role_id
	        LEFT JOIN permissions p ON rp.permission_id = p.id
	        WHERE r.id = $1
	        GROUP BY r.id, r.name
	    `

	row := r.db.QueryRowContext(ctx, query, roleId)
	var role entity.Role
	var permissionsJSON []byte

	err := row.Scan(&role.ID, &role.Name, &permissionsJSON)
	if err != nil {
		return nil, err
	}
	role.Permissions = []entity.Permission{}

	if len(permissionsJSON) > 0 && string(permissionsJSON) != "[]" {
		if err := json.Unmarshal(permissionsJSON, &role.Permissions); err != nil {
			log.Printf("Failed to unmarshal permissions: %v", err)
		}
	}
	return &role, nil

	/*roleQuery := `SELECT id, name FROM roles WHERE id = $1`
		var role entity.Role
		err := r.db.QueryRowContext(ctx, roleQuery, roleId).Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, err
		}

		// Query 2: Get permissions separately
		permQuery := `
	        SELECT p.id, p.name
	        FROM permissions p
	        JOIN role_permissions rp ON p.id = rp.permission_id
	        WHERE rp.role_id = $1
	    `

		rows, err := r.db.QueryContext(ctx, permQuery, roleId)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var permissions []entity.Permission
		for rows.Next() {
			var perm entity.Permission
			_ = rows.Scan(&perm.ID, &perm.Name)
			permissions = append(permissions, perm)
		}

		role.Permissions = permissions
		return &role, nil*/
}
func (r *RoleRepository) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	query := `SELECT id,name FROM roles WHERE name=$1`
	row := r.db.QueryRowContext(ctx, query, name)
	var role entity.Role
	if err := row.Scan(&role.ID, &role.Name); err != nil {
		return nil, err
	}
	return &role, nil
}
func (r *RoleRepository) Lists(ctx context.Context, page, perPage int64, orderBy, sortBy, search string) ([]entity.Role, entity.PaginationMeta, error) {
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

	query := fmt.Sprintf(`
    SELECT id, name 
    FROM roles 
    WHERE ($1 = '' OR name ILIKE '%%' || $1 || '%%')
    ORDER BY %s %s 
    LIMIT $2 OFFSET $3
`, orderBy, sortBy)

	rows, err := r.db.QueryContext(ctx, query, searchTerm, perPage, offset)
	if err != nil {
		return nil, entity.PaginationMeta{}, err
	}
	defer rows.Close()
	var roles []entity.Role
	for rows.Next() {
		var role entity.Role
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, entity.PaginationMeta{}, err
		}
		roles = append(roles, role)
	}

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM permissions`
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
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

	return roles, paginationMeta, nil
}
func (r *RoleRepository) Delete(ctx context.Context, roleId int64) error {
	query := `DELETE FROM roles WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, roleId)
	if err != nil {
		return err
	}
	return nil
}
