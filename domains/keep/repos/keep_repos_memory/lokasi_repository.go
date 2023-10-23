package keep_repos_memory

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"strings"
)

var lokasiEntityName = "lokasi"

func NewLokasiMemoryRepository() *LokasiMemoryRepository {
	return &LokasiMemoryRepository{}
}

type LokasiMemoryRepository struct {
	Data []*keep_entities.Lokasi
}

func (x *LokasiMemoryRepository) Get(_ context.Context, search string) []*keep_entities.Lokasi {
	search = strings.ToLower(search)
	models := helpers.Filter(x.Data, func(lokasi *keep_entities.Lokasi) bool {
		return strings.Contains(strings.ToLower(lokasi.Nama), search)
	})
	return helpers.Map(models, func(lokasi *keep_entities.Lokasi) *keep_entities.Lokasi {
		return lokasi.Copy()
	})
}
func (x *LokasiMemoryRepository) Insert(_ context.Context, lokasi *keep_entities.Lokasi) (affected int, err error) {
	x.Data = append(x.Data, lokasi)
	return 1, nil
}
func (x *LokasiMemoryRepository) Update(_ context.Context, lokasi *keep_entities.Lokasi) (affected int, err error) {
	index, err := x.findIndexByNama(lokasi.Nama)
	if err != nil {
		return 0, err
	}
	x.Data[index].LastUpdate = lokasi.LastUpdate
	return 1, nil
}
func (x *LokasiMemoryRepository) DeleteByNama(_ context.Context, nama string) (affected int, err error) {
	index, err := x.findIndexByNama(nama)
	if err != nil {
		return 0, err
	}
	x.Data = append(x.Data[0:index], x.Data[index+1:]...)
	return 1, nil
}

func (x *LokasiMemoryRepository) findIndexByNama(nama string) (index int, err error) {
	index, err = helpers.FindIndex(x.Data, func(v *keep_entities.Lokasi) bool {
		return v.Nama == nama
	})
	if err != nil {
		return 0, helpers_error.NewEntryNotFoundError(lokasiEntityName, "nama", nama)
	}
	return index, nil
}
