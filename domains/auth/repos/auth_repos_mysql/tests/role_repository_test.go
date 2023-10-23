package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_requests"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/repos/auth_repos_mysql"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_env"
	"testing"
)

func TestRoleRepositoryTests(t *testing.T) {
	helpers_env.Init(5)
	repo := auth_repos_mysql.NewRoleMysqlRepository()
	defer repo.Cleanup()
	s := RoleRepositoryTest{repo: repo}

	t.Run("InsertSuccess", func(t *testing.T) {
		s.Clear()

		input := &auth_entities.Role{
			Nama:      "Admin",
			Deskripsi: "Deskripsi role admin",
			Level:     1,
			Permissions: []string{
				"user.role.get",
				"user.role.insert",
				"user.role.update",
				"user.role.delete",
			},
		}
		model, _ := repo.Insert(context.Background(), input)

		assert.NotEmpty(t, model.Id)
	})
	t.Run("CountSuccess", func(t *testing.T) {
		s.ClearAndPopulate()
		count := repo.Count(context.Background(), auth_requests.NewGet())
		assert.Equal(t, 2, count)
	})
	t.Run("GetSuccess", func(t *testing.T) {
		s.ClearAndPopulate()
		models := repo.Get(context.Background(), auth_requests.NewGet())
		assert.Len(t, models, 2)
	})
	t.Run("UpdateSuccess", func(t *testing.T) {
		ori := s.ClearAndPopulate()[1]

		input := &auth_entities.Role{
			Id:        ori.Id,
			Nama:      "Staf",
			Deskripsi: "Deskripsi role staf",
			Level:     20,
			Permissions: []string{
				"user.role.get",
			},
		}
		affected, err := repo.Update(context.Background(), input)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)
	})
	t.Run("DeleteSuccess", func(t *testing.T) {
		ori := s.ClearAndPopulate()[0]

		count := repo.Count(context.Background(), auth_requests.NewGet())
		assert.Equal(t, 2, count)

		affected, err := repo.DeleteById(context.Background(), ori.Id)
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, affected)

		count = repo.Count(context.Background(), auth_requests.NewGet())
		assert.Equal(t, 1, count)
	})
}

type RoleRepositoryTest struct {
	repo auth_repo_interfaces.IRoleRepository
}

func (x *RoleRepositoryTest) Clear() {
	ctx := context.Background()
	models := x.repo.Get(ctx, auth_requests.NewGet())
	for _, model := range models {
		_, _ = x.repo.DeleteById(ctx, model.Id)
	}
}
func (x *RoleRepositoryTest) ClearAndPopulate() (models []*auth_entities.Role) {
	x.Clear()
	ctx := context.Background()

	input := []*auth_entities.Role{
		{
			Nama:      "Admin",
			Deskripsi: "Deskripsi role admin",
			Level:     1,
			Permissions: []string{
				"user.role.get",
				"user.role.insert",
				"user.role.update",
				"user.role.delete",
			},
		},
		{
			Nama:      "Staf",
			Deskripsi: "Deskripsi role staf",
			Level:     10,
			Permissions: []string{
				"user.role.get",
			},
		},
	}

	models = nil
	for _, in := range input {
		model, _ := x.repo.Insert(ctx, in)
		models = append(models, model)
	}
	return models
}
