package keep_handlers_restful_interfaces

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type IPosRestfulHandler interface {
	Get(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Insert(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Update(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	DeleteById(w http.ResponseWriter, r *http.Request, p httprouter.Params)

	GetTrashed(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	RestoreTrashedById(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	DeleteTrashedById(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}
