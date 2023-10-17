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
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"testing"
)

func TestPos(t *testing.T) {
	x := NewPosServicesTest()

	t.Run("GetSuccessAll", func(t *testing.T) {
		ori := x.setUpAndPopulate()
		defer x.dbCleanup()

		oriKey := helpers.KeyBy(ori, func(d *keep_entities.Pos) string {
			return d.Id
		})

		models := x.services.Get(context.Background(), keep_request.NewPosGetRequest())
		assert.Len(t, models, 3)

		for _, m := range models {
			o := oriKey[m.Id]
			assert.Equal(t, o.Id, m.Id)
			assert.Equal(t, o.Nama, m.Nama)
			assert.Equal(t, o.Urutan, m.Urutan)
			assert.Equal(t, o.Saldo, m.Saldo) // 0
			assert.Equal(t, o.ParentId, m.ParentId)
			assert.Equal(t, o.Level, m.Level)
			assert.Equal(t, true, m.IsShow)
			assert.Equal(t, "aktif", m.Status)
		}
	})
	t.Run("GetSuccessLeafOnly", func(t *testing.T) {
		ori := x.setUpAndPopulate()
		defer x.dbCleanup()

		oriKeyBy := helpers.KeyBy(ori, func(d *keep_entities.Pos) string {
			return d.Id
		})

		req := keep_request.NewPosGetRequest()
		req.IsLeafOnly = true
		models := x.services.Get(context.Background(), req)
		assert.Len(t, models, 2)

		for _, m := range models {
			o := oriKeyBy[m.Id]
			assert.Equal(t, o.Id, m.Id)
			assert.Equal(t, o.Nama, m.Nama)
			assert.Equal(t, o.Urutan, m.Urutan)
			assert.Equal(t, o.Saldo, m.Saldo) // 0
			assert.Equal(t, o.ParentId, m.ParentId)
			assert.Equal(t, o.Level, m.Level)
			assert.Equal(t, true, m.IsShow)
			assert.Equal(t, "aktif", m.Status)
		}
	})
	t.Run("GetTrashedSuccess", func(t *testing.T) {
		ori := x.setUpAndPopulate()
		defer x.dbCleanup()

		models := x.services.GetTrashed(context.Background())
		assert.Len(t, models, 1)

		m := models[0]
		o := ori[3]
		assert.Equal(t, o.Id, m.Id)
		assert.Equal(t, o.Nama, m.Nama)
		assert.Equal(t, o.Urutan, m.Urutan)
		assert.Equal(t, o.Saldo, m.Saldo) // 0
		assert.Equal(t, o.ParentId, m.ParentId)
		assert.Equal(t, o.Level, m.Level)
		assert.Equal(t, o.IsShow, m.IsShow)
		assert.Equal(t, "trashed", m.Status)
	})
	t.Run("InsertSuccess", func(t *testing.T) {
		ori := x.setUpAndPopulate()
		defer x.dbCleanup()

		nama := "Mas Luxman"
		urutan := 1
		parentId := ori[0].Id
		input := &keep_request.PosInputUpdateRequest{
			Nama:     nama,
			Urutan:   urutan,
			ParentId: parentId,
		}
		model, err := x.services.Insert(context.Background(), input)
		assert.Nil(t, err)
		assert.Equal(t, nama, model.Nama)
		assert.Equal(t, urutan, model.Urutan)
		assert.Equal(t, parentId, model.ParentId)
		assert.Equal(t, true, model.IsShow)
		assert.Equal(t, 0, model.Saldo)
		assert.Equal(t, 2, model.Level)
	})
	t.Run("UpdateSuccess", func(t *testing.T) {
		ori := x.setUpAndPopulate()
		defer x.dbCleanup()

		id := ori[0].Id
		nama := "Mas Luxman"
		urutan := 1
		parentId := ori[1].Id
		input := &keep_request.PosInputUpdateRequest{
			Id:       id,
			Nama:     nama,
			Urutan:   urutan,
			ParentId: parentId,
			IsShow:   false,
		}
		model, err := x.services.Update(context.Background(), input)
		assert.Nil(t, err)
		assert.Equal(t, nama, model.Nama)
		assert.Equal(t, urutan, model.Urutan)
		assert.Equal(t, parentId, model.ParentId)
		assert.Equal(t, false, model.IsShow)
		assert.Equal(t, 10000, model.Saldo)
		assert.Equal(t, 2, model.Level)
	})
	t.Run("DeleteSuccess", func(t *testing.T) {
		ori := x.setUpAndPopulate()
		defer x.dbCleanup()

		affected, err := x.services.DeleteById(context.Background(), ori[0].Id)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		_, err = x.repo.FindById(context.Background(), ori[0].Id)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		model, err := x.repo.FindTrashedById(context.Background(), ori[0].Id)
		assert.Nil(t, err)
		assert.Equal(t, "trashed", model.Status)
	})
	t.Run("RestoreTrashedSuccess", func(t *testing.T) {
		ori := x.setUpAndPopulate()
		defer x.dbCleanup()

		affected, err := x.services.RestoreTrashedById(context.Background(), ori[3].Id)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		model, err := x.repo.FindById(context.Background(), ori[3].Id)
		assert.Nil(t, err)
		assert.Equal(t, "aktif", model.Status)
	})
	t.Run("DeleteTrashedSuccess", func(t *testing.T) {
		ori := x.setUpAndPopulate()
		defer x.dbCleanup()

		affected, err := x.services.DeleteTrashedById(context.Background(), ori[3].Id)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		_, err = x.repo.FindById(context.Background(), ori[3].Id)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		_, err = x.repo.FindTrashedById(context.Background(), ori[3].Id)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
	t.Run("UpdateFailedCauseParentToItself", func(t *testing.T) {
		ori := x.setUpAndPopulate()
		defer x.dbCleanup()

		pos := &keep_request.PosInputUpdateRequest{
			Id:       ori[0].Id,
			ParentId: ori[0].Id,
			Nama:     "Pemasukan",
			Urutan:   1,
		}
		_, err := x.services.Update(context.Background(), pos)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, helpers_error.ValidationError)
		assert.ErrorContains(t, err, "parent_to_self")
	})
	t.Run("RestoreTrashedDeleteTrashedFailedCauseStatusAktif", func(t *testing.T) {
		ori := x.setUpAndPopulate()
		defer x.dbCleanup()

		id := ori[0].Id
		affected, err := x.services.RestoreTrashedById(context.Background(), id)
		assert.NotNil(t, err)
		assert.Empty(t, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		affected, err = x.services.DeleteTrashedById(context.Background(), id)
		assert.NotNil(t, err)
		assert.Empty(t, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
	t.Run("UpdateDeleteFailedCauseTrashed", func(t *testing.T) {
		ori := x.setUpAndPopulate()
		defer x.dbCleanup()

		id := ori[3].Id
		affected, err := x.services.DeleteById(context.Background(), id)
		assert.NotNil(t, err)
		assert.Empty(t, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		input := &keep_request.PosInputUpdateRequest{
			Id:     id,
			Nama:   "Pemasukan",
			Urutan: 1,
			IsShow: true,
		}
		_, err = x.services.Update(context.Background(), input)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
	t.Run("UpdateDeleteFailedCauseNotFound", func(t *testing.T) {
		x.setUpAndPopulate()
		defer x.dbCleanup()

		id := "9999"
		affected, err := x.services.DeleteById(context.Background(), id)
		assert.NotNil(t, err)
		assert.Empty(t, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		input := &keep_request.PosInputUpdateRequest{
			Id:     id,
			Nama:   "Pemasukan",
			Urutan: 1,
			IsShow: true,
		}
		_, err = x.services.Update(context.Background(), input)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
}

func NewPosServicesTest() *PosServicesTest {
	return &PosServicesTest{}
}

type PosServicesTest struct {
	repo      keep_repo_interfaces.IPosRepository
	services  keep_service_interfaces.IPosServices
	dbCleanup func()
}

func (x *PosServicesTest) setUp() {
	x.setUpMemoryRepository()
	x.services = keep_services.NewPosServices(x.repo)
}
func (x *PosServicesTest) setUpMemoryRepository() {
	x.repo = keep_repos_memory.NewPosMemoryRepository()
	x.dbCleanup = func() {}
}
func (x *PosServicesTest) setUpMysqlRepository() {
	repo := keep_repos_mysql.NewPosMySqlRepository()
	x.dbCleanup = repo.Cleanup
	x.repo = repo

	models := repo.Get(context.Background(), keep_request.NewPosGetRequest())
	for _, m := range models {
		_, _ = repo.SoftDeleteById(context.Background(), m.Id)
	}
	models = repo.GetTrashed(context.Background())
	for _, m := range models {
		_, _ = repo.DeleteById(context.Background(), m.Id)
	}
}
func (x *PosServicesTest) setUpAndPopulate() []*keep_entities.Pos {
	x.setUp()
	posInput := []*keep_entities.Pos{
		{
			Nama:     "Pemasukan",
			Urutan:   1,
			Saldo:    10000,
			ParentId: "",
			Level:    1,
			IsShow:   true,
			Status:   "aktif",
			IsLeaf:   true,
		},
		{
			Nama:     "Pengeluaran",
			Urutan:   2,
			Saldo:    5000,
			ParentId: "",
			Level:    1,
			IsShow:   true,
			Status:   "aktif",
			IsLeaf:   false,
		},
	}
	poses := make([]*keep_entities.Pos, 0)
	for _, pos := range posInput {
		m, _ := x.repo.Insert(context.Background(), pos)
		poses = append(poses, m)
	}
	posInput = []*keep_entities.Pos{
		{
			Nama:     "Main",
			Urutan:   1,
			Saldo:    5000,
			ParentId: poses[1].Id,
			Level:    2,
			IsShow:   true,
			Status:   "aktif",
			IsLeaf:   true,
		},
		{
			Nama:   "Trashed",
			Urutan: 1,
			Level:  1,
			Status: "trashed",
		},
	}
	for _, pos := range posInput {
		m, _ := x.repo.Insert(context.Background(), pos)
		poses = append(poses, m)
	}
	return poses
}
