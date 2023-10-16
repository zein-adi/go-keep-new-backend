package auth_services

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/helpers"
)

func NewPermissionServices(permissionRepo auth_repo_interfaces.IPermissionRepository, roleRepo auth_repo_interfaces.IRoleRepository) *PermissionServices {
	return &PermissionServices{
		repo:     permissionRepo,
		roleRepo: roleRepo,
	}
}

type PermissionServices struct {
	repo     auth_repo_interfaces.IPermissionRepository
	roleRepo auth_repo_interfaces.IRoleRepository
}

func (x *PermissionServices) Get(ctx context.Context, roleIds []string) []string {
	roleDeveloperId := "1"
	_, notFoundError := helpers.FindIndex(roleIds, func(s string) bool {
		return s == roleDeveloperId
	})
	if notFoundError == nil {
		return x.repo.Get(ctx)
	}

	validPermissions := x.repo.Get(ctx)
	validPermissionsMap := helpers.KeyBy(validPermissions, func(permission string) string {
		return permission
	})
	roles, _ := x.roleRepo.GetById(ctx, roleIds)
	var userPermissions []string
	for _, role := range roles {
		userPermissions = append(userPermissions, role.Permissions...)
	}
	userPermissions = helpers.Unique(userPermissions)

	var validUserPermissions []string
	for _, permission := range userPermissions {
		_, ok := validPermissionsMap[permission]
		if !ok {
			continue
		}
		validUserPermissions = append(validUserPermissions, permission)
	}
	return validUserPermissions
}
