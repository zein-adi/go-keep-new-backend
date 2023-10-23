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

func TestLokasi(t *testing.T) {
	helpers_env.Init(5)
	x := NewLokasiServicesTest()
	defer x.cleanup()

	t.Run("GetSuccess", func(t *testing.T) {
		ori := x.reset()
		oriKey := helpers.KeyBy(ori, func(d *keep_entities.Lokasi) string {
			return d.Nama
		})

		models := x.services.Get(context.Background(), "")
		assert.Len(t, models, 2)

		for _, m := range models {
			o := oriKey[m.Nama]
			assert.Equal(t, o.Nama, m.Nama)
			assert.Equal(t, o.LastUpdate, m.LastUpdate)
		}
	})
	t.Run("GetSearchSuccess", func(t *testing.T) {
		ori := x.reset()
		oriKey := helpers.KeyBy(ori, func(d *keep_entities.Lokasi) string {
			return d.Nama
		})

		models := x.services.Get(context.Background(), "citra")
		assert.Len(t, models, 1)

		for _, m := range models {
			o := oriKey[m.Nama]
			assert.Equal(t, o.Nama, m.Nama)
			assert.Equal(t, o.LastUpdate, m.LastUpdate)
		}

		models = x.services.Get(context.Background(), "av")
		assert.Len(t, models, 1)

		for _, m := range models {
			o := oriKey[m.Nama]
			assert.Equal(t, o.Nama, m.Nama)
			assert.Equal(t, o.LastUpdate, m.LastUpdate)
		}
	})

	/*
	 * Testing Listener
	 */
	l := keep_handlers_events.NewLokasiEventListenerHandler(x.services)
	d := helpers_events.GetDispatcher()
	_ = d.Register(keep_events.TransaksiCreated, l.TransaksiCreated)
	_ = d.Register(keep_events.TransaksiUpdated, l.TransaksiUpdated)
	_ = d.Register(keep_events.TransaksiSoftDeleted, l.TransaksiSoftDeleted)
	_ = d.Register(keep_events.TransaksiRestored, l.TransaksiRestored)

	t.Run("UpdateFromTransakasi", func(t *testing.T) {
		x.reset()
		ctx := context.Background()

		models := x.repo.Get(ctx, "")
		assert.Len(t, models, 2)
		models = x.repo.Get(ctx, "avan")
		assert.Len(t, models, 1)

		affected, err := x.services.UpdateLokasiFromTransaksi(ctx)
		assert.Nil(t, err)
		assert.Equal(t, 4, affected)

		models = x.repo.Get(ctx, "")
		assert.Len(t, models, 2)
		models = x.repo.Get(ctx, "hero")
		assert.Len(t, models, 1)
	})
	t.Run("ListenerTransaksiCreated", func(t *testing.T) {
		x.reset()
		ctx := context.Background()

		models := x.repo.Get(ctx, "")
		assert.Len(t, models, 2)
		models = x.repo.Get(ctx, "avan")
		assert.Len(t, models, 1)

		_ = d.Dispatch(keep_events.TransaksiCreated, keep_events.TransaksiCreatedEventData{})
		time.Sleep(time.Millisecond * 10)

		models = x.repo.Get(ctx, "")
		assert.Len(t, models, 2)
		models = x.repo.Get(ctx, "hero")
		assert.Len(t, models, 1)
	})
	t.Run("ListenerTransaksiUpdated", func(t *testing.T) {
		x.reset()
		ctx := context.Background()

		models := x.repo.Get(ctx, "")
		assert.Len(t, models, 2)
		models = x.repo.Get(ctx, "avan")
		assert.Len(t, models, 1)

		_ = d.Dispatch(keep_events.TransaksiUpdated, keep_events.TransaksiUpdatedEventData{})
		time.Sleep(time.Millisecond * 10)

		models = x.repo.Get(ctx, "")
		assert.Len(t, models, 2)
		models = x.repo.Get(ctx, "hero")
		assert.Len(t, models, 1)
	})
	t.Run("ListenerTransaksiSoftDeleted", func(t *testing.T) {
		x.reset()
		ctx := context.Background()

		models := x.repo.Get(ctx, "")
		assert.Len(t, models, 2)
		models = x.repo.Get(ctx, "avan")
		assert.Len(t, models, 1)

		_ = d.Dispatch(keep_events.TransaksiSoftDeleted, keep_events.TransaksiSoftDeletedEventData{})
		time.Sleep(time.Millisecond * 10)

		models = x.repo.Get(ctx, "")
		assert.Len(t, models, 2)
		models = x.repo.Get(ctx, "hero")
		assert.Len(t, models, 1)
	})
	t.Run("ListenerTransaksiRestored", func(t *testing.T) {
		x.reset()
		ctx := context.Background()

		models := x.repo.Get(ctx, "")
		assert.Len(t, models, 2)
		models = x.repo.Get(ctx, "avan")
		assert.Len(t, models, 1)

		_ = d.Dispatch(keep_events.TransaksiRestored, keep_events.TransaksiRestoredEventData{})
		time.Sleep(time.Millisecond * 10)

		models = x.repo.Get(ctx, "")
		assert.Len(t, models, 2)
		models = x.repo.Get(ctx, "hero")
		assert.Len(t, models, 1)
	})
}

func NewLokasiServicesTest() *LokasiServicesTest {
	x := &LokasiServicesTest{}
	x.setUp()
	return x
}

type LokasiServicesTest struct {
	transaksiRepo keep_repo_interfaces.ITransaksiRepository
	repo          keep_repo_interfaces.ILokasiRepository
	services      keep_service_interfaces.ILokasiServices
	truncate      func()
	cleanup       func()
}

func (x *LokasiServicesTest) setUp() {
	x.setUpMemoryRepository()

	x.services = keep_services.NewLokasiServices(x.repo, x.transaksiRepo)
}
func (x *LokasiServicesTest) setUpMemoryRepository() {
	transaksiRepo := keep_repos_memory.NewTransaksiMemoryRepository()
	x.transaksiRepo = transaksiRepo
	repo := keep_repos_memory.NewLokasiMemoryRepository()
	x.repo = repo
	x.cleanup = func() {
	}
	x.truncate = func() {
		repo.Data = make([]*keep_entities.Lokasi, 0)
		transaksiRepo.Data = make([]*keep_entities.Transaksi, 0)
	}
}
func (x *LokasiServicesTest) setUpMysqlRepository() {
	transaksiRepo := keep_repos_memory.NewTransaksiMemoryRepository()
	x.transaksiRepo = transaksiRepo
	repo := keep_repos_mysql.NewLokasiMySqlRepository()
	x.repo = repo
	x.cleanup = func() {
		repo.Cleanup()
	}
	x.truncate = func() {
		transaksiRepo.Data = make([]*keep_entities.Transaksi, 0)

		models := repo.Get(context.Background(), "")
		for _, m := range models {
			_, _ = repo.DeleteByNama(context.Background(), m.Nama)
		}
	}
}
func (x *LokasiServicesTest) reset() []*keep_entities.Lokasi {
	x.truncate()
	ctx := context.Background()

	_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
		Lokasi: "Hero",
		Waktu:  time.Now().Unix(),
		Status: "aktif",
	})
	_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
		Lokasi: "Dea Bakery",
		Waktu:  time.Now().Unix(),
		Status: "aktif",
	})
	_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
		Lokasi: "Superindo",
		Waktu:  time.Now().AddDate(0, -6, -1).Unix(),
		Status: "aktif",
	})
	_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
		Lokasi: "Superindo",
		Waktu:  time.Now().Unix(),
		Status: "trashed",
	})

	lokasis := []*keep_entities.Lokasi{
		{
			Nama:       "Citra Ken Dedes",
			LastUpdate: time.Now().Unix(),
		},
		{
			Nama:       "Avan",
			LastUpdate: time.Now().Unix(),
		},
	}
	for _, pos := range lokasis {
		_, err := x.repo.Insert(ctx, pos)
		helpers_error.PanicIfError(err)
	}
	return lokasis
}
