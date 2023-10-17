package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/zein-adi/go-keep-new-backend/app/components"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/dependency_injection"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
)

func injectAuthRoutes(r *components.Router) {
	authRestful := dependency_injection.InitUserAuthRestful()

	r.GET("/", defHandler, "")
	r.Group("/auth", "", func(r *components.Router) {

		r.POST("/login", authRestful.Login, "")

		r.Group("", "", func(r *components.Router) {
			r.GET("/profile", authRestful.Profile, "")
		}).SetMiddleware(middlewares.AuthHandle)

		r.Group("", "", func(r *components.Router) {
			r.POST("/logout", authRestful.Logout, "")
			r.POST("/refresh", authRestful.Refresh, "")
		}).SetMiddleware(middlewares.AuthRefreshHandle)
	})
}
func defHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	helpers_http.SendSingleResponse(w, http.StatusOK, "")
}
