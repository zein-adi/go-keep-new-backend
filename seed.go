package main

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/dependency_injection"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_requests"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_services"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/validator"
)

func RunSeed(username, password string) {
	validate := auth_requests.UserInputRequest{
		Username:             username,
		Password:             password,
		PasswordConfirmation: password,
		Nama:                 "placeholder",
		RoleIds:              []string{"placeholder"},
	}
	err := validator.New().ValidateStruct(validate)
	helpers_error.PanicIfError(err)

	permissionRepo := dependency_injection.InitUserPermissionMemoryRepository()
	roleRepo := dependency_injection.InitUserRoleMysqlRepository()
	userRepo := dependency_injection.InitUserUserMysqlRepository()

	ctx := context.Background()

	roleInput := &auth_entities.Role{
		Nama:        "Developer",
		Deskripsi:   "Roles khusus untuk developer / super administrator",
		Level:       1,
		Permissions: permissionRepo.Get(ctx),
	}
	role, _ := roleRepo.Insert(ctx, roleInput)

	user := &auth_entities.User{
		Nama:     username,
		Username: username,
		Password: auth_services.HashPassword(password),
		RoleIds:  []string{role.Id},
	}
	userRepo.Insert(ctx, user)
}
