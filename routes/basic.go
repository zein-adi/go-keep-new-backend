package routes

import (
	"fmt"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"github.com/zein-adi/go-keep-new-backend/app/components/gorillamux_router"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"net/http"
	"strconv"
	"strings"
)

func StartHttpServer() {
	allowedOrigins := strings.Split(viper.GetString("CORS_ALLOWED_ORIGINS"), ",")
	address := viper.GetString("HTTP_SERVER_ADDRESS")
	port := viper.GetInt("HTTP_SERVER_PORT")
	tlsCertPath := viper.GetString("HTTP_SERVER_TLS_CERT")
	tlsKeyPath := viper.GetString("HTTP_SERVER_TLS_KEY")

	opt := cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"HEAD", "OPTIONS", "GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            false,
	}
	r := gorillamux_router.NewRouter(opt)
	injectRoutes(r)

	fmt.Printf("Listening\nAddress: %s:%d \n...", address, port)
	if port == 443 {
		err := http.ListenAndServeTLS(address+":"+strconv.Itoa(port), tlsCertPath, tlsKeyPath, r)
		helpers_error.PanicIfError(err)
	} else {
		err := http.ListenAndServe(address+":"+strconv.Itoa(port), r)
		helpers_error.PanicIfError(err)
	}
}

func injectRoutes(r *gorillamux_router.Router) {
	injectAuthRoutes(r)
	injectUserRoutes(r)
	injectKeepRoutes(r)
}
