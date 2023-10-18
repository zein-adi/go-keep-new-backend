package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_services"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/repos/auth_repos_memory"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_env"
	"testing"
)

func TestPermission(t *testing.T) {
	helpers_env.Init(5)
	r := PermissionServicesTest{}
	r.setup()

	t.Run("GetPermissionForDeveloper", func(t *testing.T) {
		r.setupAndPopulate()
		permissions := r.services.Get(context.Background(), []string{"1"})
		assert.Len(t, permissions, len(r.permissionRepo.Get(context.Background())))
	})
	t.Run("GetPermissionFor2", func(t *testing.T) {
		r.setupAndPopulate()
		permissions := r.services.Get(context.Background(), []string{"2"})
		assert.Len(t, permissions, 3)
	})
	t.Run("GetPermissionFor3", func(t *testing.T) {
		r.setupAndPopulate()
		permissions := r.services.Get(context.Background(), []string{"3"})
		assert.Len(t, permissions, 2)
	})
	t.Run("GetPermissionFor2And3", func(t *testing.T) {
		r.setupAndPopulate()
		permissions := r.services.Get(context.Background(), []string{"2", "3"})
		assert.Len(t, permissions, 5)
	})
}

type PermissionServicesTest struct {
	permissionRepo auth_repo_interfaces.IPermissionRepository
	roleRepo       auth_repo_interfaces.IRoleRepository
	services       auth_service_interfaces.IPermissionServices
}

func (r *PermissionServicesTest) setup() {
	r.setMemoryRepository()
	r.services = auth_services.NewPermissionServices(r.permissionRepo, r.roleRepo)
}
func (r *PermissionServicesTest) setMemoryRepository() {
	r.permissionRepo = auth_repos_memory.NewPermissionMemoryRepository()
	r.roleRepo = auth_repos_memory.NewRoleMemoryRepository()
}
func (r *PermissionServicesTest) setupAndPopulate() []*auth_entities.Role {
	r.setup()
	input := []*auth_entities.Role{
		{
			Id:          "1", // Developer id must be 1
			Nama:        "Developer",
			Permissions: []string{},
		},
		{
			Id: "2",
			Permissions: []string{
				"user.permission.get",
				"user.role.get",
				"user.user.get",
			},
		},
		{
			Id: "3",
			Permissions: []string{
				"user.permission.insert",
				"user.role.insert",
				"user.user.insert",
			},
		},
	}

	var models []*auth_entities.Role
	for _, datum := range input {
		model, _ := r.roleRepo.Insert(context.Background(), datum)
		models = append(models, model)
	}
	return models
}
