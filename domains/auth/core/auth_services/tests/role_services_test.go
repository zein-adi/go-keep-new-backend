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
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"testing"
)

func TestRole(t *testing.T) {
	r := RoleServicesTest{}
	r.setup()

	t.Run("GetSuccess", func(t *testing.T) {
		ori := r.setupAndPopulate()

		req := auth_requests.NewGetRequest()
		models := r.services.Get(context.Background(), req)
		assert.Len(t, models, 2)

		for i := range models {
			assert.Equal(t, ori[i].Id, models[i].Id)
			assert.Equal(t, ori[i].Nama, models[i].Nama)
			assert.Equal(t, ori[i].Deskripsi, models[i].Deskripsi)
			assert.Equal(t, ori[i].Level, models[i].Level)
			assert.Equal(t, ori[i].Permissions, models[i].Permissions)
		}
	})
	t.Run("InsertSuccess", func(t *testing.T) {
		r.setup()
		ctx := context.Background()
		models := r.repo.Get(ctx, auth_requests.NewGetRequest())
		assert.Len(t, models, 0)

		input := &auth_entities.Role{
			Nama:      "Developer",
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
			},
		}

		model, err := r.services.Insert(ctx, input)
		assert.Empty(t, input.Id)
		assert.Nil(t, err)

		assert.NotEmpty(t, model.Id)
		assert.Equal(t, input.Nama, model.Nama)
		assert.Equal(t, input.Deskripsi, model.Deskripsi)
		assert.Equal(t, input.Level, model.Level)
		assert.Equal(t, input.Permissions, model.Permissions)

		models = r.repo.Get(ctx, auth_requests.NewGetRequest())
		assert.Len(t, models, 1)
		model = models[0]

		assert.NotEmpty(t, model.Id)
		assert.Equal(t, input.Nama, model.Nama)
		assert.Equal(t, input.Deskripsi, model.Deskripsi)
		assert.Equal(t, input.Level, model.Level)
		assert.Equal(t, input.Permissions, model.Permissions)
	})
	t.Run("UpdateSuccess", func(t *testing.T) {
		ctx := context.Background()
		input := r.setupAndPopulate()[0]

		input.Level = 100

		model, err := r.services.Update(ctx, input)
		assert.Nil(t, err)
		assert.Equal(t, input.Id, model.Id)
		assert.Equal(t, input.Nama, model.Nama)
		assert.Equal(t, input.Deskripsi, model.Deskripsi)
		assert.Equal(t, input.Level, model.Level)
		assert.Equal(t, 100, model.Level)
		assert.Equal(t, input.Permissions, model.Permissions)
	})
	t.Run("DeleteSuccess", func(t *testing.T) {
		ctx := context.Background()
		r.setupAndPopulate()
		models := r.repo.Get(ctx, auth_requests.NewGetRequest())
		assert.Len(t, models, 2)

		affected, err := r.services.DeleteById(ctx, "1")
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		models = r.repo.Get(ctx, auth_requests.NewGetRequest())
		assert.Len(t, models, 1)
		assert.Equal(t, "2", models[0].Id)
	})
	t.Run("InsertAndUpdateFailedValidation", func(t *testing.T) {
		ctx := context.Background()
		r.setupAndPopulate()

		input := &auth_entities.Role{}
		_, err := r.services.Insert(ctx, input)

		assert.ErrorIs(t, err, helpers_error.ValidationError)
		assert.ErrorContains(t, err, "nama.required")
		assert.ErrorContains(t, err, "level.min.1")

		input.Id = "1"
		input.Level = -1
		_, err = r.services.Update(ctx, input)
		assert.ErrorIs(t, err, helpers_error.ValidationError)
		assert.ErrorContains(t, err, "level.min.1")
	})
	t.Run("UpdateAndDeleteFailedNotFound", func(t *testing.T) {
		ctx := context.Background()
		r.setup()

		affected, err := r.services.DeleteById(ctx, "1")
		assert.Equal(t, 0, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
}

type RoleServicesTest struct {
	repo     auth_repo_interfaces.IRoleRepository
	services auth_service_interfaces.IRoleServices
}

func (r *RoleServicesTest) setup() {
	r.setMemoryRepository()
	r.services = auth_services.NewRoleServices(r.repo)
}
func (r *RoleServicesTest) setMemoryRepository() {
	r.repo = auth_repos_memory.NewRoleMemoryRepository()
}
func (r *RoleServicesTest) setupAndPopulate() []*auth_entities.Role {
	r.setup()
	input := []*auth_entities.Role{
		{
			Nama:      "Developer",
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
