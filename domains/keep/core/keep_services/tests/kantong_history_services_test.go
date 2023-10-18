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

func TestKantongHistory(t *testing.T) {
	helpers_env.Init(5)
	x := NewKantongHistoryServicesTest()
	defer x.cleanup()

	t.Run("GetSuccess", func(t *testing.T) {
		_, ori := x.reset()
		oriKey := helpers.KeyBy(ori, func(d *keep_entities.KantongHistory) string {
			return d.Id
		})

		models := x.services.Get(context.Background())
		assert.Len(t, models, 2)

		for _, m := range models {
			o := oriKey[m.Id]
			assert.Equal(t, o.Id, m.Id)
			assert.Equal(t, o.KantongId, m.KantongId)
			assert.Equal(t, o.Jumlah, m.Jumlah)
			assert.Equal(t, o.Uraian, m.Uraian)
			assert.Equal(t, o.Waktu, m.Waktu)
		}
	})
	t.Run("InsertSuccess", func(t *testing.T) {
		kantong, _ := x.reset()

		kantongId := kantong.Id
		jumlah := 10000
		uraian := ""

		input := &keep_request.KantongHistoryInsertUpdate{
			KantongId: kantongId,
			Jumlah:    jumlah,
			Uraian:    uraian,
		}
		m, err := x.services.Insert(context.Background(), input)
		assert.Nil(t, err)
		assert.NotEmpty(t, m.Id)
		assert.Equal(t, kantongId, m.KantongId)
		assert.Equal(t, jumlah, m.Jumlah)
		assert.Equal(t, uraian, m.Uraian)
		assert.NotEmpty(t, m.Waktu)
	})
	t.Run("UpdateSuccess", func(t *testing.T) {
		kantong, kantongHistories := x.reset()

		id := kantongHistories[0].Id
		kantongId := kantong.Id
		jumlah := 20000
		uraian := "edited"

		input := &keep_request.KantongHistoryInsertUpdate{
			Id:        id,
			KantongId: kantongId,
			Jumlah:    jumlah,
			Uraian:    uraian,
		}
		affected, err := x.services.Update(context.Background(), input)
		assert.Nil(t, err)
		assert.NotEmpty(t, affected)

		m, err := x.repo.FindById(context.Background(), id)
		assert.Nil(t, err)
		assert.NotEmpty(t, m.Id)
		assert.Equal(t, kantongId, m.KantongId)
		assert.Equal(t, jumlah, m.Jumlah)
		assert.Equal(t, uraian, m.Uraian)
		assert.NotEmpty(t, m.Waktu)
	})
	t.Run("DeleteSuccess", func(t *testing.T) {
		_, kantongHistories := x.reset()
		m := kantongHistories[0]

		affected, err := x.services.DeleteById(context.Background(), m.Id)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		_, err = x.repo.FindById(context.Background(), m.Id)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
	t.Run("UpdateDeleteFailedCauseNotFound", func(t *testing.T) {
		kantong, _ := x.reset()
		id := "9999"

		affected, err := x.services.DeleteById(context.Background(), id)
		assert.NotNil(t, err)
		assert.Empty(t, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		input := &keep_request.KantongHistoryInsertUpdate{
			Id:        id,
			KantongId: kantong.Id,
			Jumlah:    0,
			Uraian:    "1",
		}
		_, err = x.services.Update(context.Background(), input)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
}

func NewKantongHistoryServicesTest() *KantongHistoryServicesTest {
	t := &KantongHistoryServicesTest{}
	t.setUp()
	return t
}

type KantongHistoryServicesTest struct {
	repo        keep_repo_interfaces.IKantongHistoryRepository
	kantongRepo keep_repo_interfaces.IKantongRepository
	services    keep_service_interfaces.IKantongHistoryServices
	truncate    func()
	cleanup     func()
	posRepo     keep_repo_interfaces.IPosRepository
}

func (x *KantongHistoryServicesTest) setUp() {
	x.setUpMemoryRepository()
	x.posRepo = keep_repos_memory.NewPosMemoryRepository()
	x.services = keep_services.NewKantongHistoryServices(x.repo, x.kantongRepo)
}
func (x *KantongHistoryServicesTest) setUpMemoryRepository() {
	kantongRepo := keep_repos_memory.NewKantongMemoryRepository()
	x.kantongRepo = kantongRepo
	repo := keep_repos_memory.NewKantongHistoryMemoryRepository()
	x.repo = repo
	x.cleanup = func() {}
	x.truncate = func() {
		kantongRepo.Data = make([]*keep_entities.Kantong, 0)
		repo.Data = make([]*keep_entities.KantongHistory, 0)
	}
}
func (x *KantongHistoryServicesTest) setUpMysqlRepository() {
	kantongRepo := keep_repos_mysql.NewKantongMysqlRepository()
	x.kantongRepo = kantongRepo
	repo := keep_repos_mysql.NewKantongHistoryMysqlRepository()
	x.repo = repo

	x.cleanup = func() {
		kantongRepo.Cleanup()
		repo.Cleanup()
	}
	x.truncate = func() {
		for _, m := range kantongRepo.Get(context.Background()) {
			_, _ = kantongRepo.SoftDeleteById(context.Background(), m.Id)
		}
		for _, m := range kantongRepo.GetTrashed(context.Background()) {
			_, _ = kantongRepo.HardDeleteTrashedById(context.Background(), m.Id)
		}
		for _, m := range repo.Get(context.Background()) {
			_, _ = repo.DeleteById(context.Background(), m.Id)
		}
	}
}
func (x *KantongHistoryServicesTest) reset() (*keep_entities.Kantong, []*keep_entities.KantongHistory) {
	x.truncate()

	posInput := &keep_entities.Pos{
		Nama:   "Pemasukan",
		Urutan: 1,
		Saldo:  1000,
		Status: "aktif",
	}
	pos, _ := x.posRepo.Insert(context.Background(), posInput)

	kantongInput := &keep_entities.Kantong{
		Nama:           "BCA",
		Urutan:         1,
		Saldo:          100000,
		SaldoMengendap: 0,
		PosId:          pos.Id,
		IsShow:         true,
		Status:         "aktif",
	}
	kantong, _ := x.kantongRepo.Insert(context.Background(), kantongInput)

	kantongHistoryInput := []*keep_entities.KantongHistory{
		{
			KantongId: kantong.Id,
			Jumlah:    100000,
			Uraian:    "pemasukan",
			Waktu:     time.Now().Unix(),
		},
		{
			KantongId: kantong.Id,
			Jumlah:    -50000,
			Uraian:    "pengeluaran",
			Waktu:     time.Now().Unix(),
		},
	}
	kantongHistories := make([]*keep_entities.KantongHistory, 0)
	for _, v := range kantongHistoryInput {
		m, _ := x.repo.Insert(context.Background(), v)
		kantongHistories = append(kantongHistories, m)
	}
	return kantong, kantongHistories
}
