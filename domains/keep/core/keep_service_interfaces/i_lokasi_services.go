package keep_service_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
)

type ILokasiServices interface {
	Get(ctx context.Context, search string) []*keep_entities.Lokasi
	UpdateLokasiFromTransaksi(ctx context.Context) (affected int, err error)
}
