package keep_repo_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
)

type IKantongHistoryRepository interface {
	Get(ctx context.Context) []*keep_entities.KantongHistory
	FindById(ctx context.Context, id string) (*keep_entities.KantongHistory, error)
	Insert(ctx context.Context, kantongHistory *keep_entities.KantongHistory) (*keep_entities.KantongHistory, error)
	Update(ctx context.Context, kantongHistory *keep_entities.KantongHistory) (affected int, err error)
	DeleteById(ctx context.Context, id string) (affected int, err error)
}
