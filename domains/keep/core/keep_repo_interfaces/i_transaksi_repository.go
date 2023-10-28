package keep_repo_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
)

type ITransaksiRepository interface {
	Get(ctx context.Context, request *keep_request.GetTransaksi) []*keep_entities.Transaksi
	GetJumlahByPosId(ctx context.Context, posId string) (saldo int)
	CountByPosId(ctx context.Context, posId string) (count int)
	FindById(ctx context.Context, id string) (*keep_entities.Transaksi, error)

	Insert(ctx context.Context, transaksi *keep_entities.Transaksi) (*keep_entities.Transaksi, error)

	Update(ctx context.Context, transaksi *keep_entities.Transaksi) (affected int, err error)

	SoftDeleteById(ctx context.Context, id string) (affected int, err error)
	GetTrashed(ctx context.Context) []*keep_entities.Transaksi
	FindTrashedById(ctx context.Context, id string) (*keep_entities.Transaksi, error)
	RestoreTrashedById(ctx context.Context, id string) (affected int, err error)
	HardDeleteTrashedById(ctx context.Context, id string) (affected int, err error)
}
