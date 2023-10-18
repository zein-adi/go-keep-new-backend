package routes

import (
	"github.com/zein-adi/go-keep-new-backend/app/components"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/dependency_injection"
)

func injectUserRoutes(r *components.Router) {
	middlewareAcl := dependency_injection.InitAclMiddleware()

	r.Group("/user", "user.", func(r *components.Router) {

		userRestful := dependency_injection.InitUserUserRestful()
		r.Group("/users", "user.", func(r *components.Router) {
			r.GET("/", userRestful.Get, "get")
			r.POST("/", userRestful.Insert, "insert")
			r.PUT("/:userId/", userRestful.Update, "update")
			r.PATCH("/:userId/password/", userRestful.UpdatePassword, "update.password")
			r.DELETE("/:userId/", userRestful.DeleteById, "delete")
		})

		roleRestful := dependency_injection.InitUserRoleRestful()
		r.Group("/roles", "role.", func(r *components.Router) {
			r.GET("/", roleRestful.Get, "get")
			r.POST("/", roleRestful.Insert, "insert")
			r.PUT("/:roleId/", roleRestful.Update, "update")
			r.DELETE("/:roleId/", roleRestful.DeleteById, "delete")
		})

		permissionRestful := dependency_injection.InitUserPermissionRestful()
		r.GET("/permissions", permissionRestful.Get, "permission.get")

	}).SetMiddleware(middlewares.AuthHandle, middlewareAcl.Handle)
}
