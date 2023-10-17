package keep_services

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/validator"
)

func NewPosServices(repo keep_repo_interfaces.IPosRepository) *PosServices {
	return &PosServices{
		repo: repo,
	}
}

type PosServices struct {
	repo keep_repo_interfaces.IPosRepository
}

func (x *PosServices) Get(ctx context.Context, request *keep_request.PosGetRequest) []*keep_entities.Pos {
	return x.repo.Get(ctx, request)
}
func (x *PosServices) Insert(ctx context.Context, posRequest *keep_request.PosInputUpdateRequest) (*keep_entities.Pos, error) {
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
		IsShow:   posRequest.IsShow,
	}
	pos.IsShow = true
	pos.Status = "aktif"
	if posRequest.ParentId != "" {
		parent, err2 := x.repo.FindById(ctx, posRequest.ParentId)
		if err2 != nil {
			return pos, err2
		}
		pos.Level = parent.Level + 1
	}

	return x.repo.Insert(ctx, pos)
}
func (x *PosServices) Update(ctx context.Context, posRequest *keep_request.PosInputUpdateRequest) (*keep_entities.Pos, error) {
	model, err := x.repo.FindById(ctx, posRequest.Id)
	if err != nil {
		return model, err
	}

	pos := &keep_entities.Pos{}
	v := validator.New()
	err = v.ValidateStruct(posRequest)
	if err != nil {
		return pos, err
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
	}
	if posRequest.ParentId != "" {
		if posRequest.ParentId == posRequest.Id {
			return pos, helpers_error.NewValidationErrors("parentId", "invalid", "parent_to_self")
		}
		parent, err := x.repo.FindById(ctx, posRequest.ParentId)
		if err != nil {
			return pos, err
		}
		pos.Level = parent.Level + 1
	}
	return x.repo.Update(ctx, pos)
}
func (x *PosServices) DeleteById(ctx context.Context, id string) (affected int, err error) {
	_, err = x.repo.FindById(ctx, id)
	if err != nil {
		return 0, err
	}
	return x.repo.SoftDeleteById(ctx, id)
}
func (x *PosServices) GetTrashed(ctx context.Context) []*keep_entities.Pos {
	return x.repo.GetTrashed(ctx)
}
func (x *PosServices) RestoreTrashedById(ctx context.Context, id string) (affected int, err error) {
	_, err = x.repo.FindTrashedById(ctx, id)
	if err != nil {
		return 0, err
	}
	return x.repo.RestoreTrashedById(ctx, id)
}
func (x *PosServices) DeleteTrashedById(ctx context.Context, id string) (affected int, err error) {
	_, err = x.repo.FindTrashedById(ctx, id)
	if err != nil {
		return 0, err
	}
	return x.repo.DeleteById(ctx, id)
}
