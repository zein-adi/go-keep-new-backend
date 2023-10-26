package keep_repo_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
)

type IPosRepository interface {
	Get(ctx context.Context, request *keep_request.GetPos) []*keep_entities.Pos
	GetJumlahById(ctx context.Context, id string) (saldo int)

	CountChildren(ctx context.Context, id string) (count int)

	FindById(ctx context.Context, id string) (*keep_entities.Pos, error)

	Insert(ctx context.Context, pos *keep_entities.Pos) (*keep_entities.Pos, error)

	Update(ctx context.Context, pos *keep_entities.Pos) (affected int, err error)
	UpdateLeaf(ctx context.Context, id string, leaf bool) (affected int, err error)
	UpdateSaldo(ctx context.Context, id string, saldo int) (affected int)
	UpdateUrutan(ctx context.Context, id string, urutan int, parentId string) (affected int, err error)
	UpdateVisibility(ctx context.Context, id string, isShow bool) (affected int, err error)

	SoftDeleteById(ctx context.Context, id string) (affected int, err error)
	GetTrashed(ctx context.Context) []*keep_entities.Pos
	FindTrashedById(ctx context.Context, id string) (*keep_entities.Pos, error)
	RestoreTrashedById(ctx context.Context, id string) (affected int, err error)
	HardDeleteTrashedById(ctx context.Context, id string) (affected int, err error)
}
