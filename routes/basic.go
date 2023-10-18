package routes

import (
	"fmt"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"github.com/zein-adi/go-keep-new-backend/app/components"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"net/http"
	"strings"
)

func StartHttpServer() {
	allowedOrigins := strings.Split(viper.GetString("CORS_ALLOWED_ORIGINS"), ",")
	addr := viper.GetString("HTTP_SERVER_ADDR")

	opt := cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"HEAD", "OPTIONS", "GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            false,
	}
	r := components.NewRouter(opt)
	injectRoutes(r)

	fmt.Println("Listening " + addr + " ...")
	err := http.ListenAndServe(addr, r)
	helpers_error.PanicIfError(err)
}

func injectRoutes(r *components.Router) {
	injectAuthRoutes(r)
	injectUserRoutes(r)
	injectKeepRoutes(r)
}
