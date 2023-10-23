package keep_service_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
)

type ITransaksiServices interface {
	Get(ctx context.Context, request *keep_request.GetTransaksi) []*keep_entities.Transaksi
	Insert(ctx context.Context, transaksiRequest *keep_request.TransaksiInputUpdate) (*keep_entities.Transaksi, error)
	Update(ctx context.Context, transaksiRequest *keep_request.TransaksiInputUpdate) (affected int, err error)
	DeleteById(ctx context.Context, id string) (affected int, err error)

	GetTrashed(ctx context.Context) []*keep_entities.Transaksi
	RestoreTrashedById(ctx context.Context, id string) (affected int, err error)
	DeleteTrashedById(ctx context.Context, id string) (affected int, err error)
}
