package basic_service_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/core/basic_entities"
)

type IChangelogServices interface {
	Get(ctx context.Context, skip, take int) []*basic_entities.Changelog
	Insert(ctx context.Context, changelog *basic_entities.Changelog) (*basic_entities.Changelog, error)
	Update(ctx context.Context, changelog *basic_entities.Changelog) (affected int, err error)
	DeleteById(ctx context.Context, id string) (affected int, err error)
}
