package auth_services

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_requests"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/validator"
)

func NewRoleServices(roleRepo auth_repo_interfaces.IRoleRepository) *RoleServices {
	return &RoleServices{
		roleRepo: roleRepo,
	}
}

type RoleServices struct {
	roleRepo auth_repo_interfaces.IRoleRepository
}

func (r *RoleServices) Get(ctx context.Context, request auth_requests.GetRequest) []*auth_entities.Role {
	return r.roleRepo.Get(ctx, request)
}
func (r *RoleServices) Count(ctx context.Context, request auth_requests.GetRequest) int {
	return r.roleRepo.Count(ctx, request)
}
func (r *RoleServices) Insert(ctx context.Context, role *auth_entities.Role) (*auth_entities.Role, error) {
	err := r.validate(ctx, role)
	if err != nil {
		return &auth_entities.Role{}, err
	}
	return r.roleRepo.Insert(ctx, role)
}
func (r *RoleServices) Update(ctx context.Context, role *auth_entities.Role) (*auth_entities.Role, error) {
	_, err := r.roleRepo.FindById(ctx, role.Id)
	if err != nil {
		return &auth_entities.Role{}, err
	}

	err = r.validate(ctx, role)
	if err != nil {
		return &auth_entities.Role{}, err
	}

	_, err = r.roleRepo.Update(ctx, role)
	if err != nil {
		return &auth_entities.Role{}, err
	}
	return role.Copy(), nil
}
func (r *RoleServices) DeleteById(ctx context.Context, id string) (int, error) {
	_, err := r.roleRepo.FindById(ctx, id)
	if err != nil {
		return 0, err
	}
	return r.roleRepo.DeleteById(ctx, id)
}

func (r *RoleServices) validate(ctx context.Context, role *auth_entities.Role) error {
	v := validator.New()
	data := map[string]interface{}{
		"nama":       role.Nama,
		"deskripsi":  role.Deskripsi,
		"level":      role.Level,
		"permission": role.Permissions,
	}
	rules := map[string]interface{}{
		"nama":       "required",
		"deskripsi":  "",
		"level":      "min=1,max=65535",
		"permission": "",
	}
	err := v.ValidateMap(data, rules)
	if err != nil {
		return err
	}

	count := r.roleRepo.CountByNama(ctx, role.Nama, role.Id)
	if count > 0 {
		return helpers_error.NewValidationErrors("nama", "duplicate", "")
	}
	return nil
}
