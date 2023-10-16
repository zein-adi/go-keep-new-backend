package auth_services

import (
	"context"
)

func NewPermissionServices() *PermissionServices {
	return &PermissionServices{}
}

type PermissionServices struct {
}

func (p *PermissionServices) Get(_ context.Context) []string {
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
