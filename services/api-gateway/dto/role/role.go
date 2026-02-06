package role

type PermissionDTO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type RoleDTO struct {
	Id          int64           `json:"id"`
	Name        string          `json:"name"`
	Permissions []PermissionDTO `json:"permissions"`
}

type GetRoleResponseDTO struct {
	Role RoleDTO `json:"role"`
}
