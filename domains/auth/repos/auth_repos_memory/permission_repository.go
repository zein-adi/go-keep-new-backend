package auth_repos_memory

import "context"

func NewPermissionMemoryRepository() *PermissionMemoryRepository {
	return &PermissionMemoryRepository{}
}

type PermissionMemoryRepository struct {
	data []string
}

func (x *PermissionMemoryRepository) Get(_ context.Context) []string {
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

		"keep.pos.get",
		"keep.pos.insert",
		"keep.pos.update",
		"keep.pos.delete",
		"keep.pos.trash",
	}
}
