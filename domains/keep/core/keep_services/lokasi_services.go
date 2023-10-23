package keep_services

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"time"
)

func NewLokasiServices(repo keep_repo_interfaces.ILokasiRepository, transaksiRepo keep_repo_interfaces.ITransaksiRepository) *LokasiServices {
	return &LokasiServices{
		repo:          repo,
		transaksiRepo: transaksiRepo,
	}
}

type LokasiServices struct {
	repo          keep_repo_interfaces.ILokasiRepository
	transaksiRepo keep_repo_interfaces.ITransaksiRepository
}

func (x *LokasiServices) Get(ctx context.Context, search string) []*keep_entities.Lokasi {
	return x.repo.Get(ctx, search)
}

func (x *LokasiServices) UpdateLokasiFromTransaksi(ctx context.Context) (affected int, err error) {
	request := keep_request.NewGetTransaksi()
	request.WaktuAwal = time.Now().AddDate(0, -6, 0).Unix()
	transaksis := x.transaksiRepo.Get(ctx, request)
	transaksiMap := make(map[string]*keep_entities.Lokasi)
	for _, transaksi := range transaksis {
		l := transaksi.Lokasi
		if l == "" {
			continue
		}

		_, ok := transaksiMap[l]
		if !ok {
			transaksiMap[l] = &keep_entities.Lokasi{Nama: l}
		}

		if transaksiMap[l].LastUpdate < transaksi.CreatedAt {
			transaksiMap[l].LastUpdate = transaksi.CreatedAt
		}
	}
	lokasis := x.repo.Get(ctx, "")
	lokasiMap := helpers.KeyBy(lokasis, func(d *keep_entities.Lokasi) string {
		return d.Nama
	})

	for nama, lokasi := range transaksiMap {
		_, ok := lokasiMap[nama]
		if !ok {
			af, err2 := x.repo.Insert(ctx, lokasi)
			affected += af
			if err2 != nil {
				return affected, err2
			}
		} else {
			af, err2 := x.repo.Update(ctx, lokasi)
			affected += af
			if err2 != nil {
				return affected, err2
			}
		}
	}
	for nama := range lokasiMap {
		_, ok := transaksiMap[nama]
		if !ok {
			af, err2 := x.repo.DeleteByNama(ctx, nama)
			affected += af
			if err2 != nil {
				return affected, err2
			}
		}
	}
	return affected, nil
}
