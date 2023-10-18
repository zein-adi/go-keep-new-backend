package keep_service_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
)

type IKantongServices interface {
	Get(ctx context.Context) []*keep_entities.Kantong
	Insert(ctx context.Context, kantongRequest *keep_request.KantongInsert) (*keep_entities.Kantong, error)
	Update(ctx context.Context, kantongRequest *keep_request.KantongUpdate) (affected int, err error)
	DeleteById(ctx context.Context, id string) (affected int, err error)

	GetTrashed(ctx context.Context) []*keep_entities.Kantong
	RestoreTrashedById(ctx context.Context, id string) (affected int, err error)
	DeleteTrashedById(ctx context.Context, id string) (affected int, err error)
}
