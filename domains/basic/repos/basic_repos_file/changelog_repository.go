package basic_repos_file

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/core/basic_entities"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_directory"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_requests"
	"os"
	"sort"
	"strconv"
	"strings"
)

var changelogEntityName = "changelog"
var changelogFileName = "changelog.json"

//go:embed changelog.json
var changelogJson []byte

func NewChangelogFileRepository() *ChangelogFileRepository {
	t := &ChangelogFileRepository{
		Data: make([]*basic_entities.Changelog, 0),
	}
	return t
}

type ChangelogFileRepository struct {
	Data []*basic_entities.Changelog
}

func (x *ChangelogFileRepository) Get(_ context.Context, request *helpers_requests.Get) []*basic_entities.Changelog {
	x.loadCache()
	data := x.newQueryRequest(request)
	if request.Take > 0 {
		data = helpers.Slice(x.Data, request.Skip, request.Take)
	}
	return helpers.Map(data, func(v *basic_entities.Changelog) *basic_entities.Changelog {
		return v.Copy()
	})
}

func (x *ChangelogFileRepository) Count(_ context.Context, request *helpers_requests.Get) (count int) {
	x.loadCache()
	return len(x.newQueryRequest(request))
}

func (x *ChangelogFileRepository) FindById(_ context.Context, id string) (*basic_entities.Changelog, error) {
	index, err := x.findIndexById(id)
	if err != nil {
		return nil, err
	}
	return x.Data[index].Copy(), nil
}
func (x *ChangelogFileRepository) Insert(_ context.Context, changelog *basic_entities.Changelog) (*basic_entities.Changelog, error) {
	lastId := helpers.Reduce(x.Data, 0, func(accumulator int, v *basic_entities.Changelog) int {
		datumId, _ := strconv.Atoi(v.Id)
		return max(accumulator, datumId)
	})

	model := changelog.Copy()
	model.Id = strconv.Itoa(lastId + 1)
	x.Data = append(x.Data, model)

	x.writeToFile()

	return model, nil
}
func (x *ChangelogFileRepository) Update(_ context.Context, changelog *basic_entities.Changelog) (affected int, err error) {
	index, err := x.findIndexById(changelog.Id)
	if err != nil {
		return 0, err
	}

	model := changelog.Copy()
	x.Data[index] = model

	x.writeToFile()

	return 1, nil
}
func (x *ChangelogFileRepository) DeleteById(_ context.Context, id string) (affected int, err error) {
	index, err := x.findIndexById(id)
	if err != nil {
		return 0, err
	}
	x.Data = append(x.Data[0:index], x.Data[index+1:]...)

	x.writeToFile()

	return 1, nil
}

func (x *ChangelogFileRepository) newQueryRequest(request *helpers_requests.Get) []*basic_entities.Changelog {
	return helpers.Filter(x.Data, func(v *basic_entities.Changelog) bool {
		res := true
		if request.Search != "" {
			res = res && strings.Contains(strings.ToLower(v.Version), strings.ToLower(request.Search))
		}
		return res
	})
}
func (x *ChangelogFileRepository) findIndexById(id string) (index int, err error) {
	x.loadCache()

	index, err = helpers.FindIndex(x.Data, func(v *basic_entities.Changelog) bool {
		return v.Id == id
	})
	if err != nil {
		return index, helpers_error.NewEntryNotFoundError(changelogEntityName, "id", id)
	}
	return index, nil
}
func (x *ChangelogFileRepository) loadCache() {
	if len(x.Data) > 0 {
		return
	}

	if os.Getenv("APP_ENV") == "development" || os.Getenv("APP_ENV") == "testing" {
		// Read From File
		if !helpers_directory.FileExists(changelogFileName) {
			x.writeToFile()
			return
		}

		data, err := os.ReadFile(changelogFileName)
		helpers_error.PanicIfError(err)
		err = json.Unmarshal(data, &x.Data)
		helpers_error.PanicIfError(err)
		return
	}

	// Read From embed
	err := json.Unmarshal(changelogJson, &x.Data)
	helpers_error.PanicIfError(err)
}
func (x *ChangelogFileRepository) writeToFile() {
	// Order
	sort.Slice(x.Data, func(i, j int) bool {
		return x.Data[i].Timestamp > x.Data[j].Timestamp
	})

	// Save
	data, err := json.Marshal(x.Data)
	helpers_error.PanicIfError(err)
	err = os.WriteFile(changelogFileName, data, 0644)
	helpers_error.PanicIfError(err)

	// Reset Data
	x.Data = make([]*basic_entities.Changelog, 0)
	x.loadCache()
}
func (x *ChangelogFileRepository) DeleteFile() {
	_ = os.Remove(changelogFileName)
}
