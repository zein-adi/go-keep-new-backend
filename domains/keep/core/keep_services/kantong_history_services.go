package keep_services

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_events"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_events"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_requests"
	"github.com/zein-adi/go-keep-new-backend/helpers/validator"
	"time"
)

func NewKantongHistoryServices(repo keep_repo_interfaces.IKantongHistoryRepository, kantongRepo keep_repo_interfaces.IKantongRepository) *KantongHistoryServices {
	return &KantongHistoryServices{
		repo:        repo,
		kantongRepo: kantongRepo,
	}
}

type KantongHistoryServices struct {
	repo        keep_repo_interfaces.IKantongHistoryRepository
	kantongRepo keep_repo_interfaces.IKantongRepository
}

func (x *KantongHistoryServices) Get(ctx context.Context, kantongId string, request *helpers_requests.Get) []*keep_entities.KantongHistory {
	return x.repo.Get(ctx, kantongId, request)
}

func (x *KantongHistoryServices) Insert(ctx context.Context, kantongHistoryRequest *keep_request.KantongHistoryInsertUpdate) (*keep_entities.KantongHistory, error) {
	err := validator.New().ValidateStruct(kantongHistoryRequest)
	if err != nil {
		return nil, err
	}

	_, err = x.kantongRepo.FindById(ctx, kantongHistoryRequest.KantongId)
	if err != nil {
		return nil, err
	}

	kantongHistory := &keep_entities.KantongHistory{
		KantongId: kantongHistoryRequest.KantongId,
		Jumlah:    kantongHistoryRequest.Jumlah,
		Uraian:    kantongHistoryRequest.Uraian,
		Waktu:     time.Now().Unix(),
	}
	model, err := x.repo.Insert(ctx, kantongHistory)
	if err != nil {
		return nil, err
	}

	_ = helpers_events.GetDispatcher().Dispatch(
		keep_events.KantongHistoryCreated,
		keep_events.KantongHistoryCreatedEventData{
			Time: time.Now(),
			Data: *model,
		})

	return model, nil
}

func (x *KantongHistoryServices) Update(ctx context.Context, kantongHistoryRequest *keep_request.KantongHistoryInsertUpdate) (affected int, err error) {
	err = validator.New().ValidateStruct(kantongHistoryRequest)
	if err != nil {
		return 0, err
	}

	_, err = x.kantongRepo.FindById(ctx, kantongHistoryRequest.KantongId)
	if err != nil {
		return 0, err
	}
	model, err := x.repo.FindById(ctx, kantongHistoryRequest.Id)
	if err != nil {
		return 0, err
	}

	kantongHistory := &keep_entities.KantongHistory{
		Id:        kantongHistoryRequest.Id,
		KantongId: kantongHistoryRequest.KantongId,
		Jumlah:    kantongHistoryRequest.Jumlah,
		Uraian:    kantongHistoryRequest.Uraian,
		Waktu:     model.Waktu,
	}
	affected, err = x.repo.Update(ctx, kantongHistory)
	if err != nil {
		return 0, err
	}

	_ = helpers_events.GetDispatcher().Dispatch(keep_events.KantongHistoryUpdated, keep_events.KantongHistoryUpdatedEventData{
		Time: time.Now(),
		Old:  *model,
		New:  *kantongHistory,
	})

	return affected, nil
}

func (x *KantongHistoryServices) DeleteById(ctx context.Context, _ string, id string) (affected int, err error) {
	model, err := x.repo.FindById(ctx, id)
	if err != nil {
		return 0, err
	}
	affected, err = x.repo.DeleteById(ctx, id)
	if err != nil {
		return 0, err
	}

	_ = helpers_events.GetDispatcher().Dispatch(keep_events.KantongHistoryDeleted, keep_events.KantongHistoryDeletedEventData{
		Time: time.Now(),
		Data: *model,
	})

	return affected, nil
}
