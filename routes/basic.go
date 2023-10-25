package routes

import (
	"github.com/zein-adi/go-keep-new-backend/app/components/gorillamux_router"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/dependency_injection"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/handlers/auth_handlers_local"
	"github.com/zein-adi/go-keep-new-backend/domains/basic/handlers/basic_handlers_restful"
)

func injectBasicRoutes(r *gorillamux_router.Router) {
	middlewareAcl := middlewares.NewMiddlewareAcl(auth_handlers_local.NewRoleLocalHandler(dependency_injection.InitUserRoleServices()))
	pos := basic_handlers_restful.NewChangelogRestfulHandler(dependency_injection.InitBasicChangelogServices())

	// Bisa dilihat public / dimunculkan sebelum login
	r.Group("/changelogs", "changelog.", func(r *gorillamux_router.Router) {
		r.GET("", pos.Get, "get")
	}) //.SetMiddleware(middlewares.AuthHandle)

	r.Group("/changelogs", "changelog.", func(r *gorillamux_router.Router) {
		r.POST("", pos.Insert, "insert")
		r.PATCH("/{changelogId:[0-9]+}", pos.Update, "update")
		r.DELETE("/{changelogId:[0-9]+}", pos.DeleteById, "delete")
	}).SetMiddleware(middlewares.AuthHandle, middlewareAcl.Handle)
}
