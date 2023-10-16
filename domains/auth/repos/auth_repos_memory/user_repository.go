package auth_repos_memory

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_requests"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"strconv"
	"strings"
)

var userEntityName = "user"

func NewUserMemoryRepository() *UserRepository {
	return &UserRepository{}
}

type UserRepository struct {
	data []*auth_entities.User
}

func (r *UserRepository) Get(_ context.Context, request auth_requests.GetRequest) []*auth_entities.User {
	data := r.getDataFiltered(request)
	data = helpers.Slice(data, request.Skip, request.Take)
	return helpers.Map(data, func(d *auth_entities.User) *auth_entities.User {
		return d.Copy()
	})
}

func (r *UserRepository) Count(_ context.Context, request auth_requests.GetRequest) (count int) {
	return len(r.getDataFiltered(request))
}

func (r *UserRepository) CountByUsername(_ context.Context, username string, exceptId string) (count int) {
	matches := helpers.Filter(r.data, func(user *auth_entities.User) bool {
		return user.Username == username && user.Id != exceptId
	})
	return len(matches)
}

func (r *UserRepository) Insert(_ context.Context, user *auth_entities.User) (*auth_entities.User, error) {
	lastId := helpers.Reduce(r.data, 0, func(accumulator int, user *auth_entities.User) int {
		datumId, _ := strconv.Atoi(user.Id)
		return max(accumulator, datumId)
	})

	model := user.Copy()
	model.Id = strconv.Itoa(lastId + 1)
	r.data = append(r.data, model)
	return model, nil
}

func (r *UserRepository) Update(ctx context.Context, user *auth_entities.User) (affected int, err error) {
	_, err = r.FindById(ctx, user.Id)
	if err != nil {
		return 0, err
	}
	index, _ := helpers.FindIndex(r.data, func(user *auth_entities.User) bool {
		return user.Id == user.Id
	})
	r.data[index] = user
	return 1, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userId, password string) (affected int, err error) {
	_, err = r.FindById(ctx, userId)
	if err != nil {
		return 0, err
	}

	index, err := helpers.FindIndex(r.data, func(user *auth_entities.User) bool {
		return user.Id == userId
	})
	model := r.data[index]
	model.Password = password
	return 1, nil
}

func (r *UserRepository) FindById(_ context.Context, id string) (*auth_entities.User, error) {
	index, err := helpers.FindIndex(r.data, func(user *auth_entities.User) bool {
		return user.Id == id
	})
	if err != nil {
		return &auth_entities.User{}, helpers_error.NewEntryNotFoundError(userEntityName, "id", id)
	}
	return r.data[index].Copy(), nil
}
func (r *UserRepository) FindByUsername(_ context.Context, username string) (*auth_entities.User, error) {
	index, err := helpers.FindIndex(r.data, func(user *auth_entities.User) bool {
		return user.Username == username
	})
	if err != nil {
		return &auth_entities.User{}, helpers_error.NewEntryNotFoundError(userEntityName, "username", username)
	}
	return r.data[index].Copy(), nil
}
func (r *UserRepository) DeleteById(ctx context.Context, id string) (affected int, err error) {
	_, err = r.FindById(ctx, id)
	if err != nil {
		return 0, err
	}
	index, _ := helpers.FindIndex(r.data, func(user *auth_entities.User) bool {
		return user.Id == id
	})
	r.data = append(r.data[:index], r.data[index+1:]...)
	return 1, nil
}

func (r *UserRepository) getDataFiltered(request auth_requests.GetRequest) []*auth_entities.User {
	return helpers.Filter(r.data, func(user *auth_entities.User) bool {
		res := true
		if request.Search != "" {
			res = res && strings.Contains(user.Nama, request.Search)
		}
		return res
	})
}
