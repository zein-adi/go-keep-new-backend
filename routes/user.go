package routes

import (
	"github.com/zein-adi/go-keep-new-backend/app/components/gorillamux_router"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/dependency_injection"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/handlers/auth_handlers_local"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/handlers/auth_handlers_restful"
)

func injectUserRoutes(r *gorillamux_router.Router) {
	middlewareAcl := middlewares.NewMiddlewareAcl(auth_handlers_local.NewRoleLocalHandler(dependency_injection.InitUserRoleServices()))

	r.Group("/user", "user.", func(r *gorillamux_router.Router) {

		userRestful := auth_handlers_restful.NewUserRestfulHandler(dependency_injection.InitUserUserServices())
		r.Group("/users", "user.", func(r *gorillamux_router.Router) {
			r.GET("", userRestful.Get, "get")
			r.POST("", userRestful.Insert, "insert")
			r.PUT("/{userId:[0-9]+}", userRestful.Update, "update")
			r.PATCH("/{userId:[0-9]+}/password", userRestful.UpdatePassword, "update.password")
			r.DELETE("/{userId:[0-9]+}", userRestful.DeleteById, "delete")
		})

		roleRestful := auth_handlers_restful.NewRoleRestfulHandler(dependency_injection.InitUserRoleServices())
		r.Group("/roles", "role.", func(r *gorillamux_router.Router) {
			r.GET("", roleRestful.Get, "get")
			r.POST("", roleRestful.Insert, "insert")
			r.PUT("/{roleId:[0-9]+}", roleRestful.Update, "update")
			r.DELETE("/{roleId:[0-9]+}", roleRestful.DeleteById, "delete")
		})

		permissionRestful := auth_handlers_restful.NewPermissionRestfulHandler(dependency_injection.InitUserPermissionServices())
		r.GET("/permissions", permissionRestful.Get, "permission.get")

	}).SetMiddleware(middlewares.AuthHandle, middlewareAcl.Handle)
}
