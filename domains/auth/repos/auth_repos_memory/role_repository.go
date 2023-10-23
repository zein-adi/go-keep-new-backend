package auth_repos_memory

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_requests"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"slices"
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

func (r *RoleRepository) Get(_ context.Context, request auth_requests.Get) []*auth_entities.Role {
	data := r.getDataFiltered(request)
	if request.Take > 0 {
		data = helpers.Slice(data, request.Skip, request.Take)
	}
	return helpers.Map(data, func(d *auth_entities.Role) *auth_entities.Role {
		return d.Copy()
	})
}

func (r *RoleRepository) Count(_ context.Context, request auth_requests.Get) (count int) {
	return len(r.getDataFiltered(request))
}
func (r *RoleRepository) CountByNama(_ context.Context, nama string, exceptId string) (count int) {
	matches := helpers.Filter(r.data, func(role *auth_entities.Role) bool {
		return role.Nama == nama && role.Id != exceptId
	})
	return len(matches)
}
func (r *RoleRepository) Insert(_ context.Context, role *auth_entities.Role) (*auth_entities.Role, error) {
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
func (r *RoleRepository) FindById(_ context.Context, id string) (*auth_entities.Role, error) {
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
func (r *RoleRepository) GetById(_ context.Context, ids []string) ([]*auth_entities.Role, error) {
	matches := helpers.Filter(r.data, func(role *auth_entities.Role) bool {
		return slices.Contains(ids, role.Id)
	})
	copied := helpers.Map(matches, func(d *auth_entities.Role) *auth_entities.Role {
		return d.Copy()
	})
	expectedLen := len(ids)
	actualLen := len(copied)
	if actualLen < expectedLen {
		return copied, helpers_error.NewEntryCountMismatchError(expectedLen, actualLen)
	}
	return copied, nil
}

func (r *RoleRepository) getDataFiltered(request auth_requests.Get) []*auth_entities.Role {
	return helpers.Filter(r.data, func(role *auth_entities.Role) bool {
		res := true

		if request.Search != "" {
			res = res && strings.Contains(role.Nama, request.Search)
		}

		return res
	})
}
