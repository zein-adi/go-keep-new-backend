package auth_services

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_requests"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_responses"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/validator"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func NewUserServices(userRepo auth_repo_interfaces.IUserRepository, roleRepo auth_repo_interfaces.IRoleRepository) *UserServices {
	return &UserServices{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

type UserServices struct {
	userRepo auth_repo_interfaces.IUserRepository
	roleRepo auth_repo_interfaces.IRoleRepository
}

func (x *UserServices) Get(ctx context.Context, request auth_requests.GetRequest) []*auth_responses.UserResponse {
	models := x.userRepo.Get(ctx, request)
	return helpers.Map(models, func(d *auth_entities.User) *auth_responses.UserResponse {
		return x.newUserResponseFromUserEntity(d)
	})
}
func (x *UserServices) Count(ctx context.Context, request auth_requests.GetRequest) int {
	return x.userRepo.Count(ctx, request)
}
func (x *UserServices) Insert(ctx context.Context, request *auth_requests.UserInputRequest, currentUserRoleIds []string) (*auth_responses.UserResponse, error) {
	// Basic Input Validation
	err := validator.New().ValidateStruct(request)
	if err != nil {
		return &auth_responses.UserResponse{}, err
	}
	err = x.validateUsername(ctx, request.Username, "")
	if err != nil {
		return &auth_responses.UserResponse{}, err
	}

	// Validate Current User Access To Updated Model
	err = x.validateAccessedLevel(ctx, currentUserRoleIds, request.RoleIds)
	if err != nil {
		return &auth_responses.UserResponse{}, err
	}

	// Proses
	userEntity := x.newUserEntityFromUserInputRequest(request)
	userEntity.Username = strings.ToLower(userEntity.Username)
	model, err := x.userRepo.Insert(ctx, userEntity)
	if err != nil {
		return &auth_responses.UserResponse{}, err
	}
	return x.newUserResponseFromUserEntity(model), nil
}
func (x *UserServices) Update(ctx context.Context, request *auth_requests.UserUpdateRequest, currentUserRoleIds []string) (*auth_responses.UserResponse, error) {
	// Check Data
	_, err := x.userRepo.FindById(ctx, request.Id)
	if err != nil {
		return &auth_responses.UserResponse{}, err
	}

	// Basic Input Validation
	err = validator.New().ValidateStruct(request)
	if err != nil {
		return &auth_responses.UserResponse{}, err
	}
	err = x.validateUsername(ctx, request.Username, request.Id)
	if err != nil {
		return &auth_responses.UserResponse{}, err
	}

	// Validate Current User Access To Updated Model
	err = x.validateAccessedLevel(ctx, currentUserRoleIds, request.RoleIds)
	if err != nil {
		return &auth_responses.UserResponse{}, err
	}

	// Proses
	userEntity := x.newUserEntityFromUserUpdateRequest(request)
	userEntity.Username = strings.ToLower(userEntity.Username)
	_, err = x.userRepo.Update(ctx, userEntity)
	if err != nil {
		return &auth_responses.UserResponse{}, err
	}
	return x.newUserResponseFromUserEntity(userEntity), nil
}
func (x *UserServices) UpdatePassword(ctx context.Context, request *auth_requests.UserUpdatePasswordRequest, currentUserRoleIds []string) (affected int, err error) {
	model, err := x.userRepo.FindById(ctx, request.Id)
	if err != nil {
		return 0, err
	}

	// Basic Input Validation
	err = validator.New().ValidateStruct(request)
	if err != nil {
		return 0, err
	}

	// Validate Current User Access To Updated Model
	err = x.validateAccessedLevel(ctx, currentUserRoleIds, model.RoleIds)
	if err != nil {
		return 0, err
	}

	return x.userRepo.UpdatePassword(ctx, request.Id, HashPassword(request.Password))
}
func (x *UserServices) DeleteById(ctx context.Context, id string, currentUserRoleIds []string) (affected int, err error) {
	model, err := x.userRepo.FindById(ctx, id)
	if err != nil {
		return 0, err
	}

	// Validate Current User Access To Updated Model
	err = x.validateAccessedLevel(ctx, currentUserRoleIds, model.RoleIds)
	if err != nil {
		return 0, err
	}

	return x.userRepo.DeleteById(ctx, id)
}

func (x *UserServices) validateUsername(ctx context.Context, username string, userId string) error {
	count := x.userRepo.CountByUsername(ctx, username, userId)
	if count > 0 {
		return helpers_error.NewValidationErrors("username", "duplicate", "")
	}
	return nil
}
func (x *UserServices) newUserEntityFromUserInputRequest(user *auth_requests.UserInputRequest) *auth_entities.User {
	return &auth_entities.User{
		Username: user.Username,
		Password: HashPassword(user.Password),
		Nama:     user.Nama,
		RoleIds:  user.RoleIds,
	}
}
func (x *UserServices) newUserEntityFromUserUpdateRequest(user *auth_requests.UserUpdateRequest) *auth_entities.User {
	return &auth_entities.User{
		Id:       user.Id,
		Username: user.Username,
		Nama:     user.Nama,
		RoleIds:  user.RoleIds,
	}
}
func (x *UserServices) newUserResponseFromUserEntity(user *auth_entities.User) *auth_responses.UserResponse {
	return &auth_responses.UserResponse{
		Id:       user.Id,
		Username: user.Username,
		Nama:     user.Nama,
		RoleIds:  user.RoleIds,
	}
}
func (x *UserServices) validateAccessedLevel(ctx context.Context, currentUserRoleIds, accessedRoleIds []string) error {
	accessedRoles, err := x.roleRepo.GetById(ctx, accessedRoleIds)
	if err != nil {
		return err
	}
	currentUserRoles, err := x.roleRepo.GetById(ctx, currentUserRoleIds)
	if err != nil {
		return err
	}
	accessedMinLevel := getMinimumLevelFromRoles(accessedRoles)
	userMinLevel := getMinimumLevelFromRoles(currentUserRoles)
	err = validateAccessedLevel(userMinLevel, accessedMinLevel)
	if err != nil {
		return err
	}
	return nil
}

func HashPassword(unencryptedPassword string) string {
	passwordByte, _ := bcrypt.GenerateFromPassword([]byte(unencryptedPassword), 10)
	return string(passwordByte)
}
