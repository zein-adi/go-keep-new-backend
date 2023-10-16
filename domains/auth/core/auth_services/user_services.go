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

func NewUserServices(userRepo auth_repo_interfaces.IUserRepository) *UserServices {
	return &UserServices{
		userRepo: userRepo,
	}
}

type UserServices struct {
	userRepo auth_repo_interfaces.IUserRepository
}

func (r *UserServices) Get(ctx context.Context, request auth_requests.GetRequest) []*auth_responses.UserResponse {
	models := r.userRepo.Get(ctx, request)
	return helpers.Map(models, func(d *auth_entities.User) *auth_responses.UserResponse {
		return r.newUserResponseFromUserEntity(d)
	})
}

func (r *UserServices) Count(ctx context.Context, request auth_requests.GetRequest) int {
	return r.userRepo.Count(ctx, request)
}

func (r *UserServices) Insert(ctx context.Context, request *auth_requests.UserInputRequest) (*auth_responses.UserResponse, error) {
	// Validation
	err := validator.New().ValidateStruct(request)
	if err != nil {
		return &auth_responses.UserResponse{}, err
	}
	err = r.validateUsername(ctx, request.Username, "")
	if err != nil {
		return &auth_responses.UserResponse{}, err
	}

	// Proses
	userEntity := r.newUserEntityFromUserInputRequest(request)
	userEntity.Username = strings.ToLower(userEntity.Username)
	model, err := r.userRepo.Insert(ctx, userEntity)
	if err != nil {
		return &auth_responses.UserResponse{}, err
	}
	return r.newUserResponseFromUserEntity(model), nil
}

func (r *UserServices) Update(ctx context.Context, request *auth_requests.UserUpdateRequest) (*auth_responses.UserResponse, error) {
	// Check Data
	_, err := r.userRepo.FindById(ctx, request.Id)
	if err != nil {
		return &auth_responses.UserResponse{}, err
	}

	// Validation
	err = validator.New().ValidateStruct(request)
	if err != nil {
		return &auth_responses.UserResponse{}, err
	}
	err = r.validateUsername(ctx, request.Username, request.Id)
	if err != nil {
		return &auth_responses.UserResponse{}, err
	}

	// Proses
	userEntity := r.newUserEntityFromUserUpdateRequest(request)
	userEntity.Username = strings.ToLower(userEntity.Username)
	_, err = r.userRepo.Update(ctx, userEntity)
	if err != nil {
		return &auth_responses.UserResponse{}, err
	}
	return r.newUserResponseFromUserEntity(userEntity), nil
}

func (r *UserServices) UpdatePassword(ctx context.Context, request *auth_requests.UserUpdatePasswordRequest) (affected int, err error) {
	_, err = r.userRepo.FindById(ctx, request.Id)
	if err != nil {
		return 0, err
	}

	err = validator.New().ValidateStruct(request)
	if err != nil {
		return 0, err
	}

	return r.userRepo.UpdatePassword(ctx, request.Id, HashPassword(request.Password))
}

func (r *UserServices) DeleteById(ctx context.Context, id string) (affected int, err error) {
	_, err = r.userRepo.FindById(ctx, id)
	if err != nil {
		return 0, err
	}
	return r.userRepo.DeleteById(ctx, id)
}

func (r *UserServices) validateUsername(ctx context.Context, username string, userId string) error {
	count := r.userRepo.CountByUsername(ctx, username, userId)
	if count > 0 {
		return helpers_error.NewValidationErrors("username", "duplicate", "")
	}
	return nil
}

func (r *UserServices) newUserEntityFromUserInputRequest(user *auth_requests.UserInputRequest) *auth_entities.User {
	return &auth_entities.User{
		Username: user.Username,
		Password: HashPassword(user.Password),
		Nama:     user.Nama,
		RoleIds:  user.RoleIds,
	}
}
func (r *UserServices) newUserEntityFromUserUpdateRequest(user *auth_requests.UserUpdateRequest) *auth_entities.User {
	return &auth_entities.User{
		Id:       user.Id,
		Username: user.Username,
		Nama:     user.Nama,
		RoleIds:  user.RoleIds,
	}
}
func (r *UserServices) newUserResponseFromUserEntity(user *auth_entities.User) *auth_responses.UserResponse {
	return &auth_responses.UserResponse{
		Id:       user.Id,
		Username: user.Username,
		Nama:     user.Nama,
		RoleIds:  user.RoleIds,
	}
}
func HashPassword(unencryptedPassword string) string {
	passwordByte, _ := bcrypt.GenerateFromPassword([]byte(unencryptedPassword), 10)
	return string(passwordByte)
}
