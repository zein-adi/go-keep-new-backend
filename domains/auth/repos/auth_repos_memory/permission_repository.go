package auth_repos_memory

import "context"

func NewPermissionMemoryRepository() *PermissionRepository {
	return &PermissionRepository{}
}

type PermissionRepository struct {
	data []string
}

func (x *PermissionRepository) Get(_ context.Context) []string {
	return []string{
		"user.permission.get",
		"user.role.get",
		"user.role.insert",
		"user.role.update",
		"user.role.delete",
		"user.user.get",
		"user.user.insert",
		"user.user.update",
		"user.user.update.password",
		"user.user.delete",
	}
}
