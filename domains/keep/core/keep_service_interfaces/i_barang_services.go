package keep_service_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
)

type IBarangServices interface {
	Get(ctx context.Context, search string, lokasi string) []*keep_entities.Barang
	UpdateBarangFromTransaksi(ctx context.Context) (affected int, err error)
}
