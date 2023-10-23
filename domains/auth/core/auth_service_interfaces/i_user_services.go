package auth_service_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_requests"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_responses"
)

type IUserServices interface {
	Get(ctx context.Context, request auth_requests.Get) []*auth_responses.UserResponse
	Count(ctx context.Context, request auth_requests.Get) int
	Insert(ctx context.Context, request *auth_requests.UserInputRequest, currentUserRoleIds []string) (*auth_responses.UserResponse, error)
	Update(ctx context.Context, request *auth_requests.UserUpdateRequest, currentUserRoleIds []string) (*auth_responses.UserResponse, error)
	UpdatePassword(ctx context.Context, request *auth_requests.UserUpdatePasswordRequest, currentUserRoleIds []string) (affected int, err error)
	DeleteById(ctx context.Context, id string, currentUserRoleIds []string) (affected int, err error)
}
