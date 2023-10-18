package keep_handlers_restful_interfaces

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type IKantongHistoryRestfulHandler interface {
	Get(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Insert(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Update(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	DeleteById(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}
