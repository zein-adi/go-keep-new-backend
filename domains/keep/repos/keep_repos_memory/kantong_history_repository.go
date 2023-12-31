package keep_repos_memory

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_requests"
	"sort"
	"strconv"
	"strings"
)

var kantongHistoryEntityName = "pos"

func NewKantongHistoryMemoryRepository() *KantongHistoryMemoryRepository {
	return &KantongHistoryMemoryRepository{}
}

type KantongHistoryMemoryRepository struct {
	Data []*keep_entities.KantongHistory
}

func (x *KantongHistoryMemoryRepository) Get(_ context.Context, kantongId string, request *helpers_requests.Get) []*keep_entities.KantongHistory {
	models := x.newQueryRequest(kantongId, request)
	sort.Slice(models, func(i, j int) bool {
		return models[i].Waktu > models[j].Waktu
	})
	if request.Take > 0 {
		models = helpers.Slice(models, request.Skip, request.Take)
	}
	return helpers.Map(models, func(d *keep_entities.KantongHistory) *keep_entities.KantongHistory {
		return d.Copy()
	})
}

func (x *KantongHistoryMemoryRepository) FindById(_ context.Context, id string) (*keep_entities.KantongHistory, error) {
	index, err := x.findIndexById(id)
	if err != nil {
		return nil, err
	}
	return x.Data[index].Copy(), err
}

func (x *KantongHistoryMemoryRepository) Insert(_ context.Context, kantongHistory *keep_entities.KantongHistory) (*keep_entities.KantongHistory, error) {
	lastId := helpers.Reduce(x.Data, 0, func(accumulator int, v *keep_entities.KantongHistory) int {
		datumId, _ := strconv.Atoi(v.Id)
		return max(accumulator, datumId)
	})

	model := kantongHistory.Copy()
	model.Id = strconv.Itoa(lastId + 1)
	x.Data = append(x.Data, model)
	return model, nil
}

func (x *KantongHistoryMemoryRepository) Update(_ context.Context, kantongHistory *keep_entities.KantongHistory) (affected int, err error) {
	index, err := x.findIndexById(kantongHistory.Id)
	if err != nil {
		return 0, err
	}

	model := kantongHistory.Copy()
	x.Data[index] = model
	return 1, nil
}

func (x *KantongHistoryMemoryRepository) DeleteById(_ context.Context, id string) (affected int, err error) {
	index, err := x.findIndexById(id)
	if err != nil {
		return 0, err
	}

	x.Data = append(x.Data[0:index], x.Data[index+1:]...)
	return 1, nil
}

func (x *KantongHistoryMemoryRepository) newQueryRequest(kantongId string, request *helpers_requests.Get) []*keep_entities.KantongHistory {
	return helpers.Filter(x.Data, func(v *keep_entities.KantongHistory) bool {
		res := v.KantongId == kantongId
		if request.Search != "" {
			res = res && strings.Contains(strings.ToLower(v.Uraian), strings.ToLower(request.Search))
		}
		return res
	})
}
func (x *KantongHistoryMemoryRepository) findIndexById(id string) (index int, err error) {
	index, err = helpers.FindIndex(x.Data, func(v *keep_entities.KantongHistory) bool {
		return v.Id == id
	})
	if err != nil {
		return -1, helpers_error.NewEntryNotFoundError(kantongHistoryEntityName, "id", "id")
	}
	return index, err
}
