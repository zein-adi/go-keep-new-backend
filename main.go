package main

import (
	"flag"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/zein-adi/go-keep-new-backend/app/components"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_services"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/handlers/auth_handlers_local"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/handlers/auth_handlers_restful"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/repos/auth_repos_memory"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/repos/auth_repos_mysql"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/repos/auth_repos_redis"
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
	roleRestful := auth_handlers_restful.NewRoleRestful(roleServices)

	userRepo := auth_repos_mysql.NewUserMysqlRepository()
	userServices := auth_services.NewUserServices(userRepo, roleRepo)
	userRestful := auth_handlers_restful.NewUserRestful(userServices)

	permissionRepo := auth_repos_memory.NewPermissionMemoryRepository()
	permissionService := auth_services.NewPermissionServices(permissionRepo, roleRepo)
	permissionRestful := auth_handlers_restful.NewPermissionRestful(permissionService)

	authRepo := auth_repos_redis.NewAuthRedisRepository()
	authService := auth_services.NewAuthServices(authRepo, userRepo, roleRepo)
	authRestful := auth_handlers_restful.NewAuthRestful(authService)

	roleLocal := auth_handlers_local.NewRoleLocal(roleServices)
	middlewareAcl := middlewares.NewMiddlewareAcl(roleLocal)

	opt := cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://192.168.232.2:3000"},
		AllowedMethods:   []string{"HEAD", "OPTIONS", "GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            false,
	}
	r := components.NewRouter(opt)

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

	// Handled by ACL
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

	fmt.Println("Listening...")
	err := http.ListenAndServe("0.0.0.0:3001", r)
	helpers_error.PanicIfError(err)
}

func defHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
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
