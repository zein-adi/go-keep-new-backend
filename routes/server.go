package routes

import (
	"fmt"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"github.com/zein-adi/go-keep-new-backend/app/components/gorillamux_router"
	"github.com/zein-adi/go-keep-new-backend/dependency_injection"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"net/http"
	"strconv"
	"strings"
)

func StartHttpServer() {
	RegisterListeners()
	startServer()
}

func startServer() {
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

	serverType := ""
	if port == 443 {
		serverType = "HTTPS"
	} else {
		serverType = "HTTP"
	}
	fmt.Printf("%-25s:\n", serverType)
	fmt.Printf("%-25s: %s:%d\n", "Address", address, port)
	fmt.Printf("%-25s: %s\n", "Cors Allowed Origins", allowedOrigins)
	fmt.Printf("%-25s:\n", "Registered Routes")

	r := gorillamux_router.NewRouter(opt)
	injectRoutes(r)

	fmt.Printf("Listening ...")
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
	injectBasicRoutes(r)
}

func RegisterListeners() {
	RegisterKeepListeners(
		dependency_injection.InitKeepPosServices(),
		dependency_injection.InitKeepKantongServices(),
		dependency_injection.InitKeepLokasiServices(),
		dependency_injection.InitKeepBarangServices(),
	)
}
