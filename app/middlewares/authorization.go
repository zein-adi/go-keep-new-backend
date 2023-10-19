package middlewares

import (
	"github.com/zein-adi/go-keep-new-backend/domains/auth/handlers/auth_handlers_local/auth_handlers_local_interfaces"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"time"
)

func NewMiddlewareAcl(service auth_handlers_local_interfaces.IRoleLocalHandler) *MiddlewareAcl {
	return &MiddlewareAcl{
		services: service,
	}
}

type MiddlewareAcl struct {
	services    auth_handlers_local_interfaces.IRoleLocalHandler
	roleCache   map[string][]string
	roleCacheAt time.Time
}

func (x *MiddlewareAcl) Handle(writer http.ResponseWriter, request *http.Request, routeName string) bool {
	tokenString, _ := GetAuthorizationToken(request)
	claims, _ := GetJwtClaims(tokenString)

	roleIds := claims.RoleIds
	permissions := x.getPermissionsByRole(roleIds...)

	_, err := helpers.FindIndex(permissions, func(s string) bool {
		return s == routeName
	})
	if err != nil {
		helpers_http.SendErrorResponse(writer, http.StatusForbidden, "")
		return false
	}
	return true
}

func (x *MiddlewareAcl) getPermissionsByRole(roleIds ...string) []string {
	x.loadPermissions()
	var combinedRoles []string
	for _, roleId := range roleIds {
		permissions, ok := x.roleCache[roleId]
		if ok {
			combinedRoles = append(combinedRoles, permissions...)
		}
	}
	combinedRoles = helpers.Unique(combinedRoles)
	return combinedRoles
}
func (x *MiddlewareAcl) loadPermissions() map[string][]string {
	isExpired := time.Now().After(x.roleCacheAt.Add(time.Minute))
	if !isExpired {
		return x.roleCache
	}
	roles := x.services.Get()
	x.roleCacheAt = time.Now()
	x.roleCache = make(map[string][]string, len(roles))
	for _, role := range roles {
		x.roleCache[role.Id] = role.Permissions
	}
	return x.roleCache
}
