package keep_services

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_events"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_events"
	"github.com/zein-adi/go-keep-new-backend/helpers/validator"
	"time"
)

func NewKantongServices(kantongRepo keep_repo_interfaces.IKantongRepository, posRepo keep_repo_interfaces.IPosRepository) *KantongServices {
	return &KantongServices{
		repo:    kantongRepo,
		posRepo: posRepo,
	}
}

type KantongServices struct {
	repo    keep_repo_interfaces.IKantongRepository
	posRepo keep_repo_interfaces.IPosRepository
}

func (x *KantongServices) Get(ctx context.Context) []*keep_entities.Kantong {
	return x.repo.Get(ctx)
}
func (x *KantongServices) FindById(ctx context.Context, id string) (*keep_entities.Kantong, error) {
	return x.repo.FindById(ctx, id)
}
func (x *KantongServices) Insert(ctx context.Context, kantongRequest *keep_request.KantongInsert) (*keep_entities.Kantong, error) {
	err := validator.New().ValidateStruct(kantongRequest)
	if err != nil {
		return nil, err
	}
	_, err = x.posRepo.FindById(ctx, kantongRequest.PosId)
	if err != nil {
		return nil, err
	}

	kantong := &keep_entities.Kantong{
		Nama:           kantongRequest.Nama,
		Urutan:         kantongRequest.Urutan,
		Saldo:          kantongRequest.Saldo,
		SaldoMengendap: kantongRequest.SaldoMengendap,
		PosId:          kantongRequest.PosId,
		IsShow:         true,
		Status:         "aktif",
	}
	model, err := x.repo.Insert(ctx, kantong)
	if err != nil {
		return nil, err
	}

	_ = helpers_events.GetDispatcher().Dispatch(keep_events.KantongCreated, &keep_events.KantongCreatedEventData{
		Time: time.Now(),
		Data: *model,
	})

	return model, nil
}
func (x *KantongServices) Update(ctx context.Context, kantongRequest *keep_request.KantongUpdate) (affected int, int error) {
	err := validator.New().ValidateStruct(kantongRequest)
	if err != nil {
		return 0, err
	}
	_, err = x.posRepo.FindById(ctx, kantongRequest.PosId)
	if err != nil {
		return 0, err
	}

	model, err := x.repo.FindById(ctx, kantongRequest.Id)
	if err != nil {
		return 0, err
	}

	kantong := &keep_entities.Kantong{
		Id:             kantongRequest.Id,
		Nama:           kantongRequest.Nama,
		Urutan:         kantongRequest.Urutan,
		Saldo:          kantongRequest.Saldo,
		SaldoMengendap: kantongRequest.SaldoMengendap,
		PosId:          kantongRequest.PosId,
		IsShow:         kantongRequest.IsShow,
		Status:         "aktif",
	}
	affected, err = x.repo.Update(ctx, kantong)
	if err != nil {
		return 0, err
	}

	_ = helpers_events.GetDispatcher().Dispatch(keep_events.KantongUpdated, &keep_events.KantongUpdatedEventData{
		Time: time.Now(),
		Old:  *model,
		New:  *kantong,
	})

	return affected, nil
}
func (x *KantongServices) DeleteById(ctx context.Context, id string) (affected int, err error) {
	model, err := x.repo.FindById(ctx, id)
	if err != nil {
		return 0, err
	}
	affected, err = x.repo.SoftDeleteById(ctx, id)
	if err != nil {
		return 0, err
	}

	_ = helpers_events.GetDispatcher().Dispatch(keep_events.KantongSoftDeleted, &keep_events.KantongSoftDeletedEventData{
		Time: time.Now(),
		Data: *model,
	})

	return affected, nil
}
func (x *KantongServices) GetTrashed(ctx context.Context) []*keep_entities.Kantong {
	return x.repo.GetTrashed(ctx)
}
func (x *KantongServices) RestoreTrashedById(ctx context.Context, id string) (affected int, err error) {
	model, err := x.repo.FindTrashedById(ctx, id)
	if err != nil {
		return 0, err
	}
	affected, err = x.repo.RestoreTrashedById(ctx, id)
	if err != nil {
		return 0, err
	}

	_ = helpers_events.GetDispatcher().Dispatch(keep_events.KantongRestored, &keep_events.KantongRestoredEventData{
		Time: time.Now(),
		Data: *model,
	})

	return affected, nil
}
func (x *KantongServices) DeleteTrashedById(ctx context.Context, id string) (affected int, err error) {
	model, err := x.repo.FindTrashedById(ctx, id)
	if err != nil {
		return 0, err
	}
	affected, err = x.repo.HardDeleteTrashedById(ctx, id)
	if err != nil {
		return 0, err
	}

	_ = helpers_events.GetDispatcher().Dispatch(keep_events.KantongHardDeleted, &keep_events.KantongHardDeletedEventData{
		Time: time.Now(),
		Data: *model,
	})

	return affected, nil
}
func (x *KantongServices) UpdateSaldo(ctx context.Context, asalId, tujuanId string, jumlah int, oldAsalId, oldTujuanId string, oldJumlah int) (affected int, err error) {
	update := map[string]int{
		asalId:      0,
		tujuanId:    0,
		oldAsalId:   0,
		oldTujuanId: 0,
	}

	if oldAsalId != "" {
		// Dulu berkurang sekarang bertambah (revert)
		update[oldAsalId] -= -1 * oldJumlah
	}
	if oldTujuanId != "" {
		// Dulu bertambah sekarang berkurang (revert)
		update[oldTujuanId] += -1 * oldJumlah
	}
	if asalId != "" {
		// Berkurang
		update[asalId] -= jumlah
	}
	if tujuanId != "" {
		// Bertambah
		update[tujuanId] += jumlah
	}

	for kantongId, saldo := range update {
		if kantongId == "" {
			continue
		}
		d, err2 := x.repo.FindById(ctx, kantongId)
		if err2 != nil {
			return 0, err2
		}
		af, err2 := x.repo.UpdateSaldo(ctx, kantongId, d.Saldo+saldo)
		if err2 != nil {
			return 0, err2
		}
		affected += af
	}
	return affected, nil
}
func (x *KantongServices) UpdateUrutan(ctx context.Context, kantongRequests []*keep_request.KantongUpdateUrutanItem) (affected int, err error) {
	v := validator.New()
	for _, request := range kantongRequests {
		err = v.ValidateStruct(request)
		if err != nil {
			return 0, err
		}
		_, err = x.posRepo.FindById(ctx, request.PosId)
		if err != nil {
			return 0, err
		}
	}
	for _, request := range kantongRequests {
		aff, err2 := x.repo.UpdateUrutan(ctx, request.Id, request.Urutan, request.PosId)
		if err2 != nil {
			return affected, err2
		}
		affected += aff
	}
	return affected, nil
}
func (x *KantongServices) UpdateVisibility(ctx context.Context, kantongRequests []*keep_request.KantongUpdateVisibilityItem) (affected int, err error) {
	v := validator.New()
	for _, request := range kantongRequests {
		err = v.ValidateStruct(request)
		if err != nil {
			return 0, err
		}
	}
	for _, request := range kantongRequests {
		aff, err2 := x.repo.UpdateVisibility(ctx, request.Id, request.IsShow)
		if err2 != nil {
			return affected, err2
		}
		affected += aff
	}
	return affected, nil
}
