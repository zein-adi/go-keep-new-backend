package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_events"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_services"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/handlers/keep_handlers_events"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/repos/keep_repos_memory"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/repos/keep_repos_mysql"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_env"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_events"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_requests"
	"testing"
	"time"
)

func TestKantong(t *testing.T) {
	helpers_env.Init(5)
	x := NewKantongServicesTest()
	defer x.cleanup()

	t.Run("CalculateSaldoAktif", func(t *testing.T) {
		_, ori := x.reset()
		oriKey := helpers.KeyBy(ori, func(d *keep_entities.Kantong) string {
			return d.Id
		})

		models := x.services.Get(context.Background(), helpers_requests.NewGet())
		assert.Len(t, models, 2)

		for _, m := range models {
			o := oriKey[m.Id]
			assert.Equal(t, o.Saldo-o.SaldoMengendap, m.CalculateSaldoAktif())
		}
	})
	t.Run("GetSuccess", func(t *testing.T) {
		_, ori := x.reset()
		oriKey := helpers.KeyBy(ori, func(d *keep_entities.Kantong) string {
			return d.Id
		})

		models := x.services.Get(context.Background(), helpers_requests.NewGet())
		assert.Len(t, models, 2)

		for _, m := range models {
			o := oriKey[m.Id]
			assert.Equal(t, o.Id, m.Id)
			assert.Equal(t, o.Nama, m.Nama)
			assert.Equal(t, o.Urutan, m.Urutan)
			assert.Equal(t, o.Saldo, m.Saldo)
			assert.Equal(t, o.SaldoMengendap, m.SaldoMengendap)
			assert.Equal(t, o.PosId, m.PosId)
			assert.Equal(t, o.IsShow, m.IsShow)
			assert.Equal(t, "aktif", m.Status)
		}
	})
	t.Run("InsertSuccess", func(t *testing.T) {
		poses, _ := x.reset()

		nama := "Mandiri"
		urutan := 1
		posId := poses[0].Id
		saldo := 200000
		saldoMengendap := 100000

		input := &keep_request.KantongInsert{
			Nama:           nama,
			Urutan:         urutan,
			Saldo:          saldo,
			SaldoMengendap: saldoMengendap,
			PosId:          posId,
		}
		m, err := x.services.Insert(context.Background(), input)
		assert.Nil(t, err)
		assert.NotEmpty(t, m.Id)
		assert.Equal(t, nama, m.Nama)
		assert.Equal(t, urutan, m.Urutan)
		assert.Equal(t, saldo, m.Saldo)
		assert.Equal(t, saldoMengendap, m.SaldoMengendap)
		assert.Equal(t, posId, m.PosId)
		assert.Equal(t, true, m.IsShow)
		assert.Equal(t, "aktif", m.Status)
	})
	t.Run("UpdateSuccess", func(t *testing.T) {
		poses, kantongs := x.reset()
		pos := poses[1]
		kantong := kantongs[0]

		id := kantong.Id
		nama := "BCA"
		urutan := 2
		posId := pos.Id
		saldo := 100000
		saldoMengendap := 0
		input := &keep_request.KantongUpdate{
			Id:             id,
			Nama:           nama,
			Urutan:         urutan,
			Saldo:          saldo,
			SaldoMengendap: saldoMengendap,
			PosId:          posId,
			IsShow:         false,
		}
		affected, err := x.services.Update(context.Background(), input)
		assert.Nil(t, err)
		assert.NotEmpty(t, affected)

		m, err := x.repo.FindById(context.Background(), id)
		assert.Nil(t, err)
		assert.NotEmpty(t, m.Id)
		assert.Equal(t, nama, m.Nama)
		assert.Equal(t, urutan, m.Urutan)
		assert.Equal(t, saldo, m.Saldo)
		assert.Equal(t, saldoMengendap, m.SaldoMengendap)
		assert.Equal(t, posId, m.PosId)
		assert.Equal(t, false, m.IsShow)
		assert.Equal(t, "aktif", m.Status)
	})
	t.Run("SoftDeleteSuccess", func(t *testing.T) {
		_, kantongs := x.reset()
		m := kantongs[0]

		affected, err := x.services.DeleteById(context.Background(), m.Id)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		_, err = x.repo.FindById(context.Background(), m.Id)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		model, err := x.repo.FindTrashedById(context.Background(), m.Id)
		assert.Nil(t, err)
		assert.Equal(t, "trashed", model.Status)
	})

	t.Run("GetTrashedSuccess", func(t *testing.T) {
		_, kantongs := x.reset()

		models := x.services.GetTrashed(context.Background(), helpers_requests.NewGet())
		assert.Len(t, models, 1)

		m := models[0]
		o := kantongs[2]
		assert.NotEmpty(t, m.Id)
		assert.Equal(t, o.Nama, m.Nama)
		assert.Equal(t, o.Urutan, m.Urutan)
		assert.Equal(t, o.Saldo, m.Saldo)
		assert.Equal(t, o.SaldoMengendap, m.SaldoMengendap)
		assert.Equal(t, o.PosId, m.PosId)
		assert.Equal(t, o.IsShow, m.IsShow)
		assert.Equal(t, "trashed", m.Status)
	})
	t.Run("RestoreTrashedSuccess", func(t *testing.T) {
		_, kantongs := x.reset()
		m := kantongs[2]

		affected, err := x.services.RestoreTrashedById(context.Background(), m.Id)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		model, err := x.repo.FindById(context.Background(), m.Id)
		assert.Nil(t, err)
		assert.Equal(t, "aktif", model.Status)
	})
	t.Run("DeleteTrashedSuccess", func(t *testing.T) {
		_, kantongs := x.reset()
		m := kantongs[2]

		affected, err := x.services.DeleteTrashedById(context.Background(), m.Id)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		_, err = x.repo.FindById(context.Background(), m.Id)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		_, err = x.repo.FindTrashedById(context.Background(), m.Id)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})

	t.Run("RestoreTrashedDeleteTrashedFailedCauseStatusAktif", func(t *testing.T) {
		_, kantongs := x.reset()
		id := kantongs[0].Id

		affected, err := x.services.RestoreTrashedById(context.Background(), id)
		assert.NotNil(t, err)
		assert.Empty(t, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		affected, err = x.services.DeleteTrashedById(context.Background(), id)
		assert.NotNil(t, err)
		assert.Empty(t, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
	t.Run("UpdateSoftDeleteFailedCauseTrashed", func(t *testing.T) {
		poses, kantongs := x.reset()
		pos := poses[0]
		id := kantongs[2].Id

		affected, err := x.services.DeleteById(context.Background(), id)
		assert.NotNil(t, err)
		assert.Empty(t, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		input := &keep_request.KantongUpdate{
			Id:             id,
			Nama:           "Mandiri",
			Urutan:         2,
			Saldo:          150000,
			SaldoMengendap: 50000,
			PosId:          pos.Id,
			IsShow:         false,
		}
		_, err = x.services.Update(context.Background(), input)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
	t.Run("UpdateSoftDeleteFailedCauseNotFound", func(t *testing.T) {
		poses, _ := x.reset()
		pos := poses[0]
		id := "9999"

		affected, err := x.services.DeleteById(context.Background(), id)
		assert.NotNil(t, err)
		assert.Empty(t, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		input := &keep_request.KantongUpdate{
			Id:             id,
			Nama:           "Mandiri",
			Urutan:         2,
			Saldo:          150000,
			SaldoMengendap: 50000,
			PosId:          pos.Id,
			IsShow:         false,
		}
		_, err = x.services.Update(context.Background(), input)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})

	t.Run("InsertFailedPosNotExists", func(t *testing.T) {
		x.reset()

		nama := "Mandiri"
		urutan := 1
		posId := "999999"
		saldo := 200000
		saldoMengendap := 100000

		input := &keep_request.KantongInsert{
			Nama:           nama,
			Urutan:         urutan,
			Saldo:          saldo,
			SaldoMengendap: saldoMengendap,
			PosId:          posId,
		}
		_, err := x.services.Insert(context.Background(), input)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
	t.Run("UpdateFailedPosNotExists", func(t *testing.T) {
		_, kantongs := x.reset()
		kantong := kantongs[0]

		id := kantong.Id
		nama := "BCA"
		urutan := 2
		posId := "999999"
		saldo := 100000
		saldoMengendap := 0
		input := &keep_request.KantongUpdate{
			Id:             id,
			Nama:           nama,
			Urutan:         urutan,
			Saldo:          saldo,
			SaldoMengendap: saldoMengendap,
			PosId:          posId,
			IsShow:         false,
		}
		affected, err := x.services.Update(context.Background(), input)
		assert.Equal(t, 0, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
	t.Run("UpdateUrutan", func(t *testing.T) {
		poses, kantongs := x.reset()
		kantongBca := kantongs[0]
		posPengeluaran := poses[1]
		ctx := context.Background()

		request := []*keep_request.KantongUpdateUrutanItem{
			{
				Id:     kantongBca.Id,
				Urutan: 99,
				PosId:  posPengeluaran.Id,
			},
		}
		affected, err := x.services.UpdateUrutan(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		kantongBca, err = x.repo.FindById(ctx, kantongBca.Id)
		assert.Nil(t, err)
		assert.Equal(t, 99, kantongBca.Urutan)
		assert.Equal(t, posPengeluaran.Id, kantongBca.PosId)
	})
	t.Run("UpdateVisibility", func(t *testing.T) {
		_, kantongs := x.reset()
		kantongBca := kantongs[0]
		ctx := context.Background()

		request := []*keep_request.KantongUpdateVisibilityItem{
			{
				Id:     kantongBca.Id,
				IsShow: false,
			},
		}
		affected, err := x.services.UpdateVisibility(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		kantongBca, err = x.repo.FindById(ctx, kantongBca.Id)
		assert.Nil(t, err)
		assert.Equal(t, false, kantongBca.IsShow)

		request = []*keep_request.KantongUpdateVisibilityItem{
			{
				Id:     kantongBca.Id,
				IsShow: true,
			},
		}
		affected, err = x.services.UpdateVisibility(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		kantongBca, err = x.repo.FindById(ctx, kantongBca.Id)
		assert.Nil(t, err)
		assert.Equal(t, true, kantongBca.IsShow)
	})

	/*
	 * Testing Listener
	 */
	l := keep_handlers_events.NewKantongEventListenerHandler(x.services, x.kantongHistoryServices)
	d := helpers_events.GetDispatcher()
	_ = d.Register(keep_events.TransaksiCreated, l.TransaksiCreated)
	_ = d.Register(keep_events.KantongHistoryCreated, l.KantongHistoryCreated)
	t.Run("UpdateFromTransakasi", func(t *testing.T) {
		_, kantongs := x.reset()
		ctx := context.Background()

		bca := kantongs[0]
		mandiri := kantongs[1]

		assert.Equal(t, 100000, bca.CalculateSaldoAktif())
		assert.Equal(t, 50000, mandiri.CalculateSaldoAktif())

		asalId := bca.Id
		tujuanId := mandiri.Id
		jumlah := 10000
		oldAsalId := ""
		oldTujuanId := ""
		oldJumlah := 0
		affected, err := x.services.UpdateSaldo(ctx, asalId, tujuanId, jumlah, oldAsalId, oldTujuanId, oldJumlah)
		assert.Nil(t, err)
		assert.Equal(t, 2, affected)

		bca, _ = x.repo.FindById(ctx, bca.Id)
		mandiri, _ = x.repo.FindById(ctx, mandiri.Id)
		assert.Equal(t, 90000, bca.CalculateSaldoAktif())
		assert.Equal(t, 60000, mandiri.CalculateSaldoAktif())
	})
	t.Run("ListenerTransaksiCreated", func(t *testing.T) {
		_, kantongs := x.reset()
		ctx := context.Background()

		bca := kantongs[0]
		mandiri := kantongs[1]

		assert.Equal(t, 100000, bca.CalculateSaldoAktif())
		assert.Equal(t, 50000, mandiri.CalculateSaldoAktif())

		asalId := bca.Id
		tujuanId := mandiri.Id
		jumlah := 10000
		_ = d.Dispatch(keep_events.TransaksiCreated, keep_events.TransaksiCreatedEventData{
			Data: keep_events.TransaksiEventData{
				KantongAsalId:   asalId,
				KantongTujuanId: tujuanId,
				Jumlah:          jumlah,
				Lokasi:          "",
				Details:         nil,
			},
		})
		time.Sleep(time.Millisecond * 1000)

		bca, _ = x.repo.FindById(ctx, bca.Id)
		mandiri, _ = x.repo.FindById(ctx, mandiri.Id)
		assert.Equal(t, 90000, bca.CalculateSaldoAktif())
		assert.Equal(t, 60000, mandiri.CalculateSaldoAktif())
	})
}

func NewKantongServicesTest() *KantongServicesTest {
	t := &KantongServicesTest{}
	t.setUp()
	return t
}

type KantongServicesTest struct {
	repo                   keep_repo_interfaces.IKantongRepository
	posRepo                keep_repo_interfaces.IPosRepository
	services               keep_service_interfaces.IKantongServices
	kantongHistoryServices keep_service_interfaces.IKantongHistoryServices
	truncate               func()
	cleanup                func()
}

func (x *KantongServicesTest) setUp() {
	x.setUpMemoryRepository()
	x.services = keep_services.NewKantongServices(x.repo, x.posRepo)

	kantongHistoryRepo := keep_repos_memory.NewKantongHistoryMemoryRepository()
	x.kantongHistoryServices = keep_services.NewKantongHistoryServices(kantongHistoryRepo, x.repo)
}
func (x *KantongServicesTest) setUpMemoryRepository() {
	x.posRepo = keep_repos_memory.NewPosMemoryRepository()
	repo := keep_repos_memory.NewKantongMemoryRepository()
	x.repo = repo
	x.cleanup = func() {}
	x.truncate = func() {
		repo.Data = make([]*keep_entities.Kantong, 0)
	}
}
func (x *KantongServicesTest) setUpMysqlRepository() {
	posRepo := keep_repos_mysql.NewPosMySqlRepository()
	x.posRepo = posRepo

	repo := keep_repos_mysql.NewKantongMysqlRepository()
	x.repo = repo

	x.cleanup = func() {
		posRepo.Cleanup()
		repo.Cleanup()
	}
	x.truncate = func() {
		request := helpers_requests.NewGet()
		request.Take = 1000

		for _, m := range posRepo.Get(context.Background()) {
			_, _ = posRepo.SoftDeleteById(context.Background(), m.Id)
		}
		for _, m := range posRepo.GetTrashed(context.Background()) {
			_, _ = posRepo.HardDeleteTrashedById(context.Background(), m.Id)
		}
		for _, m := range repo.Get(context.Background(), request) {
			_, _ = repo.SoftDeleteById(context.Background(), m.Id)
		}
		for _, m := range repo.GetTrashed(context.Background(), request) {
			_, _ = repo.HardDeleteTrashedById(context.Background(), m.Id)
		}
	}
}
func (x *KantongServicesTest) reset() ([]*keep_entities.Pos, []*keep_entities.Kantong) {
	x.truncate()

	posInput := []*keep_entities.Pos{
		{
			Nama:   "Pemasukan",
			Urutan: 1,
			Saldo:  1000,
			Status: "aktif",
		},
		{
			Nama:   "Pengeluaran",
			Urutan: 2,
			Saldo:  2000,
			Status: "aktif",
		},
	}
	poses := make([]*keep_entities.Pos, 0)
	for _, v := range posInput {
		m, _ := x.posRepo.Insert(context.Background(), v)
		poses = append(poses, m)
	}

	kantongInput := []*keep_entities.Kantong{
		{
			Nama:           "BCA",
			Urutan:         1,
			Saldo:          100000,
			SaldoMengendap: 0,
			PosId:          poses[0].Id,
			IsShow:         true,
			Status:         "aktif",
		},
		{
			Nama:           "Mandiri",
			Urutan:         1,
			Saldo:          150000,
			SaldoMengendap: 100000,
			PosId:          poses[0].Id,
			IsShow:         true,
			Status:         "aktif",
		},
		{
			Nama:           "BRI",
			Urutan:         1,
			Saldo:          70000,
			SaldoMengendap: 50000,
			PosId:          poses[0].Id,
			IsShow:         true,
			Status:         "trashed",
		},
	}
	kantongs := make([]*keep_entities.Kantong, 0)
	for _, v := range kantongInput {
		m, _ := x.repo.Insert(context.Background(), v)
		kantongs = append(kantongs, m)
	}
	return poses, kantongs
}
