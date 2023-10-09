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
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
)

func TestUser(t *testing.T) {
	r := UserServicesTest{}
	r.setup()

	t.Run("InsertSuccess", func(t *testing.T) {
		r.setup()
		ctx := context.Background()

		req := &auth_requests.UserInputRequest{
			Username:             "zeinadiyusuf",
			Password:             "aA123456",
			PasswordConfirmation: "aA123456",
			Nama:                 "Zein",
			RoleIds:              []string{"1"},
		}
		response, err := r.services.Insert(ctx, req)
		assert.Nil(t, err)
		assert.NotEmpty(t, response.Id)

		model, _ := r.repo.FindById(ctx, response.Id)
		assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(model.Password), []byte(req.Password)))
	})
	t.Run("GetSuccess", func(t *testing.T) {
		r.setupAndPopulate()
		models := r.services.Get(context.Background(), auth_requests.NewGetRequest())
		assert.Len(t, models, 2)
	})
	t.Run("UpdateSuccess", func(t *testing.T) {
		id := r.setupAndPopulate()[0].Id

		req := &auth_requests.UserUpdateRequest{
			Id:       id,
			Username: "bambangsantoso",
			Nama:     "Bambang Santoso",
			RoleIds:  []string{"2"},
		}
		response, err := r.services.Update(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, response.Id, req.Id)
		assert.Equal(t, response.Username, req.Username)
		assert.Equal(t, response.Nama, req.Nama)
		assert.Equal(t, response.RoleIds, req.RoleIds)

		model, err := r.repo.FindById(context.Background(), id)
		assert.Nil(t, err)
		assert.Equal(t, model.Id, req.Id)
		assert.Equal(t, model.Username, req.Username)
		assert.Equal(t, model.Nama, req.Nama)
		assert.Equal(t, model.RoleIds, req.RoleIds)
	})
	t.Run("UpdatePasswordSuccess", func(t *testing.T) {
		id := r.setupAndPopulate()[0].Id
		req := &auth_requests.UserUpdatePasswordRequest{
			Id:                   id,
			Password:             "aA123456789",
			PasswordConfirmation: "aA123456789",
		}
		aff, err := r.services.UpdatePassword(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, 1, aff)

		model, err := r.repo.FindById(context.Background(), id)
		assert.Nil(t, err)
		assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(model.Password), []byte(req.Password)))
	})
	t.Run("DeleteSuccess", func(t *testing.T) {
		id := r.setupAndPopulate()[0].Id
		models := r.repo.Get(context.Background(), auth_requests.NewGetRequest())
		assert.Len(t, models, 2)

		aff, err := r.services.DeleteById(context.Background(), id)
		assert.Nil(t, err)
		assert.Equal(t, 1, aff)

		models = r.repo.Get(context.Background(), auth_requests.NewGetRequest())
		assert.Len(t, models, 1)
	})
	t.Run("InsertUpdateFailedValidationsUsername", func(t *testing.T) {
		id := r.setupAndPopulate()[0].Id
		tests := map[string]string{
			"":                      "username.required",
			"a":                     "username.min",
			strings.Repeat("a", 7):  "username.min",
			strings.Repeat("a", 65): "username.max",
			"a ":                    "username.alphanum",
			"a!":                    "username.alphanum",
			"a/":                    "username.alphanum",
		}
		for value, expected := range tests {
			_, err := r.services.Insert(context.Background(), &auth_requests.UserInputRequest{
				Username: value,
			})
			assert.ErrorIs(t, err, helpers_error.ValidationError)
			assert.ErrorContains(t, err, expected)
		}
		for value, expected := range tests {
			_, err := r.services.Update(context.Background(), &auth_requests.UserUpdateRequest{
				Id:       id,
				Username: value,
			})
			assert.ErrorIs(t, err, helpers_error.ValidationError)
			assert.ErrorContains(t, err, expected)
		}
	})
	t.Run("InsertUpdateFailedValidationsNama", func(t *testing.T) {
		id := r.setupAndPopulate()[0].Id
		tests := map[string]string{
			"":                       "nama.required",
			"a":                      "nama.min",
			"aa":                     "nama.min",
			strings.Repeat("a", 129): "nama.max",
		}
		for value, expected := range tests {
			_, err := r.services.Insert(context.Background(), &auth_requests.UserInputRequest{
				Nama: value,
			})
			assert.ErrorIs(t, err, helpers_error.ValidationError)
			assert.ErrorContains(t, err, expected)
		}
		for value, expected := range tests {
			_, err := r.services.Update(context.Background(), &auth_requests.UserUpdateRequest{
				Id:   id,
				Nama: value,
			})
			assert.ErrorIs(t, err, helpers_error.ValidationError)
			assert.ErrorContains(t, err, expected)
		}
	})
	t.Run("InsertUpdatePasswordFailedValidations", func(t *testing.T) {
		id := r.setupAndPopulate()[0].Id
		tests := []map[string]string{
			{
				"ex": "password.required",
				"p":  "",
				"pc": "",
			},
			{
				"ex": "password.min",
				"p":  "a",
				"pc": "a",
			},
			{
				"ex": "password.min",
				"p":  strings.Repeat("a", 7),
				"pc": strings.Repeat("a", 7),
			},
			{
				"ex": "password.valid_password",
				"p":  strings.Repeat("a", 8),
				"pc": strings.Repeat("a", 8),
			},
			{
				"ex": "password.valid_password",
				"p":  strings.Repeat("a", 72),
				"pc": strings.Repeat("a", 72),
			},
			{
				"ex": "password.max",
				"p":  strings.Repeat("a", 73),
				"pc": strings.Repeat("a", 73),
			},
			{
				"ex": "password.valid_password",
				"p":  strings.Repeat("A", 8),
				"pc": strings.Repeat("A", 8),
			},
			{
				"ex": "password.valid_password",
				"p":  strings.Repeat("aA", 8),
				"pc": strings.Repeat("aA", 8),
			},
			{
				"ex": "password.valid_password",
				"p":  strings.Repeat("a1", 8),
				"pc": strings.Repeat("a1", 8),
			},
			{
				"ex": "password.valid_password",
				"p":  strings.Repeat("A1", 8),
				"pc": strings.Repeat("A1", 8),
			},
			{
				"ex": "password_confirmation.eqfield",
				"p":  "aA123456789",
				"pc": "b",
			},
		}
		for _, value := range tests {
			_, err := r.services.Insert(context.Background(), &auth_requests.UserInputRequest{
				Password:             value["p"],
				PasswordConfirmation: value["pc"],
			})
			assert.ErrorIs(t, err, helpers_error.ValidationError)
			assert.ErrorContains(t, err, value["ex"])
		}
		for _, value := range tests {
			_, err := r.services.UpdatePassword(context.Background(), &auth_requests.UserUpdatePasswordRequest{
				Id:                   id,
				Password:             value["p"],
				PasswordConfirmation: value["pc"],
			})
			assert.ErrorIs(t, err, helpers_error.ValidationError)
			assert.ErrorContains(t, err, value["ex"])
		}
	})
	t.Run("UpdateAndDeleteFailedNotFound", func(t *testing.T) {
		ctx := context.Background()
		r.setup()

		affected, err := r.services.DeleteById(ctx, "1")
		assert.Equal(t, 0, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		_, err = r.services.Update(ctx, &auth_requests.UserUpdateRequest{Id: "1"})
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
}

type UserServicesTest struct {
	repo     auth_repo_interfaces.IUserRepository
	services auth_service_interfaces.IUserServices
}

func (r *UserServicesTest) setup() {
	r.setMemoryRepository()
	r.services = auth_services.NewUserServices(r.repo)
}
func (r *UserServicesTest) setMemoryRepository() {
	r.repo = auth_repos_memory.NewUserMemoryRepository()
}
func (r *UserServicesTest) setMysqlRepository() {
	r.repo = auth_repos_mysql.NewUserMysqlRepository()
	models := r.repo.Get(context.Background(), auth_requests.NewGetRequest())
	for _, model := range models {
		_, _ = r.repo.DeleteById(context.Background(), model.Id)
	}
}
func (r *UserServicesTest) setupAndPopulate() []*auth_entities.User {
	r.setup()
	input := []*auth_entities.User{
		{
			Id:       "1",
			Username: "zeinadimukadar",
			Password: "$2a$15$KkPamGNJDGEAD8xA/S4XfOFJn0vxbSmXYWypYoTyba3f3wljB0kfC",
			Nama:     "Zein Adi Mukadar",
			RoleIds:  []string{"1"},
		},
		{
			Id:       "2",
			Username: "rachmadyanuarianto",
			Password: "$2a$15$UU2rQCKNlVeYaIqi2CSUnO7vMWwykFQLCOoOpoNusvoU/MaxOLlR2",
			Nama:     "Rachmad Yanuarianto",
			RoleIds:  []string{"2"},
		},
	}

	var models []*auth_entities.User
	for _, datum := range input {
		model, _ := r.repo.Insert(context.Background(), datum)
		models = append(models, model)
	}
	return models
}
