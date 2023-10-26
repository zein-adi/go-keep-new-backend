package keep_repo_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
)

type IKantongRepository interface {
	Get(ctx context.Context) []*keep_entities.Kantong

	FindById(ctx context.Context, id string) (*keep_entities.Kantong, error)

	Insert(ctx context.Context, kantong *keep_entities.Kantong) (*keep_entities.Kantong, error)
	Update(ctx context.Context, kantong *keep_entities.Kantong) (affected int, err error)
	UpdateSaldo(ctx context.Context, id string, saldo int) (affected int, err error)
	UpdateUrutan(ctx context.Context, id string, urutan int, posId string) (affected int, err error)
	UpdateVisibility(ctx context.Context, id string, isShow bool) (affected int, err error)

	SoftDeleteById(ctx context.Context, id string) (affected int, err error)
	GetTrashed(ctx context.Context) []*keep_entities.Kantong
	FindTrashedById(ctx context.Context, id string) (*keep_entities.Kantong, error)
	RestoreTrashedById(ctx context.Context, id string) (affected int, err error)
	HardDeleteTrashedById(ctx context.Context, id string) (affected int, err error)
}
