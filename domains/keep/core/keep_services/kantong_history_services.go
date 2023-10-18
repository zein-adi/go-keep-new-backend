package keep_services

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
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

func (x *KantongHistoryServices) Get(ctx context.Context) []*keep_entities.KantongHistory {
	return x.repo.Get(ctx)
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
	return x.repo.Insert(ctx, kantongHistory)
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
	return x.repo.Update(ctx, kantongHistory)
}

func (x *KantongHistoryServices) DeleteById(ctx context.Context, id string) (affected int, err error) {
	_, err = x.repo.FindById(ctx, id)
	if err != nil {
		return 0, err
	}
	return x.repo.DeleteById(ctx, id)
}