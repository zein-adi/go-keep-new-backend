package auth_service_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_requests"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_responses"
)

type IUserServices interface {
	Get(ctx context.Context, request auth_requests.GetRequest) []*auth_responses.UserResponse
	Count(ctx context.Context, request auth_requests.GetRequest) int
	Insert(ctx context.Context, request *auth_requests.UserInputRequest) (*auth_responses.UserResponse, error)
	Update(ctx context.Context, request *auth_requests.UserUpdateRequest) (*auth_responses.UserResponse, error)
	UpdatePassword(ctx context.Context, request *auth_requests.UserUpdatePasswordRequest) (affected int, err error)
	DeleteById(ctx context.Context, id string) (affected int, err error)
}
