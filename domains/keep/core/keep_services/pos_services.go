package keep_services

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_events"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_events"
	"github.com/zein-adi/go-keep-new-backend/helpers/validator"
	"time"
)

func NewPosServices(repo keep_repo_interfaces.IPosRepository, transaksiRepo keep_repo_interfaces.ITransaksiRepository) *PosServices {
	return &PosServices{
		repo:          repo,
		transaksiRepo: transaksiRepo,
	}
}

type PosServices struct {
	repo          keep_repo_interfaces.IPosRepository
	transaksiRepo keep_repo_interfaces.ITransaksiRepository
}

func (x *PosServices) Get(ctx context.Context, request *keep_request.GetPos) []*keep_entities.Pos {
	return x.repo.Get(ctx, request)
}
func (x *PosServices) FindById(ctx context.Context, id string) (*keep_entities.Pos, error) {
	return x.repo.FindById(ctx, id)
}
func (x *PosServices) Insert(ctx context.Context, posRequest *keep_request.PosInputUpdate) (*keep_entities.Pos, error) {
	pos := &keep_entities.Pos{}
	v := validator.New()
	err := v.ValidateStruct(posRequest)
	if err != nil {
		return pos, err
	}

	pos = &keep_entities.Pos{
		Nama:     posRequest.Nama,
		Urutan:   posRequest.Urutan,
		ParentId: posRequest.ParentId,
		IsShow:   true,
		Level:    1,
		Status:   "aktif",
		IsLeaf:   true,
	}
	if posRequest.ParentId != "" {
		parent, err2 := x.repo.FindById(ctx, posRequest.ParentId)
		if err2 != nil {
			return pos, err2
		}
		pos.Level = parent.Level + 1
	}
	pos, err = x.repo.Insert(ctx, pos)

	_ = helpers_events.GetDispatcher().Dispatch(keep_events.PosCreated, keep_events.PosCreatedEventData{
		Time:     time.Now(),
		Id:       pos.Id,
		Nama:     pos.Nama,
		Urutan:   pos.Urutan,
		Saldo:    pos.Saldo,
		ParentId: pos.ParentId,
		Level:    pos.Level,
		IsShow:   pos.IsShow,
		IsLeaf:   pos.IsLeaf,
		Status:   pos.Status,
	})

	return pos, err
}
func (x *PosServices) Update(ctx context.Context, posRequest *keep_request.PosInputUpdate) (affected int, err error) {
	model, err := x.repo.FindById(ctx, posRequest.Id)
	if err != nil {
		return 0, err
	}

	pos := &keep_entities.Pos{}
	v := validator.New()
	err = v.ValidateStruct(posRequest)
	if err != nil {
		return 0, err
	}

	pos = &keep_entities.Pos{
		Id:       posRequest.Id,
		Nama:     posRequest.Nama,
		Urutan:   posRequest.Urutan,
		ParentId: posRequest.ParentId,
		IsShow:   posRequest.IsShow,
		Saldo:    model.Saldo,
		IsLeaf:   model.IsLeaf,
		Status:   model.Status,
		Level:    1,
	}
	if posRequest.ParentId != "" {
		if posRequest.ParentId == posRequest.Id {
			return 0, helpers_error.NewValidationErrors("parentId", "invalid", "parent_to_self")
		}
		parent, err2 := x.repo.FindById(ctx, posRequest.ParentId)
		if err2 != nil {
			return 0, err2
		}
		pos.Level = parent.Level + 1
	}
	affected, err = x.repo.Update(ctx, pos)

	_ = helpers_events.GetDispatcher().Dispatch(keep_events.PosUpdated, keep_events.PosUpdatedEventData{
		Time: time.Now(),
		Old: keep_events.PosEventData{
			Time:     time.Now(),
			Id:       model.Id,
			Nama:     model.Nama,
			Urutan:   model.Urutan,
			Saldo:    model.Saldo,
			ParentId: model.ParentId,
			Level:    model.Level,
			IsShow:   model.IsShow,
			IsLeaf:   model.IsLeaf,
			Status:   model.Status,
		},
		New: keep_events.PosEventData{
			Time:     time.Now(),
			Id:       pos.Id,
			Nama:     pos.Nama,
			Urutan:   pos.Urutan,
			Saldo:    pos.Saldo,
			ParentId: pos.ParentId,
			Level:    pos.Level,
			IsShow:   pos.IsShow,
			IsLeaf:   pos.IsLeaf,
			Status:   pos.Status,
		},
	})

	return affected, err
}
func (x *PosServices) DeleteById(ctx context.Context, id string) (affected int, err error) {
	model, err := x.repo.FindById(ctx, id)
	if err != nil {
		return 0, err
	}
	affected, err = x.repo.SoftDeleteById(ctx, id)

	_ = helpers_events.GetDispatcher().Dispatch(keep_events.PosSoftDeleted, keep_events.PosSoftDeletedEventData{
		Time:     time.Now(),
		Id:       model.Id,
		Nama:     model.Nama,
		Urutan:   model.Urutan,
		Saldo:    model.Saldo,
		ParentId: model.ParentId,
		Level:    model.Level,
		IsShow:   model.IsShow,
		IsLeaf:   model.IsLeaf,
		Status:   model.Status,
	})

	return affected, err
}
func (x *PosServices) GetTrashed(ctx context.Context) []*keep_entities.Pos {
	return x.repo.GetTrashed(ctx)
}
func (x *PosServices) RestoreTrashedById(ctx context.Context, id string) (affected int, err error) {
	model, err := x.repo.FindTrashedById(ctx, id)
	if err != nil {
		return 0, err
	}
	affected, err = x.repo.RestoreTrashedById(ctx, id)

	_ = helpers_events.GetDispatcher().Dispatch(keep_events.PosRestored, keep_events.PosRestoredEventData{
		Time:     time.Now(),
		Id:       model.Id,
		Nama:     model.Nama,
		Urutan:   model.Urutan,
		Saldo:    model.Saldo,
		ParentId: model.ParentId,
		Level:    model.Level,
		IsShow:   model.IsShow,
		IsLeaf:   model.IsLeaf,
		Status:   model.Status,
	})

	return affected, err
}
func (x *PosServices) DeleteTrashedById(ctx context.Context, id string) (affected int, err error) {
	model, err := x.repo.FindTrashedById(ctx, id)
	if err != nil {
		return 0, err
	}
	affected, err = x.repo.HardDeleteTrashedById(ctx, id)

	_ = helpers_events.GetDispatcher().Dispatch(keep_events.PosHardDeleted, keep_events.PosHardDeletedEventData{
		Time:     time.Now(),
		Id:       model.Id,
		Nama:     model.Nama,
		Urutan:   model.Urutan,
		Saldo:    model.Saldo,
		ParentId: model.ParentId,
		Level:    model.Level,
		IsShow:   model.IsShow,
		IsLeaf:   model.IsLeaf,
		Status:   model.Status,
	})

	return affected, err
}
func (x *PosServices) UpdateSaldoFromTransaksi(ctx context.Context, ids []string) (affected int, err error) {
	ids = helpers.Unique(ids)
	ids = helpers.Filter(ids, func(s string) bool { return s != "" })

	for _, id := range ids {
		_, err2 := x.repo.FindById(ctx, id)
		if err2 != nil {
			return 0, err2
		}
		jumlah := x.transaksiRepo.GetJumlahByPosId(ctx, id)
		aff := x.repo.UpdateSaldo(ctx, id, jumlah)
		affected += aff
		aff, err2 = x.updateParentSaldo(ctx, id)
		if err2 != nil {
			return 0, err2
		}
		affected += aff
	}
	return affected, nil
}
func (x *PosServices) UpdateLeafStatus(ctx context.Context, ids []string) (affected int, err error) {
	ids = helpers.Unique(ids)
	ids = helpers.Filter(ids, func(s string) bool { return s != "" })

	for _, id := range ids {
		_, err2 := x.repo.FindById(ctx, id)
		if err2 != nil {
			return 0, err2
		}
		count := x.repo.CountChildren(ctx, id)
		isLeaf := true
		if count > 0 {
			isLeaf = false
		}
		aff, err2 := x.repo.UpdateLeaf(ctx, id, isLeaf)
		if err2 != nil {
			return 0, err2
		}
		affected += aff
	}
	return affected, nil
}
func (x *PosServices) UpdateUrutan(ctx context.Context, posRequests []*keep_request.PosUpdateUrutanItem) (affected int, err error) {
	v := validator.New()
	for _, request := range posRequests {
		err = v.ValidateStruct(request)
		if err != nil {
			return 0, err
		}
		if request.ParentId != "" {
			_, err = x.repo.FindById(ctx, request.ParentId)
			if err != nil {
				return 0, err
			}
		}
	}
	for _, request := range posRequests {
		aff, err2 := x.repo.UpdateUrutan(ctx, request.Id, request.Urutan, request.ParentId)
		if err2 != nil {
			return affected, err2
		}
		affected += aff
	}
	return affected, nil
}
func (x *PosServices) UpdateVisibility(ctx context.Context, posRequests []*keep_request.PosUpdateVisibilityItem) (affected int, err error) {
	v := validator.New()
	for _, request := range posRequests {
		err = v.ValidateStruct(request)
		if err != nil {
			return 0, err
		}
	}
	for _, request := range posRequests {
		aff, err2 := x.repo.UpdateVisibility(ctx, request.Id, request.IsShow)
		if err2 != nil {
			return affected, err2
		}
		affected += aff
	}
	return affected, nil
}

func (x *PosServices) updateParentSaldo(ctx context.Context, id string) (affected int, err error) {
	m, err := x.repo.FindById(ctx, id)
	if err != nil {
		return 0, err
	}
	id = m.ParentId
	if id == "" {
		return 0, nil
	}

	for {
		parent, err2 := x.repo.FindById(ctx, id)
		if err2 != nil {
			return 0, err2
		}

		jumlah := x.repo.GetJumlahById(ctx, parent.Id)
		aff := x.repo.UpdateSaldo(ctx, parent.Id, jumlah)
		affected += aff

		id = parent.ParentId
		if id == "" {
			break
		}
	}

	return affected, nil
}
