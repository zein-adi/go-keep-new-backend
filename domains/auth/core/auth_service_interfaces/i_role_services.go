package auth_service_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_requests"
)

type IRoleServices interface {
	Get(ctx context.Context, request *helpers_requests.Get) []*auth_entities.Role
	Count(ctx context.Context, request *helpers_requests.Get) int
	Insert(ctx context.Context, role *auth_entities.Role, currentUserRoleIds []string) (*auth_entities.Role, error)
	Update(ctx context.Context, role *auth_entities.Role, currentUserRoleIds []string) (*auth_entities.Role, error)
	DeleteById(ctx context.Context, id string, currentUserRoleIds []string) (affected int, err error)
}
