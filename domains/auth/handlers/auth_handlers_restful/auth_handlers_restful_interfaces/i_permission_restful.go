package auth_handlers_restful_interfaces

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type IPermissionRestfulHandler interface {
	Get(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}
