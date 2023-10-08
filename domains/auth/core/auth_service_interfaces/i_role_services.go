package auth_service_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_requests"
)

type IRoleServices interface {
	Get(ctx context.Context, request auth_requests.GetRequest) []*auth_entities.Role
	Count(ctx context.Context, request auth_requests.GetRequest) int
	Insert(ctx context.Context, role *auth_entities.Role) (*auth_entities.Role, error)
	Update(ctx context.Context, rolePatch *auth_entities.Role) (*auth_entities.Role, error)
	DeleteById(ctx context.Context, id string) (int, error)
}
