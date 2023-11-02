package keep_service_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_requests"
)

type IKantongHistoryServices interface {
	Get(ctx context.Context, kantongId string, request *helpers_requests.Get) []*keep_entities.KantongHistory
	Insert(ctx context.Context, kantongHistoryRequest *keep_request.KantongHistoryInsertUpdate) (*keep_entities.KantongHistory, error)
	Update(ctx context.Context, kantongHistoryRequest *keep_request.KantongHistoryInsertUpdate) (affected int, err error)
	DeleteById(ctx context.Context, kantongId string, id string) (affected int, err error)
}
