package auth_repo_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_requests"
)

type IUserRepository interface {
	Get(ctx context.Context, request auth_requests.GetRequest) []*auth_entities.User
	Count(ctx context.Context, request auth_requests.GetRequest) (count int)
	FindById(ctx context.Context, id string) (*auth_entities.User, error)
	FindByUsername(ctx context.Context, username string) (*auth_entities.User, error)
	CountByUsername(ctx context.Context, username string, exceptId string) (count int)
	Insert(ctx context.Context, user *auth_entities.User) (*auth_entities.User, error)
	Update(ctx context.Context, user *auth_entities.User) (affected int, er error)
	UpdatePassword(ctx context.Context, userId, password string) (affected int, err error)
	DeleteById(ctx context.Context, id string) (affected int, err error)
}
