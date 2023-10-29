package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_services"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/repos/keep_repos_memory"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/repos/keep_repos_mysql"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_env"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"testing"
	"time"
)

func TestTransaksi(t *testing.T) {
	helpers_env.Init(5)
	x := NewTransaksiServicesTest()
	defer x.cleanup()

	t.Run("GetSuccess", func(t *testing.T) {
		// Check insert & get value
		ctx := context.Background()
		ori, _, _ := x.reset()
		oriKey := helpers.KeyBy(ori, func(d *keep_entities.Transaksi) string {
			return d.Id
		})

		models := x.services.Get(ctx, keep_request.NewGetTransaksi())
		assert.Len(t, models, 2)

		for _, m := range models {
			o := oriKey[m.Id]
			assert.Equal(t, o.Id, m.Id)
			assert.Equal(t, o.Waktu, m.Waktu)
			assert.Equal(t, o.Jenis, m.Jenis)
			assert.Equal(t, o.Jumlah, m.Jumlah)
			assert.Equal(t, o.PosAsalId, m.PosAsalId)
			assert.Equal(t, o.PosAsalNama, m.PosAsalNama)
			assert.Equal(t, o.PosTujuanId, m.PosTujuanId)
			assert.Equal(t, o.PosTujuanNama, m.PosTujuanNama)
			assert.Equal(t, o.KantongAsalId, m.KantongAsalId)
			assert.Equal(t, o.KantongAsalNama, m.KantongAsalNama)
			assert.Equal(t, o.KantongTujuanId, m.KantongTujuanId)
			assert.Equal(t, o.KantongTujuanNama, m.KantongTujuanNama)
			assert.Equal(t, o.Uraian, m.Uraian)
			assert.Equal(t, o.Keterangan, m.Keterangan)
			assert.Equal(t, o.Lokasi, m.Lokasi)
			assert.Equal(t, o.UrlFoto, m.UrlFoto)
			assert.Equal(t, o.CreatedAt, m.CreatedAt)
			assert.Equal(t, o.UpdatedAt, m.UpdatedAt)
			assert.Equal(t, o.Details, m.Details)
			assert.Equal(t, o.Status, m.Status)
		}
	})
	t.Run("SoftDeleteSuccess", func(t *testing.T) {
		ctx := context.Background()
		transaksis, _, _ := x.reset()
		m := transaksis[1]

		affected, err := x.services.DeleteById(ctx, m.Id)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		_, err = x.repo.FindById(ctx, m.Id)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		model, err := x.repo.FindTrashedById(ctx, m.Id)
		assert.Nil(t, err)
		assert.Equal(t, "trashed", model.Status)
	})
	t.Run("UpdateDeleteFailedCauseNotFound", func(t *testing.T) {
		_, poses, _ := x.reset()
		posPemasukan := poses[0]
		id := "999999"

		input := &keep_request.TransaksiInputUpdate{
			Id:          id,
			Jenis:       "pemasukan",
			Jumlah:      1000,
			PosTujuanId: posPemasukan.Id,
			Uraian:      "Gajian",
		}
		_, err := x.services.Update(context.Background(), input)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		affected, err := x.services.DeleteById(context.Background(), id)
		assert.Equal(t, 0, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
	t.Run("UpdateDeleteFailedCauseTrashed", func(t *testing.T) {
		transaksis, poses, _ := x.reset()
		posPemasukan := poses[0]
		m := transaksis[2]

		input := &keep_request.TransaksiInputUpdate{
			Id:          m.Id,
			Jenis:       "pemasukan",
			Jumlah:      1000,
			PosTujuanId: posPemasukan.Id,
			Uraian:      "Gajian",
		}
		_, err := x.services.Update(context.Background(), input)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		affected, err := x.services.DeleteById(context.Background(), m.Id)
		assert.Equal(t, 0, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
	t.Run("RestoreHardDeleteFailedCauseStatusActive", func(t *testing.T) {
		transaksis, _, _ := x.reset()
		m := transaksis[0]

		affected, err := x.services.RestoreTrashedById(context.Background(), m.Id)
		assert.Equal(t, 0, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		affected, err = x.services.DeleteTrashedById(context.Background(), m.Id)
		assert.Equal(t, 0, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
	t.Run("GetTrashedSuccess", func(t *testing.T) {
		x.reset()
		time.Sleep(time.Millisecond * 10)
		now := time.Now()

		models := x.services.GetTrashed(context.Background())
		assert.Len(t, models, 1)

		for _, m := range models {
			assert.True(t, now.After(time.Unix(m.CreatedAt, 0)))
			assert.True(t, now.After(time.Unix(m.UpdatedAt, 0)))
			assert.Equal(t, "trashed", m.Status)
		}
	})
	t.Run("RestoreTrashedSuccess", func(t *testing.T) {
		ctx := context.Background()
		transaksis, _, _ := x.reset()
		m := transaksis[2]

		affected, err := x.services.RestoreTrashedById(ctx, m.Id)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		model, err := x.repo.FindById(ctx, m.Id)
		assert.Nil(t, err)
		assert.Equal(t, "aktif", model.Status)

		_, err = x.repo.FindTrashedById(ctx, m.Id)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
	t.Run("HardDeleteTrashedSuccess", func(t *testing.T) {
		transaksis, _, _ := x.reset()
		m := transaksis[2]

		affected, err := x.services.DeleteTrashedById(context.Background(), m.Id)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		_, err = x.repo.FindById(context.Background(), m.Id)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		_, err = x.repo.FindTrashedById(context.Background(), m.Id)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
	t.Run("InsertPemasukanSuccess", func(t *testing.T) {
		ctx := context.Background()
		_, poses, kantongs := x.reset()
		posPemasukan := poses[0]
		kantongMandiri := kantongs[0]

		assert.Equal(t, 100000, posPemasukan.Saldo)
		assert.Equal(t, 100000, kantongMandiri.Saldo)

		before := time.Now().Add(-1 * time.Second)
		jenis := "pemasukan"
		jumlah := 1000000
		posTujuanId := posPemasukan.Id
		kantongTujuanId := kantongMandiri.Id
		uraian := "Gajian"
		keterangan := "November 2023"
		urlFoto := "https://test.com"
		waktu := time.Now().Unix()
		input := &keep_request.TransaksiInputUpdate{
			Waktu:           waktu,
			Jenis:           jenis,
			Jumlah:          jumlah,
			PosTujuanId:     posTujuanId,
			KantongTujuanId: kantongTujuanId,
			Uraian:          uraian,
			Keterangan:      keterangan,
			UrlFoto:         urlFoto,
		}
		m, err := x.services.Insert(ctx, input)
		assert.Nil(t, err)
		assert.NotEmpty(t, m.Id)
		assert.True(t, before.Before(time.Unix(m.Waktu, 0)))
		assert.Equal(t, jenis, m.Jenis)
		assert.Equal(t, jumlah, m.Jumlah)
		assert.Empty(t, m.PosAsalId)
		assert.Empty(t, m.PosAsalNama)
		assert.Equal(t, posTujuanId, m.PosTujuanId)
		assert.Equal(t, posPemasukan.Nama, m.PosTujuanNama)
		assert.Empty(t, m.KantongAsalId)
		assert.Empty(t, m.KantongAsalNama)
		assert.Equal(t, kantongTujuanId, m.KantongTujuanId)
		assert.Equal(t, kantongMandiri.Nama, m.KantongTujuanNama)
		assert.Equal(t, uraian, m.Uraian)
		assert.Equal(t, keterangan, m.Keterangan)
		assert.Empty(t, m.Lokasi)
		assert.Equal(t, urlFoto, m.UrlFoto)
		assert.True(t, before.Before(time.Unix(m.CreatedAt, 0)))
		assert.True(t, before.Before(time.Unix(m.UpdatedAt, 0)))
		assert.Empty(t, m.Details)
		assert.Equal(t, "aktif", m.Status)
	})
	t.Run("InsertMutasiSuccess", func(t *testing.T) {
		ctx := context.Background()
		_, poses, kantongs := x.reset()
		posPemasukan := poses[0]
		posMain := poses[2]
		kantongMandiri := kantongs[0]
		kantongBca := kantongs[1]

		assert.Equal(t, 100000, posPemasukan.Saldo)
		assert.Equal(t, 0, posMain.Saldo)
		assert.Equal(t, 100000, posPemasukan.Saldo)
		assert.Equal(t, 0, kantongBca.Saldo)

		jenis := "mutasi"
		jumlah := 40000
		posAsalId := posPemasukan.Id
		posTujuanId := posMain.Id
		kantongAsalId := kantongMandiri.Id
		kantongTujuanId := kantongBca.Id
		uraian := "Gajian"
		waktu := time.Now().Unix()
		input := &keep_request.TransaksiInputUpdate{
			Waktu:           waktu,
			Jenis:           jenis,
			Jumlah:          jumlah,
			PosAsalId:       posAsalId,
			PosTujuanId:     posTujuanId,
			KantongAsalId:   kantongAsalId,
			KantongTujuanId: kantongTujuanId,
			Uraian:          uraian,
		}
		_, err := x.services.Insert(ctx, input)
		assert.Nil(t, err)
	})
	t.Run("InsertPengeluaranSimpleSuccess", func(t *testing.T) {
		ctx := context.Background()
		_, poses, kantongs := x.reset()
		posPemasukan := poses[0]
		kantongMandiri := kantongs[0]

		assert.Equal(t, 100000, posPemasukan.Saldo)
		assert.Equal(t, 100000, posPemasukan.Saldo)

		jenis := "pengeluaran"
		jumlah := 10000
		posAsalId := posPemasukan.Id
		kantongAsalId := kantongMandiri.Id
		uraian := "Jajan"
		keterangan := "Cilok"
		lokasi := "Citra Ken Dedes"
		input := &keep_request.TransaksiInputUpdate{
			Waktu:         time.Now().Unix(),
			Jenis:         jenis,
			Jumlah:        jumlah,
			PosAsalId:     posAsalId,
			KantongAsalId: kantongAsalId,
			Uraian:        uraian,
			Keterangan:    keterangan,
			Lokasi:        lokasi,
		}
		_, err := x.services.Insert(ctx, input)
		assert.Nil(t, err)
	})
	t.Run("InsertPengeluaranDetailSuccess", func(t *testing.T) {
		ctx := context.Background()
		_, poses, kantongs := x.reset()
		posPemasukan := poses[0]
		kantongMandiri := kantongs[0]

		assert.Equal(t, 100000, posPemasukan.Saldo)
		assert.Equal(t, 100000, posPemasukan.Saldo)

		jenis := "pengeluaran"
		posAsalId := posPemasukan.Id
		kantongAsalId := kantongMandiri.Id
		uraian := "Jajan"
		lokasi := "Citra Ken Dedes"

		nama := "Bloat Cake Special"
		harga := float64(13000)
		jumlah := float64(1)
		diskon := float64(1000)
		satuanJumlah := float64(1)
		satuanNama := "pcs"
		keteranganDetail := "lagi diskon"

		input := &keep_request.TransaksiInputUpdate{
			Waktu:         time.Now().Unix(),
			Jenis:         jenis,
			PosAsalId:     posAsalId,
			KantongAsalId: kantongAsalId,
			Uraian:        uraian,
			Keterangan:    nama,
			Lokasi:        lokasi,
			Details: []*keep_request.TransaksiInputUpdateDetail{
				{
					Uraian:       nama,
					Harga:        harga,
					Jumlah:       jumlah,
					Diskon:       diskon,
					SatuanJumlah: satuanJumlah,
					SatuanNama:   satuanNama,
					Keterangan:   keteranganDetail,
				},
			},
		}
		m, err := x.services.Insert(ctx, input)
		d := m.Details[0]
		assert.Nil(t, err)
		assert.Equal(t, 12000, m.Jumlah)
		assert.Equal(t, nama, d.Uraian)
		assert.Equal(t, harga, d.Harga)
		assert.Equal(t, jumlah, d.Jumlah)
		assert.Equal(t, diskon, d.Diskon)
		assert.Equal(t, satuanJumlah, d.SatuanJumlah)
		assert.Equal(t, satuanNama, d.SatuanNama)
		assert.Equal(t, keteranganDetail, d.Keterangan)
	})
	t.Run("UpdatePemasukanKePengeluaranSuccess", func(t *testing.T) {
		ctx := context.Background()
		_, poses, kantongs := x.reset()
		posPemasukan := poses[0]
		kantongMandiri := kantongs[0]

		// Prepare Data
		input := &keep_request.TransaksiInputUpdate{
			Waktu:           time.Now().Unix(),
			Jenis:           "pemasukan",
			Jumlah:          10000,
			PosTujuanId:     posPemasukan.Id,
			KantongTujuanId: kantongMandiri.Id,
			Uraian:          "test",
		}
		m, err := x.services.Insert(ctx, input)
		assert.Nil(t, err)

		v, _ := x.repo.FindById(ctx, m.Id)
		assert.Equal(t, "pemasukan", v.Jenis)

		// Progress
		input = &keep_request.TransaksiInputUpdate{
			Id:            m.Id,
			Waktu:         time.Now().Unix(),
			Jenis:         "pengeluaran",
			Jumlah:        10000,
			PosAsalId:     posPemasukan.Id,
			KantongAsalId: kantongMandiri.Id,
			Uraian:        "test",
		}
		affected, err := x.services.Update(ctx, input)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		v, _ = x.repo.FindById(ctx, m.Id)
		assert.Equal(t, "pengeluaran", v.Jenis)
	})
	t.Run("UpdatePemasukanKeMutasiSuccess", func(t *testing.T) {
		ctx := context.Background()
		_, poses, kantongs := x.reset()
		posPemasukan := poses[0]
		posMain := poses[2]
		kantongMandiri := kantongs[0]
		kantongBca := kantongs[1]

		// Prepare Data
		input := &keep_request.TransaksiInputUpdate{
			Waktu:           time.Now().Unix(),
			Jenis:           "pemasukan",
			Jumlah:          10000,
			PosTujuanId:     posPemasukan.Id,
			KantongTujuanId: kantongMandiri.Id,
			Uraian:          "test",
		}
		m, err := x.services.Insert(ctx, input)
		assert.Nil(t, err)

		v, _ := x.repo.FindById(ctx, m.Id)
		assert.Equal(t, "pemasukan", v.Jenis)

		// Progress
		input = &keep_request.TransaksiInputUpdate{
			Id:              m.Id,
			Waktu:           time.Now().Unix(),
			Jenis:           "mutasi",
			Jumlah:          10000,
			PosAsalId:       posPemasukan.Id,
			PosTujuanId:     posMain.Id,
			KantongAsalId:   kantongMandiri.Id,
			KantongTujuanId: kantongBca.Id,
			Uraian:          "test",
		}
		affected, err := x.services.Update(ctx, input)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		v, _ = x.repo.FindById(ctx, m.Id)
		assert.Equal(t, "mutasi", v.Jenis)
	})
	t.Run("UpdatePengeluaranKePemasukanSuccess", func(t *testing.T) {
		ctx := context.Background()
		_, poses, kantongs := x.reset()
		posPemasukan := poses[0]
		kantongMandiri := kantongs[0]

		// Prepare Data
		input := &keep_request.TransaksiInputUpdate{
			Waktu:         time.Now().Unix(),
			Jenis:         "pengeluaran",
			Jumlah:        10000,
			PosAsalId:     posPemasukan.Id,
			KantongAsalId: kantongMandiri.Id,
			Uraian:        "test",
		}
		m, err := x.services.Insert(ctx, input)
		assert.Nil(t, err)

		v, _ := x.repo.FindById(ctx, m.Id)
		assert.Equal(t, "pengeluaran", v.Jenis)

		// Progress
		input = &keep_request.TransaksiInputUpdate{
			Id:              m.Id,
			Waktu:           time.Now().Unix(),
			Jenis:           "pemasukan",
			Jumlah:          10000,
			PosTujuanId:     posPemasukan.Id,
			KantongTujuanId: kantongMandiri.Id,
			Uraian:          "test",
		}
		affected, err := x.services.Update(ctx, input)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		v, _ = x.repo.FindById(ctx, m.Id)
		assert.Equal(t, "pemasukan", v.Jenis)
	})
	t.Run("UpdatePengeluaranKeMutasiSuccess", func(t *testing.T) {
		ctx := context.Background()
		_, poses, kantongs := x.reset()
		posPemasukan := poses[0]
		posMain := poses[2]
		kantongMandiri := kantongs[0]
		kantongBca := kantongs[1]

		// Prepare Data
		input := &keep_request.TransaksiInputUpdate{
			Waktu:         time.Now().Unix(),
			Jenis:         "pengeluaran",
			Jumlah:        10000,
			PosAsalId:     posPemasukan.Id,
			KantongAsalId: kantongMandiri.Id,
			Uraian:        "test",
		}
		m, err := x.services.Insert(ctx, input)
		assert.Nil(t, err)

		v, _ := x.repo.FindById(ctx, m.Id)
		assert.Equal(t, "pengeluaran", v.Jenis)

		// Progress
		input = &keep_request.TransaksiInputUpdate{
			Id:              m.Id,
			Waktu:           time.Now().Unix(),
			Jenis:           "mutasi",
			Jumlah:          10000,
			PosAsalId:       posPemasukan.Id,
			PosTujuanId:     posMain.Id,
			KantongAsalId:   kantongMandiri.Id,
			KantongTujuanId: kantongBca.Id,
			Uraian:          "test",
		}
		affected, err := x.services.Update(ctx, input)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		v, _ = x.repo.FindById(ctx, m.Id)
		assert.Equal(t, "mutasi", v.Jenis)
	})
	t.Run("UpdateMutasiKePemasukanSuccess", func(t *testing.T) {
		ctx := context.Background()
		_, poses, kantongs := x.reset()
		posPemasukan := poses[0]
		posMain := poses[2]
		kantongMandiri := kantongs[0]
		kantongBca := kantongs[1]

		// Prepare Data
		input := &keep_request.TransaksiInputUpdate{
			Waktu:           time.Now().Unix(),
			Jenis:           "mutasi",
			Jumlah:          10000,
			PosAsalId:       posPemasukan.Id,
			PosTujuanId:     posMain.Id,
			KantongAsalId:   kantongMandiri.Id,
			KantongTujuanId: kantongBca.Id,
			Uraian:          "test",
		}
		m, err := x.services.Insert(ctx, input)
		assert.Nil(t, err)

		v, _ := x.repo.FindById(ctx, m.Id)
		assert.Equal(t, "mutasi", v.Jenis)

		// Progress
		input = &keep_request.TransaksiInputUpdate{
			Id:              m.Id,
			Waktu:           time.Now().Unix(),
			Jenis:           "pemasukan",
			Jumlah:          10000,
			PosTujuanId:     posPemasukan.Id,
			KantongTujuanId: kantongMandiri.Id,
			Uraian:          "test",
		}
		affected, err := x.services.Update(ctx, input)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		v, _ = x.repo.FindById(ctx, m.Id)
		assert.Equal(t, "pemasukan", v.Jenis)
	})
	t.Run("UpdateMutasiKePengeluaranSuccess", func(t *testing.T) {
		ctx := context.Background()
		_, poses, kantongs := x.reset()
		posPemasukan := poses[0]
		posMain := poses[2]
		kantongMandiri := kantongs[0]
		kantongBca := kantongs[1]

		// Prepare Data
		input := &keep_request.TransaksiInputUpdate{
			Waktu:           time.Now().Unix(),
			Jenis:           "mutasi",
			Jumlah:          10000,
			PosAsalId:       posPemasukan.Id,
			PosTujuanId:     posMain.Id,
			KantongAsalId:   kantongMandiri.Id,
			KantongTujuanId: kantongBca.Id,
			Uraian:          "test",
		}
		m, err := x.services.Insert(ctx, input)
		assert.Nil(t, err)

		v, _ := x.repo.FindById(ctx, m.Id)
		assert.Equal(t, "mutasi", v.Jenis)

		// Progress
		input = &keep_request.TransaksiInputUpdate{
			Id:            m.Id,
			Waktu:         time.Now().Unix(),
			Jenis:         "pengeluaran",
			Jumlah:        10000,
			PosAsalId:     posPemasukan.Id,
			KantongAsalId: kantongMandiri.Id,
			Uraian:        "test",
		}
		affected, err := x.services.Update(ctx, input)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		v, _ = x.repo.FindById(ctx, m.Id)
		assert.Equal(t, "pengeluaran", v.Jenis)
	})
	t.Run("InsertUpdateFailedCauseBasicValidation", func(t *testing.T) {
		ctx := context.Background()
		transaksis, _, _ := x.reset()
		trx := transaksis[0]

		tests := []map[string]string{
			// Jenis Tests
			{
				"err":   "jenis.required",
				"jenis": "",
			},
			{
				"err":   "jenis.oneof",
				"jenis": "pemasukann",
			},
			{
				"err":   "jenis.oneof",
				"jenis": "pengeluarann",
			},
			{
				"err":   "jenis.oneof",
				"jenis": "mutasis",
			},

			// Uraian Test
			{
				"err":    "uraian.required",
				"uraian": "",
			},

			// Pos
			{
				"err":         "pos_asal_id.required",
				"jenis":       "pengeluaran",
				"posAsalId":   "",
				"posTujuanId": "",
			},
			{
				"err":         "pos_asal_id.required",
				"jenis":       "mutasi",
				"posAsalId":   "",
				"posTujuanId": "",
			},
			{
				"err":         "pos_tujuan_id.excluded_if",
				"jenis":       "pengeluaran",
				"posAsalId":   "",
				"posTujuanId": "1",
			},
			{
				"err":         "pos_tujuan_id.required",
				"jenis":       "pemasukan",
				"posAsalId":   "",
				"posTujuanId": "",
			},
			{
				"err":         "pos_tujuan_id.required",
				"jenis":       "mutasi",
				"posAsalId":   "",
				"posTujuanId": "",
			},
			{
				"err":         "pos_asal_id.excluded_if",
				"jenis":       "pemasukan",
				"posAsalId":   "1",
				"posTujuanId": "",
			},
		}
		for _, v := range tests {

			input := &keep_request.TransaksiInputUpdate{
				Id:              trx.Id,
				Jenis:           getMap(v, "jenis", ""),
				PosAsalId:       getMap(v, "posAsalId", ""),
				PosTujuanId:     getMap(v, "posTujuanId", ""),
				Uraian:          getMap(v, "uraian", ""),
				KantongAsalId:   getMap(v, "kantongAsalId", ""),
				KantongTujuanId: getMap(v, "kantongTujuanId", ""),
				Keterangan:      getMap(v, "keterangan", ""),
				UrlFoto:         getMap(v, "urlFoto", ""),
				Lokasi:          getMap(v, "lokasi", ""),
			}
			_, err := x.services.Insert(ctx, input)
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, helpers_error.ValidationError)
			assert.ErrorContains(t, err, v["err"])
			_, err = x.services.Update(ctx, input)
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, helpers_error.ValidationError)
			assert.ErrorContains(t, err, v["err"])
		}
	})
	t.Run("InsertUpdateFailedCausePosIsNotLeaf", func(t *testing.T) {
		ctx := context.Background()
		transaksis, poses, _ := x.reset()
		trx := transaksis[1]
		posPengeluaranRoot := poses[1]

		jenis := "pemasukan"
		jumlah := 10000
		posTujuanId := posPengeluaranRoot.Id
		uraian := "test"
		input := &keep_request.TransaksiInputUpdate{
			Id:          trx.Id,
			Waktu:       time.Now().Unix(),
			Jenis:       jenis,
			Jumlah:      jumlah,
			PosTujuanId: posTujuanId,
			Uraian:      uraian,
		}
		_, err := x.services.Insert(ctx, input)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, helpers_error.ValidationError)
		assert.ErrorContains(t, err, "has children")

		affected, err := x.services.Update(ctx, input)
		assert.NotNil(t, err)
		assert.Equal(t, 0, affected)
		assert.ErrorIs(t, err, helpers_error.ValidationError)
		assert.ErrorContains(t, err, "has children")
	})
}

func NewTransaksiServicesTest() *TransaksiServicesTest {
	t := &TransaksiServicesTest{}
	t.setUp()
	return t
}

type TransaksiServicesTest struct {
	posRepo     keep_repo_interfaces.IPosRepository
	kantongRepo keep_repo_interfaces.IKantongRepository
	repo        keep_repo_interfaces.ITransaksiRepository
	services    keep_service_interfaces.ITransaksiServices
	truncate    func()
	cleanup     func()
}

func (x *TransaksiServicesTest) setUp() {
	x.setUpMemoryRepository()
	x.services = keep_services.NewTransaksiServices(x.repo, x.posRepo, x.kantongRepo)
}
func (x *TransaksiServicesTest) setUpMemoryRepository() {
	repo := keep_repos_memory.NewTransaksiMemoryRepository()
	posRepo := keep_repos_memory.NewPosMemoryRepository()
	kantongRepo := keep_repos_memory.NewKantongMemoryRepository()

	x.truncate = func() {
		repo.Data = make([]*keep_entities.Transaksi, 0)
		posRepo.Data = make([]*keep_entities.Pos, 0)
		kantongRepo.Data = make([]*keep_entities.Kantong, 0)
	}
	x.cleanup = func() {}
	x.repo = repo
	x.posRepo = posRepo
	x.kantongRepo = kantongRepo
}
func (x *TransaksiServicesTest) setUpMysqlRepository() {
	repo := keep_repos_mysql.NewTransaksiMySqlRepository()
	posRepo := keep_repos_memory.NewPosMemoryRepository()
	kantongRepo := keep_repos_memory.NewKantongMemoryRepository()

	x.truncate = func() {
		posRepo.Data = make([]*keep_entities.Pos, 0)
		kantongRepo.Data = make([]*keep_entities.Kantong, 0)

		models := repo.Get(context.Background(), keep_request.NewGetTransaksi())
		for _, m := range models {
			_, err := repo.SoftDeleteById(context.Background(), m.Id)
			helpers_error.PanicIfError(err)
		}
		models = repo.GetTrashed(context.Background())
		for _, m := range models {
			_, err := repo.HardDeleteTrashedById(context.Background(), m.Id)
			helpers_error.PanicIfError(err)
		}
	}
	x.cleanup = func() {
		repo.Cleanup()
	}
	x.repo = repo
	x.posRepo = posRepo
	x.kantongRepo = kantongRepo
}
func (x *TransaksiServicesTest) reset() ([]*keep_entities.Transaksi, []*keep_entities.Pos, []*keep_entities.Kantong) {
	x.truncate()
	ctx := context.Background()

	pInput := []*keep_entities.Pos{
		{
			Nama:   "Pemasukan",
			Urutan: 1,
			Saldo:  100000,
			IsShow: true,
			Status: "aktif",
		},
		{
			Nama:   "Pengeluaran",
			Urutan: 2,
			Saldo:  0,
			IsShow: true,
			Status: "aktif",
		},
	}
	poses := make([]*keep_entities.Pos, 0)
	for _, v := range pInput {
		m, _ := x.posRepo.Insert(ctx, v)
		poses = append(poses, m)
	}
	pInput = []*keep_entities.Pos{
		{
			Nama:     "Main",
			Urutan:   3,
			Saldo:    0,
			ParentId: poses[1].Id,
			IsShow:   true,
			Status:   "aktif",
		},
		{
			Nama:   "ZAM",
			Urutan: 4,
			Saldo:  0,
			IsShow: true,
			Status: "aktif",
		},
	}
	for _, v := range pInput {
		m, _ := x.posRepo.Insert(ctx, v)
		poses = append(poses, m)
	}

	kInput := []*keep_entities.Kantong{
		{
			Nama:   "Mandiri",
			Urutan: 1,
			Saldo:  100000,
			PosId:  poses[0].Id,
			IsShow: true,
			Status: "aktif",
		},
		{
			Nama:   "BCA",
			Urutan: 1,
			Saldo:  0,
			PosId:  poses[2].Id,
			IsShow: true,
			Status: "aktif",
		},
		{
			Nama:   "Jago",
			Urutan: 1,
			Saldo:  100000,
			PosId:  poses[3].Id,
			IsShow: true,
			Status: "aktif",
		},
	}
	kantongs := make([]*keep_entities.Kantong, 0)
	for _, v := range kInput {
		m, _ := x.kantongRepo.Insert(ctx, v)
		kantongs = append(kantongs, m)
	}

	tInput := []*keep_entities.Transaksi{
		{
			Waktu:             time.Now().Unix(),
			Jenis:             "pemasukan",
			Jumlah:            110000,
			PosTujuanId:       poses[0].Id,
			PosTujuanNama:     poses[0].Nama,
			KantongTujuanId:   kantongs[0].Id,
			KantongTujuanNama: kantongs[0].Nama,
			Uraian:            "Gajian",
			CreatedAt:         time.Now().Unix(),
			UpdatedAt:         time.Now().Unix(),
			Status:            "aktif",
		},
		{
			Waktu:           time.Now().Unix(),
			Jenis:           "pengeluaran",
			Jumlah:          10000,
			PosAsalId:       poses[0].Id,
			PosAsalNama:     poses[0].Nama,
			KantongAsalId:   kantongs[0].Id,
			KantongAsalNama: kantongs[0].Nama,
			Uraian:          "test",
			CreatedAt:       time.Now().Unix(),
			UpdatedAt:       time.Now().Unix(),
			Status:          "aktif",
		},
		{
			Waktu:           time.Now().Unix(),
			Jenis:           "pengeluaran",
			Jumlah:          10000,
			PosAsalId:       poses[0].Id,
			PosAsalNama:     poses[0].Nama,
			KantongAsalId:   kantongs[0].Id,
			KantongAsalNama: kantongs[0].Nama,
			Uraian:          "trashed",
			CreatedAt:       time.Now().Unix(),
			UpdatedAt:       time.Now().Unix(),
			Status:          "trashed",
		},
	}
	transaksis := make([]*keep_entities.Transaksi, 0)
	for _, v := range tInput {
		m, err := x.repo.Insert(ctx, v)
		helpers_error.PanicIfError(err)
		transaksis = append(transaksis, m)
	}

	return transaksis, poses, kantongs
}
func getMap(mp map[string]string, key, def string) string {
	v, ok := mp[key]
	if !ok {
		return def
	}
	return v
}
