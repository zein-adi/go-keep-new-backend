package keep_repo_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
)

type ILokasiRepository interface {
	Get(ctx context.Context, search string) []*keep_entities.Lokasi
	Insert(ctx context.Context, lokasi *keep_entities.Lokasi) (affected int, err error)
	Update(ctx context.Context, lokasi *keep_entities.Lokasi) (affected int, err error)
	DeleteByNama(ctx context.Context, nama string) (affected int, err error)
}
