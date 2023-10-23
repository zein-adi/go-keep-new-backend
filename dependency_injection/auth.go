//go:build wireinject
// +build wireinject

package dependency_injection

import (
	"github.com/google/wire"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_services"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/repos/auth_repos_memory"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/repos/auth_repos_mysql"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/repos/auth_repos_redis"
)

var (
	UserRoleSet = wire.NewSet(
		auth_services.NewRoleServices,
		wire.Bind(new(auth_repo_interfaces.IRoleRepository), new(*auth_repos_mysql.RoleMysqlRepository)),
		auth_repos_mysql.NewRoleMysqlRepository,
	)
	UserPermissionSet = wire.NewSet(
		auth_services.NewPermissionServices,
		wire.Bind(new(auth_repo_interfaces.IPermissionRepository), new(*auth_repos_memory.PermissionMemoryRepository)),
		auth_repos_memory.NewPermissionMemoryRepository,
	)
	UserAuthSet = wire.NewSet(
		auth_services.NewAuthServices,
		wire.Bind(new(auth_repo_interfaces.IAuthRepository), new(*auth_repos_redis.AuthMysqlRepository)),
		auth_repos_redis.NewAuthRedisRepository,
	)
	UserUserSet = wire.NewSet(
		auth_services.NewUserServices,
		wire.Bind(new(auth_repo_interfaces.IUserRepository), new(*auth_repos_mysql.UserMysqlRepository)),
		auth_repos_mysql.NewUserMysqlRepository,
	)
)

// User
func InitUserUserMysqlRepository() *auth_repos_mysql.UserMysqlRepository {
	wire.Build(UserUserSet)
	return nil
}
func InitUserUserServices() *auth_services.UserServices {
	wire.Build(UserUserSet, UserRoleSet)
	return nil
}

// Role
func InitUserRoleServices() *auth_services.RoleServices {
	wire.Build(UserRoleSet)
	return nil
}
func InitUserRoleMysqlRepository() *auth_repos_mysql.RoleMysqlRepository {
	wire.Build(UserRoleSet)
	return nil
}

// Permission
func InitUserPermissionMemoryRepository() *auth_repos_memory.PermissionMemoryRepository {
	wire.Build(UserPermissionSet)
	return nil
}
func InitUserPermissionServices() *auth_services.PermissionServices {
	wire.Build(UserPermissionSet, UserRoleSet)
	return nil
}

// Auth
func InitUserAuthServices() *auth_services.AuthServices {
	wire.Build(UserAuthSet, UserUserSet, UserRoleSet)
	return nil
}
