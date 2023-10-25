package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/core/basic_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/core/basic_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/core/basic_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/core/basic_services"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/repos/basic_repos_file"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_env"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_requests"
	"testing"
	"time"
)

func TestChangelog(t *testing.T) {
	helpers_env.Init(5)
	x := NewChangelogServicesTest()
	defer x.cleanup()

	t.Run("InsertSuccess", func(t *testing.T) {
		x.truncate()
		ctx := context.Background()

		version := "1.0.0"
		timestamp := time.Now().Unix()
		body := "- Penambahakan fitur xxxx"
		changelog := &basic_entities.Changelog{
			Version:   version,
			Timestamp: timestamp,
			Body:      body,
		}
		model, err := x.services.Insert(ctx, changelog)
		assert.Nil(t, err)

		assert.Equal(t, version, model.Version)
		assert.Equal(t, timestamp, model.Timestamp)
		assert.Equal(t, body, model.Body)
	})
	t.Run("UpdateSuccess", func(t *testing.T) {
		x.truncate()
		ctx := context.Background()

		changelog := &basic_entities.Changelog{
			Version:   "1.0.0",
			Timestamp: time.Now().Add(time.Hour * -100).Unix(),
			Body:      "- Penambahakan fitur xxxx",
		}
		model, err := x.repo.Insert(ctx, changelog)
		assert.Nil(t, err)

		version := "1.0.0"
		timestamp := time.Now().Unix()
		body := "- Penambahan fitur yyyy "
		affected, err := x.services.Update(ctx, &basic_entities.Changelog{
			Id:        model.Id,
			Version:   version,
			Timestamp: timestamp,
			Body:      body,
		})
		assert.Nil(t, err)

		model, _ = x.repo.FindById(ctx, model.Id)
		assert.Equal(t, 1, affected)
		assert.Equal(t, version, model.Version)
		assert.Equal(t, timestamp, model.Timestamp)
		assert.Equal(t, body, model.Body)
	})
	t.Run("GetSuccess", func(t *testing.T) {
		x.truncate()
		ctx := context.Background()

		input := []*basic_entities.Changelog{
			{
				Version:   "1.0.0",
				Timestamp: time.Now().Add(time.Hour * -100).Unix(),
				Body:      "- Penambahakan fitur xxxx",
			},
			{
				Version:   "1.0.1",
				Timestamp: time.Now().Add(time.Hour * -50).Unix(),
				Body:      "- Penambahakan fitur yyyy",
			},
			{
				Version:   "1.0.2",
				Timestamp: time.Now().Add(time.Hour * -40).Unix(),
				Body:      "- Penambahakan fitur zzzz",
			},
		}

		ori := make([]*basic_entities.Changelog, 0)
		oriMap := make(map[string]*basic_entities.Changelog)
		for _, v := range input {
			m, err := x.repo.Insert(ctx, v)
			helpers_error.PanicIfError(err)
			ori = append(ori, m)
			oriMap[m.Id] = m
		}

		request := helpers_requests.NewGet()
		models := x.services.Get(ctx, request)
		for _, m := range models {
			o := oriMap[m.Id]
			assert.Equal(t, o.Id, m.Id)
			assert.Equal(t, o.Version, m.Version)
			assert.Equal(t, o.Timestamp, m.Timestamp)
			assert.Equal(t, o.Body, m.Body)
		}

		request = helpers_requests.NewGet()
		request.Skip = 0
		request.Take = 1
		models = x.services.Get(ctx, request)
		for _, m := range models {
			o := oriMap[m.Id]
			assert.Equal(t, "3", m.Id)
			assert.Equal(t, o.Version, m.Version)
			assert.Equal(t, o.Timestamp, m.Timestamp)
			assert.Equal(t, o.Body, m.Body)
		}

		request = helpers_requests.NewGet()
		request.Skip = 1
		request.Take = 1
		models = x.services.Get(ctx, request)
		for _, m := range models {
			o := oriMap[m.Id]
			assert.Equal(t, "2", m.Id)
			assert.Equal(t, o.Version, m.Version)
			assert.Equal(t, o.Timestamp, m.Timestamp)
			assert.Equal(t, o.Body, m.Body)
		}
	})
	t.Run("DeleteSuccess", func(t *testing.T) {
		x.truncate()
		ctx := context.Background()

		input := &basic_entities.Changelog{
			Version:   "1.0.0",
			Timestamp: time.Now().Add(time.Hour * -100).Unix(),
			Body:      "- Penambahakan fitur xxxx",
		}
		m, err := x.repo.Insert(ctx, input)
		helpers_error.PanicIfError(err)

		affected, err := x.services.DeleteById(ctx, m.Id)
		assert.Nil(t, err)
		assert.Equal(t, 1, affected)

		_, err = x.repo.FindById(ctx, m.Id)
		assert.ErrorIs(t, err, helpers_error.EntryNotFoundError)
	})
}

func NewChangelogServicesTest() *ChangelogServicesTest {
	x := &ChangelogServicesTest{}
	x.setUp()
	return x
}

type ChangelogServicesTest struct {
	repo     basic_repo_interfaces.IChangelogRepository
	services basic_service_interfaces.IChangelogServices
	truncate func()
	cleanup  func()
}

func (x *ChangelogServicesTest) setUp() {
	x.setUpRepository()

	x.services = basic_services.NewChangelogServices(x.repo)
}
func (x *ChangelogServicesTest) setUpRepository() {
	repo := basic_repos_file.NewChangelogFileRepository()
	x.repo = repo
	x.cleanup = func() {
		repo.DeleteFile()
	}
	x.truncate = func() {
		models := repo.Get(context.Background(), helpers_requests.NewGet())
		for _, model := range models {
			_, err := repo.DeleteById(context.Background(), model.Id)
			helpers_error.PanicIfError(err)
		}
	}
}
func (x *ChangelogServicesTest) reset() []*basic_entities.Changelog {
	x.truncate()
	//ctx := context.Background()

	barangs := []*basic_entities.Changelog{}
	return barangs
}
