package routes

import (
	"github.com/zein-adi/go-keep-new-backend/app/components"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/dependency_injection"
)

func injectUserRoutes(r *components.Router) {
	middlewareAcl := dependency_injection.InitAclMiddleware()
	roleRestful := dependency_injection.InitUserRoleRestful()
	userRestful := dependency_injection.InitUserUserRestful()
	permissionRestful := dependency_injection.InitUserPermissionRestful()

	r.Group("/user", "user.", func(r *components.Router) {

		r.GET("/users", userRestful.Get, "user.get")
		r.POST("/users", userRestful.Insert, "user.insert")
		r.PUT("/users/:userId", userRestful.Update, "user.update")
		r.PATCH("/users/:userId/password", userRestful.UpdatePassword, "user.update.password")
		r.DELETE("/users/:userId", userRestful.DeleteById, "user.delete")

		r.GET("/roles", roleRestful.Get, "role.get")
		r.POST("/roles", roleRestful.Insert, "role.insert")
		r.PUT("/roles/:roleId", roleRestful.Update, "role.update")
		r.DELETE("/roles/:roleId", roleRestful.DeleteById, "role.delete")

		r.GET("/permissions", permissionRestful.Get, "permission.get")

	}).SetMiddleware(middlewares.AuthHandle, middlewareAcl.Handle)
}
