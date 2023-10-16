package auth_handlers_restful_interfaces

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type IAuthRestful interface {
	Login(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Refresh(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Profile(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}
