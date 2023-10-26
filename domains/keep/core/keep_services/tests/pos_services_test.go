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
	"testing"
	"time"
)

func TestPos(t *testing.T) {
	helpers_env.Init(5)
	x := NewPosServicesTest()
	defer x.cleanup()

	t.Run("GetSuccessAll", func(t *testing.T) {
		ori := x.reset()

		oriKey := helpers.KeyBy(ori, func(d *keep_entities.Pos) string {
			return d.Id
		})

		models := x.services.Get(context.Background(), keep_request.NewGetPos())
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
		ori := x.reset()

		oriKeyBy := helpers.KeyBy(ori, func(d *keep_entities.Pos) string {
			return d.Id
		})

		req := keep_request.NewGetPos()
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
		ori := x.reset()

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
		ori := x.reset()

		nama := "Mas Luxman"
		urutan := 1
		parentId := ori[0].Id
		input := &keep_request.PosInputUpdate{
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
		ori := x.reset()

		id := ori[0].Id
		nama := "Mas Luxman"
		urutan := 1
		parentId := ori[1].Id
		input := &keep_request.PosInputUpdate{
			Id:       id,
			Nama:     nama,
			Urutan:   urutan,
			ParentId: parentId,
			IsShow:   false,
		}
		affected, err := x.services.Update(context.Background(), input)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		model, err := x.repo.FindById(context.Background(), id)
		assert.Nil(t, err)
		assert.Equal(t, nama, model.Nama)
		assert.Equal(t, urutan, model.Urutan)
		assert.Equal(t, parentId, model.ParentId)
		assert.Equal(t, false, model.IsShow)
		assert.Equal(t, 10000, model.Saldo)
		assert.Equal(t, 2, model.Level)
	})
	t.Run("DeleteSuccess", func(t *testing.T) {
		ori := x.reset()

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
		ori := x.reset()

		affected, err := x.services.RestoreTrashedById(context.Background(), ori[3].Id)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		model, err := x.repo.FindById(context.Background(), ori[3].Id)
		assert.Nil(t, err)
		assert.Equal(t, "aktif", model.Status)
	})
	t.Run("DeleteTrashedSuccess", func(t *testing.T) {
		ori := x.reset()

		affected, err := x.services.DeleteTrashedById(context.Background(), ori[3].Id)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		_, err = x.repo.FindById(context.Background(), ori[3].Id)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		_, err = x.repo.FindTrashedById(context.Background(), ori[3].Id)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
	t.Run("UpdateFailedCauseParentToItself", func(t *testing.T) {
		ori := x.reset()

		pos := &keep_request.PosInputUpdate{
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
		ori := x.reset()

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
		ori := x.reset()

		id := ori[3].Id
		affected, err := x.services.DeleteById(context.Background(), id)
		assert.NotNil(t, err)
		assert.Empty(t, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		input := &keep_request.PosInputUpdate{
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
		x.reset()

		id := "9999"
		affected, err := x.services.DeleteById(context.Background(), id)
		assert.NotNil(t, err)
		assert.Empty(t, affected)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)

		input := &keep_request.PosInputUpdate{
			Id:     id,
			Nama:   "Pemasukan",
			Urutan: 1,
			IsShow: true,
		}
		_, err = x.services.Update(context.Background(), input)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
	t.Run("InsertSuccessLevel1", func(t *testing.T) {
		x.reset()

		nama := "Mas Luxman"
		urutan := 1
		input := &keep_request.PosInputUpdate{
			Nama:     nama,
			Urutan:   urutan,
			ParentId: "",
		}
		model, err := x.services.Insert(context.Background(), input)
		assert.Nil(t, err)
		assert.Equal(t, true, model.IsShow)
		assert.Equal(t, 1, model.Level)

		model, _ = x.repo.FindById(context.Background(), model.Id)
		assert.Equal(t, true, model.IsShow)
		assert.Equal(t, 1, model.Level)

		input = &keep_request.PosInputUpdate{
			Id:       model.Id,
			Nama:     nama,
			Urutan:   urutan,
			ParentId: "",
			IsShow:   true,
		}
		_, err = x.services.Update(context.Background(), input)
		assert.Nil(t, err)

		model, _ = x.repo.FindById(context.Background(), model.Id)
		assert.Equal(t, true, model.IsShow)
		assert.Equal(t, 1, model.Level)
	})
	t.Run("InsertSuccessLevel2", func(t *testing.T) {
		poses := x.reset()
		pos := poses[0]

		nama := "Mas Luxman"
		urutan := 1
		input := &keep_request.PosInputUpdate{
			Nama:     nama,
			Urutan:   urutan,
			ParentId: pos.Id,
		}
		model, err := x.services.Insert(context.Background(), input)
		assert.Nil(t, err)
		assert.Equal(t, true, model.IsShow)
		assert.Equal(t, 2, model.Level)

		model, _ = x.repo.FindById(context.Background(), model.Id)
		assert.Equal(t, true, model.IsShow)
		assert.Equal(t, 2, model.Level)
	})
	t.Run("InsertSuccessLevel3", func(t *testing.T) {
		poses := x.reset()
		pos := poses[2]

		nama := "Mas Luxman"
		urutan := 1
		input := &keep_request.PosInputUpdate{
			Nama:     nama,
			Urutan:   urutan,
			ParentId: pos.Id,
		}
		model, err := x.services.Insert(context.Background(), input)
		assert.Nil(t, err)
		assert.Equal(t, true, model.IsShow)
		assert.Equal(t, 3, model.Level)

		model, _ = x.repo.FindById(context.Background(), model.Id)
		assert.Equal(t, true, model.IsShow)
		assert.Equal(t, 3, model.Level)
	})
	t.Run("InsertSuccessIsLeafTrue", func(t *testing.T) {
		poses := x.reset()
		pos := poses[0]

		nama := "Mas Luxman"
		urutan := 1
		input := &keep_request.PosInputUpdate{
			Nama:     nama,
			Urutan:   urutan,
			ParentId: "",
		}
		model, err := x.services.Insert(context.Background(), input)
		assert.Nil(t, err)
		assert.Equal(t, true, model.IsLeaf)

		model, _ = x.repo.FindById(context.Background(), model.Id)
		assert.Equal(t, true, model.IsLeaf)

		x.reset()

		input = &keep_request.PosInputUpdate{
			Nama:     nama,
			Urutan:   urutan,
			ParentId: pos.ParentId,
		}
		model, err = x.services.Insert(context.Background(), input)
		assert.Nil(t, err)
		assert.Equal(t, true, model.IsLeaf)

		model, _ = x.repo.FindById(context.Background(), model.Id)
		assert.Equal(t, true, model.IsLeaf)
	})
	t.Run("UpdateUrutan", func(t *testing.T) {
		poses := x.reset()
		posMain := poses[2]
		ctx := context.Background()

		request := []*keep_request.PosUpdateUrutanItem{
			{
				Id:       posMain.Id,
				Urutan:   99,
				ParentId: "",
			},
		}
		affected, err := x.services.UpdateUrutan(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		posMain, err = x.repo.FindById(ctx, posMain.Id)
		assert.Nil(t, err)
		assert.Equal(t, 99, posMain.Urutan)
		assert.Equal(t, "", posMain.ParentId)
	})
	t.Run("UpdateUrutan", func(t *testing.T) {
		poses := x.reset()
		posMain := poses[2]
		ctx := context.Background()

		request := []*keep_request.PosUpdateVisibilityItem{
			{
				Id:     posMain.Id,
				IsShow: false,
			},
		}
		affected, err := x.services.UpdateVisibility(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		posMain, err = x.repo.FindById(ctx, posMain.Id)
		assert.Nil(t, err)
		assert.Equal(t, false, posMain.IsShow)

		request = []*keep_request.PosUpdateVisibilityItem{
			{
				Id:     posMain.Id,
				IsShow: true,
			},
		}
		affected, err = x.services.UpdateVisibility(ctx, request)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		posMain, err = x.repo.FindById(ctx, posMain.Id)
		assert.Nil(t, err)
		assert.Equal(t, true, posMain.IsShow)
	})

	/*
	 * Testing Listener
	 */
	l := keep_handlers_events.NewPosEventListenerHandler(x.services)
	d := helpers_events.GetDispatcher()
	_ = d.Register(keep_events.TransaksiCreated, l.TransaksiCreated)
	_ = d.Register(keep_events.TransaksiUpdated, l.TransaksiUpdated)
	_ = d.Register(keep_events.TransaksiSoftDeleted, l.TransaksiSoftDeleted)
	_ = d.Register(keep_events.TransaksiRestored, l.TransaksiRestored)

	t.Run("UpdateSaldoFromTransaksi", func(t *testing.T) {
		poses := x.reset()
		posPemasukan := poses[0]
		posPengeluaran := poses[1]
		posMain := poses[2]
		ctx := context.Background()

		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:       "pemasukan",
			Jumlah:      1000000,
			PosTujuanId: posPemasukan.Id,
			Status:      "aktif",
		})
		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:       "pemasukan",
			Jumlah:      100000,
			PosTujuanId: posPemasukan.Id,
			Status:      "aktif",
		})
		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:     "pengeluaran",
			Jumlah:    10000,
			PosAsalId: posPemasukan.Id,
			Status:    "aktif",
		})
		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:     "pengeluaran",
			Jumlah:    1000000,
			PosAsalId: posPemasukan.Id,
			Status:    "trashed",
		})
		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:       "mutasi",
			Jumlah:      10000,
			PosAsalId:   posPemasukan.Id,
			PosTujuanId: posMain.Id,
			Status:      "aktif",
		})

		affected, err := x.services.UpdateSaldoFromTransaksi(ctx, []string{posPemasukan.Id, posMain.Id})
		assert.Nil(t, err)
		assert.Equal(t, 3, affected)

		model, err := x.repo.FindById(ctx, posPemasukan.Id)
		assert.Nil(t, err)
		assert.Equal(t, 1080000, model.Saldo)

		model, err = x.repo.FindById(ctx, posMain.Id)
		assert.Nil(t, err)
		assert.Equal(t, 10000, model.Saldo)

		model, err = x.repo.FindById(ctx, posPengeluaran.Id)
		assert.Nil(t, err)
		assert.Equal(t, 10000, model.Saldo)
	})
	t.Run("ListenerTransaksiCreated", func(t *testing.T) {
		poses := x.reset()
		posPemasukan := poses[0]
		posPengeluaran := poses[1]
		posMain := poses[2]
		ctx := context.Background()

		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:       "pemasukan",
			Jumlah:      1000000,
			PosTujuanId: posPemasukan.Id,
			Status:      "aktif",
		})
		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:       "pemasukan",
			Jumlah:      100000,
			PosTujuanId: posPemasukan.Id,
			Status:      "aktif",
		})
		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:     "pengeluaran",
			Jumlah:    10000,
			PosAsalId: posPemasukan.Id,
			Status:    "aktif",
		})
		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:     "pengeluaran",
			Jumlah:    1000000,
			PosAsalId: posPemasukan.Id,
			Status:    "trashed",
		})
		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:       "mutasi",
			Jumlah:      10000,
			PosAsalId:   posPemasukan.Id,
			PosTujuanId: posMain.Id,
			Status:      "aktif",
		})

		_ = d.Dispatch(keep_events.TransaksiCreated, keep_events.TransaksiCreatedEventData{
			Data: keep_events.TransaksiEventData{
				PosAsalId:   posPemasukan.Id,
				PosTujuanId: posMain.Id,
			},
		})
		time.Sleep(time.Millisecond * 10)

		model, err := x.repo.FindById(ctx, posPemasukan.Id)
		assert.Nil(t, err)
		assert.Equal(t, 1080000, model.Saldo)

		model, err = x.repo.FindById(ctx, posMain.Id)
		assert.Nil(t, err)
		assert.Equal(t, 10000, model.Saldo)

		model, err = x.repo.FindById(ctx, posPengeluaran.Id)
		assert.Nil(t, err)
		assert.Equal(t, 10000, model.Saldo)
	})
	t.Run("ListenerTransaksiUpdated", func(t *testing.T) {
		poses := x.reset()
		posPemasukan := poses[0]
		posPengeluaran := poses[1]
		posMain := poses[2]
		ctx := context.Background()

		// Basic
		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:       "pemasukan",
			Jumlah:      1000000,
			PosTujuanId: posPemasukan.Id,
			Status:      "aktif",
		})
		transaksi, _ := x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:       "mutasi",
			Jumlah:      10000,
			PosAsalId:   posPemasukan.Id,
			PosTujuanId: posMain.Id,
			Status:      "aktif",
		})
		_, err := x.services.UpdateSaldoFromTransaksi(ctx, []string{posPemasukan.Id, posMain.Id})
		assert.Nil(t, err)

		model, _ := x.repo.FindById(ctx, posPemasukan.Id)
		assert.Equal(t, 990000, model.Saldo)
		model, _ = x.repo.FindById(ctx, posMain.Id)
		assert.Equal(t, 10000, model.Saldo)
		model, _ = x.repo.FindById(ctx, posPengeluaran.Id)
		assert.Equal(t, 10000, model.Saldo)

		// Update
		transaksi.Jenis = "pengeluaran"
		transaksi.PosAsalId = posPemasukan.Id
		transaksi.PosTujuanId = ""
		_ = d.Dispatch(keep_events.TransaksiUpdated, keep_events.TransaksiUpdatedEventData{
			Old: keep_events.TransaksiEventData{
				PosAsalId:   posPemasukan.Id,
				PosTujuanId: posMain.Id,
			},
			New: keep_events.TransaksiEventData{
				PosAsalId:   posPemasukan.Id,
				PosTujuanId: posMain.Id,
			},
		})
		time.Sleep(time.Millisecond * 10)

		model, _ = x.repo.FindById(ctx, posPemasukan.Id)
		assert.Equal(t, 990000, model.Saldo)
		model, _ = x.repo.FindById(ctx, posMain.Id)
		assert.Equal(t, 0, model.Saldo)
		model, _ = x.repo.FindById(ctx, posPengeluaran.Id)
		assert.Equal(t, 0, model.Saldo)
	})
	t.Run("ListenerTransaksiSoftDeleted", func(t *testing.T) {
		poses := x.reset()
		posPemasukan := poses[0]
		posPengeluaran := poses[1]
		posMain := poses[2]
		ctx := context.Background()

		// Basic
		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:       "pemasukan",
			Jumlah:      1000000,
			PosTujuanId: posPemasukan.Id,
			Status:      "aktif",
		})
		transaksi, _ := x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:       "mutasi",
			Jumlah:      10000,
			PosAsalId:   posPemasukan.Id,
			PosTujuanId: posMain.Id,
			Status:      "aktif",
		})
		_, err := x.services.UpdateSaldoFromTransaksi(ctx, []string{posPemasukan.Id, posMain.Id})
		assert.Nil(t, err)

		model, _ := x.repo.FindById(ctx, posPemasukan.Id)
		assert.Equal(t, 990000, model.Saldo)
		model, _ = x.repo.FindById(ctx, posMain.Id)
		assert.Equal(t, 10000, model.Saldo)
		model, _ = x.repo.FindById(ctx, posPengeluaran.Id)
		assert.Equal(t, 10000, model.Saldo)

		// Update
		transaksi.Status = "trashed"
		_ = d.Dispatch(keep_events.TransaksiUpdated, keep_events.TransaksiUpdatedEventData{
			Old: keep_events.TransaksiEventData{
				PosAsalId:   posPemasukan.Id,
				PosTujuanId: posMain.Id,
			},
			New: keep_events.TransaksiEventData{
				PosAsalId:   posPemasukan.Id,
				PosTujuanId: posMain.Id,
			},
		})
		time.Sleep(time.Millisecond * 10)

		model, _ = x.repo.FindById(ctx, posPemasukan.Id)
		assert.Equal(t, 1000000, model.Saldo)
		model, _ = x.repo.FindById(ctx, posMain.Id)
		assert.Equal(t, 0, model.Saldo)
		model, _ = x.repo.FindById(ctx, posPengeluaran.Id)
		assert.Equal(t, 0, model.Saldo)
	})
	t.Run("ListenerTransaksiRestored", func(t *testing.T) {
		poses := x.reset()
		posPemasukan := poses[0]
		posPengeluaran := poses[1]
		posMain := poses[2]
		ctx := context.Background()

		// Basic
		_, _ = x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:       "pemasukan",
			Jumlah:      1000000,
			PosTujuanId: posPemasukan.Id,
			Status:      "aktif",
		})
		transaksi, _ := x.transaksiRepo.Insert(ctx, &keep_entities.Transaksi{
			Jenis:       "mutasi",
			Jumlah:      10000,
			PosAsalId:   posPemasukan.Id,
			PosTujuanId: posMain.Id,
			Status:      "trashed",
		})
		_, err := x.services.UpdateSaldoFromTransaksi(ctx, []string{posPemasukan.Id, posMain.Id})
		assert.Nil(t, err)

		model, _ := x.repo.FindById(ctx, posPemasukan.Id)
		assert.Equal(t, 1000000, model.Saldo)
		model, _ = x.repo.FindById(ctx, posMain.Id)
		assert.Equal(t, 0, model.Saldo)
		model, _ = x.repo.FindById(ctx, posPengeluaran.Id)
		assert.Equal(t, 0, model.Saldo)

		// Update
		transaksi.Status = "aktif"
		_ = d.Dispatch(keep_events.TransaksiUpdated, keep_events.TransaksiUpdatedEventData{
			Old: keep_events.TransaksiEventData{
				PosAsalId:   posPemasukan.Id,
				PosTujuanId: posMain.Id,
			},
			New: keep_events.TransaksiEventData{
				PosAsalId:   posPemasukan.Id,
				PosTujuanId: posMain.Id,
			},
		})
		time.Sleep(time.Millisecond * 10)

		model, _ = x.repo.FindById(ctx, posPemasukan.Id)
		assert.Equal(t, 990000, model.Saldo)
		model, _ = x.repo.FindById(ctx, posMain.Id)
		assert.Equal(t, 10000, model.Saldo)
		model, _ = x.repo.FindById(ctx, posPengeluaran.Id)
		assert.Equal(t, 10000, model.Saldo)
	})

	_ = d.Register(keep_events.PosCreated, l.PosCreated)
	_ = d.Register(keep_events.PosUpdated, l.PosUpdated)
	_ = d.Register(keep_events.PosSoftDeleted, l.PosSoftDeleted)
	_ = d.Register(keep_events.PosRestored, l.PosRestored)
	t.Run("ListenerPosCreated", func(t *testing.T) {
		poses := x.reset()
		ctx := context.Background()

		posPemasukan := poses[0]
		assert.True(t, posPemasukan.IsLeaf)

		_, err := x.services.Insert(ctx, &keep_request.PosInputUpdate{
			Nama:     "Mas Luxman",
			Urutan:   1,
			ParentId: posPemasukan.Id,
			IsShow:   true,
		})
		assert.Nil(t, err)
		time.Sleep(time.Millisecond * 10)

		posPemasukan, err = x.repo.FindById(ctx, posPemasukan.Id)
		assert.Nil(t, err)
		assert.False(t, posPemasukan.IsLeaf)
	})
	t.Run("ListenerPosUpdated", func(t *testing.T) {
		poses := x.reset()
		ctx := context.Background()

		posPemasukan := poses[0]
		posPengeluaran := poses[1]
		posMain := poses[2]
		assert.True(t, posPemasukan.IsLeaf)
		assert.False(t, posPengeluaran.IsLeaf)
		assert.True(t, posMain.IsLeaf)

		_, err := x.services.Update(ctx, &keep_request.PosInputUpdate{
			Id:       posMain.Id,
			Nama:     "Mas Luxman",
			Urutan:   1,
			ParentId: posPemasukan.Id,
			IsShow:   true,
		})
		assert.Nil(t, err)
		time.Sleep(time.Millisecond * 10)

		posPemasukan, err = x.repo.FindById(ctx, posPemasukan.Id)
		assert.Nil(t, err)
		assert.False(t, posPemasukan.IsLeaf)

		posPengeluaran, err = x.repo.FindById(ctx, posPengeluaran.Id)
		assert.Nil(t, err)
		assert.True(t, posPengeluaran.IsLeaf)
	})
	t.Run("ListenerPosSoftDeleted", func(t *testing.T) {
		poses := x.reset()
		ctx := context.Background()

		posPemasukan := poses[0]
		posPengeluaran := poses[1]
		posMain := poses[2]
		assert.True(t, posPemasukan.IsLeaf)
		assert.False(t, posPengeluaran.IsLeaf)
		assert.True(t, posMain.IsLeaf)

		_, err := x.services.DeleteById(ctx, posMain.Id)
		assert.Nil(t, err)
		time.Sleep(time.Millisecond * 10)

		posPemasukan, _ = x.repo.FindById(ctx, posPemasukan.Id)
		assert.True(t, posPemasukan.IsLeaf)
		posPengeluaran, _ = x.repo.FindById(ctx, posPengeluaran.Id)
		assert.True(t, posPengeluaran.IsLeaf)
	})
	t.Run("ListenerPosRestored", func(t *testing.T) {
		poses := x.reset()
		ctx := context.Background()

		posPengeluaran := poses[1]
		posMain := poses[2]

		_, _ = x.services.DeleteById(ctx, posMain.Id)
		time.Sleep(time.Millisecond * 10)

		posPengeluaran, _ = x.repo.FindById(ctx, posPengeluaran.Id)
		assert.True(t, posPengeluaran.IsLeaf)
		posMain, _ = x.repo.FindTrashedById(ctx, posMain.Id)
		assert.True(t, posMain.IsLeaf)
		assert.Equal(t, "trashed", posMain.Status)

		_, err := x.services.RestoreTrashedById(ctx, posMain.Id)
		assert.Nil(t, err)
		time.Sleep(time.Millisecond * 10)

		posPengeluaran, _ = x.repo.FindById(ctx, posPengeluaran.Id)
		assert.False(t, posPengeluaran.IsLeaf)
		posMain, _ = x.repo.FindById(ctx, posMain.Id)
		assert.True(t, posMain.IsLeaf)
	})
}

func NewPosServicesTest() *PosServicesTest {
	t := &PosServicesTest{}
	t.setUp()
	return t
}

type PosServicesTest struct {
	repo          keep_repo_interfaces.IPosRepository
	transaksiRepo keep_repo_interfaces.ITransaksiRepository
	services      keep_service_interfaces.IPosServices
	truncate      func()
	cleanup       func()
}

func (x *PosServicesTest) setUp() {
	x.setUpMemoryRepository()
	x.services = keep_services.NewPosServices(x.repo, x.transaksiRepo)
}
func (x *PosServicesTest) setUpMemoryRepository() {
	transaksiRepo := keep_repos_memory.NewTransaksiMemoryRepository()
	x.transaksiRepo = transaksiRepo
	repo := keep_repos_memory.NewPosMemoryRepository()
	x.repo = repo

	x.cleanup = func() {
	}
	x.truncate = func() {
		transaksiRepo.Data = make([]*keep_entities.Transaksi, 0)
		repo.Data = make([]*keep_entities.Pos, 0)
	}
}
func (x *PosServicesTest) setUpMysqlRepository() {
	transaksiRepo := keep_repos_memory.NewTransaksiMemoryRepository()
	x.transaksiRepo = transaksiRepo
	repo := keep_repos_mysql.NewPosMySqlRepository()
	x.repo = repo

	x.cleanup = func() {
		repo.Cleanup()
	}
	x.truncate = func() {
		transaksiRepo.Data = make([]*keep_entities.Transaksi, 0)

		models := repo.Get(context.Background(), keep_request.NewGetPos())
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
}
func (x *PosServicesTest) reset() []*keep_entities.Pos {
	x.truncate()

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
