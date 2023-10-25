package keep_repo_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
)

type IPosRepository interface {
	Get(ctx context.Context, request *keep_request.GetPos) []*keep_entities.Pos
	FindById(ctx context.Context, id string) (*keep_entities.Pos, error)
	Insert(ctx context.Context, pos *keep_entities.Pos) (*keep_entities.Pos, error)
	Update(ctx context.Context, pos *keep_entities.Pos) (affected int, err error)

	SoftDeleteById(ctx context.Context, id string) (affected int, err error)
	DeleteById(ctx context.Context, id string) (affected int, err error)

	GetTrashed(ctx context.Context) []*keep_entities.Pos
	FindTrashedById(ctx context.Context, id string) (*keep_entities.Pos, error)
	RestoreTrashedById(ctx context.Context, id string) (affected int, err error)

	CountChildren(ctx context.Context, id string) (count int)
	UpdateLeaf(ctx context.Context, id string, leaf bool) (affected int, err error)

	GetJumlahById(ctx context.Context, id string) (saldo int)
	UpdateSaldo(ctx context.Context, id string, saldo int) (affected int)
}
