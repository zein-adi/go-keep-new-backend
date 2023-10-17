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

func NewPosMemoryRepository() *PosRepository {
	return &PosRepository{}
}

type PosRepository struct {
	data []*keep_entities.Pos
}

func (x *PosRepository) Get(_ context.Context, request *keep_request.PosGetRequest) []*keep_entities.Pos {
	models := x.newQueryRequest(request, "aktif")
	return helpers.Map(models, func(d *keep_entities.Pos) *keep_entities.Pos {
		return d.Copy()
	})
}
func (x *PosRepository) FindById(_ context.Context, id string) (*keep_entities.Pos, error) {
	model, err := x.findById(id, "aktif")
	return model.Copy(), err
}
func (x *PosRepository) Insert(_ context.Context, pos *keep_entities.Pos) (*keep_entities.Pos, error) {
	lastId := helpers.Reduce(x.data, 0, func(accumulator int, pos *keep_entities.Pos) int {
		datumId, _ := strconv.Atoi(pos.Id)
		return max(accumulator, datumId)
	})

	model := pos.Copy()
	model.Id = strconv.Itoa(lastId + 1)
	x.data = append(x.data, model)
	return model, nil
}
func (x *PosRepository) Update(_ context.Context, pos *keep_entities.Pos) (*keep_entities.Pos, error) {
	model := &keep_entities.Pos{}
	_, err := x.findById(pos.Id, "aktif")
	if err != nil {
		return model, err
	}
	index, _ := helpers.FindIndex(x.data, func(p *keep_entities.Pos) bool {
		return p.Id == pos.Id
	})
	model = pos.Copy()
	x.data[index] = model
	return model, nil
}
func (x *PosRepository) SoftDeleteById(_ context.Context, id string) (affected int, err error) {
	model, _ := x.findById(id, "aktif")
	model.Status = "trashed"
	return 1, nil
}
func (x *PosRepository) DeleteById(_ context.Context, id string) (affected int, err error) {
	index, _ := helpers.FindIndex(x.data, func(p *keep_entities.Pos) bool {
		return p.Id == id
	})
	x.data = append(x.data[:index], x.data[index+1:]...)
	return 1, nil
}
func (x *PosRepository) GetTrashed(_ context.Context) []*keep_entities.Pos {
	models := x.newQueryRequest(keep_request.NewPosGetRequest(), "trashed")
	return helpers.Map(models, func(d *keep_entities.Pos) *keep_entities.Pos {
		return d.Copy()
	})
}
func (x *PosRepository) FindTrashedById(_ context.Context, id string) (*keep_entities.Pos, error) {
	model, err := x.findById(id, "trashed")
	return model.Copy(), err
}
func (x *PosRepository) RestoreTrashedById(_ context.Context, id string) (affected int, err error) {
	model, _ := x.findById(id, "trashed")
	model.Status = "aktif"
	return 1, nil
}

func (x *PosRepository) newQueryRequest(request *keep_request.PosGetRequest, status string) []*keep_entities.Pos {
	return helpers.Filter(x.data, func(pos *keep_entities.Pos) bool {
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
func (x *PosRepository) findById(id string, status string) (*keep_entities.Pos, error) {
	models := helpers.Filter(x.data, func(pos *keep_entities.Pos) bool {
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
