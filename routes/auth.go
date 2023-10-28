package routes

import (
	"github.com/zein-adi/go-keep-new-backend/app/components/gorillamux_router"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/dependency_injection"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/handlers/auth_handlers_restful"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
)

func injectAuthRoutes(r *gorillamux_router.Router) {
	authRestful := auth_handlers_restful.NewAuthRestfulHandler(dependency_injection.InitUserAuthServices())
	r.GET("/", defHandler, "")
	r.Group("/auth", "", func(r *gorillamux_router.Router) {

		r.POST("/login", authRestful.Login, "")

		r.New().SetMiddleware(middlewares.AuthHandle).
			Group("", "", func(r *gorillamux_router.Router) {
				r.GET("/profile", authRestful.Profile, "")
			})

		r.New().SetMiddleware(middlewares.AuthRefreshHandle).
			Group("", "", func(r *gorillamux_router.Router) {
				r.POST("/logout", authRestful.Logout, "")
				r.POST("/refresh", authRestful.Refresh, "")
			})
	})
}
func defHandler(w http.ResponseWriter, _ *http.Request) {
	helpers_http.SendSingleResponse(w, http.StatusOK, "")
}
