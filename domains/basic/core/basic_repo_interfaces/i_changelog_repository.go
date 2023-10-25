package basic_repo_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/core/basic_entities"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_requests"
)

type IChangelogRepository interface {
	Get(ctx context.Context, request *helpers_requests.Get) []*basic_entities.Changelog
	Count(ctx context.Context, request *helpers_requests.Get) (count int)
	FindById(ctx context.Context, id string) (*basic_entities.Changelog, error)
	Insert(ctx context.Context, changelog *basic_entities.Changelog) (*basic_entities.Changelog, error)
	Update(ctx context.Context, changelog *basic_entities.Changelog) (affected int, err error)
	DeleteById(ctx context.Context, id string) (affected int, err error)
}
