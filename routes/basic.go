package routes

import (
	"fmt"
	"github.com/rs/cors"
	"github.com/zein-adi/go-keep-new-backend/app/components"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"net/http"
)

func StartHttpServer() {
	//TODO get from env
	opt := cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://192.168.232.2:3000"},
		AllowedMethods:   []string{"HEAD", "OPTIONS", "GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            false,
	}
	r := components.NewRouter(opt)
	injectRoutes(r)

	fmt.Println("Listening...")
	err := http.ListenAndServe("0.0.0.0:3001", r)
	helpers_error.PanicIfError(err)
}

func injectRoutes(r *components.Router) {
	injectAuthRoutes(r)
	injectUserRoutes(r)
	injectKeepRoutes(r)
}
