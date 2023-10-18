package keep_service_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
)

type IKantongHistoryServices interface {
	Get(ctx context.Context) []*keep_entities.KantongHistory
	Insert(ctx context.Context, kantongHistoryRequest *keep_request.KantongHistoryInsertUpdate) (*keep_entities.KantongHistory, error)
	Update(ctx context.Context, kantongHistoryRequest *keep_request.KantongHistoryInsertUpdate) (affected int, err error)
	DeleteById(ctx context.Context, id string) (affected int, err error)
}
