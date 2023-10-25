package auth_services

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_requests"
	"github.com/zein-adi/go-keep-new-backend/helpers/validator"
)

var (
	RoleAccessUnauthorizedError = errors.New("role access unauthorized")
)

func NewRoleServices(roleRepo auth_repo_interfaces.IRoleRepository) *RoleServices {
	return &RoleServices{
		roleRepo: roleRepo,
	}
}

type RoleServices struct {
	roleRepo auth_repo_interfaces.IRoleRepository
}

func (x *RoleServices) Get(ctx context.Context, request *helpers_requests.Get) []*auth_entities.Role {
	return x.roleRepo.Get(ctx, request)
}
func (x *RoleServices) Count(ctx context.Context, request *helpers_requests.Get) int {
	return x.roleRepo.Count(ctx, request)
}
func (x *RoleServices) Insert(ctx context.Context, role *auth_entities.Role, currentUserRoleIds []string) (*auth_entities.Role, error) {
	err := x.validate(ctx, role)
	if err != nil {
		return &auth_entities.Role{}, err
	}
	currentUserRoles, err := x.roleRepo.GetById(ctx, currentUserRoleIds)
	if err != nil {
		return &auth_entities.Role{}, err
	}
	err = x.validateLevelAndPermission(role, currentUserRoles)
	if err != nil {
		return &auth_entities.Role{}, err
	}

	return x.roleRepo.Insert(ctx, role)
}
func (x *RoleServices) Update(ctx context.Context, role *auth_entities.Role, currentUserRoleIds []string) (*auth_entities.Role, error) {
	wanToUpdateModel, err := x.roleRepo.FindById(ctx, role.Id)
	if err != nil {
		return &auth_entities.Role{}, err
	}

	// Validate Patch
	err = x.validate(ctx, role)
	if err != nil {
		return &auth_entities.Role{}, err
	}

	isNotDeveloper := currentUserRoleIds[0] != "1"
	if isNotDeveloper {
		currentUserRoles, err := x.roleRepo.GetById(ctx, currentUserRoleIds)
		if err != nil {
			return &auth_entities.Role{}, err
		}
		err = x.validateLevelAndPermission(role, currentUserRoles)
		if err != nil {
			return nil, err
		}
		// Validate Current User Access To Updated Model
		userMinLevel := getMinimumLevelFromRoles(currentUserRoles)
		err = validateAccessedLevel(userMinLevel, wanToUpdateModel.Level)
		if err != nil {
			return &auth_entities.Role{}, err
		}
	}

	// Process
	_, err = x.roleRepo.Update(ctx, role)
	if err != nil {
		return &auth_entities.Role{}, err
	}
	return role.Copy(), nil
}
func (x *RoleServices) DeleteById(ctx context.Context, id string, currentUserRoleIds []string) (affected int, err error) {
	wantToDeleteModel, err := x.roleRepo.FindById(ctx, id)
	if err != nil {
		return 0, err
	}

	// Validate Current User Access To Updated Model
	currentUserRoles, err := x.roleRepo.GetById(ctx, currentUserRoleIds)
	if err != nil {
		return 0, err
	}
	userMinLevel := getMinimumLevelFromRoles(currentUserRoles)
	err = validateAccessedLevel(userMinLevel, wantToDeleteModel.Level)
	if err != nil {
		return 0, err
	}

	// Process
	return x.roleRepo.DeleteById(ctx, id)
}

func (x *RoleServices) validate(ctx context.Context, role *auth_entities.Role) error {
	v := validator.New()
	err := v.ValidateStruct(role)
	if err != nil {
		return err
	}
	err = x.validateRoleName(ctx, role)
	if err != nil {
		return err
	}
	return nil
}
func (x *RoleServices) validateLevelAndPermission(role *auth_entities.Role, currentUserRoles []*auth_entities.Role) error {
	userMinLevel := getMinimumLevelFromRoles(currentUserRoles)
	if role.Level <= userMinLevel {
		return helpers_error.NewValidationErrors("level", "invalid", "")
	}
	err := x.validatePermission(role.Permissions, currentUserRoles)
	if err != nil {
		return err
	}
	return nil
}
func (x *RoleServices) validateRoleName(ctx context.Context, role *auth_entities.Role) error {
	count := x.roleRepo.CountByNama(ctx, role.Nama, role.Id)
	if count > 0 {
		return helpers_error.NewValidationErrors("nama", "duplicate", "")
	}
	return nil
}
func (x *RoleServices) validatePermission(requestedPermissions []string, currentUserRoles []*auth_entities.Role) error {
	var currentUserPermissions []string
	for _, userRole := range currentUserRoles {
		currentUserPermissions = append(currentUserPermissions, userRole.Permissions...)
	}
	currentUserPermissions = helpers.Unique(currentUserPermissions)
	currentUserPermissionsMap := helpers.KeyByMap(
		currentUserPermissions,
		func(permission string) string {
			return permission
		},
		func(permission string) bool {
			return true
		})

	for _, permission := range requestedPermissions {
		_, ok := currentUserPermissionsMap[permission]
		if !ok {
			return helpers_error.NewValidationErrors("permissions", "invalid", "")
		}
	}
	return nil
}

func validateAccessedLevel(userMinLevel, accessedMinLevel int) error {
	if accessedMinLevel <= userMinLevel {
		return errors.Wrapf(RoleAccessUnauthorizedError,
			"role level %d want to update/delete role with level %d you can only delete level higher than your level",
			userMinLevel,
			accessedMinLevel,
		)
	}
	return nil
}
func getMinimumLevelFromRoles(currentUserRoles []*auth_entities.Role) int {
	userMinLevel := helpers.Reduce(currentUserRoles, 10000, func(currentSmallest int, role *auth_entities.Role) int {
		if role.Level < currentSmallest {
			return role.Level
		}
		return currentSmallest
	})
	return userMinLevel
}
