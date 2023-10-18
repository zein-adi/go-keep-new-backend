package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_requests"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_services"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/repos/auth_repos_memory"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/repos/auth_repos_mysql"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_env"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"testing"
)

func TestRole(t *testing.T) {
	helpers_env.Init(5)
	r := RoleServicesTest{}
	r.setup()
	defer r.dbCleanup()

	t.Run("GetSuccess", func(t *testing.T) {
		ori := r.setupAndPopulate()

		req := auth_requests.NewGetRequest()
		models := r.services.Get(context.Background(), req)
		assert.Len(t, models, 3)

		for i := range models {
			assert.Equal(t, ori[i].Id, models[i].Id)
			assert.Equal(t, ori[i].Nama, models[i].Nama)
			assert.Equal(t, ori[i].Deskripsi, models[i].Deskripsi)
			assert.Equal(t, ori[i].Level, models[i].Level)
			assert.Equal(t, ori[i].Permissions, models[i].Permissions)
		}
	})
	t.Run("DeleteSuccess", func(t *testing.T) {
		ctx := context.Background()
		ori := r.setupAndPopulate()
		models := r.repo.Get(ctx, auth_requests.NewGetRequest())
		assert.Len(t, models, 3)
		currentRoleIds := []string{ori[0].Id}

		affected, err := r.services.DeleteById(ctx, ori[1].Id, currentRoleIds)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		models = r.repo.Get(ctx, auth_requests.NewGetRequest())
		assert.Len(t, models, 2)
		assert.Equal(t, ori[2].Id, models[1].Id)
	})
	t.Run("UpdateAndDeleteFailedNotFound", func(t *testing.T) {
		ctx := context.Background()
		r.setupAndPopulate()

		currentUserRoleIds := []string{"1"}

		affected, err := r.services.DeleteById(ctx, "4", currentUserRoleIds)
		assert.Equal(t, 0, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		_, err = r.services.Update(ctx, &auth_entities.Role{Id: "4"}, currentUserRoleIds)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
	t.Run("UpdateDeleteFailedCauseCurrentUserUnauthorizedToAccessRole", func(t *testing.T) {
		ctx := context.Background()
		ori := r.setupAndPopulate()

		currentUserRoleIds := []string{ori[2].Id}
		accessedId := ori[1].Id

		affected, err := r.services.DeleteById(ctx, accessedId, currentUserRoleIds)
		assert.Equal(t, 0, affected)
		assert.ErrorIs(t, err, auth_services.RoleAccessUnauthorizedError)

		input := &auth_entities.Role{
			Id:    accessedId,
			Nama:  "a",
			Level: 100,
		}
		_, err = r.services.Update(ctx, input, currentUserRoleIds)
		assert.ErrorIs(t, err, auth_services.RoleAccessUnauthorizedError)
	})
	t.Run("InsertSuccess", func(t *testing.T) {
		ori := r.setupAndPopulate()
		ctx := context.Background()
		models := r.repo.Get(ctx, auth_requests.NewGetRequest())
		assert.Len(t, models, 3)

		nama := "Staf 2"
		deskripsi := "deskripsi baru"
		level := 2
		permissions := []string{
			"user.user.get",
			"user.role.get",
		}

		input := &auth_entities.Role{
			Nama:        nama,
			Deskripsi:   deskripsi,
			Level:       level,
			Permissions: permissions,
		}
		currentUserRoleIds := []string{ori[0].Id}
		model, err := r.services.Insert(ctx, input, currentUserRoleIds)
		assert.Empty(t, input.Id)
		assert.Nil(t, err)

		assert.NotEmpty(t, model.Id)
		assert.Equal(t, nama, model.Nama)
		assert.Equal(t, deskripsi, model.Deskripsi)
		assert.Equal(t, level, model.Level)
		assert.Equal(t, permissions, model.Permissions)

		models = r.repo.Get(ctx, auth_requests.NewGetRequest())
		assert.Len(t, models, 4)
		model = models[3]

		assert.NotEmpty(t, model.Id)
		assert.Equal(t, nama, model.Nama)
		assert.Equal(t, deskripsi, model.Deskripsi)
		assert.Equal(t, level, model.Level)
		assert.Equal(t, permissions, model.Permissions)
	})
	t.Run("UpdateSuccess", func(t *testing.T) {
		ori := r.setupAndPopulate()
		ctx := context.Background()
		models := r.repo.Get(ctx, auth_requests.NewGetRequest())
		assert.Len(t, models, 3)

		currentUserRoleIds := []string{ori[0].Id}
		id := ori[2].Id
		nama := "Staf 2"
		deskripsi := "deskripsi baru"
		level := 2
		permissions := []string{
			"user.user.get",
			"user.role.get",
		}

		input := &auth_entities.Role{
			Id:          id,
			Nama:        nama,
			Deskripsi:   deskripsi,
			Level:       level,
			Permissions: permissions,
		}
		model, err := r.services.Update(ctx, input, currentUserRoleIds)
		assert.Nil(t, err)
		assert.Equal(t, id, model.Id)
		assert.Equal(t, nama, model.Nama)
		assert.Equal(t, deskripsi, model.Deskripsi)
		assert.Equal(t, level, model.Level)
		assert.Equal(t, permissions, model.Permissions)

		model, err = r.repo.FindById(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, id, model.Id)
		assert.Equal(t, nama, model.Nama)
		assert.Equal(t, deskripsi, model.Deskripsi)
		assert.Equal(t, level, model.Level)
		assert.Equal(t, permissions, model.Permissions)
	})
	t.Run("InsertUpdateDeleteFailedRoleIdNotFound", func(t *testing.T) {
		ctx := context.Background()
		ori := r.setupAndPopulate()

		id := ori[1].Id
		input := &auth_entities.Role{
			Id:    id,
			Nama:  "a",
			Level: 100,
		}

		currentUserRoleIds := []string{"100"}
		affected, err := r.services.DeleteById(ctx, id, currentUserRoleIds)
		assert.Equal(t, 0, affected)
		assert.ErrorIs(t, err, helpers_error.EntryCountMismatchError)

		_, err = r.services.Insert(ctx, input, currentUserRoleIds)
		assert.ErrorIs(t, err, helpers_error.EntryCountMismatchError)

		_, err = r.services.Update(ctx, input, currentUserRoleIds)
		assert.ErrorIs(t, err, helpers_error.EntryCountMismatchError)
	})
	t.Run("InsertUpdateFailedCauseValidationErrorRequestedLevelLowerThanEqualUserLevel", func(t *testing.T) {
		ctx := context.Background()
		ori := r.setupAndPopulate()

		id := ori[2].Id
		nama := "User"
		level := 15

		currentUserRoleIds := []string{id}
		input := &auth_entities.Role{
			Id:    id,
			Nama:  nama,
			Level: level,
		}
		_, err := r.services.Insert(ctx, input, currentUserRoleIds)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, helpers_error.ValidationError)

		_, err = r.services.Update(ctx, input, currentUserRoleIds)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, helpers_error.ValidationError)
	})
	t.Run("InsertUpdateFailedCauseRequestedPermissionsOutsideUserPermissions", func(t *testing.T) {
		ctx := context.Background()
		ori := r.setupAndPopulate()

		currentUserRoleIds := []string{ori[1].Id}
		id := ori[2].Id
		nama := "Guru 2"
		level := 100
		permissions := []string{
			"user.permission.get",
		}

		input := &auth_entities.Role{
			Id:          id,
			Nama:        nama,
			Level:       level,
			Permissions: permissions,
		}
		_, err := r.services.Insert(ctx, input, currentUserRoleIds)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, helpers_error.ValidationError)

		_, err = r.services.Update(ctx, input, currentUserRoleIds)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, helpers_error.ValidationError)
	})
	t.Run("InsertAndUpdateFailedValidation", func(t *testing.T) {
		ctx := context.Background()
		id := r.setupAndPopulate()[0].Id

		currentUserRoleIds := []string{id}
		input := &auth_entities.Role{}
		_, err := r.services.Insert(ctx, input, currentUserRoleIds)

		assert.ErrorIs(t, err, helpers_error.ValidationError)
		assert.ErrorContains(t, err, "nama.required")
		assert.ErrorContains(t, err, "level.min.1")

		input.Id = id
		input.Level = -1
		_, err = r.services.Update(ctx, input, currentUserRoleIds)
		assert.ErrorIs(t, err, helpers_error.ValidationError)
		assert.ErrorContains(t, err, "level.min.1")
	})
}

type RoleServicesTest struct {
	repo      auth_repo_interfaces.IRoleRepository
	services  auth_service_interfaces.IRoleServices
	dbCleanup func()
}

func (r *RoleServicesTest) setup() {
	r.setMemoryRepository()
	r.services = auth_services.NewRoleServices(r.repo)
}
func (r *RoleServicesTest) setMemoryRepository() {
	r.repo = auth_repos_memory.NewRoleMemoryRepository()
	r.dbCleanup = func() {
	}
}
func (r *RoleServicesTest) setMysqlRepository() {
	repo := auth_repos_mysql.NewRoleMysqlRepository()
	r.dbCleanup = repo.Cleanup
	r.repo = repo
	models := r.repo.Get(context.Background(), auth_requests.NewGetRequest())
	for _, model := range models {
		_, _ = r.repo.DeleteById(context.Background(), model.Id)
	}
}
func (r *RoleServicesTest) setupAndPopulate() []*auth_entities.Role {
	r.setup()
	input := []*auth_entities.Role{
		{
			Nama:      "Aa Developer",
			Deskripsi: "Deskripsi role developer",
			Level:     1,
			Permissions: []string{
				"user.user.get",
				"user.user.insert",
				"user.user.update",
				"user.user.delete",
				"user.role.get",
				"user.role.insert",
				"user.role.update",
				"user.role.delete",
				"user.permission.get",
			},
		},
		{
			Nama:      "Admin IT Sekolah",
			Deskripsi: "Deskripsi role admin IT sekolah",
			Level:     10,
			Permissions: []string{
				"user.user.get",
				"user.user.insert",
				"user.user.update",
				"user.user.delete",
				"user.role.get",
				"user.role.insert",
				"user.role.update",
				"user.role.delete",
			},
		},
		{
			Nama:      "Staf",
			Deskripsi: "Deskripsi staf",
			Level:     15,
			Permissions: []string{
				"user.user.get",
			},
		},
	}

	var models []*auth_entities.Role
	for _, datum := range input {
		model, _ := r.repo.Insert(context.Background(), datum)
		models = append(models, model)
	}
	return models
}
