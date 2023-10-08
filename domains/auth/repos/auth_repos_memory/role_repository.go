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

var roleEntityName = "role"

func NewRoleMemoryRepository() *RoleRepository {
	return &RoleRepository{}
}

type RoleRepository struct {
	data []*auth_entities.Role
}

func (r *RoleRepository) Get(ctx context.Context, request auth_requests.GetRequest) []*auth_entities.Role {
	data := r.getDataFiltered(request)
	data = helpers.Slice(data, request.Skip, request.Take)
	return helpers.Map(data, func(d *auth_entities.Role) *auth_entities.Role {
		return d.Copy()
	})
}

func (r *RoleRepository) Count(ctx context.Context, request auth_requests.GetRequest) (count int) {
	return len(r.getDataFiltered(request))
}

func (r *RoleRepository) Insert(ctx context.Context, role *auth_entities.Role) (*auth_entities.Role, error) {
	lastId := helpers.Reduce(r.data, 0, func(accumulator int, role *auth_entities.Role) int {
		datumId, _ := strconv.Atoi(role.Id)
		return max(accumulator, datumId)
	})

	model := role.Copy()
	model.Id = strconv.Itoa(lastId + 1)
	r.data = append(r.data, model)
	return model, nil
}

func (r *RoleRepository) Update(ctx context.Context, role *auth_entities.Role) (affected int, err error) {
	_, err = r.FindById(ctx, role.Id)
	if err != nil {
		return 0, err
	}
	index, _ := helpers.FindIndex(r.data, func(role *auth_entities.Role) bool {
		return role.Id == role.Id
	})
	r.data[index] = role
	return 1, nil
}

func (r *RoleRepository) FindById(ctx context.Context, id string) (*auth_entities.Role, error) {
	index, err := helpers.FindIndex(r.data, func(role *auth_entities.Role) bool {
		return role.Id == id
	})
	if err != nil {
		return &auth_entities.Role{}, helpers_error.NewEntryNotFoundError(roleEntityName, "id", id)
	}
	return r.data[index].Copy(), nil
}

func (r *RoleRepository) DeleteById(ctx context.Context, id string) (affected int, err error) {
	_, err = r.FindById(ctx, id)
	if err != nil {
		return 0, err
	}
	index, _ := helpers.FindIndex(r.data, func(role *auth_entities.Role) bool {
		return role.Id == id
	})
	r.data = append(r.data[:index], r.data[index+1:]...)
	return 1, nil
}

func (r *RoleRepository) getDataFiltered(request auth_requests.GetRequest) []*auth_entities.Role {
	return helpers.Filter(r.data, func(role *auth_entities.Role) bool {
		res := true
		if request.Search != "" {
			res = res && strings.Contains(role.Nama, request.Search)
		}
		return res
	})
}
