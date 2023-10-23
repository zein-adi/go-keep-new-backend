package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_events"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_services"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/handlers/keep_handlers_events"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/repos/keep_repos_memory"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/repos/keep_repos_mysql"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_env"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_events"
	"testing"
	"time"
)

func TestBarang(t *testing.T) {
	helpers_env.Init(5)
	x := NewBarangServicesTest()
	defer x.cleanup()

	t.Run("GetSuccess", func(t *testing.T) {
		ori := x.reset()
		oriKey := helpers.KeyBy(ori, func(d *keep_entities.Barang) string {
			return d.Nama
		})

		models := x.services.Get(context.Background(), "", "")
		assert.Len(t, models, 2)

		for _, m := range models {
			o := oriKey[m.Nama]
			assert.Equal(t, o.Nama, m.Nama)
			assert.Equal(t, o.Harga, m.Harga)
			assert.Equal(t, o.Diskon, m.Diskon)
			assert.Equal(t, o.SatuanNama, m.SatuanNama)
			assert.Equal(t, o.SatuanJumlah, m.SatuanJumlah)
			assert.Equal(t, o.SatuanHarga, m.SatuanHarga)
			assert.Equal(t, o.Keterangan, m.Keterangan)
			assert.Equal(t, o.LastUpdate, m.LastUpdate)
			assert.Equal(t, o.Details, m.Details)
		}
	})
	t.Run("GetSearchSuccess", func(t *testing.T) {
		ori := x.reset()
		oriKey := helpers.KeyBy(ori, func(d *keep_entities.Barang) string {
			return d.Nama
		})

		models := x.services.Get(context.Background(), "sque", "")
		assert.Len(t, models, 1)

		for _, m := range models {
			o := oriKey[m.Nama]
			assert.Equal(t, o.Nama, m.Nama)
			assert.Equal(t, o.LastUpdate, m.LastUpdate)
			assert.Len(t, o.Details, 2)
		}

		models = x.services.Get(context.Background(), "sque", "indo")
		assert.Len(t, models, 1)

		for _, m := range models {
			o := oriKey[m.Nama]
			assert.Equal(t, o.Nama, m.Nama)
			assert.Equal(t, o.LastUpdate, m.LastUpdate)
			assert.Len(t, o.Details, 2)
			assert.Len(t, m.Details, 1)
		}
	})

	/*
	 * Testing Listener
	 */
	l := keep_handlers_events.NewBarangEventListenerHandler(x.services)
	d := helpers_events.GetDispatcher()
	_ = d.Register(keep_events.TransaksiCreated, l.TransaksiCreated)
	_ = d.Register(keep_events.TransaksiUpdated, l.TransaksiUpdated)
	_ = d.Register(keep_events.TransaksiSoftDeleted, l.TransaksiSoftDeleted)
	_ = d.Register(keep_events.TransaksiRestored, l.TransaksiRestored)

	t.Run("UpdateFromTransakasi", func(t *testing.T) {
		x.reset()
		ctx := context.Background()

		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:  "pengeluaran",
			Lokasi: "Hero",
			Waktu:  time.Now().Add(-time.Second).Unix(),
			Status: "aktif",
			Details: []*keep_entities.TransaksiDetail{
				{
					Uraian:       "Le Minerale 600ml",
					Harga:        3300,
					Jumlah:       1,
					Diskon:       0,
					SatuanNama:   "ml",
					SatuanJumlah: 600,
					SatuanHarga:  5.5,
					Keterangan:   "",
				},
			},
		})
		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:  "pengeluaran",
			Lokasi: "Avan",
			Waktu:  time.Now().Unix(),
			Status: "aktif",
			Details: []*keep_entities.TransaksiDetail{
				{
					Uraian:       "Le Minerale 600ml",
					Harga:        3300,
					Jumlah:       1,
					Diskon:       300,
					SatuanNama:   "ml",
					SatuanJumlah: 600,
					SatuanHarga:  5,
					Keterangan:   "",
				},
				{
					Uraian:       "Aqua 600ml",
					Harga:        3500,
					Jumlah:       1,
					Diskon:       0,
					SatuanNama:   "ml",
					SatuanJumlah: 600,
					SatuanHarga:  5.8333333333,
					Keterangan:   "",
				},
			},
		})
		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Waktu: time.Now().AddDate(0, -6, -1).Unix(),

			Jenis:  "pengeluaran",
			Lokasi: "Avan",
			Status: "aktif",
			Details: []*keep_entities.TransaksiDetail{
				{
					Uraian:       "Le Minerale 600ml",
					Harga:        2000,
					Jumlah:       1,
					Diskon:       0,
					SatuanNama:   "ml",
					SatuanJumlah: 600,
					SatuanHarga:  3.33,
					Keterangan:   "",
				},
			},
		})
		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Status: "trashed",
			Jenis:  "pengeluaran",
			Lokasi: "Avan",
			Waktu:  time.Now().Unix(),
			Details: []*keep_entities.TransaksiDetail{
				{
					Uraian:       "Le Minerale 600ml",
					Harga:        1000,
					Jumlah:       1,
					Diskon:       0,
					SatuanNama:   "ml",
					SatuanJumlah: 600,
					SatuanHarga:  1.67,
					Keterangan:   "",
				},
			},
		})

		models := x.repo.Get(ctx, "", "")
		assert.Len(t, models, 2)
		models = x.repo.Get(ctx, "yoghurt", "")
		assert.Len(t, models, 1)
		assert.Len(t, models[0].Details, 2)
		models = x.repo.Get(ctx, "yoghurt", "alfa")
		assert.Len(t, models, 1)
		assert.Len(t, models[0].Details, 1)

		affected, err := x.services.UpdateBarangFromTransaksi(ctx)
		assert.Nil(t, err)
		assert.Equal(t, 4, affected)

		models = x.repo.Get(ctx, "", "")
		assert.Len(t, models, 2)
		models = x.repo.Get(ctx, "miner", "")
		assert.Len(t, models, 1)
		assert.EqualValues(t, 5, models[0].SatuanHarga)
		models = x.repo.Get(ctx, "miner", "avan")
		assert.Len(t, models, 1)
		assert.Len(t, models[0].Details, 1)
	})
	t.Run("ListenerTransaksiCreated", func(t *testing.T) {
		x.truncate()
		ctx := context.Background()

		models := x.repo.Get(ctx, "", "")
		assert.Len(t, models, 0)

		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:  "pengeluaran",
			Lokasi: "Hero",
			Waktu:  time.Now().Add(-time.Second).Unix(),
			Status: "aktif",
			Details: []*keep_entities.TransaksiDetail{
				{
					Uraian:       "Le Minerale 600ml",
					Harga:        3300,
					Jumlah:       1,
					Diskon:       0,
					SatuanNama:   "ml",
					SatuanJumlah: 600,
					SatuanHarga:  5.5,
					Keterangan:   "",
				},
			},
		})

		_ = d.Dispatch(keep_events.TransaksiCreated, keep_events.TransaksiCreatedEventData{})
		time.Sleep(time.Millisecond * 10)

		models = x.repo.Get(ctx, "", "")
		assert.Len(t, models, 1)
		models = x.repo.Get(ctx, "miner", "")
		assert.Len(t, models, 1)
		assert.Len(t, models[0].Details, 1)
		assert.EqualValues(t, 5.5, models[0].SatuanHarga)
		models = x.repo.Get(ctx, "miner", "hero")
		assert.Len(t, models, 1)
		assert.Len(t, models[0].Details, 1)
	})
	t.Run("ListenerTransaksiUpdated", func(t *testing.T) {
		x.truncate()
		ctx := context.Background()

		transaksi, _ := x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:  "pengeluaran",
			Lokasi: "Avan",
			Waktu:  time.Now().Add(-time.Second).Unix(),
			Status: "aktif",
			Details: []*keep_entities.TransaksiDetail{
				{
					Uraian:       "Le Minerale 600ml",
					Harga:        3300,
					Jumlah:       1,
					Diskon:       0,
					SatuanNama:   "ml",
					SatuanJumlah: 600,
					SatuanHarga:  5.5,
					Keterangan:   "",
				},
			},
		})
		_, _ = x.services.UpdateBarangFromTransaksi(ctx)
		models := x.repo.Get(ctx, "minerale", "")
		assert.Len(t, models, 1)

		transaksi.Details[0].Uraian = "Aqua 600ml"
		transaksi.Details[0].Harga = 3000
		transaksi.Details[0].SatuanHarga = 5
		_ = d.Dispatch(keep_events.TransaksiUpdated, keep_events.TransaksiUpdatedEventData{})
		time.Sleep(time.Millisecond * 10)

		models = x.repo.Get(ctx, "", "")
		assert.Len(t, models, 1)
		models = x.repo.Get(ctx, "aqua", "")
		assert.Len(t, models, 1)
		assert.Len(t, models[0].Details, 1)
		assert.EqualValues(t, 5, models[0].SatuanHarga)
		models = x.repo.Get(ctx, "aqua", "avan")
		assert.Len(t, models, 1)
		assert.Len(t, models[0].Details, 1)
	})
	t.Run("ListenerTransaksiSoftDeleted", func(t *testing.T) {
		x.reset()
		ctx := context.Background()

		transaksi, _ := x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:  "pengeluaran",
			Lokasi: "Avan",
			Waktu:  time.Now().Add(-time.Second).Unix(),
			Status: "aktif",
			Details: []*keep_entities.TransaksiDetail{
				{
					Uraian:       "Le Minerale 600ml",
					Harga:        3300,
					Jumlah:       1,
					Diskon:       0,
					SatuanNama:   "ml",
					SatuanJumlah: 600,
					SatuanHarga:  5.5,
					Keterangan:   "",
				},
			},
		})
		_, _ = x.services.UpdateBarangFromTransaksi(ctx)
		models := x.repo.Get(ctx, "minerale", "")
		assert.Len(t, models, 1)

		transaksi.Status = "trashed"
		_ = d.Dispatch(keep_events.TransaksiUpdated, keep_events.TransaksiUpdatedEventData{})
		time.Sleep(time.Millisecond * 10)

		models = x.repo.Get(ctx, "", "")
		assert.Len(t, models, 0)
	})
	t.Run("ListenerTransaksiRestored", func(t *testing.T) {
		x.reset()
		ctx := context.Background()

		transaksi, _ := x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:  "pengeluaran",
			Lokasi: "Avan",
			Waktu:  time.Now().Add(-time.Second).Unix(),
			Status: "trashed",
			Details: []*keep_entities.TransaksiDetail{
				{
					Uraian:       "Le Minerale 600ml",
					Harga:        3300,
					Jumlah:       1,
					Diskon:       0,
					SatuanNama:   "ml",
					SatuanJumlah: 600,
					SatuanHarga:  5.5,
					Keterangan:   "",
				},
			},
		})
		_, _ = x.services.UpdateBarangFromTransaksi(ctx)
		models := x.repo.Get(ctx, "", "")
		assert.Len(t, models, 0)

		transaksi.Status = "aktif"
		_ = d.Dispatch(keep_events.TransaksiUpdated, keep_events.TransaksiUpdatedEventData{})
		time.Sleep(time.Millisecond * 10)

		models = x.repo.Get(ctx, "", "")
		assert.Len(t, models, 1)
		models = x.repo.Get(ctx, "minerale", "")
		assert.Len(t, models, 1)
		assert.Len(t, models[0].Details, 1)
		assert.EqualValues(t, 5.5, models[0].SatuanHarga)
		models = x.repo.Get(ctx, "minerale", "avan")
		assert.Len(t, models, 1)
		assert.Len(t, models[0].Details, 1)
	})
}

func NewBarangServicesTest() *BarangServicesTest {
	x := &BarangServicesTest{}
	x.setUp()
	return x
}

type BarangServicesTest struct {
	transaksiRepo keep_repo_interfaces.ITransaksiRepository
	repo          keep_repo_interfaces.IBarangRepository
	services      keep_service_interfaces.IBarangServices
	truncate      func()
	cleanup       func()
}

func (x *BarangServicesTest) setUp() {
	x.setUpMemoryRepository()

	x.services = keep_services.NewBarangServices(x.repo, x.transaksiRepo)
}
func (x *BarangServicesTest) setUpMemoryRepository() {
	transaksiRepo := keep_repos_memory.NewTransaksiMemoryRepository()
	x.transaksiRepo = transaksiRepo
	repo := keep_repos_memory.NewBarangMemoryRepository()
	x.repo = repo
	x.cleanup = func() {
	}
	x.truncate = func() {
		repo.Data = make([]*keep_entities.Barang, 0)
		transaksiRepo.Data = make([]*keep_entities.Transaksi, 0)
	}
}
func (x *BarangServicesTest) setUpMysqlRepository() {
	transaksiRepo := keep_repos_memory.NewTransaksiMemoryRepository()
	x.transaksiRepo = transaksiRepo
	repo := keep_repos_mysql.NewBarangMySqlRepository()
	x.repo = repo
	x.cleanup = func() {
		repo.Cleanup()
	}
	x.truncate = func() {
		transaksiRepo.Data = make([]*keep_entities.Transaksi, 0)

		models := repo.Get(context.Background(), "", "")
		for _, m := range models {
			_, _ = repo.DeleteByNama(context.Background(), m.Nama)
		}
	}
}
func (x *BarangServicesTest) reset() []*keep_entities.Barang {
	x.truncate()
	ctx := context.Background()

	barangs := []*keep_entities.Barang{
		{
			Nama:         "Cimory Yoghurt Drink Strawberry 120ml",
			Harga:        9100,
			Diskon:       100,
			SatuanNama:   "ml",
			SatuanJumlah: 120,
			SatuanHarga:  75,
			Keterangan:   "",
			LastUpdate:   time.Now().Unix(),
			Details: []*keep_entities.BarangDetail{
				{
					Lokasi:      "Alfamart",
					Harga:       9100,
					Diskon:      100,
					SatuanHarga: 75,
					Keterangan:  "",
				},
				{
					Lokasi:      "Indomaret",
					Harga:       9200,
					Diskon:      0,
					SatuanHarga: 76.67,
					Keterangan:  "",
				},
			},
		},
		{
			Nama:         "Cimory Squeeze Blueberry 150ml",
			Harga:        10500,
			Diskon:       400,
			SatuanNama:   "ml",
			SatuanJumlah: 150,
			SatuanHarga:  67.33,
			Keterangan:   "",
			LastUpdate:   time.Now().Unix(),
			Details: []*keep_entities.BarangDetail{
				{
					Lokasi:      "Alfamart",
					Harga:       10200,
					Diskon:      0,
					SatuanHarga: 68,
					Keterangan:  "",
				},
				{
					Lokasi:      "Indomaret",
					Harga:       10500,
					Diskon:      400,
					SatuanHarga: 67.33,
					Keterangan:  "",
				},
			},
		},
	}
	for _, pos := range barangs {
		_, err := x.repo.Insert(ctx, pos)
		helpers_error.PanicIfError(err)
	}
	return barangs
}
