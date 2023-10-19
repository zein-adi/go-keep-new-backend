package keep_handlers_restful_interfaces

import (
	"net/http"
)

type IKantongHistoryRestfulHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Insert(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	DeleteById(w http.ResponseWriter, r *http.Request)
}
