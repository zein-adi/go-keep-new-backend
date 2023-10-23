package keep_repo_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
)

type IBarangRepository interface {
	Get(ctx context.Context, search string, lokasi string) []*keep_entities.Barang
	Insert(ctx context.Context, barang *keep_entities.Barang) (affected int, err error)
	Update(ctx context.Context, barang *keep_entities.Barang) (affected int, err error)
	DeleteByNama(ctx context.Context, nama string) (affected int, err error)
}
