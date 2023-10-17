//go:build wireinject
// +build wireinject

package dependency_injection

import (
	"github.com/google/wire"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_services"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/handlers/auth_handlers_local"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/handlers/auth_handlers_local/auth_handlers_local_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/handlers/auth_handlers_restful"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/handlers/auth_handlers_restful/auth_handlers_restful_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/repos/auth_repos_memory"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/repos/auth_repos_mysql"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/repos/auth_repos_redis"
)

var (
	UserRoleSet = wire.NewSet(
		auth_repos_mysql.NewRoleMysqlRepository,
		wire.Bind(new(auth_repo_interfaces.IRoleRepository), new(*auth_repos_mysql.RoleMysqlRepository)),
		auth_services.NewRoleServices,
		wire.Bind(new(auth_service_interfaces.IRoleServices), new(*auth_services.RoleServices)),
		auth_handlers_restful.NewRoleRestfulHandler,
		wire.Bind(new(auth_handlers_restful_interfaces.IRoleRestfulHandler), new(*auth_handlers_restful.RoleRestfulHandler)),
		auth_handlers_local.NewRoleLocalHandler,
		wire.Bind(new(auth_handlers_local_interfaces.IRoleLocalHandler), new(*auth_handlers_local.RoleLocalHandler)),
		middlewares.NewMiddlewareAcl,
	)
	UserPermissionSet = wire.NewSet(
		auth_repos_memory.NewPermissionMemoryRepository,
		wire.Bind(new(auth_repo_interfaces.IPermissionRepository), new(*auth_repos_memory.PermissionMemoryRepository)),
		auth_services.NewPermissionServices,
		wire.Bind(new(auth_service_interfaces.IPermissionServices), new(*auth_services.PermissionServices)),
		auth_handlers_restful.NewPermissionRestfulHandler,
		wire.Bind(new(auth_handlers_restful_interfaces.IPermissionRestfulHandler), new(*auth_handlers_restful.PermissionRestfulHandler)),
	)
	UserAuthSet = wire.NewSet(
		auth_repos_redis.NewAuthRedisRepository,
		wire.Bind(new(auth_repo_interfaces.IAuthRepository), new(*auth_repos_redis.AuthMysqlRepository)),
		auth_services.NewAuthServices,
		wire.Bind(new(auth_service_interfaces.IAuthServices), new(*auth_services.AuthServices)),
		auth_handlers_restful.NewAuthRestfulHandler,
		wire.Bind(new(auth_handlers_restful_interfaces.IAuthRestfulHandler), new(*auth_handlers_restful.AuthRestfulHandler)),
	)
	UserUserSet = wire.NewSet(
		auth_repos_mysql.NewUserMysqlRepository,
		wire.Bind(new(auth_repo_interfaces.IUserRepository), new(*auth_repos_mysql.UserMysqlRepository)),
		auth_services.NewUserServices,
		wire.Bind(new(auth_service_interfaces.IUserServices), new(*auth_services.UserServices)),
		auth_handlers_restful.NewUserRestfulHandler,
		wire.Bind(new(auth_handlers_restful_interfaces.IUserRestfulHandler), new(*auth_handlers_restful.UserRestfulHandler)),
	)
)

// User
func InitUserUserMysqlRepository() *auth_repos_mysql.UserMysqlRepository {
	wire.Build(UserUserSet)
	return nil
}
func InitUserUserRestful() *auth_handlers_restful.UserRestfulHandler {
	wire.Build(UserUserSet, UserRoleSet)
	return nil
}

// Role
func InitUserRoleRestful() *auth_handlers_restful.RoleRestfulHandler {
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
func InitUserPermissionRestful() *auth_handlers_restful.PermissionRestfulHandler {
	wire.Build(UserPermissionSet, UserRoleSet)
	return nil
}

// Auth
func InitUserAuthRestful() *auth_handlers_restful.AuthRestfulHandler {
	wire.Build(UserAuthSet, UserUserSet, UserRoleSet)
	return nil
}
func InitAclMiddleware() *middlewares.MiddlewareAcl {
	wire.Build(UserRoleSet)
	return nil
}
