package repository_interface

import "context"

type RBACRepository interface {
	AssignRoleToUser(ctx context.Context, userId, roleId int64) error
	AssignPermissionToRole(ctx context.Context, permissionID, roleID int64) error
	GetPermissionsByUserID(ctx context.Context, userID int64) ([]string, error)
}
