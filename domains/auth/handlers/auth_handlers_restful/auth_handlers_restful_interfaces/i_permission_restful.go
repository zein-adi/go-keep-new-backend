package auth_handlers_restful_interfaces

import (
	"net/http"
)

type IPermissionRestfulHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
}
