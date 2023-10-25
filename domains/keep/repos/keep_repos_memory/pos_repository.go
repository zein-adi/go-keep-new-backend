package keep_repos_memory

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"strconv"
)

var posEntityName = "pos"

func NewPosMemoryRepository() *PosMemoryRepository {
	return &PosMemoryRepository{}
}

type PosMemoryRepository struct {
	Data []*keep_entities.Pos
}

func (x *PosMemoryRepository) Get(_ context.Context, request *keep_request.GetPos) []*keep_entities.Pos {
	models := x.newQueryRequest(request, "aktif")
	return helpers.Map(models, func(d *keep_entities.Pos) *keep_entities.Pos {
		return d.Copy()
	})
}
func (x *PosMemoryRepository) FindById(_ context.Context, id string) (*keep_entities.Pos, error) {
	index, err := x.findIndexById(id, "aktif")
	if err != nil {
		return nil, err
	}
	return x.Data[index].Copy(), nil
}
func (x *PosMemoryRepository) Insert(_ context.Context, pos *keep_entities.Pos) (*keep_entities.Pos, error) {
	lastId := helpers.Reduce(x.Data, 0, func(accumulator int, pos *keep_entities.Pos) int {
		datumId, _ := strconv.Atoi(pos.Id)
		return max(accumulator, datumId)
	})

	model := pos.Copy()
	model.Id = strconv.Itoa(lastId + 1)
	x.Data = append(x.Data, model)
	return model, nil
}
func (x *PosMemoryRepository) Update(_ context.Context, pos *keep_entities.Pos) (affected int, err error) {
	model := &keep_entities.Pos{}
	_, err = x.findById(pos.Id, "aktif")
	if err != nil {
		return 0, err
	}
	index, _ := helpers.FindIndex(x.Data, func(p *keep_entities.Pos) bool {
		return p.Id == pos.Id
	})
	model = pos.Copy()
	x.Data[index] = model
	return 1, nil
}
func (x *PosMemoryRepository) SoftDeleteById(_ context.Context, id string) (affected int, err error) {
	model, _ := x.findById(id, "aktif")
	model.Status = "trashed"
	return 1, nil
}
func (x *PosMemoryRepository) DeleteById(_ context.Context, id string) (affected int, err error) {
	index, _ := helpers.FindIndex(x.Data, func(p *keep_entities.Pos) bool {
		return p.Id == id
	})
	x.Data = append(x.Data[:index], x.Data[index+1:]...)
	return 1, nil
}
func (x *PosMemoryRepository) GetTrashed(_ context.Context) []*keep_entities.Pos {
	models := x.newQueryRequest(keep_request.NewGetPos(), "trashed")
	return helpers.Map(models, func(d *keep_entities.Pos) *keep_entities.Pos {
		return d.Copy()
	})
}
func (x *PosMemoryRepository) FindTrashedById(_ context.Context, id string) (*keep_entities.Pos, error) {
	index, err := x.findIndexById(id, "trashed")
	if err != nil {
		return nil, err
	}
	return x.Data[index].Copy(), nil
}
func (x *PosMemoryRepository) RestoreTrashedById(_ context.Context, id string) (affected int, err error) {
	model, _ := x.findById(id, "trashed")
	model.Status = "aktif"
	return 1, nil
}
func (x *PosMemoryRepository) UpdateSaldo(_ context.Context, id string, saldo int) (affected int) {
	model, _ := x.findById(id, "aktif")
	model.Saldo = saldo
	return 1
}

func (x *PosMemoryRepository) GetJumlahById(_ context.Context, id string) (saldo int) {
	models := helpers.Filter(x.Data, func(v *keep_entities.Pos) bool {
		return v.ParentId == id
	})
	return helpers.Reduce(models, 0, func(accumulator int, v *keep_entities.Pos) int {
		return accumulator + v.Saldo
	})
}
func (x *PosMemoryRepository) CountChildren(_ context.Context, id string) (count int) {
	models := helpers.Filter(x.Data, func(v *keep_entities.Pos) bool {
		return v.Status == "aktif" && v.ParentId == id
	})
	return len(models)
}
func (x *PosMemoryRepository) UpdateLeaf(_ context.Context, id string, leaf bool) (affected int, err error) {
	index, _ := x.findIndexById(id, "aktif")
	x.Data[index].IsLeaf = leaf
	return 1, nil
}

func (x *PosMemoryRepository) newQueryRequest(request *keep_request.GetPos, status string) []*keep_entities.Pos {
	return helpers.Filter(x.Data, func(pos *keep_entities.Pos) bool {
		res := true
		if request.IsLeafOnly == true {
			res = res && pos.IsLeaf == true
		}
		if status != "" {
			res = res && pos.Status == status
		}
		return res
	})
}
func (x *PosMemoryRepository) findIndexById(id string, status string) (index int, err error) {
	index, err = helpers.FindIndex(x.Data, func(v *keep_entities.Pos) bool {
		res := v.Id == id
		if status != "" {
			res = res && v.Status == status
		}
		return res
	})
	if err != nil {
		return index, helpers_error.NewEntryNotFoundError(posEntityName, "id", "id")
	}
	return index, nil
}
func (x *PosMemoryRepository) findById(id string, status string) (*keep_entities.Pos, error) {
	models := helpers.Filter(x.Data, func(pos *keep_entities.Pos) bool {
		res := pos.Id == id
		if status != "" {
			res = res && pos.Status == status
		}
		return res
	})
	model := &keep_entities.Pos{}
	if len(models) == 0 {
		return model, helpers_error.NewEntryNotFoundError(posEntityName, "id", "id")
	}
	return models[0], nil
}
