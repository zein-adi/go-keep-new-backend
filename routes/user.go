package routes

import (
	"github.com/zein-adi/go-keep-new-backend/app/components/gorillamux_router"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/dependency_injection"
)

func injectUserRoutes(r *gorillamux_router.Router) {
	middlewareAcl := dependency_injection.InitAclMiddleware()

	r.Group("/user", "user.", func(r *gorillamux_router.Router) {

		userRestful := dependency_injection.InitUserUserRestful()
		r.Group("/users", "user.", func(r *gorillamux_router.Router) {
			r.GET("", userRestful.Get, "get")
			r.POST("", userRestful.Insert, "insert")
			r.PUT("/{userId:[0-9]+}", userRestful.Update, "update")
			r.PATCH("/{userId:[0-9]+}/password", userRestful.UpdatePassword, "update.password")
			r.DELETE("/{userId:[0-9]+}", userRestful.DeleteById, "delete")
		})

		roleRestful := dependency_injection.InitUserRoleRestful()
		r.Group("/roles", "role.", func(r *gorillamux_router.Router) {
			r.GET("", roleRestful.Get, "get")
			r.POST("", roleRestful.Insert, "insert")
			r.PUT("/{roleId:[0-9]+}", roleRestful.Update, "update")
			r.DELETE("/{roleId:[0-9]+}", roleRestful.DeleteById, "delete")
		})

		permissionRestful := dependency_injection.InitUserPermissionRestful()
		r.GET("/permissions", permissionRestful.Get, "permission.get")

	}).SetMiddleware(middlewares.AuthHandle, middlewareAcl.Handle)
}
