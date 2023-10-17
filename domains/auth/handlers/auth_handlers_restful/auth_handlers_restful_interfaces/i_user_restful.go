package auth_handlers_restful_interfaces

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type IUserRestfulHandler interface {
	Get(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Insert(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Update(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	UpdatePassword(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	DeleteById(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}
