package main

import (
	"flag"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/zein-adi/go-keep-new-backend/app/components"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_services"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/handlers/auth_handlers_restful"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/repos/auth_repos_mysql"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"strings"
)

func main() {
	cliHandler()
}

func startHttpServer() {
	roleRepo := auth_repos_mysql.NewRoleMysqlRepository()
	roleServices := auth_services.NewRoleServices(roleRepo)
	role := auth_handlers_restful.NewRoleRestful(roleServices)

	opt := cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://192.168.232.2:3000"},
		AllowedMethods:   []string{"HEAD", "OPTIONS", "GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            false,
	}
	router := components.NewRouter(opt)

	router.GET("/", sh, "")
	router.Group("/auth", "auth", func(router *components.Router) {
		router.GET("/login", sh, "")
		router.POST("/logout", sh, "")
		router.Group("", "", func(router *components.Router) {
			router.POST("/refresh", sh, "")
		}).SetMiddleware(middlewares.AuthRefresh)
	})
	router.Group("/user", "user.", func(router *components.Router) {
		router.GET("/roles", role.Get, "role.get")
		router.POST("/roles", role.Insert, "role.insert")
		router.PUT("/roles/:roleId", role.Update, "role.update")
		router.DELETE("/roles/:roleId", role.DeleteById, "role.delete")
	}) // .SetMiddleware(middlewares.Auth, middlewares.Acl)

	err := http.ListenAndServe("0.0.0.0:3001", router)
	helpers_error.PanicIfError(err)
}

func sh(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	helpers_http.SendSingleResponse(w, http.StatusOK, "")
}

func cliHandler() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		startHttpServer()
	} else {
		action := strings.ToLower(flag.Arg(0))
		if action == "migrate" {
			RunMigration()
		}
	}
}
