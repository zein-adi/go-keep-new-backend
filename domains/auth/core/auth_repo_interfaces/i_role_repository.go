package auth_repo_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_requests"
)

type IRoleRepository interface {
	Get(ctx context.Context, request *helpers_requests.Get) []*auth_entities.Role
	Count(ctx context.Context, request *helpers_requests.Get) (count int)
	// GetById throws NewEntryCountMismatchError
	GetById(ctx context.Context, ids []string) ([]*auth_entities.Role, error)
	FindById(ctx context.Context, id string) (*auth_entities.Role, error)
	CountByNama(ctx context.Context, nama string, exceptId string) (count int)
	Insert(ctx context.Context, role *auth_entities.Role) (*auth_entities.Role, error)
	Update(ctx context.Context, role *auth_entities.Role) (affected int, er error)
	DeleteById(ctx context.Context, id string) (affected int, err error)
}
