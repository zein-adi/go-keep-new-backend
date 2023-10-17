package keep_service_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
)

type IPosServices interface {
	Get(ctx context.Context, request *keep_request.PosGetRequest) []*keep_entities.Pos
	Insert(ctx context.Context, posRequest *keep_request.PosInputUpdateRequest) (*keep_entities.Pos, error)
	Update(ctx context.Context, posRequest *keep_request.PosInputUpdateRequest) (*keep_entities.Pos, error)
	DeleteById(ctx context.Context, id string) (affected int, err error)

	GetTrashed(ctx context.Context) []*keep_entities.Pos
	RestoreTrashedById(ctx context.Context, id string) (affected int, err error)
	DeleteTrashedById(ctx context.Context, id string) (affected int, err error)
}
