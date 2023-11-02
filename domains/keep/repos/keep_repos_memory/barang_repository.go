package keep_repos_memory

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"strings"
)

var barangEntityName = "barang"

func NewBarangMemoryRepository() *BarangMemoryRepository {
	return &BarangMemoryRepository{}
}

type BarangMemoryRepository struct {
	Data []*keep_entities.Barang
}

func (x *BarangMemoryRepository) Get(_ context.Context, search string, lokasi string) []*keep_entities.Barang {
	search = strings.ToLower(search)
	searchArray := strings.Split(search, " ")
	lokasi = strings.ToLower(lokasi)
	models := helpers.Filter(x.Data, func(v *keep_entities.Barang) bool {
		res := true

		if search != "" {
			nama := strings.ToLower(v.Nama)
			tmpRes := true
			for _, s := range searchArray {
				if !strings.Contains(nama, s) {
					tmpRes = false
					break
				}
			}
			res = res && tmpRes
		}
		if lokasi != "" {
			_, err := helpers.FindIndex(v.Details, func(v *keep_entities.BarangDetail) bool {
				return strings.Contains(strings.ToLower(v.Lokasi), lokasi)
			})
			if err != nil {
				res = false
			}
		}

		return res
	})
	return helpers.Map(models, func(v *keep_entities.Barang) *keep_entities.Barang {
		model := v.Copy()
		model.Details = helpers.Filter(model.Details, func(d *keep_entities.BarangDetail) bool {
			return strings.Contains(strings.ToLower(d.Lokasi), lokasi)
		})
		return model
	})
}

func (x *BarangMemoryRepository) Insert(_ context.Context, barang *keep_entities.Barang) (affected int, err error) {
	x.Data = append(x.Data, barang.Copy())
	return 1, nil
}

func (x *BarangMemoryRepository) Update(_ context.Context, barang *keep_entities.Barang) (affected int, err error) {
	index, err := x.findIndexByNama(barang.Nama)
	if err != nil {
		return 0, err
	}
	x.Data[index] = barang.Copy()
	return 1, nil
}

func (x *BarangMemoryRepository) DeleteByNama(_ context.Context, nama string) (affected int, err error) {
	index, err := x.findIndexByNama(nama)
	if err != nil {
		return 0, err
	}
	x.Data = append(x.Data[0:index], x.Data[index+1:]...)
	return 1, nil
}

func (x *BarangMemoryRepository) findIndexByNama(nama string) (index int, err error) {
	index, err = helpers.FindIndex(x.Data, func(v *keep_entities.Barang) bool {
		return v.Nama == nama
	})
	if err != nil {
		return 0, helpers_error.NewEntryNotFoundError(barangEntityName, "nama", nama)
	}
	return index, nil
}
