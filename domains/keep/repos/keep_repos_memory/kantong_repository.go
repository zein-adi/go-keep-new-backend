package keep_repos_memory

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"sort"
	"strconv"
)

var kantongEntityName = "pos"

func NewKantongMemoryRepository() *KantongMemoryRepository {
	return &KantongMemoryRepository{}
}

type KantongMemoryRepository struct {
	Data []*keep_entities.Kantong
}

func (x *KantongMemoryRepository) Get(_ context.Context) []*keep_entities.Kantong {
	models := x.newQueryRequest("aktif")
	sort.Slice(models, func(i, j int) bool {
		return models[i].Urutan < models[i].Urutan
	})
	return helpers.Map(models, func(d *keep_entities.Kantong) *keep_entities.Kantong {
		return d.Copy()
	})
}
func (x *KantongMemoryRepository) FindById(_ context.Context, id string) (*keep_entities.Kantong, error) {
	index, err := x.findIndexById(id, "aktif")
	if err != nil {
		return nil, err
	}
	return x.Data[index].Copy(), err
}
func (x *KantongMemoryRepository) Insert(_ context.Context, kantong *keep_entities.Kantong) (*keep_entities.Kantong, error) {
	lastId := helpers.Reduce(x.Data, 0, func(accumulator int, pos *keep_entities.Kantong) int {
		datumId, _ := strconv.Atoi(pos.Id)
		return max(accumulator, datumId)
	})

	model := kantong.Copy()
	model.Id = strconv.Itoa(lastId + 1)
	x.Data = append(x.Data, model)
	return model, nil
}
func (x *KantongMemoryRepository) Update(_ context.Context, kantong *keep_entities.Kantong) (affected int, err error) {
	index, err := x.findIndexById(kantong.Id, "aktif")
	if err != nil {
		return 0, err
	}

	model := kantong.Copy()
	x.Data[index] = model
	return 1, nil
}
func (x *KantongMemoryRepository) UpdateSaldo(_ context.Context, id string, saldo int) (affected int, err error) {
	index, err := x.findIndexById(id, "aktif")
	if err != nil {
		return 0, err
	}

	x.Data[index].Saldo = saldo
	return 1, nil
}
func (x *KantongMemoryRepository) SoftDeleteById(_ context.Context, id string) (affected int, err error) {
	index, err := x.findIndexById(id, "aktif")
	if err != nil {
		return 0, err
	}

	x.Data[index].Status = "trashed"
	return 1, nil
}
func (x *KantongMemoryRepository) GetTrashed(_ context.Context) []*keep_entities.Kantong {
	models := x.newQueryRequest("trashed")
	return helpers.Map(models, func(d *keep_entities.Kantong) *keep_entities.Kantong {
		return d.Copy()
	})
}
func (x *KantongMemoryRepository) FindTrashedById(_ context.Context, id string) (*keep_entities.Kantong, error) {
	index, err := x.findIndexById(id, "trashed")
	if err != nil {
		return nil, err
	}
	return x.Data[index].Copy(), err
}
func (x *KantongMemoryRepository) RestoreTrashedById(_ context.Context, id string) (affected int, err error) {
	index, err := x.findIndexById(id, "trashed")
	if err != nil {
		return 0, err
	}

	x.Data[index].Status = "aktif"
	return 1, nil
}
func (x *KantongMemoryRepository) HardDeleteTrashedById(_ context.Context, id string) (affected int, err error) {
	index, err := x.findIndexById(id, "trashed")
	if err != nil {
		return 0, err
	}

	x.Data = append(x.Data[0:index], x.Data[index+1:]...)
	return 1, nil
}
func (x *KantongMemoryRepository) UpdateUrutan(_ context.Context, id string, urutan int, posId string) (affected int, err error) {
	index, err := x.findIndexById(id, "aktif")
	if err != nil {
		return 0, err
	}

	x.Data[index].Urutan = urutan
	x.Data[index].PosId = posId
	return 1, nil
}
func (x *KantongMemoryRepository) UpdateVisibility(_ context.Context, id string, isShow bool) (affected int, err error) {
	index, err := x.findIndexById(id, "aktif")
	if err != nil {
		return 0, err
	}

	x.Data[index].IsShow = isShow
	return 1, nil
}

func (x *KantongMemoryRepository) newQueryRequest(status string) []*keep_entities.Kantong {
	return helpers.Filter(x.Data, func(pos *keep_entities.Kantong) bool {
		res := true
		if status != "" {
			res = res && pos.Status == status
		}
		return res
	})
}
func (x *KantongMemoryRepository) findIndexById(id string, status string) (index int, err error) {
	index, err = helpers.FindIndex(x.Data, func(pos *keep_entities.Kantong) bool {
		res := pos.Id == id
		if status != "" {
			res = res && pos.Status == status
		}
		return res
	})
	if err != nil {
		return -1, helpers_error.NewEntryNotFoundError(kantongEntityName, "id", "id")
	}
	return index, err
}
