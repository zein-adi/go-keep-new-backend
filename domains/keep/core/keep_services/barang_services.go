package keep_services

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"math"
	"time"
)

func NewBarangServices(repo keep_repo_interfaces.IBarangRepository, transaksiRepo keep_repo_interfaces.ITransaksiRepository) *BarangServices {
	return &BarangServices{
		repo:          repo,
		transaksiRepo: transaksiRepo,
	}
}

type BarangServices struct {
	repo          keep_repo_interfaces.IBarangRepository
	transaksiRepo keep_repo_interfaces.ITransaksiRepository
}

func (x *BarangServices) Get(ctx context.Context, search string, lokasi string) []*keep_entities.Barang {
	return x.repo.Get(ctx, search, lokasi)
}

func (x *BarangServices) UpdateBarangFromTransaksi(ctx context.Context) (affected int, err error) {
	request := keep_request.NewGetTransaksi()
	request.WaktuAwal = time.Now().AddDate(0, -6, 0).Unix()
	request.Jenis = "pengeluaran"
	transaksis := x.transaksiRepo.Get(ctx, request)
	transaksiMap := make(map[string]*keep_entities.Barang)
	for _, transaksi := range transaksis {
		for _, detail := range transaksi.Details {
			k := detail.Uraian
			if k == "" {
				continue
			}

			// Mapping Per Barang
			_, ok := transaksiMap[k]
			if !ok {
				transaksiMap[k] = &keep_entities.Barang{
					Nama:         detail.Uraian,
					Harga:        detail.Harga,
					Diskon:       detail.Diskon,
					SatuanNama:   detail.SatuanNama,
					SatuanJumlah: detail.SatuanJumlah,
					SatuanHarga:  detail.SatuanHarga,
					Keterangan:   detail.Keterangan,
					LastUpdate:   0,
					Details:      make([]*keep_entities.BarangDetail, 0),
				}
			}

			t := transaksiMap[k]
			if t.LastUpdate < transaksi.Waktu {
				t.LastUpdate = transaksi.Waktu
			}

			// Mapping Per Detail Barang / Lokasi
			_, err2 := helpers.FindIndex(t.Details, func(detail *keep_entities.BarangDetail) bool {
				return detail.Lokasi == transaksi.Lokasi && math.Ceil(detail.SatuanHarga) == math.Ceil(detail.SatuanHarga)
			})
			if err2 != nil {
				barangDetail := &keep_entities.BarangDetail{
					Lokasi:      transaksi.Lokasi,
					Harga:       detail.Harga,
					Diskon:      detail.Diskon,
					SatuanHarga: detail.SatuanHarga,
					Keterangan:  detail.Keterangan,
				}
				t.Details = append(t.Details, barangDetail)

				// Bila detail barang lebih murah dari yang tercatat saat ini, maka update
				if barangDetail.SatuanHarga < t.SatuanHarga {
					t.Harga = barangDetail.Harga
					t.Diskon = barangDetail.Diskon
					t.SatuanHarga = barangDetail.SatuanHarga
				}
			}
		}
	}

	barangs := x.repo.Get(ctx, "", "")
	barangMap := helpers.KeyBy(barangs, func(d *keep_entities.Barang) string {
		return d.Nama
	})

	for nama, barang := range transaksiMap {
		_, ok := barangMap[nama]
		if !ok {
			af, err2 := x.repo.Insert(ctx, barang)
			affected += af
			if err2 != nil {
				return affected, err2
			}
		} else {
			af, err2 := x.repo.Update(ctx, barang)
			affected += af
			if err2 != nil {
				return affected, err2
			}
		}
	}
	for nama := range barangMap {
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
